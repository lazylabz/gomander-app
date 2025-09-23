import type { domain } from "../../wailsjs/go/models.ts";
import { event } from "../../wailsjs/go/models.ts";

// Types
export type Command = domain.Command;
export type UserConfig = domain.Config;
export type CommandGroup = domain.CommandGroup;
export type Project = domain.Project;
export type ProjectExport = domain.ProjectExportJSONv1;

// Enums
export enum Event {
  NEW_LOG_ENTRY = event.Event.NEW_LOG_ENTRY,
  PROCESS_FINISHED = event.Event.PROCESS_FINISHED,
  PROCESS_STARTED = event.Event.PROCESS_STARTED,
  COMMAND_GROUP_DELETED = event.Event.COMMAND_GROUP_DELETED,
  COMMAND_FAILED = event.Event.COMMAND_FAILED
}

export type EventData = {
  [Event.NEW_LOG_ENTRY]: {
    id: string;
    line: string;
  };
  [Event.PROCESS_FINISHED]: string;
  [Event.PROCESS_STARTED]: string;
  [Event.COMMAND_GROUP_DELETED]: string;
  [Event.COMMAND_FAILED]: {
    id: string;
    line: string;
    pattern: string[];
  };
};
