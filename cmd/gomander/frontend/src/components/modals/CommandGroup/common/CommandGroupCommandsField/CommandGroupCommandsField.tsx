import {
	closestCorners,
	DndContext,
	type DragEndEvent,
	type DragOverEvent,
	PointerSensor,
	useSensor,
	useSensors,
} from "@dnd-kit/core";
import {
	SortableContext,
	verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { ChevronRight, X } from "lucide-react";
import { useMemo } from "react";
import { useFormContext } from "react-hook-form";
import { useTranslation } from "react-i18next";

import {
	addCommand,
	type Drag,
	moveAcrossContainers,
	removeCommand,
	reorderAdded,
} from "@/components/modals/CommandGroup/common/CommandGroupCommandsField/commandGroupSelection.ts";
import {
	ADDED_COMMANDS,
	AVAILABLE_COMMANDS,
} from "@/components/modals/CommandGroup/common/CommandGroupCommandsField/constants.ts";
import { DraggableCommandItem } from "@/components/modals/CommandGroup/common/CommandGroupCommandsField/DraggableCommandItem.tsx";
import { DroppableContainer } from "@/components/modals/CommandGroup/common/CommandGroupCommandsField/DroppableContainer.tsx";
import type { FormSchemaType } from "@/components/modals/CommandGroup/common/formSchema.ts";
import {
	FormDescription,
	FormField,
	FormItem,
	FormMessage,
} from "@/design-system/components/ui/form.tsx";
import { resolveCommands } from "@/helpers/commandHelpers.ts";
import { useCommandStore } from "@/store/commandStore.ts";

export const CommandGroupCommandsField = () => {
	const { t } = useTranslation();
	const allCommands = useCommandStore((state) => state.commands);
	const form = useFormContext<FormSchemaType>();

	const formCommands = form.watch("commands");
	const selectedCommandIds = useMemo(() => {
		return formCommands || [];
	}, [formCommands]);

	const { availableCommands, addedCommands } = useMemo(() => {
		const selectedSet = new Set(selectedCommandIds);

		return {
			availableCommands: allCommands.filter((cmd) => !selectedSet.has(cmd.id)),
			addedCommands: resolveCommands(selectedCommandIds, allCommands),
		};
	}, [allCommands, selectedCommandIds]);

	const sensors = useSensors(
		useSensor(PointerSensor, {
			activationConstraint: {
				distance: 8,
			},
		}),
	);

	const allCommandIds = allCommands.map((cmd) => cmd.id);

	const setSelected = (next: string[]) => {
		if (next !== selectedCommandIds) {
			form.setValue("commands", next);
		}
	};

	const applyDrag =
		(
			resolve: (
				selectedIds: string[],
				allCommandIds: string[],
				drag: Drag,
			) => string[],
		) =>
		({ active, over }: DragOverEvent | DragEndEvent) => {
			if (!over) {
				return;
			}

			setSelected(
				resolve(selectedCommandIds, allCommandIds, {
					activeId: active.id.toString(),
					overId: over.id.toString(),
				}),
			);
		};

	return (
		<FormField
			control={form.control}
			name="commands"
			render={() => (
				<FormItem>
					<DndContext
						sensors={sensors}
						collisionDetection={closestCorners}
						onDragOver={applyDrag(moveAcrossContainers)}
						onDragEnd={applyDrag(reorderAdded)}
					>
						<div className="flex gap-6 select-none">
							<SortableContext
								items={availableCommands.map((cmd) => cmd.id)}
								strategy={verticalListSortingStrategy}
							>
								<DroppableContainer
									variant={AVAILABLE_COMMANDS}
									className="flex-1"
								>
									{availableCommands.map((command) => (
										<DraggableCommandItem
											key={command.id}
											command={command}
											rightComponent={
												<button
													type="button"
													className="cursor-pointer flex items-center justify-center p-2 rounded text-neutral-900 shadow-xs dark:text-neutral-50 bg-accent group-hover:bg-neutral-200 hover:bg-neutral-300/80 dark:group-hover:bg-card/60 dark:hover:bg-card"
													onClick={() =>
														setSelected(
															addCommand(selectedCommandIds, command.id),
														)
													}
												>
													<ChevronRight className="size-4" />
												</button>
											}
										/>
									))}
								</DroppableContainer>
							</SortableContext>

							<SortableContext
								items={addedCommands.map((cmd) => cmd.id)}
								strategy={verticalListSortingStrategy}
							>
								<DroppableContainer variant={ADDED_COMMANDS} className="flex-1">
									{addedCommands.map((command) => (
										<DraggableCommandItem
											key={command.id}
											command={command}
											rightComponent={
												<button
													type="button"
													className="cursor-pointer flex items-center justify-center p-2 rounded text-neutral-900 shadow-xs dark:text-neutral-50 bg-accent group-hover:bg-neutral-200 hover:bg-neutral-300/80 dark:group-hover:bg-card/60 dark:hover:bg-card"
													onClick={() =>
														setSelected(
															removeCommand(selectedCommandIds, command.id),
														)
													}
												>
													<X className="size-4" />
												</button>
											}
										/>
									))}
								</DroppableContainer>
							</SortableContext>
						</div>
					</DndContext>
					<FormDescription className="text-xs mt-1">
						{t("commandGroupForm.commandsDescription")}
					</FormDescription>

					<FormMessage />
				</FormItem>
			)}
		/>
	);
};
