import { toast } from "sonner";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { appendCommandOutput } from "@/commandOutput/commandOutput.ts";
import { MISSING_ENV_PATH_TOAST_ID } from "@/components/MissingEnvironmentPath/MissingEnvironmentPathToast.tsx";
import {
	installRecordingTerminals,
	resetTerminals,
} from "@/testing/terminals.ts";
import { detectMissingEnvironmentPathFailure } from "@/useCases/command/detectMissingEnvironmentPathFailure.tsx";

describe("detectMissingEnvironmentPathFailure", () => {
	const sut = detectMissingEnvironmentPathFailure;

	const warning = vi.spyOn(toast, "warning");

	beforeEach(() => {
		vi.useFakeTimers();
		installRecordingTerminals();
		warning.mockClear();
	});

	afterEach(() => {
		resetTerminals();
		vi.useRealTimers();
	});

	it.each([
		"bash: line 1: nope: command not found",
		"sh: 1: nope: not found",
		"fish: Unknown command: nope",
		"'nope' is not recognized as an internal or external command",
		'exec: "nope": executable file not found in $PATH',
	])("Should warn when the output ends with %s", (line) => {
		// Arrange
		appendCommandOutput("cmd-1", [line]);

		// Act
		sut("cmd-1");

		// Assert
		expect(warning).toHaveBeenCalledOnce();
		expect(warning.mock.calls[0][1]).toMatchObject({
			id: MISSING_ENV_PATH_TOAST_ID,
		});
	});

	it("Should warn on a line that has not reached the terminal yet", () => {
		// Arrange the failing line without ever letting the flush loop run: the
		// process finishes before the 30 ms buffer is drained.
		appendCommandOutput("cmd-1", ["bash: nope: command not found"]);

		// Act
		sut("cmd-1");

		// Assert
		expect(warning).toHaveBeenCalledOnce();
	});

	it("Should stay quiet when the output holds no missing-path error", () => {
		// Arrange
		appendCommandOutput("cmd-1", ["Server listening on :3000"]);

		// Act
		sut("cmd-1");

		// Assert
		expect(warning).not.toHaveBeenCalled();
	});

	it("Should stay quiet when the command produced no output at all", () => {
		// Act
		sut("cmd-1");

		// Assert
		expect(warning).not.toHaveBeenCalled();
	});

	it("Should stay quiet when the failure belongs to another command", () => {
		// Arrange
		appendCommandOutput("cmd-2", ["bash: nope: command not found"]);

		// Act
		sut("cmd-1");

		// Assert
		expect(warning).not.toHaveBeenCalled();
	});

	it("Should stay quiet once the failing line has scrolled out of the tail", () => {
		// Arrange
		appendCommandOutput("cmd-1", ["bash: nope: command not found"]);
		appendCommandOutput(
			"cmd-1",
			Array.from({ length: 20 }, (_, i) => `line ${i}`),
		);

		// Act
		sut("cmd-1");

		// Assert
		expect(warning).not.toHaveBeenCalled();
	});
});
