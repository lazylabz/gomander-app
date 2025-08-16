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
  arrayMove,
  SortableContext,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { ChevronRight, X } from "lucide-react";
import { useMemo } from "react";
import { useFormContext } from "react-hook-form";

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
} from "@/components/ui/form.tsx";
import { isDefined } from "@/helpers/mapHelpers.ts";
import { useCommandStore } from "@/store/commandStore.ts";

export const CommandGroupCommandsField = () => {
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
      addedCommands: selectedCommandIds
        .map((id) => allCommands.find((cmd) => cmd.id === id))
        .filter(isDefined),
    };
  }, [allCommands, selectedCommandIds]);

  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 8,
      },
    }),
  );

  const findContainer = (id: string): string | null => {
    if (id === AVAILABLE_COMMANDS || id === ADDED_COMMANDS) {
      return id;
    }

    if (availableCommands.find((cmd) => cmd.id === id)) {
      return AVAILABLE_COMMANDS;
    }
    if (addedCommands.find((cmd) => cmd.id === id)) {
      return ADDED_COMMANDS;
    }

    return null;
  };

  const addCommand = (commandId: string) => {
    const newSelected = [...selectedCommandIds, commandId];
    form.setValue("commands", newSelected);
  };

  const removeCommand = (commandId: string) => {
    const newSelected = selectedCommandIds.filter((id) => id !== commandId);
    form.setValue("commands", newSelected);
  };

  const handleDragOver = (event: DragOverEvent) => {
    const { active, over } = event;

    if (!over) {
      return;
    }

    const activeContainer = findContainer(active.id.toString());
    const overContainer = findContainer(over.id.toString());

    if (
      !activeContainer ||
      !overContainer ||
      activeContainer === overContainer
    ) {
      return;
    }

    const commandId = active.id.toString();

    // Move between containers
    if (
      activeContainer === AVAILABLE_COMMANDS &&
      overContainer === ADDED_COMMANDS
    ) {
      addCommand(commandId);
    } else if (
      activeContainer === ADDED_COMMANDS &&
      overContainer === AVAILABLE_COMMANDS
    ) {
      removeCommand(commandId);
    }
  };

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;

    if (!over) {
      return;
    }

    const activeContainer = findContainer(active.id.toString());
    const overContainer = findContainer(over.id.toString());

    // Reorder within added commands
    if (
      activeContainer === ADDED_COMMANDS &&
      overContainer === ADDED_COMMANDS
    ) {
      const oldIndex = selectedCommandIds.findIndex(
        (id) => id === active.id.toString(),
      );
      const newIndex = selectedCommandIds.findIndex(
        (id) => id === over.id.toString(),
      );

      if (oldIndex !== -1 && newIndex !== -1 && oldIndex !== newIndex) {
        const reorderedIds = arrayMove(selectedCommandIds, oldIndex, newIndex);
        form.setValue("commands", reorderedIds);
      }
    }
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
            onDragOver={handleDragOver}
            onDragEnd={handleDragEnd}
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
                          onClick={() => addCommand(command.id)}
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
                          onClick={() => removeCommand(command.id)}
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
            Drag commands from left to right to add them, and reorder them as
            needed
          </FormDescription>

          <FormMessage />
        </FormItem>
      )}
    />
  );
};
