import { commandGroupStore } from "@/store/commandGroupStore.ts";
import { commandStore } from "@/store/commandStore.ts";

export const cleanCommandGroupFromErrors = (commandGroupId: string) => {
  const { commandIdsWithErrors, setCommandIdsWithErrors } =
    commandStore.getState();

  const { commandGroups } = commandGroupStore.getState();

  const commandGroup = commandGroups.find(
    (group) => group.id === commandGroupId,
  );

  if (!commandGroup) {
    return;
  }

  setCommandIdsWithErrors(
    commandIdsWithErrors.filter((id) =>
      commandGroup.commands.map((c) => c.id).includes(id),
    ),
  );
};
