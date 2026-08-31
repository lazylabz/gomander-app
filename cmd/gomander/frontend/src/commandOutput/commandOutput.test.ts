import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import type { RecordingTerminals } from "@/commandOutput/adapters/recording.ts";
import {
	appendCommandLogEntry,
	appendCommandOutput,
	attachCommandOutput,
	commandOutputTail,
	disposeCommandOutput,
	FLUSH_INTERVAL_MS,
	LOG_TAIL_SIZE,
	resetCommandOutput,
	resetCommandOutputForNewRun,
	setCommandOutputTheme,
} from "@/commandOutput/commandOutput.ts";
import { TERMINAL_SCROLLBACK } from "@/commandOutput/ports.ts";
import { LogEntryKind } from "@/contracts/types.ts";
import {
	installRecordingTerminals,
	resetTerminals,
} from "@/testing/terminals.ts";

describe("the command output pipeline", () => {
	let recording: RecordingTerminals;

	const terminalOf = (commandId: string) => {
		const terminal = recording.terminals.get(commandId);
		if (!terminal) {
			throw new Error(`No terminal was created for ${commandId}`);
		}
		return terminal;
	};

	const flush = () => vi.advanceTimersByTime(FLUSH_INTERVAL_MS);

	beforeEach(() => {
		vi.useFakeTimers();
		vi.setSystemTime(new Date("2024-03-05T09:07:02"));
		recording = installRecordingTerminals();
	});

	afterEach(() => {
		resetTerminals();
		vi.useRealTimers();
	});

	describe("flushing", () => {
		it("Should hold an appended line until the flush interval elapses", () => {
			// Arrange
			attachCommandOutput("cmd-1", document.createElement("div"));

			// Act
			appendCommandOutput("cmd-1", ["first"]);
			vi.advanceTimersByTime(FLUSH_INTERVAL_MS - 1);

			// Assert
			expect(terminalOf("cmd-1").written).toEqual([]);
		});

		it("Should write every line appended within one interval in a single flush", () => {
			// Arrange
			attachCommandOutput("cmd-1", document.createElement("div"));

			// Act
			appendCommandOutput("cmd-1", ["first"]);
			appendCommandOutput("cmd-1", ["second", "third"]);
			flush();

			// Assert
			expect(terminalOf("cmd-1").written).toEqual([
				"\x1b[90m[09:07:02]\x1b[0m first",
				"\x1b[90m[09:07:02]\x1b[0m second",
				"\x1b[90m[09:07:02]\x1b[0m third",
			]);
		});

		it("Should keep the output of each command on its own terminal", () => {
			// Arrange
			attachCommandOutput("cmd-1", document.createElement("div"));
			attachCommandOutput("cmd-2", document.createElement("div"));

			// Act
			appendCommandOutput("cmd-1", ["mine"]);
			appendCommandOutput("cmd-2", ["yours"]);
			flush();

			// Assert
			expect(terminalOf("cmd-1").written).toEqual([
				"\x1b[90m[09:07:02]\x1b[0m mine",
			]);
			expect(terminalOf("cmd-2").written).toEqual([
				"\x1b[90m[09:07:02]\x1b[0m yours",
			]);
		});

		it("Should ignore an append that carries no lines", () => {
			// Arrange
			attachCommandOutput("cmd-1", document.createElement("div"));

			// Act
			appendCommandOutput("cmd-1", []);
			flush();

			// Assert
			expect(terminalOf("cmd-1").written).toEqual([]);
			expect(commandOutputTail("cmd-1")).toEqual([]);
		});

		it("Should keep flushing after a quiet spell", () => {
			// Arrange the loop into whatever state a long silence leaves it in
			attachCommandOutput("cmd-1", document.createElement("div"));
			appendCommandOutput("cmd-1", ["first"]);
			vi.advanceTimersByTime(FLUSH_INTERVAL_MS * 100);

			// Act
			appendCommandOutput("cmd-1", ["second"]);
			flush();

			// Assert
			expect(terminalOf("cmd-1").written).toEqual([
				"\x1b[90m[09:07:02]\x1b[0m first",
				"\x1b[90m[09:07:05]\x1b[0m second",
			]);
		});

		it("Should stamp a later flush with the time it happened", () => {
			// Arrange
			attachCommandOutput("cmd-1", document.createElement("div"));

			// Act
			appendCommandOutput("cmd-1", ["first"]);
			flush();
			vi.setSystemTime(new Date("2024-03-05T09:07:41"));
			appendCommandOutput("cmd-1", ["second"]);
			flush();

			// Assert
			expect(terminalOf("cmd-1").written).toEqual([
				"\x1b[90m[09:07:02]\x1b[0m first",
				"\x1b[90m[09:07:41]\x1b[0m second",
			]);
		});
	});

	describe("log entries", () => {
		it("Should write a command log entry in bold cyan", () => {
			// Arrange
			attachCommandOutput("cmd-1", document.createElement("div"));

			// Act
			appendCommandLogEntry("cmd-1", "pnpm dev", LogEntryKind.COMMAND);
			flush();

			// Assert
			expect(terminalOf("cmd-1").written).toEqual([
				"\x1b[90m[09:07:02]\x1b[0m \x1b[1;36mpnpm dev\x1b[0m",
			]);
		});

		it("Should write an output log entry as it arrived", () => {
			// Arrange
			attachCommandOutput("cmd-1", document.createElement("div"));

			// Act
			appendCommandLogEntry("cmd-1", "ready in 42 ms", LogEntryKind.OUTPUT);
			flush();

			// Assert
			expect(terminalOf("cmd-1").written).toEqual([
				"\x1b[90m[09:07:02]\x1b[0m ready in 42 ms",
			]);
		});
	});

	describe("the tail", () => {
		it("Should hold an appended line before it has been flushed", () => {
			// Act
			appendCommandOutput("cmd-1", ["not flushed yet"]);

			// Assert
			expect(commandOutputTail("cmd-1")).toEqual(["not flushed yet"]);
		});

		it("Should hold the raw lines, without the flush timestamp", () => {
			// Act
			appendCommandOutput("cmd-1", ["raw line"]);
			flush();

			// Assert
			expect(commandOutputTail("cmd-1")).toEqual(["raw line"]);
		});

		it("Should keep only the last lines of a long run", () => {
			// Act
			appendCommandOutput(
				"cmd-1",
				Array.from({ length: LOG_TAIL_SIZE + 5 }, (_, i) => `line ${i}`),
			);

			// Assert
			expect(commandOutputTail("cmd-1")).toHaveLength(LOG_TAIL_SIZE);
			expect(commandOutputTail("cmd-1").at(-1)).toBe(
				`line ${LOG_TAIL_SIZE + 4}`,
			);
		});

		it("Should be empty for a command that has produced no output", () => {
			// Assert
			expect(commandOutputTail("cmd-1")).toEqual([]);
		});
	});

	describe("attaching", () => {
		it("Should write the lines that arrived before the terminal was attached", () => {
			// Arrange
			appendCommandOutput("cmd-1", ["early"]);
			flush();

			// Act
			attachCommandOutput("cmd-1", document.createElement("div"));

			// Assert
			expect(terminalOf("cmd-1").written).toEqual([
				"\x1b[90m[09:07:02]\x1b[0m early",
			]);
		});

		it("Should backfill a terminal only once", () => {
			// Arrange
			appendCommandOutput("cmd-1", ["early"]);
			flush();
			const attached = attachCommandOutput(
				"cmd-1",
				document.createElement("div"),
			);
			attached.detach();

			// Act
			attachCommandOutput("cmd-1", document.createElement("div"));

			// Assert
			expect(terminalOf("cmd-1").written).toEqual([
				"\x1b[90m[09:07:02]\x1b[0m early",
			]);
		});

		it("Should backfill only once the terminal is on screen", () => {
			// Arrange
			appendCommandOutput("cmd-1", ["early"]);
			flush();

			// Act
			attachCommandOutput("cmd-1", document.createElement("div"));

			// Assert
			expect(terminalOf("cmd-1").writtenWhenAttached).toBe(0);
		});

		it("Should hold at most a terminal's worth of lines for a command never attached", () => {
			// Arrange
			appendCommandOutput(
				"cmd-1",
				Array.from({ length: TERMINAL_SCROLLBACK + 5 }, (_, i) => `line ${i}`),
			);
			flush();

			// Act
			attachCommandOutput("cmd-1", document.createElement("div"));

			// Assert
			const { written } = terminalOf("cmd-1");
			expect(written).toHaveLength(TERMINAL_SCROLLBACK);
			expect(written.at(-1)).toContain(`line ${TERMINAL_SCROLLBACK + 4}`);
		});

		it("Should attach the terminal to the given element", () => {
			// Arrange
			const element = document.createElement("div");

			// Act
			attachCommandOutput("cmd-1", element);

			// Assert
			expect(terminalOf("cmd-1").attachedTo).toBe(element);
		});
	});

	describe("resetting", () => {
		it("Should reset the terminal", () => {
			// Arrange
			attachCommandOutput("cmd-1", document.createElement("div"));
			appendCommandOutput("cmd-1", ["from the previous run"]);
			flush();

			// Act
			resetCommandOutput("cmd-1");

			// Assert
			expect(terminalOf("cmd-1").resets).toBe(1);
			expect(terminalOf("cmd-1").written).toEqual([]);
		});

		it("Should drop the lines of the previous run that had not been flushed", () => {
			// Arrange
			attachCommandOutput("cmd-1", document.createElement("div"));
			appendCommandOutput("cmd-1", ["from the previous run"]);

			// Act
			resetCommandOutput("cmd-1");
			flush();

			// Assert
			expect(terminalOf("cmd-1").written).toEqual([]);
		});

		it("Should drop the backfill of a terminal that was never attached", () => {
			// Arrange
			appendCommandOutput("cmd-1", ["from the previous run"]);
			flush();

			// Act
			resetCommandOutput("cmd-1");
			attachCommandOutput("cmd-1", document.createElement("div"));

			// Assert
			expect(terminalOf("cmd-1").written).toEqual([]);
		});

		it("Should clear the tail", () => {
			// Arrange
			appendCommandOutput("cmd-1", ["from the previous run"]);

			// Act
			resetCommandOutput("cmd-1");

			// Assert
			expect(commandOutputTail("cmd-1")).toEqual([]);
		});

		it("Should drop an unflushed opening line too, since the screen was asked to be empty", () => {
			// Arrange
			attachCommandOutput("cmd-1", document.createElement("div"));
			appendCommandLogEntry("cmd-1", "pnpm dev", LogEntryKind.COMMAND);

			// Act
			resetCommandOutput("cmd-1");
			flush();

			// Assert
			expect(terminalOf("cmd-1").written).toEqual([]);
		});

		it("Should leave the output of the other commands alone", () => {
			// Arrange
			attachCommandOutput("cmd-2", document.createElement("div"));
			appendCommandOutput("cmd-2", ["untouched"]);

			// Act
			resetCommandOutput("cmd-1");
			flush();

			// Assert
			expect(commandOutputTail("cmd-2")).toEqual(["untouched"]);
			expect(terminalOf("cmd-2").written).toEqual([
				"\x1b[90m[09:07:02]\x1b[0m untouched",
			]);
		});
	});

	describe("resetting for a new run", () => {
		const openingLine = "\x1b[90m[09:07:02]\x1b[0m \x1b[1;36mpnpm dev\x1b[0m";

		it("Should keep the line naming the command that is starting", () => {
			// Arrange: the runner sends the opening line before it announces the
			// process started, so it is buffered when the reset arrives.
			attachCommandOutput("cmd-1", document.createElement("div"));
			appendCommandLogEntry("cmd-1", "pnpm dev", LogEntryKind.COMMAND);

			// Act
			resetCommandOutputForNewRun("cmd-1");
			flush();

			// Assert
			expect(terminalOf("cmd-1").written).toEqual([openingLine]);
		});

		it("Should drop the unflushed output of the run that just ended", () => {
			// Arrange
			attachCommandOutput("cmd-1", document.createElement("div"));
			appendCommandOutput("cmd-1", ["from the previous run"]);
			appendCommandLogEntry("cmd-1", "pnpm dev", LogEntryKind.COMMAND);

			// Act
			resetCommandOutputForNewRun("cmd-1");
			flush();

			// Assert
			expect(terminalOf("cmd-1").written).toEqual([openingLine]);
		});

		it("Should keep it for a terminal that is not on screen yet", () => {
			// Arrange
			appendCommandLogEntry("cmd-1", "pnpm dev", LogEntryKind.COMMAND);

			// Act
			resetCommandOutputForNewRun("cmd-1");
			flush();
			attachCommandOutput("cmd-1", document.createElement("div"));

			// Assert
			expect(terminalOf("cmd-1").written).toEqual([openingLine]);
		});

		it("Should reset the terminal of the run that just ended", () => {
			// Arrange
			attachCommandOutput("cmd-1", document.createElement("div"));
			appendCommandOutput("cmd-1", ["from the previous run"]);
			flush();

			// Act
			resetCommandOutputForNewRun("cmd-1");

			// Assert
			expect(terminalOf("cmd-1").resets).toBe(1);
			expect(terminalOf("cmd-1").written).toEqual([]);
		});
	});

	describe("disposing", () => {
		it("Should dispose the terminal and forget the tail", () => {
			// Arrange
			attachCommandOutput("cmd-1", document.createElement("div"));
			appendCommandOutput("cmd-1", ["some output"]);
			flush();

			// Act
			disposeCommandOutput("cmd-1");

			// Assert
			expect(terminalOf("cmd-1").disposed).toBe(true);
			expect(commandOutputTail("cmd-1")).toEqual([]);
		});

		it("Should build a fresh terminal if the command is attached again", () => {
			// Arrange
			attachCommandOutput("cmd-1", document.createElement("div"));
			disposeCommandOutput("cmd-1");

			// Act
			attachCommandOutput("cmd-1", document.createElement("div"));

			// Assert
			expect(terminalOf("cmd-1").disposed).toBe(false);
		});
	});

	describe("theming", () => {
		it("Should apply a new theme to the terminals that already exist", () => {
			// Arrange
			attachCommandOutput("cmd-1", document.createElement("div"));

			// Act
			setCommandOutputTheme("light");

			// Assert
			expect(terminalOf("cmd-1").theme).toBe("light");
		});

		it("Should build later terminals with the current theme", () => {
			// Arrange
			setCommandOutputTheme("light");

			// Act
			attachCommandOutput("cmd-1", document.createElement("div"));

			// Assert
			expect(terminalOf("cmd-1").theme).toBe("light");
		});
	});
});
