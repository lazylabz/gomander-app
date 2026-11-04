// Package notifications sends OS notifications through the Wails service.
//
// That service refuses to start without a macOS bundle identifier or a Linux
// notification daemon, and registered directly the refusal would abort the app,
// so WailsNotifier owns its lifecycle and only logs it.
package notifications

import (
	"context"
	"sync/atomic"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/services/notifications"

	"gomander/internal/command/application/handlers"
)

type Logger interface {
	Error(message string)
}

// notificationCentre is the part of the Wails service the notifier drives, so
// its lifecycle can be tested without an OS notification centre.
type notificationCentre interface {
	ServiceStartup(ctx context.Context, options application.ServiceOptions) error
	ServiceShutdown() error
	RequestNotificationAuthorization() (bool, error)
	SendNotification(options notifications.NotificationOptions) error
}

type WailsNotifier struct {
	service notificationCentre
	logger  Logger
	started atomic.Bool
}

func NewWailsNotifier(logger Logger) *WailsNotifier {
	return newWailsNotifier(notifications.New(), logger)
}

func newWailsNotifier(service notificationCentre, logger Logger) *WailsNotifier {
	return &WailsNotifier{service: service, logger: logger}
}

func (n *WailsNotifier) Notify(notification handlers.Notification) error {
	// Why it never started was logged once at startup, not on every failure
	if !n.started.Load() {
		return nil
	}

	return n.service.SendNotification(notifications.NotificationOptions{
		ID:    notification.Id,
		Title: notification.Title,
		Body:  notification.Body,
	})
}

// Lifecycle is the Wails service that starts the notifier with the app. It is a
// type of its own so the frontend bindings expose none of the notifier.
func (n *WailsNotifier) Lifecycle() *Lifecycle {
	return &Lifecycle{notifier: n}
}

type Lifecycle struct {
	notifier *WailsNotifier
}

func (l *Lifecycle) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	n := l.notifier

	if err := n.service.ServiceStartup(ctx, options); err != nil {
		n.logger.Error("[ERROR - Starting notifications]: " + err.Error())
		return nil
	}
	n.started.Store(true)

	// On macOS this waits on the user's answer, so it must not hold up startup.
	go func() {
		if _, err := n.service.RequestNotificationAuthorization(); err != nil {
			n.logger.Error("[ERROR - Requesting notification authorization]: " + err.Error())
		}
	}()

	return nil
}

func (l *Lifecycle) ServiceShutdown() error {
	if !l.notifier.started.Load() {
		return nil
	}

	return l.notifier.service.ServiceShutdown()
}
