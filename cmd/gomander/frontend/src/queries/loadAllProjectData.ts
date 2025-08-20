import { fetchCommands } from "@/queries/fetchCommands.ts";
import { fetchProject } from "@/queries/fetchProject.ts";

export const loadAllProjectData = async () => {
  await fetchCommands();
  await fetchProject();
};
