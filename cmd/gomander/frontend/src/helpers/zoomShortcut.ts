import type { ZoomDirection } from "@/contracts/types.ts";

const DIRECTION_BY_KEY: Record<string, ZoomDirection> = {
	"=": "in",
	"+": "in",
	"-": "out",
	"0": "reset",
};

// Ctrl alone: outside macOS that is the zoom modifier, and Ctrl+Alt is how
// several layouts type a character rather than a shortcut.
export const zoomShortcutOf = (event: KeyboardEvent): ZoomDirection | null => {
	if (!event.ctrlKey || event.altKey || event.metaKey) {
		return null;
	}

	return DIRECTION_BY_KEY[event.key] ?? null;
};
