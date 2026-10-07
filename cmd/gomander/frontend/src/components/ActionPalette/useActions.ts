import { useTranslation } from "react-i18next";

import type { Action } from "@/components/ActionPalette/action.ts";
import { commandActions } from "@/components/ActionPalette/actions/commandActions.ts";
import { commandGroupActions } from "@/components/ActionPalette/actions/commandGroupActions.ts";
import { useCommandGroupStore } from "@/store/commandGroupStore.ts";
import { useCommandStore } from "@/store/commandStore.ts";

export const useActions = (): Action[] => {
	// Labels are translated outside React, so this re-renders the palette on a language change.
	useTranslation();

	const commands = useCommandStore((state) => state.commands);
	const commandsStatus = useCommandStore((state) => state.commandsStatus);
	const commandGroups = useCommandGroupStore((state) => state.commandGroups);

	return [
		...commandActions(commands, commandsStatus),
		...commandGroupActions(commandGroups, commands, commandsStatus),
	];
};
