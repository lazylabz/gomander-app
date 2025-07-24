import { useEffect } from "react";

import { fetchCommandGroups } from "@/queries/fetchCommandGroups.ts";
import { fetchCommands } from "@/queries/fetchCommands.ts";
import { fetchUserConfig } from "@/queries/fetchUserConfig.ts";

export const useFetchInitialData = () => {
  useEffect(() => {
    fetchCommandGroups();
    fetchCommands();
    fetchUserConfig();
  }, []);
};
