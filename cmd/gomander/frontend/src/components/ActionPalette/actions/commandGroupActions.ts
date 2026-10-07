import { Folder, Play, Square } from "lucide-react";

import type { Action, ActionKind } from "@/components/ActionPalette/action.ts";
import type { Command, CommandGroup } from "@/contracts/types.ts";
import i18n from "@/design-system/lib/i18n.ts";
import { resolveCommands } from "@/helpers/commandHelpers.ts";
import { CommandStatus } from "@/types/CommandStatus.ts";
import { runCommandGroup } from "@/useCases/commandGroup/runCommandGroup.ts";
import { stopCommandGroup } from "@/useCases/commandGroup/stopCommandGroup.ts";

// A partly running group offers both, as its row in the sidebar does.
export const commandGroupActions = (
	commandGroups: CommandGroup[],
	commands: Command[],
	commandsStatus: Record<string, CommandStatus>,
): Action[] =>
	commandGroups.flatMap((commandGroup): Action[] => {
		const commandsHeld = resolveCommands(commandGroup.commandIds, commands);
		const someCommandHas = (status: CommandStatus) =>
			commandsHeld.some((command) => commandsStatus[command.id] === status);

		const keywords = [commandGroup.name];
		const kind: ActionKind = {
			icon: Folder,
			label: i18n.t("actionPalette.kinds.commandGroup"),
		};
		const actions: Action[] = [];

		if (someCommandHas(CommandStatus.IDLE)) {
			actions.push({
				id: `commandGroup.run.${commandGroup.id}`,
				section: "commands",
				label: i18n.t("actionPalette.run", { name: commandGroup.name }),
				keywords,
				icon: Play,
				kind,
				perform: () => runCommandGroup(commandGroup.id),
			});
		}

		if (someCommandHas(CommandStatus.RUNNING)) {
			actions.push({
				id: `commandGroup.stop.${commandGroup.id}`,
				section: "commands",
				label: i18n.t("actionPalette.stop", { name: commandGroup.name }),
				keywords,
				icon: Square,
				kind,
				perform: () => stopCommandGroup(commandGroup.id),
			});
		}

		return actions;
	});
