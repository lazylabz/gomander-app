import {
  AddCommand,
  CloseProject,
  CreateCommandGroup,
  CreateProject,
  DeleteCommandGroup,
  DeleteProject,
  DuplicateCommand,
  EditCommand,
  EditProject,
  ExportProject,
  GetAvailableProjects,
  GetCommandGroups,
  GetCommands,
  GetCurrentProject,
  GetCurrentRelease,
  GetProjectToImport,
  GetUserConfig,
  ImportProject,
  IsThereANewRelease,
  OpenProject,
  RemoveCommand,
  RemoveCommandFromCommandGroup,
  ReorderCommandGroups,
  ReorderCommands,
  RunCommand,
  SaveUserConfig,
  StopCommand,
  UpdateCommandGroup,
} from "../../wailsjs/go/app/App";
import type { domain } from "../../wailsjs/go/models.ts";
import { GetComputedPath } from "../../wailsjs/go/path/UiPathHelper";
import { BrowserOpenURL, EventsOff, EventsOn } from "../../wailsjs/runtime";

export const dataService = {
  addCommand: AddCommand,
  duplicateCommand: DuplicateCommand,
  editCommand: EditCommand,
  reorderCommands: ReorderCommands,
  getAvailableProjects: GetAvailableProjects,
  editCommandGroup: UpdateCommandGroup,
  createCommandGroup: CreateCommandGroup,
  deleteCommandGroup: DeleteCommandGroup,
  reorderCommandGroups: ReorderCommandGroups,
  removeCommandFromGroup: RemoveCommandFromCommandGroup,
  getCommandGroups: GetCommandGroups,
  getCommands: GetCommands,
  getCurrentProject: GetCurrentProject as () => Promise<domain.Project | null>,
  getUserConfig: GetUserConfig,
  removeCommand: RemoveCommand,
  runCommand: RunCommand,
  saveUserConfig: SaveUserConfig,
  createProject: CreateProject,
  stopCommand: StopCommand,
  openProject: OpenProject,
  closeProject: CloseProject,
  deleteProject: DeleteProject,
  exportProject: ExportProject,
  importProject: ImportProject,
  getProjectToImport: GetProjectToImport,
  editProject: EditProject,
};

export const helpersService = {
  getComputedPath: GetComputedPath,
  isThereANewRelease: IsThereANewRelease,
  getCurrentRelease: GetCurrentRelease,
};

export const eventService = {
  eventsOn: EventsOn,
  eventsOff: EventsOff,
};

export const externalBrowserService = {
  browserOpenURL: BrowserOpenURL,
};
