import { dataService } from "@/contracts/service.ts";
import { commandGroupStore } from "@/store/commandGroupStore.ts";

export const reorderCommandGroups = async (
  reorderedCommandGroupIds: string[],
): Promise<void> => {
  const { commandGroups, setCommandGroups } = commandGroupStore.getState();

  // Optimistic approach to avoid flickering
  setCommandGroups(
    commandGroups.sort(
      (a, b) =>
        reorderedCommandGroupIds.indexOf(a.id) -
        reorderedCommandGroupIds.indexOf(b.id),
    ),
  );

  await dataService.reorderCommandGroups(reorderedCommandGroupIds);
};
