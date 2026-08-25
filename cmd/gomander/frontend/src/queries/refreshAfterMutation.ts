// A mutation refreshes once it has already reported its own outcome, and no call
// site keeps a `catch` around it any more - so a refresh that fails must not
// turn a finished mutation into a rejection.
export const refreshAfterMutation = async (
	...refreshes: (() => Promise<void>)[]
): Promise<void> => {
	await Promise.allSettled(refreshes.map((refresh) => refresh()));
};
