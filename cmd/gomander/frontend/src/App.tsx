import { useCallback, useEffect, useState } from "react";
import { Route, Routes } from "react-router";

import { LoadFailure } from "@/components/LoadFailure/LoadFailure.tsx";
import { AppSidebarLayout } from "@/components/layout/AppSidebarLayout/AppSidebarLayout.tsx";
import { MissingEnvironmentPathDialog } from "@/components/MissingEnvironmentPath/MissingEnvironmentPathDialog.tsx";
import { EventListenersContainer } from "@/components/utility/EventListenersContainer.tsx";
import { ZoomShortcuts } from "@/components/utility/ZoomShortcuts.tsx";
import { ThemeProvider } from "@/contexts/theme.tsx";
import { Toaster } from "@/design-system/components/ui/sonner.tsx";
import { initFallbackI18n, initI18n } from "@/design-system/lib/i18n.ts";
import { useLoadFailure } from "@/hooks/useLoadFailure.ts";
import { fetchCurrentRelease } from "@/queries/fetchCurrentRelease.ts";
import { fetchUserConfig } from "@/queries/fetchUserConfig.ts";
import { loadAllProjectData } from "@/queries/loadAllProjectData.ts";
import { ScreenRoutes } from "@/routes.ts";
import { LogsScreen } from "@/screens/LogsScreen/LogsScreen.tsx";
import { ProjectSelectionScreen } from "@/screens/ProjectSelectionScreen/ProjectSelectionScreen.tsx";
import { SettingsContextProvider } from "@/screens/SettingsScreen/context/settingsContext.tsx";
import { SettingsScreen } from "@/screens/SettingsScreen/SettingsScreen.tsx";
import { useProjectStore } from "@/store/projectStore.ts";
import { useUserConfigurationStore } from "@/store/userConfigurationStore.ts";
import { checkForNewRelease } from "@/useCases/release/checkForNewRelease.ts";

function App() {
	const projectIsLoaded = useProjectStore((state) => state.isLoaded);
	const userConfigIsLoaded = useUserConfigurationStore(
		(state) => state.isLoaded,
	);

	const [i18nReady, setI18nReady] = useState(false);
	const { hasFailed, isPending, attempt, retry } = useLoadFailure();

	const initialFetchesAreDone = projectIsLoaded && userConfigIsLoaded;

	const initializeApp = useCallback(
		() =>
			attempt(async () => {
				try {
					await initI18n();
				} catch (error) {
					await initFallbackI18n();
					throw error;
				}

				setI18nReady(true);
				await Promise.all([loadAllProjectData(), fetchUserConfig()]);
			}),
		[attempt],
	);

	useEffect(() => {
		initializeApp();
	}, [initializeApp]);

	const appIsReady = i18nReady && initialFetchesAreDone;

	// Waits for the Toaster to be mounted, or a failed check would toast into nothing.
	useEffect(() => {
		if (!appIsReady) {
			return;
		}

		fetchCurrentRelease();
		checkForNewRelease();
	}, [appIsReady]);

	if (hasFailed) {
		return (
			<ThemeProvider defaultTheme="system" storageKey="vite-ui-theme">
				<div className="w-full h-full flex items-center justify-center">
					<LoadFailure onRetry={retry} isRetrying={isPending} />
				</div>
			</ThemeProvider>
		);
	}

	if (!appIsReady) {
		return null;
	}

	return (
		<ThemeProvider defaultTheme="system" storageKey="vite-ui-theme">
			<EventListenersContainer />
			<ZoomShortcuts />
			<MissingEnvironmentPathDialog />
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
	);
}

export default App;
