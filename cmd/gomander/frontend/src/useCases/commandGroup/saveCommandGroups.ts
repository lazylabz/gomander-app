import { toast } from "sonner";

import { dataService } from "@/contracts/service.ts";
import type { CommandGroup } from "@/contracts/types.ts";
import { commandGroupStore } from "@/store/commandGroupStore.ts";

export const saveCommandGroups = async (
  groups: CommandGroup[],
): Promise<void> => {
  const { setCommandGroups, commandGroups } = commandGroupStore.getState();

  // Optimistic save to avoid flickering while drag and dropping
  const prev = commandGroups;
  setCommandGroups(groups);
  try {
    await dataService.saveCommandGroups(groups);
  } catch {
    // If saving fails, revert to previous state
    setCommandGroups(prev);
    toast.error("Failed to save command groups");
  }
};
