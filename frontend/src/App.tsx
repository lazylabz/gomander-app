import { AppSidebar } from "@/components/layout/AppSidebar/AppSidebar.tsx";
import { SidebarProvider } from "@/components/ui/sidebar.tsx";
import { Toaster } from "@/components/ui/sonner.tsx";
import { DataContextProvider } from "@/contexts/DataContext.tsx";

import { MainScreen } from "./screens/MainScreen";

function App() {
  return (
    <DataContextProvider>
      <SidebarProvider>
        <nav>
          <AppSidebar />
        </nav>
        <main className="w-full h-screen bg-white">
          <MainScreen />
        </main>
        <Toaster />
      </SidebarProvider>
    </DataContextProvider>
  );
}

export default App;
