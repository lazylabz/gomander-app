import { dataService } from "@/contracts/service.ts";
import { EXPECTED_VALIDATION_ERROR } from "@/helpers/errorHelpers.ts";
import { commandGroupStore } from "@/store/commandGroupStore.ts";

export const removeCommandFromGroup = async (
  commandId: string,
  groupId: string,
) => {
  const { commandGroups } = commandGroupStore.getState();
  const group = commandGroups.find((g) => g.id === groupId);
  if (!group) {
    throw new Error("Command group not found");
  }
  if (group.commands.length === 1) {
    throw new Error(
      "Cannot remove the last command from the group. Delete the group instead.",
      { cause: EXPECTED_VALIDATION_ERROR },
    );
  }

  await dataService.removeCommandFromGroup(commandId, groupId);
};
