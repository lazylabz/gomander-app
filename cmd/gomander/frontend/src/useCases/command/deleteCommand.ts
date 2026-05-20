import { dataService } from "@/contracts/service.ts";
import { terminalStore } from "@/store/terminalStore.ts";

export const deleteCommand = async (commandId: string) => {
	await dataService.removeCommand(commandId);

	// Dispose existing terminal
	terminalStore.getState().dispose(commandId);
};
