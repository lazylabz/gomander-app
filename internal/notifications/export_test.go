package notifications

type NotificationCentre = notificationCentre

func NewWailsNotifierOver(service NotificationCentre, logger Logger) *WailsNotifier {
	return newWailsNotifier(service, logger)
}
