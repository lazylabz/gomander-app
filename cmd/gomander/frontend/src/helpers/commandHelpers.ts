import type { Command } from "@/contracts/types.ts";
import { isDefined } from "@/helpers/mapHelpers.ts";

// Walks the id list rather than the commands, so the order is the one that named
// them. An id no loaded command answers to is dropped, so a command that is gone
// leaves no hole behind.
export const resolveCommands = (
	commandIds: string[],
	commands: Command[],
): Command[] =>
	commandIds
		.map((commandId) => commands.find((command) => command.id === commandId))
		.filter(isDefined);
