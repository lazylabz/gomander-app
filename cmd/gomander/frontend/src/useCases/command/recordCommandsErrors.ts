import { commandStore } from "@/store/commandStore.ts";

export const recordCommandsErrors = (commandIds: string[]) => {
  const { commandIdsWithErrors, setCommandIdsWithErrors } =
    commandStore.getState();

  const updatedCommandIdsWithErrors = new Set([
    ...commandIdsWithErrors,
    ...commandIds,
  ]);

  setCommandIdsWithErrors(Array.from(updatedCommandIdsWithErrors));
};
