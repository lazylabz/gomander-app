import { useEffect, useState } from "react";
import { Route, Routes } from "react-router";

import { AppSidebarLayout } from "@/components/layout/AppSidebarLayout/AppSidebarLayout.tsx";
import { EventListenersContainer } from "@/components/utility/EventListenersContainer.tsx";
import { ThemeProvider } from "@/contexts/theme.tsx";
import { VersionProvider } from "@/contexts/version.tsx";
import { Toaster } from "@/design-system/components/ui/sonner.tsx";
import { initI18n } from "@/lib/i18n.ts";
import { fetchUserConfig } from "@/queries/fetchUserConfig.ts";
import { loadAllProjectData } from "@/queries/loadAllProjectData.ts";
import { ScreenRoutes } from "@/routes.ts";
import { ProjectSelectionScreen } from "@/screens/ProjectSelectionScreen/ProjectSelectionScreen.tsx";
import { SettingsContextProvider } from "@/screens/SettingsScreen/contexts/settingsContext.tsx";
import { SettingsScreen } from "@/screens/SettingsScreen/SettingsScreen.tsx";

import { LogsScreen } from "./screens/LogsScreen/LogsScreen.tsx";

function App() {
  const [i18nReady, setI18nReady] = useState(false);

  useEffect(() => {
    const initializeApp = async () => {
      await initI18n();
      setI18nReady(true);
      loadAllProjectData();
      fetchUserConfig();
    };

    initializeApp();
  }, []);

  if (!i18nReady) {
    return <div>Loading...</div>;
  }

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
