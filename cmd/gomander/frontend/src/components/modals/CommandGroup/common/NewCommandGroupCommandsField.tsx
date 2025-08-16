import {
  closestCorners,
  DndContext,
  type DragEndEvent,
  type DragOverEvent,
  PointerSensor,
  useDroppable,
  useSensor,
  useSensors,
} from "@dnd-kit/core";
import {
  arrayMove,
  SortableContext,
  useSortable,
  verticalListSortingStrategy,
} from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { ChevronRight, X } from "lucide-react";
import { useMemo } from "react";
import React from "react";
import { useFormContext } from "react-hook-form";

import type { FormSchemaType } from "@/components/modals/CommandGroup/common/formSchema.ts";
import {
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form.tsx";
import type { Command } from "@/contracts/types.ts";
import { isDefined } from "@/helpers/mapHelpers.ts";
import { useCommandStore } from "@/store/commandStore.ts";

const DraggableCommandItem = ({
  command,
  rightComponent,
}: {
  command: Command;
  rightComponent?: React.ReactNode;
}) => {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: command.id });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
    opacity: isDragging ? 0.5 : 1,
  };

  return (
    <div
      ref={setNodeRef}
      style={style}
      className="flex items-center p-1.5 pl-3 bg-card border rounded hover:bg-accent transition-colors"
    >
      <div {...attributes} {...listeners} className="flex-1 cursor-grab">
        <span className="text-sm">{command.name}</span>
      </div>

      {rightComponent}
    </div>
  );
};

const DroppableContainer = ({
  id,
  title,
  children,
  className,
}: {
  id: string;
  title: string;
  children: React.ReactNode;
  className?: string;
}) => {
  const { setNodeRef } = useDroppable({ id });

  return (
    <div className={className}>
      <h4 className="font-medium text-sm mb-2 text-muted-foreground">
        {title}
      </h4>
      <div
        ref={setNodeRef}
        className="h-80 max-h-[calc(100vh-400px)] p-3 border rounded-lg overflow-y-auto overflow-x-hidden"
      >
        <div className="space-y-2">{children}</div>
        {React.Children.count(children) === 0 && (
          <div className="flex items-center justify-center h-full text-muted-foreground text-sm">
            {id === AVAILABLE_COMMANDS
              ? "All commands are added"
              : "Drag commands here"}
          </div>
        )}
      </div>
    </div>
  );
};

const AVAILABLE_COMMANDS = "available-commands";
const ADDED_COMMANDS = "added-commands";

export const NewCommandGroupCommandsField = () => {
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
          <div className="flex flex-col gap-1 mb-4">
            <FormLabel>Commands</FormLabel>
            <FormDescription className="text-xs">
              Drag commands from left to right to add them, and reorder them as
              needed
            </FormDescription>
          </div>

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
                  id={AVAILABLE_COMMANDS}
                  title="Available Commands"
                  className="flex-1"
                >
                  {availableCommands.map((command) => (
                    <DraggableCommandItem
                      key={command.id}
                      command={command}
                      rightComponent={
                        <button
                          type="button"
                          className="cursor-pointer flex items-center justify-center p-2 rounded bg-neutral-100 text-neutral-900 shadow-xs hover:bg-neutral-100/80 dark:bg-neutral-800 dark:text-neutral-50 dark:hover:bg-neutral-800/80 transition-colors"
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
                <DroppableContainer
                  id={ADDED_COMMANDS}
                  title="Added Commands"
                  className="flex-1"
                >
                  {addedCommands.map((command) => (
                    <DraggableCommandItem
                      key={command.id}
                      command={command}
                      rightComponent={
                        <button
                          type="button"
                          className="cursor-pointer flex items-center justify-center p-2 rounded bg-neutral-100 text-neutral-900 shadow-xs hover:bg-neutral-100/80 dark:bg-neutral-800 dark:text-neutral-50 dark:hover:bg-neutral-800/80 transition-colors"
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

          <FormMessage />
        </FormItem>
      )}
    />
  );
};
