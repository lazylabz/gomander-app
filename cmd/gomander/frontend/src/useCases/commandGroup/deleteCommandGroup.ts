import { getCommandGroupSectionOpenLocalStorageKey } from "@/constants/localStorage.ts";
import { dataService } from "@/contracts/service.ts";
import { removeKeyFromLocalStorage } from "@/helpers/localStorage.ts";

export const deleteCommandGroup = async (commandGroupId: string) => {
  await dataService.deleteCommandGroup(commandGroupId);
  removeKeyFromLocalStorage(
    getCommandGroupSectionOpenLocalStorageKey(commandGroupId),
  );
};
