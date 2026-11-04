package handlers

import (
	"regexp"
	"strings"

	commanddomain "gomander/internal/command/domain"
	commandevent "gomander/internal/command/domain/event"
	configdomain "gomander/internal/config/domain"
	"gomander/internal/eventbus"
	"gomander/internal/localization"
)

const (
	fallbackLocale = "en"
	maxBodyRunes   = 200
)

// The runner forces colour on, so output lines carry ANSI escapes that would
// show up raw in a notification.
var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;?]*[ -/]*[@-~]`)

type Notification struct {
	Id    string
	Title string
	Body  string
}

type Notifier interface {
	Notify(notification Notification) error
}

type Window interface {
	IsFocused() bool
}

// NotifyOnErrorPatternMatched tells the user a Command reported an error when
// they are not looking at Gomander: with the window focused the logs show it.
type NotifyOnErrorPatternMatched struct {
	commandRepository commanddomain.Repository
	configRepository  configdomain.Repository
	catalogue         localization.Catalogue
	window            Window
	notifier          Notifier
}

func NewNotifyOnErrorPatternMatched(
	commandRepository commanddomain.Repository,
	configRepository configdomain.Repository,
	catalogue localization.Catalogue,
	window Window,
	notifier Notifier,
) *NotifyOnErrorPatternMatched {
	return &NotifyOnErrorPatternMatched{
		commandRepository: commandRepository,
		configRepository:  configRepository,
		catalogue:         catalogue,
		window:            window,
		notifier:          notifier,
	}
}

func (h *NotifyOnErrorPatternMatched) GetEvent() eventbus.Event {
	return commandevent.ErrorPatternMatchedEvent{}
}

func (h *NotifyOnErrorPatternMatched) Execute(e eventbus.Event) error {
	match, ok := e.(commandevent.ErrorPatternMatchedEvent)
	if !ok {
		return nil
	}

	config, err := h.configRepository.GetOrCreate()
	if err != nil {
		return err
	}

	if !config.NotificationsEnabled || h.window.IsFocused() {
		return nil
	}

	command, err := h.commandRepository.Get(match.CommandId)
	if err != nil {
		return err
	}

	texts, err := h.translations(config.Locale)
	if err != nil {
		return err
	}

	return h.notifier.Notify(Notification{
		// One per Command, so each run replaces the notification of the last
		Id:    "error-pattern-" + command.Id,
		Title: strings.ReplaceAll(texts.NotificationsErrorPatternTitle, "{{name}}", command.Name),
		Body:  notificationBody(match.Line),
	})
}

// translations falls back to English the way the frontend does, so a locale
// never saved yet still reads.
func (h *NotifyOnErrorPatternMatched) translations(locale string) (*localization.Localization, error) {
	if locale != "" {
		if texts, err := h.catalogue.Translations(locale); err == nil {
			return texts, nil
		}
	}

	return h.catalogue.Translations(fallbackLocale)
}

func notificationBody(line string) string {
	body := []rune(strings.TrimSpace(ansiEscape.ReplaceAllString(line, "")))
	if len(body) <= maxBodyRunes {
		return string(body)
	}

	return string(body[:maxBodyRunes-1]) + "…"
}
