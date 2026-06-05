import { toast } from "sonner";

import {
	MISSING_ENV_PATH_TOAST_ID,
	MissingEnvironmentPathToast,
} from "@/components/MissingEnvironmentPath/MissingEnvironmentPathToast.tsx";
import { commandStore } from "@/store/commandStore.ts";

// Shell/OS specific wording for "this executable could not be located".
// "No such file or directory" is intentionally excluded as it is too noisy.
const MISSING_PATH_PATTERNS = [
	/: command not found/i, // bash, zsh
	/: not found/i, // dash, sh, busybox ("sh: 1: foo: not found")
	/unknown command/i, // fish ("fish: Unknown command: foo")
	/is not recognized as an internal or external command/i, // Windows cmd.exe
	/executable file not found in \$?path/i, // Go exec
];

const isMissingPathLine = (line: string): boolean =>
	MISSING_PATH_PATTERNS.some((pattern) => pattern.test(line));

// The error is emitted as the process exits, so it lands in the last log lines.
const LINES_TO_SCAN = 20;

// `bufferedLines` are lines not yet flushed from the buffer into the store; we
// scan them alongside stored logs so the final error line is never missed.
export const detectMissingEnvironmentPathFailure = (
	commandId: string,
	bufferedLines: string[] = [],
): void => {
	const storedLines = commandStore.getState().commandsLogs[commandId] ?? [];
	const lastLines = [...storedLines, ...bufferedLines].slice(-LINES_TO_SCAN);

	if (!lastLines.some(isMissingPathLine)) {
		return;
	}

	toast.warning(<MissingEnvironmentPathToast />, {
		id: MISSING_ENV_PATH_TOAST_ID,
		icon: null,
	});
};
