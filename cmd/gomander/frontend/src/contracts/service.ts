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
  GetCommandGroups,
  GetCommands,
  GetCurrentRelease,
  GetProjectToImport,
  ImportProject,
  IsThereANewRelease,
  RemoveCommand,
  RemoveCommandFromCommandGroup,
  ReorderCommandGroups,
  ReorderCommands,
  RunCommand,
  StopCommand,
  UpdateCommandGroup,
} from "../../wailsjs/go/app/App";
import {
  GetAvailableProjectsController,
  GetCurrentProjectController,
  GetUserConfigController,
  OpenProjectController,
  SaveUserConfigController,
} from "../../wailsjs/go/main/WailsControllers";
import type { domain } from "../../wailsjs/go/models.ts";
import { GetComputedPath } from "../../wailsjs/go/path/UiPathHelper";
import { BrowserOpenURL, EventsOff, EventsOn } from "../../wailsjs/runtime";

export const dataService = {
  addCommand: AddCommand,
  duplicateCommand: DuplicateCommand,
  editCommand: EditCommand,
  reorderCommands: ReorderCommands,
  getAvailableProjects: GetAvailableProjectsController,
  editCommandGroup: UpdateCommandGroup,
  createCommandGroup: CreateCommandGroup,
  deleteCommandGroup: DeleteCommandGroup,
  reorderCommandGroups: ReorderCommandGroups,
  removeCommandFromGroup: RemoveCommandFromCommandGroup,
  getCommandGroups: GetCommandGroups,
  getCommands: GetCommands,
  getCurrentProject:
    GetCurrentProjectController as () => Promise<domain.Project | null>,
  getUserConfig: GetUserConfigController,
  removeCommand: RemoveCommand,
  runCommand: RunCommand,
  saveUserConfig: SaveUserConfigController,
  createProject: CreateProject,
  stopCommand: StopCommand,
  openProject: OpenProjectController,
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
