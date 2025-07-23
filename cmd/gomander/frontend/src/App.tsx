import { AppSidebar } from "@/components/layout/AppSidebar/AppSidebar.tsx";
import { SidebarProvider } from "@/components/ui/sidebar.tsx";
import { Toaster } from "@/components/ui/sonner.tsx";
import { EventListenersContainer } from "@/components/utility/EventListenersContainer.tsx";
import { DataContextProvider } from "@/contexts/DataContext.tsx";
import { UserConfigDataContextProvider } from "@/contexts/UserConfigDataContext.tsx";

import { LogsScreen } from "./screens/LogsScreen.tsx";

function App() {
  return (
    <DataContextProvider>
      <UserConfigDataContextProvider>
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
      </UserConfigDataContextProvider>
    </DataContextProvider>
  );
}

export default App;
