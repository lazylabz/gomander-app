import { useDroppable } from "@dnd-kit/core";
import React from "react";

import {
  ADDED_COMMANDS,
  AVAILABLE_COMMANDS,
} from "@/components/modals/CommandGroup/common/CommandGroupCommandsField/constants.ts";
import { FormLabel } from "@/design-system/components/ui/form.tsx";

export const DroppableContainer = ({
  variant,
  children,
  className,
}: {
  variant: typeof AVAILABLE_COMMANDS | typeof ADDED_COMMANDS;
  children: React.ReactNode;
  className?: string;
}) => {
  const { setNodeRef } = useDroppable({ id: variant });

  return (
    <div className={className}>
      {variant === AVAILABLE_COMMANDS && (
        <h4 className="font-medium text-sm mb-2">Available Commands</h4>
      )}
      {variant === ADDED_COMMANDS && (
        <FormLabel className="font-medium text-sm mb-2">
          Group Commands
        </FormLabel>
      )}
      <div
        ref={setNodeRef}
        className="h-80 max-h-[calc(100vh-400px)] p-3 border rounded-lg overflow-y-auto overflow-x-hidden"
      >
        <div className="space-y-2">{children}</div>
        {React.Children.count(children) === 0 && (
          <div className="flex items-center justify-center h-full text-muted-foreground text-sm text-center">
            {variant === AVAILABLE_COMMANDS
              ? "No commands available"
              : "Drag commands here"}
          </div>
        )}
      </div>
    </div>
  );
};
