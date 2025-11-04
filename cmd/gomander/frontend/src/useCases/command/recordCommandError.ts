import { commandStore } from "@/store/commandStore.ts";

export const recordCommandError = (commandId: string) => {
  const { commandIdsWithErrors, setCommandIdsWithErrors } =
    commandStore.getState();

  if (!commandIdsWithErrors.includes(commandId)) {
    setCommandIdsWithErrors([...commandIdsWithErrors, commandId]);
  }
};
