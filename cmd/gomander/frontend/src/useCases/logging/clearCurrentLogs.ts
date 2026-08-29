import { resetCommandOutput } from "@/commandOutput/commandOutput.ts";
import { commandStore } from "@/store/commandStore.ts";

export const clearCurrentLogs = () => {
	const { activeCommandId } = commandStore.getState();

	if (!activeCommandId) {
		return;
	}

	resetCommandOutput(activeCommandId);
};
