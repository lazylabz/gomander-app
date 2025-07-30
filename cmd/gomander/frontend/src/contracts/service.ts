import {
  AddCommand,
  EditCommand,
  GetAvailableProjects,
  GetCommandGroups,
  GetCommands,
  GetCurrentProject,
  GetUserConfig,
  RemoveCommand,
  RunCommand,
  SaveCommandGroups,
  SaveExtraPaths,
  StopCommand,
} from "../../wailsjs/go/app/App";
import { project } from "../../wailsjs/go/models.ts";
import { EventsOff, EventsOn } from "../../wailsjs/runtime";

export const dataService = {
  addCommand: AddCommand,
  editCommand: EditCommand,
  getAvailableProjects: GetAvailableProjects,
  getCommandGroups: GetCommandGroups,
  getCommands: GetCommands,
  getCurrentProject: GetCurrentProject as () => Promise<project.Project | null>,
  getUserConfig: GetUserConfig,
  removeCommand: RemoveCommand,
  runCommand: RunCommand,
  saveCommandGroups: SaveCommandGroups,
  saveExtraPaths: SaveExtraPaths,
  stopCommand: StopCommand,
};

export const eventService = {
  eventsOn: EventsOn,
  eventsOff: EventsOff,
};
