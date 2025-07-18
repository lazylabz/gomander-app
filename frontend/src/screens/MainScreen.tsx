import AnsiToHtml from "ansi-to-html";

import { useDataContext } from "@/contexts/DataContext.tsx";

const ansiConverter = new AnsiToHtml();

export const MainScreen = () => {
  const { currentLogs } = useDataContext();
  
  return (
    <div className="p-4 overflow-y-auto h-full">
      {currentLogs.map((log, index) => (
        <div
          dangerouslySetInnerHTML={{ __html: ansiConverter.toHtml(log) }}
          key={index}
        />
      ))}
    </div>
  );
};
