import type { config, project } from "../../wailsjs/go/models.ts";
import { event } from "../../wailsjs/go/models.ts";

// Types
export type Command = project.Command;
export type UserConfig = config.UserConfig;
export type CommandGroup = project.CommandGroup;
export type Project = project.Project;
export type ProjectInfo = Pick<Project, "id" | "name" | "baseWorkingDirectory">;

// Enums
export enum Event {
  GET_COMMANDS = event.Event.GET_COMMANDS,
  NEW_LOG_ENTRY = event.Event.NEW_LOG_ENTRY,
  SUCCESS_NOTIFICATION = event.Event.SUCCESS_NOTIFICATION,
  ERROR_NOTIFICATION = event.Event.ERROR_NOTIFICATION,
  PROCESS_FINISHED = event.Event.PROCESS_FINISHED,
  PROCESS_STARTED = event.Event.PROCESS_STARTED,
  GET_USER_CONFIG = event.Event.GET_USER_CONFIG,
  GET_COMMAND_GROUPS = event.Event.GET_COMMAND_GROUPS,
}

export type EventData = {
  [Event.GET_COMMANDS]: null;
  [Event.NEW_LOG_ENTRY]: {
    id: string;
    line: string;
  };
  [Event.ERROR_NOTIFICATION]: string;
  [Event.PROCESS_FINISHED]: string;
  [Event.PROCESS_STARTED]: string;
  [Event.SUCCESS_NOTIFICATION]: string;
  [Event.GET_USER_CONFIG]: null;
  [Event.GET_COMMAND_GROUPS]: null;
};
