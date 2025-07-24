import AnsiToHtml from "ansi-to-html";
import { BrushCleaning } from "lucide-react";

import { useCurrentLogs } from "@/hooks/useCurrentLogs.ts";
import { clearCurrentLogs } from "@/useCases/logging/clearCurrentLogs.ts";

const ansiConverter = new AnsiToHtml({
  fg: "#000000",
  bg: "#ffffff",
  escapeXML: true,
  newline: true,
  stream: false,
});

export const LogsScreen = () => {
  const { currentLogs } = useCurrentLogs();

  return (
    <div className="p-4 overflow-y-auto h-full w-full flex flex-col font-mono justify-end">
      <BrushCleaning
        onClick={clearCurrentLogs}
        className="fixed top-3 right-6 z-1 text-foreground opacity-25 hover:opacity-100 transition-opacity cursor-pointer"
      />
      {currentLogs.map((log, index) => (
        <p
          className="w-full wrap-anywhere"
          dangerouslySetInnerHTML={{ __html: ansiConverter.toHtml(log) }}
          key={index}
        />
      ))}
    </div>
  );
};
