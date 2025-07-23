import { createContext, useContext, useEffect, useState } from "react";

import { type UserConfig } from "@/types/contracts.ts";

import { GetUserConfig, SaveUserConfig } from "../../wailsjs/go/app/App";

type UserConfigDataContextValue = {
  userConfig: UserConfig;
  saveUserConfig: (config: UserConfig) => Promise<void>;
  fetchUserConfig: () => Promise<void>;
};

export const userConfigDataContext = createContext<UserConfigDataContextValue>(
  {} as UserConfigDataContextValue,
);

export const UserConfigDataContextProvider = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  const [userConfig, setUserConfig] = useState<UserConfig>({
    extraPaths: [],
  });

  // User config operations
  const fetchUserConfig = async (): Promise<void> => {
    const config = await GetUserConfig();

    setUserConfig(config);
  };

  const saveUserConfig = async (config: UserConfig): Promise<void> => {
    // Assuming there's a function to set user config
    // This function should be implemented in the backend
    await SaveUserConfig(config);
  };

  useEffect(() => {
    fetchUserConfig();
  }, []);

  const value = { userConfig, saveUserConfig, fetchUserConfig };
  
  return (
    <userConfigDataContext.Provider value={value}>
      {children}
    </userConfigDataContext.Provider>
  );
};

export const useUserConfigDataContext = () => {
  return useContext(userConfigDataContext);
};
