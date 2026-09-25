import { arrayMove } from "@dnd-kit/sortable";

import {
	ADDED_COMMANDS,
	AVAILABLE_COMMANDS,
} from "@/components/modals/CommandGroup/common/CommandGroupCommandsField/constants.ts";

type Container = typeof AVAILABLE_COMMANDS | typeof ADDED_COMMANDS;

export type Drag = { activeId: string; overId: string };

// Each function hands back the very same array when the drag changes nothing,
// so the caller can skip touching the form.

const containerOf = (
	id: string,
	selectedIds: string[],
	allCommandIds: string[],
): Container | null => {
	if (id === AVAILABLE_COMMANDS || id === ADDED_COMMANDS) {
		return id;
	}
	if (!allCommandIds.includes(id)) {
		return null;
	}

	return selectedIds.includes(id) ? ADDED_COMMANDS : AVAILABLE_COMMANDS;
};

export const addCommand = (selectedIds: string[], commandId: string) => [
	...selectedIds,
	commandId,
];

export const removeCommand = (selectedIds: string[], commandId: string) =>
	selectedIds.filter((id) => id !== commandId);

export const moveAcrossContainers = (
	selectedIds: string[],
	allCommandIds: string[],
	{ activeId, overId }: Drag,
): string[] => {
	const from = containerOf(activeId, selectedIds, allCommandIds);
	const to = containerOf(overId, selectedIds, allCommandIds);

	if (from === AVAILABLE_COMMANDS && to === ADDED_COMMANDS) {
		return addCommand(selectedIds, activeId);
	}
	if (from === ADDED_COMMANDS && to === AVAILABLE_COMMANDS) {
		return removeCommand(selectedIds, activeId);
	}

	return selectedIds;
};

export const reorderAdded = (
	selectedIds: string[],
	allCommandIds: string[],
	{ activeId, overId }: Drag,
): string[] => {
	const from = containerOf(activeId, selectedIds, allCommandIds);
	const to = containerOf(overId, selectedIds, allCommandIds);
	if (from !== ADDED_COMMANDS || to !== ADDED_COMMANDS) {
		return selectedIds;
	}

	const oldIndex = selectedIds.indexOf(activeId);
	const newIndex = selectedIds.indexOf(overId);
	if (oldIndex === -1 || newIndex === -1 || oldIndex === newIndex) {
		return selectedIds;
	}

	return arrayMove(selectedIds, oldIndex, newIndex);
};
