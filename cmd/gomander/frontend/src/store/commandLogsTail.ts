// Deliberately a plain module rather than a Zustand store: nothing renders this
// state, and it is written on every log flush. A store would either churn a new
// Map ~33 times a second or need setters that never call `set`.

// The missing-executable error is emitted as the process exits, so only the
// tail of a command's output is ever inspected.
export const LOG_TAIL_SIZE = 20;

const tailByCommandId = new Map<string, string[]>();

export const recordLogTail = (commandId: string, lines: string[]): void => {
	const existing = tailByCommandId.get(commandId) ?? [];
	tailByCommandId.set(commandId, [...existing, ...lines].slice(-LOG_TAIL_SIZE));
};

export const getLogTail = (commandId: string): string[] =>
	tailByCommandId.get(commandId) ?? [];

export const clearLogTail = (commandId: string): void => {
	tailByCommandId.delete(commandId);
};
