// A mutation refreshes once it has already reported its own outcome, and no call
// site keeps a `catch` around it any more - so a refresh that fails must not
// turn a finished mutation into a rejection. The `async` wrapper is what keeps a
// refresh that throws synchronously inside `allSettled`.
export const refreshAfterMutation = async (
	...refreshes: (() => Promise<void>)[]
): Promise<void> => {
	const results = await Promise.allSettled(
		refreshes.map(async (refresh) => refresh()),
	);

	for (const result of results) {
		if (result.status === "rejected") {
			console.error("Failed to refresh after a mutation:", result.reason);
		}
	}
};
