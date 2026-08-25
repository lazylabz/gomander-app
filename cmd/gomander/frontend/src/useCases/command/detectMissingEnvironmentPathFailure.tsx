import { toast } from "sonner";

import { commandOutputTail } from "@/commandOutput/commandOutput.ts";
import {
	MISSING_ENV_PATH_TOAST_ID,
	MissingEnvironmentPathToast,
} from "@/components/MissingEnvironmentPath/MissingEnvironmentPathToast.tsx";

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

export const detectMissingEnvironmentPathFailure = (
	commandId: string,
): void => {
	if (!commandOutputTail(commandId).some(isMissingPathLine)) {
		return;
	}

	toast.warning(<MissingEnvironmentPathToast />, {
		id: MISSING_ENV_PATH_TOAST_ID,
		icon: null,
	});
};
