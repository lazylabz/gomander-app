import { dataService } from "@/contracts/service.ts";
import type { ProjectExport } from "@/contracts/types.ts";

export const importProject = async (
  project: ProjectExport,
  name: string,
  workingDirectory: string,
) => {
  await dataService.importProject(project, name, workingDirectory);
};
