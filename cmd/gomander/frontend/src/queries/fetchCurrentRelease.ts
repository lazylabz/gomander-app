import { dataService } from "@/contracts/service.ts";
import { releaseStore } from "@/store/releaseStore.ts";

export const fetchCurrentRelease = async (): Promise<void> => {
	releaseStore.setState({
		currentVersion: await dataService.getCurrentRelease(),
	});
};
