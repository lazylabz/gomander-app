import { useCallback, useEffect } from "react";

type Key =
	| "A"
	| "B"
	| "C"
	| "D"
	| "E"
	| "F"
	| "G"
	| "H"
	| "I"
	| "J"
	| "K"
	| "L"
	| "M"
	| "N"
	| "O"
	| "P"
	| "Q"
	| "R"
	| "S"
	| "T"
	| "U"
	| "V"
	| "W"
	| "X"
	| "Y"
	| "Z";

type Modifier = "Control" | "Shift" | "Alt" | "Mod";

const ModifierMap: Record<Modifier, (keyof KeyboardEvent)[]> = {
	Control: ["ctrlKey"],
	Shift: ["shiftKey"],
	Alt: ["altKey"],
	Mod: ["metaKey", "ctrlKey"],
};

type Shortcut = `${Modifier}-${Key}` | Key;

export const matchesShortcut = (
	shortcut: Shortcut,
	event: KeyboardEvent,
): boolean => {
	const [modifier, key] = shortcut.includes("-")
		? (shortcut.split("-") as [Modifier, Key])
		: [null, shortcut as Key];

	const keyMatches = event.key.toLowerCase() === key.toLowerCase();
	const modifierMatches =
		!modifier || ModifierMap[modifier].some((mod) => event[mod]);

	return keyMatches && modifierMatches;
};

export const useShortcut = (shortcut: Shortcut, callback: () => void) => {
	const handleKeyDown = useCallback(
		(event: KeyboardEvent) => {
			if (matchesShortcut(shortcut, event)) {
				event.preventDefault();
				callback();
			}
		},
		[shortcut, callback],
	);

	useEffect(() => {
		window.addEventListener("keydown", handleKeyDown);
		return () => {
			window.removeEventListener("keydown", handleKeyDown);
		};
	}, [handleKeyDown]);
};
