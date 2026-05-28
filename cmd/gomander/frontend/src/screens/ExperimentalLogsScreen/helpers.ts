const pad = (n: number) => n.toString().padStart(2, "0");

export const formatLogTimestamp = (d: Date): string =>
	`${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`;

// Dim-gray ANSI prefix; raw line keeps its own formatting.
export const prependTimestamp = (line: string, ts: string): string =>
	`\x1b[90m[${ts}]\x1b[0m ${line}`;
