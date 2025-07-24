import { AppSidebar } from "@/components/layout/AppSidebar/AppSidebar.tsx";
import { SidebarProvider } from "@/components/ui/sidebar.tsx";
import { Toaster } from "@/components/ui/sonner.tsx";
import { EventListenersContainer } from "@/components/utility/EventListenersContainer.tsx";
import { useFetchInitialData } from "@/hooks/useFetchInitialData.ts";

import { LogsScreen } from "./screens/LogsScreen.tsx";

function App() {
  useFetchInitialData();

  return (
    <>
      <EventListenersContainer />
      <SidebarProvider>
        <nav>
          <AppSidebar />
        </nav>
        <main className="w-full h-screen bg-white">
          <LogsScreen />
        </main>
        <Toaster richColors />
      </SidebarProvider>
    </>
  );
}

export default App;
