import { dataService } from "@/contracts/service.ts";
import { projectStore } from "@/store/projectStore.ts";

export const fetchAvailableProjects = async (): Promise<void> => {
	const { setAvailableProjects } = projectStore.getState();

	setAvailableProjects(await dataService.getAvailableProjects());
};
