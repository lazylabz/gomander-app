import { useEffect } from "react";
import { util } from "zod/v3";

import { Toaster } from "@/components/ui/sonner.tsx";
import { EventListenersContainer } from "@/components/utility/EventListenersContainer.tsx";
import { fetchProject } from "@/queries/fetchProject.ts";
import { fetchUserConfig } from "@/queries/fetchUserConfig.ts";
import { ProjectSelectionScreen } from "@/screens/ProjectSelectionScreen/ProjectSelectionScreen.tsx";
import { SettingsScreen } from "@/screens/SettingsScreen/SettingsScreen.tsx";
import { useProjectStore } from "@/store/projectStore.ts";
import { Screen, useScreensStore } from "@/store/screensStore.ts";

import { LogsScreen } from "./screens/LogsScreen/LogsScreen.tsx";
import assertNever = util.assertNever;
import { AppSidebarLayout } from "@/components/layout/AppSidebarLayout/AppSidebarLayout.tsx";

function App() {
  const project = useProjectStore((state) => state.projectInfo);
  const currentScreen = useScreensStore((state) => state.currentScreen);
  const setCurrentScreen = useScreensStore((state) => state.setCurrentScreen);

  useEffect(() => {
    if (!project) {
      setCurrentScreen(Screen.ProjectSelection);
    } else {
      setCurrentScreen(Screen.Logs);
    }
  }, [project, setCurrentScreen]);

  useEffect(() => {
    fetchProject();
    fetchUserConfig();
  }, []);

  const getScreen = () => {
    switch (currentScreen) {
      case Screen.Logs:
        return (
          <AppSidebarLayout>
            <LogsScreen />
          </AppSidebarLayout>
        );
      case Screen.ProjectSelection:
        return <ProjectSelectionScreen />;
      case Screen.Settings:
        return <SettingsScreen />;
      default:
        assertNever(currentScreen);
    }
  };

  return (
    <>
      <EventListenersContainer />
      <Toaster richColors />
      {getScreen()}
    </>
  );
}

export default App;
