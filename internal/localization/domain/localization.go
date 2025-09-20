package domain

type Localization struct {
	// PoC: Only a few keys for testing
	SidebarTitle         string `json:"sidebar.title"`
	SidebarMinimize      string `json:"sidebar.minimize"`
	SidebarCommandsTitle string `json:"sidebar.commands.title"`

	ActionsCancel string `json:"actions.cancel"`
	ActionsSave   string `json:"actions.save"`

	ToastTestSuccess string `json:"toast.test.success"`
	ToastTestError   string `json:"toast.test.error"`
}
