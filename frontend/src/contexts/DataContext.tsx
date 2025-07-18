import { createContext, useContext, useEffect, useMemo, useState } from "react";
import { toast } from "sonner";

import { Event, type EventData } from "@/types/contracts.ts";

import {
  AddCommand,
  ExecCommand,
  GetCommands,
} from "../../wailsjs/go/main/LogServer";
import { EventsOff, EventsOn } from "../../wailsjs/runtime";
import type { Command } from "../types/contracts";

type DataContextValue = {
  commands: Record<string, Command>;
  activeCommandId: string | null;
  setActiveCommandId: (commandId: string | null) => void;
  createCommand: (command: Command) => Promise<void>;
  execCommand: (commandId: string) => Promise<void>;
  currentLogs: string[];
};

export const dataContext = createContext<DataContextValue>(
  {} as DataContextValue,
);

export const DataContextProvider = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  const [commands, setCommands] = useState<Record<string, Command>>({});
  const [logs, setLogs] = useState<Record<string, string[]>>({});
  const [activeCommandId, setActiveCommandId] = useState<string | null>(null);

  const refreshCommands = async () => {
    const commandsData = await GetCommands();

    setCommands(commandsData);
  };

  // Handlers

  const createCommand = async (command: Command) => {
    await AddCommand(command);
  };

  const execCommand = async (commandId: string) => {
    await ExecCommand(commandId);
  };

  // Register events listeners
  useEffect(() => {
    EventsOn(Event.GET_COMMANDS, () => {
      refreshCommands();
    });

    EventsOn(Event.NEW_LOG_ENTRY, (data: EventData[Event.NEW_LOG_ENTRY]) => {
      const { id, line } = data;

      setLogs((prevLogs) => {
        const newLogs = { ...prevLogs };
        if (!newLogs[id]) {
          newLogs[id] = [];
        }
        newLogs[id].push(line);
        return newLogs;
      });
    });

    EventsOn(
      Event.ERROR_NOTIFICATION,
      (data: EventData[Event.ERROR_NOTIFICATION]) => {
        toast("Error", {
          description: data,
        });
      },
    );

    // Clean listeners on unmount
    return () =>
      EventsOff(
        Event.GET_COMMANDS,
        Event.NEW_LOG_ENTRY,
        Event.ERROR_NOTIFICATION,
        Event.PROCESS_FINISHED,
      );
  });

  // Initial fetch of commands
  useEffect(() => {
    refreshCommands();
  }, []);

  const currentLogs = useMemo(() => {
    return logs[activeCommandId ?? ""] || [];
  }, [logs, activeCommandId]);

  const value: DataContextValue = {
    commands,
    activeCommandId,
    setActiveCommandId,
    createCommand,
    execCommand,
    currentLogs,
  };

  return <dataContext.Provider value={value}>{children}</dataContext.Provider>;
};

export const useDataContext = () => {
  return useContext(dataContext);
};
