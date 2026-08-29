import type { Command, CommandGroup } from "@/contracts/types.ts";
import { isDefined } from "@/helpers/mapHelpers.ts";
import { commandStore } from "@/store/commandStore.ts";

export interface CommandGroupWithCommandIds
	extends Omit<CommandGroup, "commands"> {
	commands: string[];
}

// A group carries a copy of every command it holds, but a form only knows the ids
// the user ticked. An id the store does not know about is dropped: the group is
// saved with what could be resolved.
export const resolveGroupCommands = (commandIds: string[]): Command[] => {
	const { commands } = commandStore.getState();

	return commandIds
		.map((commandId) => commands.find((c) => c.id === commandId))
		.filter(isDefined);
};
