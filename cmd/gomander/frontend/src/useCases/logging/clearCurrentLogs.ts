import { clearLogTail } from "@/store/commandLogsTail.ts";
import { commandStore } from "@/store/commandStore.ts";
import { terminalStore } from "@/store/terminalStore.ts";

export const clearCurrentLogs = () => {
	const { activeCommandId } = commandStore.getState();

	if (!activeCommandId) {
		return;
	}

	clearLogTail(activeCommandId);

	const { resetTerminal } = terminalStore.getState();
	resetTerminal(activeCommandId);
};
