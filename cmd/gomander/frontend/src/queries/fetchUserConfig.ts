import { dataService } from "@/contracts/service.ts";
import { userConfigurationStore } from "@/store/userConfigurationStore.ts";

export const fetchUserConfig = async (): Promise<void> => {
  const { setUserConfig } = userConfigurationStore.getState();
  const config = await dataService.getUserConfig();

  setUserConfig(config);
};
