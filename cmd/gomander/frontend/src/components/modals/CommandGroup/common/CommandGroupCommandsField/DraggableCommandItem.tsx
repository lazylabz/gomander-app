import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import React from "react";

import type { Command } from "@/contracts/types.ts";

export const DraggableCommandItem = ({
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
      className="group flex items-center p-1.5 pl-3 bg-card hover:bg-accent border rounded"
    >
      <div {...attributes} {...listeners} className="flex-1 cursor-grab">
        <span className="text-sm">{command.name}</span>
      </div>

      {rightComponent}
    </div>
  );
};
