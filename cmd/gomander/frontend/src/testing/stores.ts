import type { StoreApi } from "zustand/vanilla";

import { commandGroupStore } from "@/store/commandGroupStore.ts";
import { commandStore } from "@/store/commandStore.ts";
import { missingEnvironmentPathStore } from "@/store/missingEnvironmentPathStore.ts";
import { projectStore } from "@/store/projectStore.ts";
import { userConfigurationStore } from "@/store/userConfigurationStore.ts";

const reset = <T>(store: StoreApi<T>) =>
	store.setState(store.getInitialState(), true);

export const resetStores = (): void => {
	reset(commandStore);
	reset(commandGroupStore);
	reset(projectStore);
	reset(userConfigurationStore);
	reset(missingEnvironmentPathStore);
};
