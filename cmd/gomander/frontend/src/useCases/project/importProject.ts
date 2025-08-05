import { dataService } from "@/contracts/service.ts";
import type { ExportableProject } from "@/contracts/types.ts";

export const importProject = async (
  project: ExportableProject,
  newBaseWorkingDir: string,
) => {
  await dataService.importProject(project, newBaseWorkingDir);
};
