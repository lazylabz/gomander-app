import { createContext, useContext } from "react";
import { useLocation } from "react-router";

export enum SettingsTab {
  User = "user",
  Project = "project",
}

// Define context
export interface SettingsContextData {
  initialTab: SettingsTab;
}

export const settingsContext = createContext<SettingsContextData>({
  initialTab: SettingsTab.User,
});

// Define provider
export const SettingsContextProvider = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  const { state } = useLocation();

  const initialTab = state?.tab || SettingsTab.User;

  const value: SettingsContextData = {
    initialTab,
  };

  return (
    <settingsContext.Provider value={value}>
      {children}
    </settingsContext.Provider>
  );
};

// Custom hook to use the settings context
export const useSettingsContext = () => {
  return useContext(settingsContext);
};
