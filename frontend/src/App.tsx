import { AppSidebar } from "@/components/layout/AppSidebar/AppSidebar.tsx";
import { SidebarProvider } from "@/components/ui/sidebar.tsx";
import { Toaster } from "@/components/ui/sonner.tsx";
import { DataContextProvider } from "@/contexts/DataContext.tsx";
import { ModalsContextProvider } from "@/contexts/ModalsContext.tsx";

import { MainScreen } from "./screens/MainScreen";

function App() {
  return (
    <DataContextProvider>
      <ModalsContextProvider>
        <SidebarProvider>
          <nav>
            <AppSidebar />
          </nav>
          <main className="w-full h-full bg-white">
            <MainScreen />
          </main>
          <Toaster />
        </SidebarProvider>
      </ModalsContextProvider>
    </DataContextProvider>
  );
}

export default App;
