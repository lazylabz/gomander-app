import { useMemo } from "react";

import { useCommandStore } from "@/store/commandStore.ts";

export const useCurrentLogs = () => {
  const activeCommandId = useCommandStore((state) => state.activeCommandId);
  const commandsLogs = useCommandStore((state) => state.commandsLogs);

  const currentLogs = useMemo(
    () => (activeCommandId ? commandsLogs[activeCommandId] || [] : []),
    [activeCommandId, commandsLogs],
  );

  return {
    currentLogs,
  };
};
