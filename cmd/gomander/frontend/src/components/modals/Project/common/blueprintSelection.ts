import type { ProjectBlueprint } from "@/contracts/types.ts";

type BlueprintGroup = ProjectBlueprint["commandGroups"][number];

// A group only makes sense to import alongside at least one of its commands.
export const isGroupSelectable = (
	group: BlueprintGroup,
	selectedCommandIds: string[],
): boolean =>
	group.commandIds.some((commandId) => selectedCommandIds.includes(commandId));

export const keepSelectableGroups = (
	blueprint: ProjectBlueprint,
	selectedGroupIds: string[],
	selectedCommandIds: string[],
): string[] =>
	selectedGroupIds.filter((groupId) => {
		const group = blueprint.commandGroups.find((cg) => cg.id === groupId);
		return !!group && isGroupSelectable(group, selectedCommandIds);
	});

export const narrowBlueprint = (
	blueprint: ProjectBlueprint,
	selectedCommandIds: string[],
	selectedGroupIds: string[],
): ProjectBlueprint => ({
	...blueprint,
	commands: blueprint.commands.filter((c) => selectedCommandIds.includes(c.id)),
	commandGroups: blueprint.commandGroups.filter((cg) =>
		selectedGroupIds.includes(cg.id),
	),
});
