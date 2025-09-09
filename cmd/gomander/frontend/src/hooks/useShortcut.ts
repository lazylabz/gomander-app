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
  Mod: ["metaKey", "ctrlKey"], // Meta is Command on Mac, Ctrl on Windows
};

type Shortcut = `${Modifier}-${Key}` | Key;

export const useShortcut = (shortCut: Shortcut, callback: () => void) => {
  const hasModifier = shortCut.includes("-");

  const modifier = hasModifier ? (shortCut.split("-")[0] as Modifier) : null;
  const key = hasModifier ? (shortCut.split("-")[1] as Key) : (shortCut as Key);

  const handleKeyDown = useCallback(
    (event: KeyboardEvent) => {
      // Check if the key is a modifier key
      const keyMatches = event.key.toLowerCase() === key.toLowerCase();
      const modifierMatches =
        !modifier || ModifierMap[modifier].some((mod) => event[mod]);

      if (keyMatches && modifierMatches) {
        event.preventDefault();
        callback();
      }
    },
    [key, modifier, callback],
  );

  useEffect(() => {
    window.addEventListener("keydown", handleKeyDown);
    return () => {
      window.removeEventListener("keydown", handleKeyDown);
    };
  }, [shortCut, callback, handleKeyDown]);
};
