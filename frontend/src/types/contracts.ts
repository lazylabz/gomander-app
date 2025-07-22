import { main } from "../../wailsjs/go/models";

// Types
export type Command = main.Command;
export type UserConfig = main.UserConfig;
export type CommandGroup = main.CommandGroup;

// Enums
export enum Event {
  GET_COMMANDS = main.Event.GET_COMMANDS,
  NEW_LOG_ENTRY = main.Event.NEW_LOG_ENTRY,
  SUCCESS_NOTIFICATION = main.Event.SUCCESS_NOTIFICATION,
  ERROR_NOTIFICATION = main.Event.ERROR_NOTIFICATION,
  PROCESS_FINISHED = main.Event.PROCESS_FINISHED,
  GET_USER_CONFIG = main.Event.GET_USER_CONFIG,
  GET_COMMAND_GROUPS = main.Event.GET_COMMAND_GROUPS,
}

export type EventData = {
  [Event.GET_COMMANDS]: null;
  [Event.NEW_LOG_ENTRY]: {
    id: string;
    line: string;
  };
  [Event.ERROR_NOTIFICATION]: string;
  [Event.PROCESS_FINISHED]: string;
  [Event.SUCCESS_NOTIFICATION]: string;
  [Event.GET_USER_CONFIG]: null;
  [Event.GET_COMMAND_GROUPS]: null;
};
