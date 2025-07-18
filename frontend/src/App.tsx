import { AppSidebar } from "@/components/layout/AppSidebar/AppSidebar.tsx";
import { SidebarProvider } from "@/components/ui/sidebar.tsx";
import { Toaster } from "@/components/ui/sonner.tsx";
import { DataContextProvider } from "@/contexts/DataContext.tsx";

import { LogsScreen } from "./screens/LogsScreen.tsx";

function App() {
  return (
    <DataContextProvider>
      <SidebarProvider>
        <nav>
          <AppSidebar />
        </nav>
        <main className="w-full h-screen bg-white">
          <LogsScreen />
        </main>
        <Toaster />
      </SidebarProvider>
    </DataContextProvider>
  );
}

export default App;
