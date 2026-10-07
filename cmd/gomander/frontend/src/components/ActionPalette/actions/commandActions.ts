import { Play, Square, Terminal } from "lucide-react";

import type { Action, ActionKind } from "@/components/ActionPalette/action.ts";
import type { Command } from "@/contracts/types.ts";
import i18n from "@/design-system/lib/i18n.ts";
import { commandStore } from "@/store/commandStore.ts";
import { CommandStatus } from "@/types/CommandStatus.ts";
import { startCommand } from "@/useCases/command/startCommand.ts";
import { stopCommand } from "@/useCases/command/stopCommand.ts";

const showOutputOf = (commandId: string) =>
	commandStore.getState().setActiveCommandId(commandId);

export const commandActions = (
	commands: Command[],
	commandsStatus: Record<string, CommandStatus>,
): Action[] =>
	commands.flatMap((command): Action[] => {
		const keywords = [command.name, command.command];
		const kind: ActionKind = {
			icon: Terminal,
			label: i18n.t("actionPalette.kinds.command"),
		};

		switch (commandsStatus[command.id]) {
			case CommandStatus.IDLE:
				return [
					{
						id: `command.start.${command.id}`,
						section: "commands",
						label: i18n.t("actionPalette.run", {
							name: command.name,
						}),
						keywords,
						icon: Play,
						kind,
						perform: () => {
							showOutputOf(command.id);
							return startCommand(command.id);
						},
					},
				];
			case CommandStatus.RUNNING:
				return [
					{
						id: `command.stop.${command.id}`,
						section: "commands",
						label: i18n.t("actionPalette.stop", {
							name: command.name,
						}),
						keywords,
						icon: Square,
						kind,
						perform: async () => {
							await stopCommand(command.id);
							showOutputOf(command.id);
						},
					},
				];
			default:
				return [];
		}
	});
