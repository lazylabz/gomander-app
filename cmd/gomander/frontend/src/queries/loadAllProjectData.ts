import { fetchCommandGroups } from "@/queries/fetchCommandGroups.ts";
import { fetchCommands } from "@/queries/fetchCommands.ts";
import { fetchProject } from "@/queries/fetchProject.ts";

export const loadAllProjectData = async () => {
  await Promise.all([fetchProject(), fetchCommands(), fetchCommandGroups()]);
};
