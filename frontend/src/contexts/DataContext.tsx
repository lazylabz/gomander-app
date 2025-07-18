import { createContext, useContext, useEffect, useMemo, useState } from "react";
import { toast } from "sonner";

import { Event, type EventData } from "@/types/contracts.ts";

import {
  AddCommand,
  ExecCommand,
  GetCommands,
  RemoveCommand,
} from "../../wailsjs/go/main/App";
import { EventsOff, EventsOn } from "../../wailsjs/runtime";
import type { Command } from "../types/contracts";

type DataContextValue = {
  // State
  commands: Record<string, Command>;
  activeCommandId: string | null;
  setActiveCommandId: (commandId: string | null) => void;
  currentLogs: string[];
  // Handlers
  createCommand: (command: Command) => Promise<void>;
  execCommand: (commandId: string) => Promise<void>;
  deleteCommand: (commandId: string) => Promise<void>;
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
    setLogs((prev) => ({
      ...prev,
      [commandId]: [], // Reset logs for the command being executed
    }));
    await ExecCommand(commandId);
  };

  const deleteCommand = async (commandId: string) => {
    await RemoveCommand(commandId);
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

    // Clean listeners on all events
    return () =>
      EventsOff(Object.keys(Event)[0], ...Object.values(Event).slice(1));
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
    deleteCommand,
  };

  return <dataContext.Provider value={value}>{children}</dataContext.Provider>;
};

export const useDataContext = () => {
  return useContext(dataContext);
};
