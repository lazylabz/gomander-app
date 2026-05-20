import type { ExperimentalFeats } from "@/store/experimentalFeatsStore.ts";
import { experimentalFeatsStore } from "@/store/experimentalFeatsStore.ts";

export const changeExperimentalFeat = <K extends keyof ExperimentalFeats>(
	key: K,
	value: ExperimentalFeats[K],
): void => {
	experimentalFeatsStore.getState().setExperimentalFeat(key, value);
};
