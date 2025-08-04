import { dataService } from "@/contracts/service.ts";

export const importProject = async () => {
  await dataService.importProject();
};
