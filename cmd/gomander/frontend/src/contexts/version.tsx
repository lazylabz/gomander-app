import { createContext, useContext, useEffect, useState } from "react";
import { toast } from "sonner";

import { helpersService } from "@/contracts/service.ts";

type VersionContext = {
  currentVersion: string;
  newVersion: string | null;
  errorLoadingNewVersion: Error | null;
};

export const versionContext = createContext<VersionContext>({
  currentVersion: "",
  newVersion: null,
  errorLoadingNewVersion: null,
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

  return (
    <versionContext.Provider
      value={{
        currentVersion: currentRelease,
        newVersion: newRelease,
        errorLoadingNewVersion,
      }}
    >
      {children}
    </versionContext.Provider>
  );
};

export const useVersionContext = () => useContext(versionContext);
