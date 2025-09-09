import parse, { type DOMNode, domToReact, Element } from "html-react-parser";
import { BrushCleaning, ChevronLeft, ChevronRight, X } from "lucide-react";
import {
  type ChangeEvent,
  type KeyboardEventHandler,
  useEffect,
  useRef,
  useState,
} from "react";

import { Input } from "@/components/ui/input.tsx";
import { SidebarTrigger, useSidebar } from "@/components/ui/sidebar.tsx";
import { useTheme } from "@/contexts/theme.tsx";
import { useCurrentLogs } from "@/hooks/useCurrentLogs.ts";
import { useShortcut } from "@/hooks/useShortcut.ts";
import { cn } from "@/lib/utils.ts";
import { extractMatchesIds, parseLog } from "@/screens/LogsScreen/helpers.ts";
import { clearCurrentLogs } from "@/useCases/logging/clearCurrentLogs.ts";

const focusElementByMatchId = (id: string) => {
  const matchElement = document.querySelector(`[data-match="${id}"]`);
  if (matchElement) {
    matchElement.scrollIntoView({ behavior: "smooth", block: "center" });
  }
};

export const LogsScreen = () => {
  const { currentLogs } = useCurrentLogs();

  const { theme } = useTheme();
  const { open, isMobile } = useSidebar();
  const isDesktopSidebarOpen = !isMobile && open;

  const searchInput = useRef<HTMLInputElement | null>(null);

  const [searchOpen, setSearchOpen] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");
  const [focusedMatchIndex, setFocusedMatchIndex] = useState(0);

  const openSearch = () => setSearchOpen(true);
  const closeSearch = () => setSearchOpen(false);

  const handleSearchInputChange = (event: ChangeEvent<HTMLInputElement>) => {
    setFocusedMatchIndex(0);
    setSearchQuery(event.target.value);
  };

  const parsedLogs = currentLogs.map((log) =>
    parseLog(log, searchOpen ? searchQuery : "", theme),
  );

  const matchesIds =
    searchQuery && searchOpen ? parsedLogs.flatMap(extractMatchesIds) : [];

  useShortcut("Mod-F", () => {
    openSearch();
  });

  const nextMatch = () => {
    const newIndex = focusedMatchIndex + 1;
    const correctedIndex = newIndex >= matchesIds.length ? 0 : newIndex;
    setFocusedMatchIndex(correctedIndex);
    focusElementByMatchId(matchesIds[correctedIndex]);
  };
  const prevMatch = () => {
    const newIndex = focusedMatchIndex - 1;
    const correctedIndex = newIndex < 0 ? matchesIds.length - 1 : newIndex;
    setFocusedMatchIndex(correctedIndex);
    focusElementByMatchId(matchesIds[correctedIndex]);
  };

  const handleInputKeyPress: KeyboardEventHandler = (event) => {
    if (event.key === "Escape") {
      setSearchOpen(false);
    }
    if (event.key === "ArrowDown") {
      event.preventDefault();
      nextMatch();
    }
    if (event.key === "ArrowUp") {
      event.preventDefault();
      prevMatch();
    }
  };

  const focusedMatchId = matchesIds[focusedMatchIndex];

  useEffect(() => {
    if (searchOpen) {
      searchInput.current?.focus();
    }
  }, [searchOpen]);

  return (
    <div className="p-4 overflow-y-auto h-full w-full flex flex-col font-mono justify-end bg-background text-foreground">
      <div
        className={cn(
          "fixed top-3.5 z-10 bg-background rounded-sm transition-[left] duration-200 ease-linear",
          isDesktopSidebarOpen ? "left-[16.5rem]" : "left-2",
        )}
      >
        <SidebarTrigger className="opacity-25 hover:opacity-100 transition-opacity" />
      </div>
      <div className="fixed top-3 right-3 z-1 flex items-center gap-2">
        {searchOpen && (
          <div className="flex flex-col bg-background gap-1.5">
            <Input
              ref={searchInput}
              autoComplete="off"
              autoCorrect="off"
              autoCapitalize="off"
              className="w-64"
              value={searchQuery}
              onChange={handleSearchInputChange}
              onKeyDown={handleInputKeyPress}
            />
            <span className="text-xs text-muted-foreground pl-2 flex items-center gap-2 pb-1 justify-between select-none">
              <div className="flex flex-row items-center gap-2">
                <div className="flex flex-row items-center">
                  <ChevronLeft
                    className="text-muted-foreground hover:text-foreground cursor-pointer"
                    onClick={prevMatch}
                    size={14}
                  />
                  <ChevronRight
                    className="text-muted-foreground hover:text-foreground cursor-pointer"
                    onClick={nextMatch}
                    size={14}
                  />
                </div>
                {matchesIds.length} matches
              </div>
              <X
                size={14}
                onClick={closeSearch}
                className="text-muted-foreground hover:text-foreground cursor-pointer"
              />
            </span>
          </div>
        )}
        <BrushCleaning
          onClick={clearCurrentLogs}
          className="text-foreground opacity-25 hover:opacity-100 transition-opacity cursor-pointer self-start mt-1.5 ml-2 size-5"
        />
      </div>
      {parsedLogs.map((log, index) => (
        <p className="w-full wrap-anywhere" key={index}>
          {parse(log, {
            replace: (domNode) => {
              if (!(domNode instanceof Element && domNode.attribs)) {
                return;
              }
              if (!domNode.attribs.class?.includes("match")) {
                return;
              }

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
            },
          })}
        </p>
      ))}
    </div>
  );
};
