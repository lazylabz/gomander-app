import {
  AddCommand,
  CloseProject,
  CreateProject,
  DeleteProject,
  EditCommand,
  ExportProject,
  GetAvailableProjects,
  GetCommandGroups,
  GetCommands,
  GetCurrentProject,
  GetUserConfig,
  ImportProject,
  OpenProject,
  RemoveCommand,
  RunCommand,
  SaveCommandGroups,
  SaveUserConfig,
  StopCommand,
} from "../../wailsjs/go/app/App";
import type { project } from "../../wailsjs/go/models.ts";
import { GetComputedPath } from "../../wailsjs/go/uihelpers/UiPathHelper";
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
  saveUserConfig: SaveUserConfig,
  createProject: CreateProject,
  stopCommand: StopCommand,
  openProject: OpenProject,
  closeProject: CloseProject,
  deleteProject: DeleteProject,
  exportProject: ExportProject,
  importProject: ImportProject,
};

export const helpersService = {
  getComputedPath: GetComputedPath,
};

export const eventService = {
  eventsOn: EventsOn,
  eventsOff: EventsOff,
};
