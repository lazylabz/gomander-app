package notifications_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wailsapp/wails/v3/pkg/application"
	wailsnotifications "github.com/wailsapp/wails/v3/pkg/services/notifications"

	"gomander/internal/command/application/handlers"
	"gomander/internal/notifications"
)

type notificationCentreFake struct {
	mutex                sync.Mutex
	startupFailure       error
	authorizationFailure error
	sent                 []wailsnotifications.NotificationOptions
	authorizationAsked   bool
	shutDown             bool
}

func (f *notificationCentreFake) ServiceStartup(context.Context, application.ServiceOptions) error {
	return f.startupFailure
}

func (f *notificationCentreFake) ServiceShutdown() error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.shutDown = true

	return nil
}

func (f *notificationCentreFake) RequestNotificationAuthorization() (bool, error) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.authorizationAsked = true

	return f.authorizationFailure == nil, f.authorizationFailure
}

func (f *notificationCentreFake) SendNotification(options wailsnotifications.NotificationOptions) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.sent = append(f.sent, options)

	return nil
}

func (f *notificationCentreFake) state() (sent []wailsnotifications.NotificationOptions, authorizationAsked, shutDown bool) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	return append([]wailsnotifications.NotificationOptions(nil), f.sent...), f.authorizationAsked, f.shutDown
}

type loggerFake struct {
	mutex  sync.Mutex
	errors []string
}

func (l *loggerFake) Error(message string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.errors = append(l.errors, message)
}

func (l *loggerFake) logged() []string {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	return append([]string(nil), l.errors...)
}

var notification = handlers.Notification{Id: "error-pattern-1", Title: "api reported an error", Body: "ERROR: boom"}

func TestWailsNotifier(t *testing.T) {
	t.Run("Should send through the notification centre once it has started", func(t *testing.T) {
		// Arrange
		centre := &notificationCentreFake{}
		sut := notifications.NewWailsNotifierOver(centre, &loggerFake{})
		assert.NoError(t, sut.Lifecycle().ServiceStartup(context.Background(), application.ServiceOptions{}))

		// Act
		err := sut.Notify(notification)

		// Assert
		assert.NoError(t, err)
		sent, _, _ := centre.state()
		assert.Equal(t, []wailsnotifications.NotificationOptions{{
			ID:    "error-pattern-1",
			Title: "api reported an error",
			Body:  "ERROR: boom",
		}}, sent)
	})

	t.Run("Should ask for permission to notify when it starts", func(t *testing.T) {
		// Arrange
		centre := &notificationCentreFake{}
		sut := notifications.NewWailsNotifierOver(centre, &loggerFake{})

		// Act
		assert.NoError(t, sut.Lifecycle().ServiceStartup(context.Background(), application.ServiceOptions{}))

		// Assert
		assert.Eventually(t, func() bool {
			_, asked, _ := centre.state()
			return asked
		}, time.Second, 10*time.Millisecond)
	})

	t.Run("Should log a refused permission and keep running", func(t *testing.T) {
		// Arrange
		centre := &notificationCentreFake{authorizationFailure: errors.New("denied")}
		logger := &loggerFake{}
		sut := notifications.NewWailsNotifierOver(centre, logger)

		// Act
		err := sut.Lifecycle().ServiceStartup(context.Background(), application.ServiceOptions{})

		// Assert
		assert.NoError(t, err)
		assert.Eventually(t, func() bool { return len(logger.logged()) == 1 }, time.Second, 10*time.Millisecond)
	})

	t.Run("Should let the app start when the notification centre refuses to, and then notify nothing", func(t *testing.T) {
		// Arrange
		centre := &notificationCentreFake{startupFailure: errors.New("notifications require a valid bundle identifier")}
		logger := &loggerFake{}
		sut := notifications.NewWailsNotifierOver(centre, logger)

		// Act
		startupErr := sut.Lifecycle().ServiceStartup(context.Background(), application.ServiceOptions{})
		notifyErr := sut.Notify(notification)
		shutdownErr := sut.Lifecycle().ServiceShutdown()

		// Assert
		assert.NoError(t, startupErr)
		assert.NoError(t, notifyErr)
		assert.NoError(t, shutdownErr)
		assert.Equal(t, []string{"[ERROR - Starting notifications]: notifications require a valid bundle identifier"}, logger.logged())
		sent, asked, shutDown := centre.state()
		assert.Empty(t, sent)
		assert.False(t, asked)
		assert.False(t, shutDown, "a centre that never started is not shut down")
	})

	t.Run("Should shut the notification centre down with the app", func(t *testing.T) {
		// Arrange
		centre := &notificationCentreFake{}
		sut := notifications.NewWailsNotifierOver(centre, &loggerFake{})
		assert.NoError(t, sut.Lifecycle().ServiceStartup(context.Background(), application.ServiceOptions{}))

		// Act
		err := sut.Lifecycle().ServiceShutdown()

		// Assert
		assert.NoError(t, err)
		_, _, shutDown := centre.state()
		assert.True(t, shutDown)
	})
}
