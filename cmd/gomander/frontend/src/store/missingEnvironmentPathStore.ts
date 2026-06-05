import { useStore } from "zustand/react";
import { createStore } from "zustand/vanilla";

type MissingEnvironmentPathStore = {
	dialogOpen: boolean;
	setDialogOpen: (open: boolean) => void;
};

// To be used in use cases
export const missingEnvironmentPathStore =
	createStore<MissingEnvironmentPathStore>()((set) => ({
		dialogOpen: false,
		setDialogOpen: (open) => set({ dialogOpen: open }),
	}));

// To be used in react components
export const useMissingEnvironmentPathStore = <T>(
	selector: (state: MissingEnvironmentPathStore) => T,
): T => useStore(missingEnvironmentPathStore, selector);
