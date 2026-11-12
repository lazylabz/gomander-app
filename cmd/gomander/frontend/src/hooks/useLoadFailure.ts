import { useCallback, useRef, useState } from "react";

type Load = () => Promise<void>;

export const useLoadFailure = () => {
	const [failedLoad, setFailedLoad] = useState<Load | null>(null);
	const [isPending, setIsPending] = useState(false);

	const latestAttempt = useRef(0);

	const attempt = useCallback(async (load: Load) => {
		latestAttempt.current += 1;
		const thisAttempt = latestAttempt.current;
		const isLatest = () => thisAttempt === latestAttempt.current;
		setIsPending(true);

		try {
			await load();
			if (isLatest()) {
				setFailedLoad(null);
			}
		} catch (error) {
			console.error("Failed to load:", error);
			if (isLatest()) {
				setFailedLoad(() => load);
			}
		} finally {
			if (isLatest()) {
				setIsPending(false);
			}
		}
	}, []);

	const retry = useCallback(async () => {
		if (failedLoad) {
			await attempt(failedLoad);
		}
	}, [attempt, failedLoad]);

	return { hasFailed: failedLoad !== null, isPending, attempt, retry };
};
