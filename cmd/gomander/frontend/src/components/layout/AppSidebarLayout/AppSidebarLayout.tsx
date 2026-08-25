import { AppSidebar } from "@/components/layout/AppSidebarLayout/components/AppSidebar/AppSidebar.tsx";
import { RunningIndicator } from "@/components/layout/AppSidebarLayout/components/AppSidebar/components/RunningIndicator/RunningIndicator.tsx";
import {
	SidebarProvider,
	useSidebar,
} from "@/design-system/components/ui/sidebar.tsx";

// Shown only when the sidebar is collapsed, since the footer indicator
// disappears together with the sidebar.
const CollapsedRunningIndicator = () => {
	const { open, openMobile, isMobile } = useSidebar();
	const isSidebarOpen = isMobile ? openMobile : open;

	if (isSidebarOpen) {
		return null;
	}

	return (
		<RunningIndicator className="fixed bottom-3 left-3 z-10 opacity-25 transition-opacity hover:opacity-100" />
	);
};

export const AppSidebarLayout = ({
	children,
}: {
	children: React.ReactNode;
}) => {
	return (
		<SidebarProvider>
			<nav>
				<AppSidebar />
			</nav>

			<main className="w-full h-screen bg-white">{children}</main>

			<CollapsedRunningIndicator />
		</SidebarProvider>
	);
};
