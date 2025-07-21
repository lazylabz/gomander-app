import { createContext, useContext, useEffect, useMemo, useState } from "react";
import { toast } from "sonner";

import { Event, type EventData, type UserConfig } from "@/types/contracts.ts";

import {
  AddCommand,
  EditCommand,
  GetCommands,
  GetUserConfig,
  RemoveCommand,
  RunCommand,
  SaveUserConfig,
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
  // User config
  userConfig: UserConfig;
  saveUserConfig: (config: UserConfig) => Promise<void>;
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

  const [userConfig, setUserConfig] = useState<UserConfig>({
    extraPaths: [],
  });

  // Computed values
  const currentLogs = useMemo(() => {
    return logs[activeCommandId ?? ""] || [];
  }, [logs, activeCommandId]);

  // Command CRUD operations
  const fetchCommands = async () => {
    const commandsData = await GetCommands();

    setCommands(commandsData);
    setCommandsStatus((prevStatus) => {
      return Object.fromEntries(
        Object.keys(commandsData).map((id) => [
          id,
          prevStatus[id] || CommandStatus.IDLE,
        ]),
      );
    });
  };

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

  // User config operations
  const fetchUserConfig = async (): Promise<void> => {
    const config = await GetUserConfig();
    
    setUserConfig(config);
  };

  const saveUserConfig = async (config: UserConfig): Promise<void> => {
    // Assuming there's a function to set user config
    // This function should be implemented in the backend
    await SaveUserConfig(config);
  };

  // Handlers
  const clearCurrentLogs = () => {
    if (!activeCommandId) {
      return;
    }
    setLogs((prev) => ({
      ...prev,
      [activeCommandId]: [],
    }));
  };

  const setCommandStatus = (commandId: string, status: CommandStatus) => {
    setCommandsStatus((prevStatus) => ({
      ...prevStatus,
      [commandId]: status,
    }));
  };

  // Register events listeners
  useEffect(() => {
    EventsOn(Event.GET_COMMANDS, () => {
      fetchCommands();
    });

    EventsOn(Event.GET_USER_CONFIG, () => {
      fetchUserConfig();
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
        toast.error("Error", {
          description: data,
        });
      },
    );

    EventsOn(
      Event.SUCCESS_NOTIFICATION,
      (data: EventData[Event.SUCCESS_NOTIFICATION]) => {
        toast.success(data);
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

  // Initial fetch of data
  useEffect(() => {
    fetchCommands();
    fetchUserConfig();
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
    // User config
    userConfig,
    saveUserConfig,
  };

  return <dataContext.Provider value={value}>{children}</dataContext.Provider>;
};

export const useDataContext = () => {
  return useContext(dataContext);
};
