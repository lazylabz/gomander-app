import AnsiToHtml from "ansi-to-html";
import { BrushCleaning } from "lucide-react";

import { useDataContext } from "@/contexts/DataContext.tsx";

const ansiConverter = new AnsiToHtml({
  fg: "#000000",
  bg: "#ffffff",
  escapeXML: true,
  newline: true,
  stream: false,
});

export const LogsScreen = () => {
  const { currentLogs, clearCurrentLogs } = useDataContext();

  return (
    <div className="p-4 overflow-y-auto h-full">
      <BrushCleaning
        onClick={clearCurrentLogs}
        className="fixed top-3 right-6 z-1 text-foreground opacity-25 hover:opacity-100 transition-opacity cursor-pointer"
      />
      {currentLogs.map((log, index) => (
        <div
          dangerouslySetInnerHTML={{ __html: ansiConverter.toHtml(log) }}
          key={index}
        />
      ))}
    </div>
  );
};
