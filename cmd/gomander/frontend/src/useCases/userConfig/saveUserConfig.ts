import { dataService } from "@/contracts/service.ts";
import type { UserConfig } from "@/contracts/types.ts";

export const saveUserConfig = async (config: UserConfig): Promise<void> => {
  await dataService.saveUserConfig(config);
};
