import { useEffect } from "react";

import { AppSidebar } from "@/components/layout/AppSidebar/AppSidebar.tsx";
import { SidebarProvider } from "@/components/ui/sidebar.tsx";
import { Toaster } from "@/components/ui/sonner.tsx";
import { EventListenersContainer } from "@/components/utility/EventListenersContainer.tsx";
import { fetchProject } from "@/queries/fetchProject.ts";
import { fetchUserConfig } from "@/queries/fetchUserConfig.ts";
import { ProjectSelectionScreen } from "@/screens/ProjectSelectionScreen/ProjectSelectionScreen.tsx";
import { useProjectStore } from "@/store/projectStore.ts";

import { LogsScreen } from "./screens/LogsScreen/LogsScreen.tsx";

function App() {
  const project = useProjectStore((state) => state.projectInfo);

  useEffect(() => {
    fetchProject();
    fetchUserConfig();
  }, []);

  return (
    <>
      <EventListenersContainer />
      <Toaster richColors />
      {project && (
        <SidebarProvider>
          <nav>
            <AppSidebar />
          </nav>

          <main className="w-full h-screen bg-white">
            <LogsScreen />
          </main>
        </SidebarProvider>
      )}
      {!project && <ProjectSelectionScreen />}
    </>
  );
}

export default App;
