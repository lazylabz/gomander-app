import parse, { type DOMNode, domToReact, Element } from "html-react-parser";
import { BrushCleaning } from "lucide-react";
import { type ChangeEvent, type KeyboardEventHandler, useState } from "react";

import { Input } from "@/components/ui/input.tsx";
import { useCurrentLogs } from "@/hooks/useCurrentLogs.ts";
import { useShortcut } from "@/hooks/useShortcut.ts";
import { cn } from "@/lib/utils.ts";
import { extractMatchesIds, parseLog } from "@/screens/LogsScreen/helpers.ts";
import { clearCurrentLogs } from "@/useCases/logging/clearCurrentLogs.ts";

export const LogsScreen = () => {
  const { currentLogs } = useCurrentLogs();

  const [searchOpen, setSearchOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const [focusedMatchIndex, setFocusedMatchIndex] = useState(0);

  const handleSearchInputChange = (event: ChangeEvent<HTMLInputElement>) => {
    setFocusedMatchIndex(0);
    setSearchQuery(event.target.value);
  };

  const parsedLogs = currentLogs.map((log) =>
    parseLog(log, searchOpen ? searchQuery : ""),
  );

  const matchesIds =
    searchQuery && searchOpen
      ? parsedLogs.flatMap((log) => {
          return extractMatchesIds(log);
        })
      : [];

  useShortcut("Meta-F", () => {
    setSearchOpen(true);
  });

  const focusElementByMatchId = (id: string) => {
    const matchElement = document.querySelector(`[data-match="${id}"]`);
    if (matchElement) {
      matchElement.scrollIntoView({ behavior: "smooth", block: "center" });
    }
  };

  const handleInputKeyPress: KeyboardEventHandler = (event) => {
    if (event.key === "Escape") {
      setSearchOpen(false);
    }
    if (event.key === "ArrowDown") {
      event.preventDefault();
      const newIndex = focusedMatchIndex + 1;
      const correctedIndex = newIndex >= matchesIds.length ? 0 : newIndex;
      setFocusedMatchIndex(correctedIndex);
      focusElementByMatchId(matchesIds[correctedIndex]);
    }
    if (event.key === "ArrowUp") {
      event.preventDefault();
      const newIndex = focusedMatchIndex - 1;
      const correctedIndex = newIndex < 0 ? matchesIds.length - 1 : newIndex;
      setFocusedMatchIndex(correctedIndex);
      focusElementByMatchId(matchesIds[correctedIndex]);
    }
  };

  const focusedMatchId = matchesIds[focusedMatchIndex];

  return (
    <div className="p-4 overflow-y-auto h-full w-full flex flex-col font-mono justify-end">
      <div className="fixed top-3 right-6 z-1 flex items-center gap-2">
        {searchOpen && (
          <div className="flex flex-col">
            <Input
              autoCorrect="off"
              autoComplete="off"
              className="bg-background"
              value={searchQuery}
              onChange={handleSearchInputChange}
              onKeyDown={handleInputKeyPress}
            />
            {searchQuery && (
              <span className="text-xs text-muted-foreground">
                {matchesIds.length} matches
              </span>
            )}
          </div>
        )}
        <BrushCleaning
          onClick={clearCurrentLogs}
          className="text-foreground opacity-25 hover:opacity-100 transition-opacity cursor-pointer"
        />
      </div>
      {parsedLogs.map((log, index) => (
        <p className="w-full wrap-anywhere" key={index}>
          {parse(log, {
            replace: (domNode) => {
              if (domNode instanceof Element && domNode.attribs) {
                if (domNode.attribs.class?.includes("match")) {
                  const dataMatch = domNode.attribs["data-match"] || "";
                  const isMatch = dataMatch === focusedMatchId;

                  return (
                    <span
                      data-match={dataMatch}
                      className={cn(
                        "match bg-yellow-100",
                        isMatch && "bg-yellow-300",
                      )}
                    >
                      {domToReact(domNode.children as DOMNode[])}
                    </span>
                  );
                }
              }
            },
          })}
        </p>
      ))}
    </div>
  );
};
