import { AppSidebar } from "@/components/layout/AppSidebarLayout/components/AppSidebar/AppSidebar.tsx";
import { SidebarProvider } from "@/components/ui/sidebar.tsx";

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
    </SidebarProvider>
  );
};
