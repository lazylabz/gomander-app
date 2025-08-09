import { useEffect } from "react";
import { Route, Routes } from "react-router";

import { Toaster } from "@/components/ui/sonner.tsx";
import { EventListenersContainer } from "@/components/utility/EventListenersContainer.tsx";
import { fetchProject } from "@/queries/fetchProject.ts";
import { fetchUserConfig } from "@/queries/fetchUserConfig.ts";
import { ScreenRoutes } from "@/routes.ts";
import { ProjectSelectionScreen } from "@/screens/ProjectSelectionScreen/ProjectSelectionScreen.tsx";
import { SettingsScreen } from "@/screens/SettingsScreen/SettingsScreen.tsx";

import { LogsScreen } from "./screens/LogsScreen/LogsScreen.tsx";

function App() {
  useEffect(() => {
    fetchProject();
    fetchUserConfig();
  }, []);

  return (
    <>
      <EventListenersContainer />
      <Toaster richColors />
      <Routes>
        <Route
          path={ScreenRoutes.ProjectSelection}
          element={<ProjectSelectionScreen />}
        />
        <Route path={ScreenRoutes.Logs} element={<LogsScreen />} />
        <Route path={ScreenRoutes.Settings} element={<SettingsScreen />} />
      </Routes>
    </>
  );
}

export default App;
