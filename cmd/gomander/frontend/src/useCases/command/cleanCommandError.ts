import { commandStore } from "@/store/commandStore.ts";

export const cleanCommandError = (commandId: string) => {
  const { commandIdsWithErrors, setCommandIdsWithErrors } =
    commandStore.getState();

  if (commandIdsWithErrors.includes(commandId)) {
    setCommandIdsWithErrors(
      commandIdsWithErrors.filter((id) => id !== commandId),
    );
  }
};
