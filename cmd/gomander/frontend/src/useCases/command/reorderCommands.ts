import { dataService } from "@/contracts/service.ts";

export const reorderCommands = async (newOrder: string[]) => {
  await dataService.reorderCommands(newOrder);
};
