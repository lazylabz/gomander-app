import { afterEach, beforeEach, describe, expect, it } from "vitest";

import type { RecordingTerminals } from "@/commandOutput/adapters/recording.ts";

import {
	appendCommandOutput,
	attachCommandOutput,
	commandOutputTail,
} from "@/commandOutput/commandOutput.ts";
import { commandStore } from "@/store/commandStore.ts";
import { resetStores } from "@/testing/stores.ts";
import {
	installRecordingTerminals,
	resetTerminals,
} from "@/testing/terminals.ts";
import { clearCurrentLogs } from "@/useCases/logging/clearCurrentLogs.ts";

describe("clearCurrentLogs", () => {
	const sut = clearCurrentLogs;

	let recording: RecordingTerminals;

	beforeEach(() => {
		resetStores();
		recording = installRecordingTerminals();
	});

	afterEach(() => {
		resetTerminals();
	});

	it("Should clear the terminal and the buffered logs of the active command", () => {
		// Arrange
		commandStore.setState({ activeCommandId: "cmd-1" });
		attachCommandOutput("cmd-1", document.createElement("div"));
		appendCommandOutput("cmd-1", ["a log line"]);

		// Act
		sut();

		// Assert
		expect(recording.terminals.get("cmd-1")?.resets).toBe(1);
		expect(commandOutputTail("cmd-1")).toEqual([]);
	});

	it("Should leave the logs of the other commands alone", () => {
		// Arrange
		commandStore.setState({ activeCommandId: "cmd-1" });
		attachCommandOutput("cmd-2", document.createElement("div"));
		appendCommandOutput("cmd-2", ["another log line"]);

		// Act
		sut();

		// Assert
		expect(recording.terminals.get("cmd-2")?.resets).toBe(0);
		expect(commandOutputTail("cmd-2")).toEqual(["another log line"]);
	});

	it("Should clear nothing when no command is active", () => {
		// Arrange
		attachCommandOutput("cmd-1", document.createElement("div"));
		appendCommandOutput("cmd-1", ["a log line"]);

		// Act
		sut();

		// Assert
		expect(recording.terminals.get("cmd-1")?.resets).toBe(0);
		expect(commandOutputTail("cmd-1")).toEqual(["a log line"]);
	});
});
