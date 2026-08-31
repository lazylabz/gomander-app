import { xtermTerminal } from "@/commandOutput/adapters/xterm.ts";
import {
	type AttachedTerminal,
	type OutputTerminal,
	TERMINAL_SCROLLBACK,
	type TerminalFactory,
	type TerminalTheme,
} from "@/commandOutput/ports.ts";
import { LogEntryKind } from "@/contracts/types.ts";

export const FLUSH_INTERVAL_MS = 30;

// The missing-executable error is emitted as the process exits, so only the
// tail of a command's output is ever inspected.
export const LOG_TAIL_SIZE = 20;

const DEFAULT_THEME: TerminalTheme = "dark";

// A raw line appended since the last flush, with the kind that decides whether
// a reset for a new run keeps it.
type BufferedLine = {
	line: string;
	kind: LogEntryKind;
};

type CommandOutput = {
	terminal: OutputTerminal | null;
	buffered: BufferedLine[];
	// Stamped lines waiting for a terminal to exist.
	backlog: string[];
	tail: string[];
};

// The emulator sits on the same kind of seam as the backend contract: a live
// binding a test swaps for a recording factory. Only tests should set it.
let createTerminal: TerminalFactory = xtermTerminal;

const outputs = new Map<string, CommandOutput>();

let currentTheme: TerminalTheme = DEFAULT_THEME;
let flushTimer: ReturnType<typeof setInterval> | null = null;

const pad = (n: number) => n.toString().padStart(2, "0");

const formatTimestamp = (d: Date): string =>
	`${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`;

// Dim-gray ANSI prefix; the raw line keeps its own formatting.
const prependTimestamp = (line: string, timestamp: string): string =>
	`\x1b[90m[${timestamp}]\x1b[0m ${line}`;

// Bold cyan, so the command the user launched reads apart from the output it
// goes on to print. The runner sends the command text plain; the styling is
// decided here, where the timestamp above is decided too.
const asCommandLine = (line: string): string => `\x1b[1;36m${line}\x1b[0m`;

const outputFor = (commandId: string): CommandOutput => {
	const existing = outputs.get(commandId);
	if (existing) {
		return existing;
	}

	const output: CommandOutput = {
		terminal: null,
		buffered: [],
		backlog: [],
		tail: [],
	};
	outputs.set(commandId, output);
	return output;
};

const flush = (): void => {
	const timestamp = formatTimestamp(new Date());
	let wrote = false;

	for (const output of outputs.values()) {
		if (output.buffered.length === 0) {
			continue;
		}
		wrote = true;

		const stamped = output.buffered.map(({ line }) =>
			prependTimestamp(line, timestamp),
		);
		output.buffered = [];

		if (!output.terminal) {
			// Capped: a chatty command whose terminal is never opened would otherwise
			// hold every line it ever printed, and the emulator would drop the older
			// ones on write anyway.
			output.backlog = [...output.backlog, ...stamped].slice(
				-TERMINAL_SCROLLBACK,
			);
			continue;
		}
		for (const line of stamped) {
			output.terminal.writeln(line);
		}
	}

	if (!wrote) {
		stopFlushing();
	}
};

// The loop only runs while there is something to write: a command's output is
// bursty, and an always-on 30 ms interval wakes the app up for nothing.
const startFlushing = (): void => {
	if (flushTimer !== null) {
		return;
	}
	flushTimer = setInterval(flush, FLUSH_INTERVAL_MS);
};

const stopFlushing = (): void => {
	if (flushTimer === null) {
		return;
	}
	clearInterval(flushTimer);
	flushTimer = null;
};

const append = (commandId: string, entries: BufferedLine[]): void => {
	if (entries.length === 0) {
		return;
	}
	const output = outputFor(commandId);

	// The tail is recorded now rather than on flush: PROCESS_FINISHED arrives
	// between two flushes often enough, and the line that explains the failure is
	// the last one written.
	output.tail = [...output.tail, ...entries.map(({ line }) => line)].slice(
		-LOG_TAIL_SIZE,
	);

	output.buffered.push(...entries);
	startFlushing();
};

export const appendCommandOutput = (commandId: string, lines: string[]): void =>
	append(
		commandId,
		lines.map((line) => ({ line, kind: LogEntryKind.OUTPUT })),
	);

export const appendCommandLogEntry = (
	commandId: string,
	line: string,
	kind: LogEntryKind,
): void =>
	append(commandId, [
		{ line: kind === LogEntryKind.COMMAND ? asCommandLine(line) : line, kind },
	]);

export const attachCommandOutput = (
	commandId: string,
	element: HTMLElement,
): AttachedTerminal => {
	const output = outputFor(commandId);
	output.terminal ??= createTerminal(commandId, currentTheme);

	const attached = output.terminal.attach(element);

	// Backfilled after the emulator is on screen: written earlier, the lines
	// would be laid out against the default terminal size rather than this one.
	for (const line of output.backlog) {
		output.terminal.writeln(line);
	}
	output.backlog = [];

	return attached;
};

const clear = (output: CommandOutput, kept: BufferedLine[]): void => {
	output.buffered = kept;
	output.backlog = [];
	output.tail = [];
	output.terminal?.reset();
};

// Buffered lines are dropped rather than flushed: the user asked for the screen
// to be empty, so nothing written before may land in the cleared terminal.
export const resetCommandOutput = (commandId: string): void => {
	const output = outputs.get(commandId);
	if (output) {
		clear(output, []);
	}
};

// The runner names the command it is about to launch before it announces the
// process started, so that line is already buffered by the time this runs. It
// belongs to the run being started rather than the one that just ended, and is
// the only thing this keeps; everything else is the finished run's output.
export const resetCommandOutputForNewRun = (commandId: string): void => {
	const output = outputs.get(commandId);
	if (output) {
		clear(
			output,
			output.buffered.filter(({ kind }) => kind === LogEntryKind.COMMAND),
		);
	}
};

export const commandOutputTail = (commandId: string): string[] =>
	outputs.get(commandId)?.tail ?? [];

export const disposeCommandOutput = (commandId: string): void => {
	outputs.get(commandId)?.terminal?.dispose();
	outputs.delete(commandId);
};

export const setCommandOutputTheme = (theme: TerminalTheme): void => {
	currentTheme = theme;
	for (const output of outputs.values()) {
		output.terminal?.setTheme(theme);
	}
};

export const setTerminalFactory = (factory: TerminalFactory): void => {
	createTerminal = factory;
};

export const resetTerminalFactory = (): void => {
	createTerminal = xtermTerminal;
};

// Tests only.
export const resetCommandOutputPipeline = (): void => {
	for (const output of outputs.values()) {
		output.terminal?.dispose();
	}
	outputs.clear();
	currentTheme = DEFAULT_THEME;
	stopFlushing();
};
