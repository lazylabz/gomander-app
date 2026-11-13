// Controls sit above the row's stretched button as siblings: nested inside it,
// they take no focus and screen readers may drop them.
export const stretchedButtonClassName =
	"outline-hidden after:absolute after:inset-0 after:rounded-md after:ring-sidebar-ring focus-visible:after:ring-2";

export const sidebarControlClassName =
	"relative z-10 flex items-center justify-center rounded-sm text-muted-foreground cursor-pointer hover:text-primary outline-hidden ring-sidebar-ring focus-visible:ring-2";
