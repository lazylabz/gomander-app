import { useDataContext } from "@/contexts/DataContext.tsx";

export const MainScreen = () => {
  const { currentLogs } = useDataContext();
  return (
    <div>
      <div>
        {currentLogs.map((log, index) => (
          <div key={index}>{log}</div>
        ))}
      </div>
    </div>
  );
};
