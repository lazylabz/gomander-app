import { dataService } from "@/contracts/service.ts";
import { cleanCommandGroupFromErrors } from "@/useCases/command/cleanCommandGroupFromErrors.ts";

export const runCommandGroup = async (groupId: string) => {
  cleanCommandGroupFromErrors(groupId);

  await dataService.runCommandGroup(groupId);
};
