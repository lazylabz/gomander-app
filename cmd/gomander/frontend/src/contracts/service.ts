import {
  AddCommand,
  DuplicateCommand,
  EditCommand,
  ExportProject,
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
} from "../../wailsjs/go/app/App";
import {
  CloseProjectController,
  CreateCommandGroupController,
  CreateProjectController,
  DeleteCommandGroupController,
  DeleteProjectController,
  EditProjectController,
  GetAvailableProjectsController,
  GetCommandGroupsController,
  GetCurrentProjectController,
  GetUserConfigController,
  OpenProjectController,
  SaveUserConfigController,
  UpdateCommandGroupController,
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
  editCommandGroup: UpdateCommandGroupController,
  createCommandGroup: CreateCommandGroupController,
  deleteCommandGroup: DeleteCommandGroupController,
  reorderCommandGroups: ReorderCommandGroups,
  removeCommandFromGroup: RemoveCommandFromCommandGroup,
  getCommandGroups: GetCommandGroupsController,
  getCommands: GetCommands,
  getCurrentProject:
    GetCurrentProjectController as () => Promise<domain.Project | null>,
  getUserConfig: GetUserConfigController,
  removeCommand: RemoveCommand,
  runCommand: RunCommand,
  saveUserConfig: SaveUserConfigController,
  createProject: CreateProjectController,
  stopCommand: StopCommand,
  openProject: OpenProjectController,
  closeProject: CloseProjectController,
  deleteProject: DeleteProjectController,
  exportProject: ExportProject,
  importProject: ImportProject,
  getProjectToImport: GetProjectToImport,
  editProject: EditProjectController,
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
