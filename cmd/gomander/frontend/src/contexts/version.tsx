import { createContext, useContext, useEffect, useState } from "react";
import { toast } from "sonner";

import { externalBrowserService, helpersService } from "@/contracts/service.ts";

type VersionContext = {
  currentVersion: string;
  newVersion: string | null;
  errorLoadingNewVersion: Error | null;
  openLatestReleasePage: () => void;
};

export const versionContext = createContext<VersionContext>({
  currentVersion: "",
  newVersion: null,
  errorLoadingNewVersion: null,
  openLatestReleasePage: () => {},
});

export const VersionProvider = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  const [currentRelease, setCurrentRelease] = useState<string>("");
  const [newRelease, setNewRelease] = useState<string | null>(null);
  const [errorLoadingNewVersion, setErrorLoadingNewVersion] =
    useState<Error | null>(null);

  const fetchCurrentRelease = async () => {
    const release = await helpersService.getCurrentRelease();
    setCurrentRelease(release);
  };

  const checkNewRelease = async () => {
    setErrorLoadingNewVersion(null);
    try {
      const release = await helpersService.isThereANewRelease();
      if (release) {
        setNewRelease(release);
      }
    } catch (err) {
      setErrorLoadingNewVersion(err as Error);
      console.error("Error checking for new releases:", err);
      toast.error(`Error checking for new releases`);
    }
  };

  useEffect(() => {
    fetchCurrentRelease();
    checkNewRelease();
  }, []);

  const openLatestReleasePage = () => {
    const url = `https://github.com/lazylabz/gomander-app/releases/latest`;

    externalBrowserService.browserOpenURL(url);
  };

  return (
    <versionContext.Provider
      value={{
        currentVersion: currentRelease,
        newVersion: newRelease,
        errorLoadingNewVersion,
        openLatestReleasePage,
      }}
    >
      {children}
    </versionContext.Provider>
  );
};

export const useVersionContext = () => useContext(versionContext);
