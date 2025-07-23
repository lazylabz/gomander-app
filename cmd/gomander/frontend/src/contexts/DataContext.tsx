import {
  createContext,
  type Dispatch,
  type SetStateAction,
  useContext,
  useEffect,
  useMemo,
  useState,
} from "react";
import { toast } from "sonner";

import { type CommandGroup } from "@/types/contracts.ts";

import {
  AddCommand,
  EditCommand,
  GetCommandGroups,
  GetCommands,
  RemoveCommand,
  RunCommand,
  SaveCommandGroups,
  StopCommand,
} from "../../wailsjs/go/app/App";
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
  duplicateCommand: (command: Command) => Promise<void>;
  // Command groups
  commandGroups: CommandGroup[];
  saveCommandGroups: (groups: CommandGroup[]) => Promise<void>;
  runCommandGroup: (groupId: string) => Promise<void>;
  stopCommandGroup: (groupId: string) => Promise<void>;
  fetchCommands: () => Promise<void>;
  fetchCommandGroups: () => Promise<void>;
  setLogs: Dispatch<SetStateAction<Record<string, string[]>>>;
  setCommandsStatus: Dispatch<SetStateAction<Record<string, CommandStatus>>>;
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
  const [commandGroups, setCommandGroups] = useState<CommandGroup[]>([]);

  const [logs, setLogs] = useState<Record<string, string[]>>({});
  const [activeCommandId, setActiveCommandId] = useState<string | null>(null);

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
    setCommandsStatus((prevStatus) => ({
      ...prevStatus,
      [commandId]: CommandStatus.RUNNING,
    }));
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

  const duplicateCommand = async (command: Command) => {
    const newCommand = {
      ...command,
      name: `${command.name} (copy)`,
      id: crypto.randomUUID(),
    };
    await AddCommand(newCommand);
  };

  // Command groups operations
  const fetchCommandGroups = async (): Promise<void> => {
    const groups = await GetCommandGroups();
    setCommandGroups(groups);
  };

  const saveCommandGroups = async (groups: CommandGroup[]): Promise<void> => {
    // Optimistic save to avoid flickering while drag and dropping
    const prev = commandGroups;
    setCommandGroups(groups);
    console.log({
      prev,
      groups,
    });
    try {
      await SaveCommandGroups(groups);
    } catch {
      // If saving fails, revert to previous state
      setCommandGroups(prev);
      toast.error("Failed to save command groups");
    }
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

  const runCommandGroup = async (groupId: string) => {
    const group = commandGroups.find((g) => g.id === groupId);
    if (!group) {
      toast.error("Command group not found");
      return;
    }

    const notRunningCommands = group.commands.filter(
      (cmdId) => commandsStatus[cmdId] !== CommandStatus.RUNNING,
    );

    await Promise.all(
      notRunningCommands.map(async (cmdId) => {
        await execCommand(cmdId);
      }),
    );
  };

  const stopCommandGroup = async (groupId: string) => {
    const group = commandGroups.find((g) => g.id === groupId);
    if (!group) {
      toast.error("Command group not found");
      return;
    }

    const runningCommands = group.commands.filter(
      (cmdId) => commandsStatus[cmdId] === CommandStatus.RUNNING,
    );

    await Promise.all(
      runningCommands.map(async (cmdId) => {
        await stopRunningCommand(cmdId);
      }),
    );
  };

  // Initial fetch of data
  useEffect(() => {
    fetchCommands();
    fetchCommandGroups();
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
    duplicateCommand,
    // Command groups
    commandGroups,
    saveCommandGroups,
    runCommandGroup,
    stopCommandGroup,
    fetchCommands,
    fetchCommandGroups,
    setLogs,
    setCommandsStatus,
  };

  return <dataContext.Provider value={value}>{children}</dataContext.Provider>;
};

export const useDataContext = () => {
  return useContext(dataContext);
};
