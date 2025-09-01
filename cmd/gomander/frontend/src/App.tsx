import { useEffect } from "react";
import { Route, Routes } from "react-router";

import { AppSidebarLayout } from "@/components/layout/AppSidebarLayout/AppSidebarLayout.tsx";
import { Toaster } from "@/components/ui/sonner.tsx";
import { EventListenersContainer } from "@/components/utility/EventListenersContainer.tsx";
import { ThemeProvider } from "@/contexts/theme.tsx";
import { VersionProvider } from "@/contexts/version.tsx";
import { fetchUserConfig } from "@/queries/fetchUserConfig.ts";
import { loadAllProjectData } from "@/queries/loadAllProjectData.ts";
import { ScreenRoutes } from "@/routes.ts";
import { ProjectSelectionScreen } from "@/screens/ProjectSelectionScreen/ProjectSelectionScreen.tsx";
import { SettingsContextProvider } from "@/screens/SettingsScreen/contexts/settingsContext.tsx";
import { SettingsScreen } from "@/screens/SettingsScreen/SettingsScreen.tsx";

import { LogsScreen } from "./screens/LogsScreen/LogsScreen.tsx";

function App() {
  useEffect(() => {
    loadAllProjectData();
    fetchUserConfig();
  }, []);

  return (
    <VersionProvider>
      <ThemeProvider defaultTheme="system" storageKey="vite-ui-theme">
        <EventListenersContainer />
        <Toaster richColors position="top-right" />
        <Routes>
          <Route
            path={ScreenRoutes.ProjectSelection}
            element={<ProjectSelectionScreen />}
          />
          <Route
            path={ScreenRoutes.Logs}
            element={
              <AppSidebarLayout>
                <LogsScreen />
              </AppSidebarLayout>
            }
          />

          <Route
            path={ScreenRoutes.Settings}
            element={
              <SettingsContextProvider>
                <SettingsScreen />
              </SettingsContextProvider>
            }
          />
        </Routes>
      </ThemeProvider>
    </VersionProvider>
  );
}

export default App;
