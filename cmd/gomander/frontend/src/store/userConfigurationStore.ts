import { useStore } from "zustand/react";
import { createStore } from "zustand/vanilla";

import { LOGS_FONT_SIZE } from "@/constants/logsFontSize.ts";
import type { UserConfig } from "@/contracts/types.ts";

type UserConfigurationStore = {
	userConfig: UserConfig;
	setUserConfig: (config: UserConfig) => void;
	isLoaded: boolean;
};

export const userConfigurationStore = createStore<UserConfigurationStore>()(
	(set) => ({
		isLoaded: false,
		userConfig: {
			environmentPaths: [],
			lastOpenedProjectId: "",
			locale: "en",
			logsFontSize: LOGS_FONT_SIZE.default,
			zoomLevel: 1,
			notificationsEnabled: true,
		},
		setUserConfig: (config: UserConfig) => {
			set({ userConfig: config, isLoaded: true });
		},
	}),
);

export const useUserConfigurationStore = <T>(
	selector: (state: UserConfigurationStore) => T,
): T => {
	return useStore(userConfigurationStore, selector);
};
