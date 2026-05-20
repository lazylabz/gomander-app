import { commandStore } from "@/store/commandStore.ts";
import { terminalStore } from "@/store/terminalStore.ts";
import { cleanCommandLogs } from "@/useCases/command/cleanCommandLogs.ts";

export const clearCurrentLogs = () => {
	const { activeCommandId } = commandStore.getState();

	if (!activeCommandId) {
		return;
	}

	cleanCommandLogs(activeCommandId);

	const { resetTerminal } = terminalStore.getState();
	resetTerminal(activeCommandId);
};
