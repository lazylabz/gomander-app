import { useStore } from "zustand/react";
import { createStore } from "zustand/vanilla";

import type { Command } from "@/contracts/types.ts";
import type { CommandStatus } from "@/types/CommandStatus.ts";

type CommandStore = {
	commands: Command[];
	setCommands: (commands: Command[]) => void;

	commandsStatus: Record<string, CommandStatus>;
	setCommandsStatus: (commandsStatus: Record<string, CommandStatus>) => void;

	commandIdsWithErrors: string[];
	setCommandIdsWithErrors: (commandsWithErrors: string[]) => void;

	activeCommandId: string | null;
	setActiveCommandId: (commandId: string | null) => void;
};

// To be used in use cases
export const commandStore = createStore<CommandStore>()((set) => ({
	commands: [],
	setCommands: (commands) => set({ commands }),

	commandsStatus: {},
	setCommandsStatus: (commandsStatus) => set({ commandsStatus }),

	commandIdsWithErrors: [],
	setCommandIdsWithErrors: (commandsWithErrors) =>
		set({ commandIdsWithErrors: commandsWithErrors }),

	activeCommandId: null,
	setActiveCommandId: (commandId) => set({ activeCommandId: commandId }),
}));

// To be used in react components
export const useCommandStore = <T>(selector: (state: CommandStore) => T): T =>
	useStore(commandStore, selector);
