import { useCallback, useState } from "react";

import { dataService } from "@/contracts/service.ts";
import type { Project } from "@/contracts/types.ts";

export const useGetAvailableProjects = () => {
  const [availableProjects, setAvailableProjects] = useState<Project[]>([]);

  const fetchAvailableProjects = useCallback(async () => {
    const projects = await dataService.getAvailableProjects();
    setAvailableProjects(projects);
  }, []);

  return {
    availableProjects,
    fetchAvailableProjects,
  };
};
