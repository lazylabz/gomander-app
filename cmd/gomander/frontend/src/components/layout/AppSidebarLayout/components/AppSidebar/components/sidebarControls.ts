// A row's own button stretches over the whole row through ::after, and its
// controls sit above it as sibling buttons: a control nested inside the row's
// button would take no focus of its own and could be hidden from screen readers.
export const stretchedButtonClassName =
	"outline-hidden after:absolute after:inset-0 after:rounded-md after:ring-sidebar-ring focus-visible:after:ring-2";

export const sidebarControlClassName =
	"relative z-10 flex items-center justify-center rounded-sm text-muted-foreground cursor-pointer hover:text-primary outline-hidden ring-sidebar-ring focus-visible:ring-2";
