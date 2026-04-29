import { useStore } from "zustand/react";
import { createStore } from "zustand/vanilla";

import { EXPERIMENTAL_XTERMJS } from "@/constants/localStorage.ts";

export type ExperimentalFeats = {
	xtermjs: boolean;
};

type ExperimentalFeatsStore = {
	experimentalFeats: ExperimentalFeats;
	setExperimentalFeat: <K extends keyof ExperimentalFeats>(
		key: K,
		value: ExperimentalFeats[K],
	) => void;
};

const KEYS: Record<keyof ExperimentalFeats, string> = {
	xtermjs: EXPERIMENTAL_XTERMJS,
};

const readBool = (key: string, fallback: boolean): boolean => {
	try {
		const v = window.localStorage.getItem(key);
		return v !== null ? (JSON.parse(v) as boolean) : fallback;
	} catch {
		return fallback;
	}
};

export const experimentalFeatsStore = createStore<ExperimentalFeatsStore>()(
	(set) => ({
		experimentalFeats: {
			xtermjs: readBool(EXPERIMENTAL_XTERMJS, false),
		},
		setExperimentalFeat: (key, value) => {
			try {
				window.localStorage.setItem(KEYS[key], JSON.stringify(value));
			} catch {
				// Ignore write errors
			}
			set((state) => ({
				experimentalFeats: { ...state.experimentalFeats, [key]: value },
			}));
		},
	}),
);

export const useExperimentalFeatsStore = <T>(
	selector: (state: ExperimentalFeatsStore) => T,
): T => useStore(experimentalFeatsStore, selector);
