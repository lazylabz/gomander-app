import { createContext, useContext, useEffect, useMemo, useState } from "react";
import { toast } from "sonner";

import { Event, type EventData } from "@/types/contracts.ts";

import {
  AddCommand,
  EditCommand,
  GetCommands,
  RemoveCommand,
  RunCommand,
  StopCommand,
} from "../../wailsjs/go/main/App";
import { EventsOff, EventsOn } from "../../wailsjs/runtime";
import type { Command } from "../types/contracts";

export enum CommandStatus {
  IDLE,
  RUNNING,
}

type DataContextValue = {
  // State
  commands: Record<string, Command>;
  commandsStatus: Record<string, CommandStatus>;
  activeCommandId: string | null;
  setActiveCommandId: (commandId: string | null) => void;
  // Command status updates
  setCommandStatus: (commandId: string, status: CommandStatus) => void;
  // Logs
  currentLogs: string[];
  clearCurrentLogs: () => void;
  // Handlers
  createCommand: (command: Command) => Promise<void>;
  editCommand: (command: Command) => Promise<void>;
  execCommand: (commandId: string) => Promise<void>;
  deleteCommand: (commandId: string) => Promise<void>;
  stopRunningCommand: (commandId: string) => Promise<void>;
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
  const [commandsStatus, setCommandsStatus] = useState<
    Record<string, CommandStatus>
  >({});

  const [logs, setLogs] = useState<Record<string, string[]>>({});
  const [activeCommandId, setActiveCommandId] = useState<string | null>(null);

  const refreshCommands = async () => {
    const commandsData = await GetCommands();

    setCommands(commandsData);
    setCommandsStatus(
      Object.fromEntries(
        Object.keys(commandsData).map((id) => [
          id,
          commandsStatus[id] || CommandStatus.IDLE,
        ]),
      ),
    );
  };

  // Command CRUD operations

  const createCommand = async (command: Command) => {
    await AddCommand(command);
  };

  const execCommand = async (commandId: string) => {
    setLogs((prev) => ({
      ...prev,
      [commandId]: [], // Reset logs for the command being executed
    }));
    await RunCommand(commandId);
  };

  const deleteCommand = async (commandId: string) => {
    await RemoveCommand(commandId);
  };

  const editCommand = async (command: Command) => {
    await EditCommand(command);
  };

  const stopRunningCommand = async (commandId: string) => {
    await StopCommand(commandId);
  };

  // Log handlers
  const currentLogs = useMemo(() => {
    return logs[activeCommandId ?? ""] || [];
  }, [logs, activeCommandId]);

  const clearCurrentLogs = () => {
    if (!activeCommandId) {
      return;
    }
    setLogs((prev) => ({
      ...prev,
      [activeCommandId]: [],
    }));
  };

  // Command status updates
  const setCommandStatus = (commandId: string, status: CommandStatus) => {
    setCommandsStatus((prevStatus) => ({
      ...prevStatus,
      [commandId]: status,
    }));
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

    EventsOn(
      Event.PROCESS_FINISHED,
      (data: EventData[Event.PROCESS_FINISHED]) => {
        setCommandsStatus((prevStatus) => ({
          ...prevStatus,
          [data]: CommandStatus.IDLE, // Reset status to IDLE when process finishes
        }));
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

  const value: DataContextValue = {
    // State
    commands,
    commandsStatus,
    activeCommandId,
    setActiveCommandId,
    // Command status updates
    setCommandStatus,
    // Logs
    currentLogs,
    clearCurrentLogs,
    // Handlers
    createCommand,
    editCommand,
    execCommand,
    deleteCommand,
    stopRunningCommand,
  };

  return <dataContext.Provider value={value}>{children}</dataContext.Provider>;
};

export const useDataContext = () => {
  return useContext(dataContext);
};
