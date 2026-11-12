import { fetchCommandGroups } from "@/queries/fetchCommandGroups.ts";
import { fetchCommands } from "@/queries/fetchCommands.ts";
import { fetchProject } from "@/queries/fetchProject.ts";

// Waits for both fetches even when one rejects: a fetch left running would write
// its project's data into the stores after the caller has moved on to another one.
export const loadAllProjectData = async () => {
	const results = await Promise.allSettled([
		fetchCommands(),
		fetchCommandGroups(),
	]);

	const failure = results.find(
		(result): result is PromiseRejectedResult => result.status === "rejected",
	);
	if (failure) {
		throw failure.reason;
	}

	await fetchProject();
};
