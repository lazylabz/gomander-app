import { dataService } from "@/contracts/service.ts";
import i18n from "@/design-system/lib/i18n.ts";
import { EXPECTED_VALIDATION_ERROR } from "@/helpers/errorHelpers.ts";
import { commandGroupStore } from "@/store/commandGroupStore.ts";

export const removeCommandFromGroup = async (
  commandId: string,
  groupId: string,
) => {
  const { commandGroups } = commandGroupStore.getState();
  const group = commandGroups.find((g) => g.id === groupId);
  if (!group) {
    throw new Error(i18n.t('toast.commandGroup.notFound'));
  }
  if (group.commands.length === 1) {
    throw new Error(i18n.t('toast.commandGroup.cannotRemoveLast'), {
      cause: EXPECTED_VALIDATION_ERROR,
    });
  }

  await dataService.removeCommandFromGroup(commandId, groupId);
};
