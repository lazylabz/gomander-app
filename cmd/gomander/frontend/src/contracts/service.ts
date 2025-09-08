import {
  ExportProject,
  GetCurrentRelease,
  GetProjectToImport,
  ImportProject,
  IsThereANewRelease,
  RunCommand,
  StopCommand,
} from "../../wailsjs/go/app/App";
import {
  AddCommandController,
  CloseProjectController,
  CreateCommandGroupController,
  CreateProjectController,
  DeleteCommandGroupController,
  DeleteProjectController,
  DuplicateCommandController,
  EditCommandController,
  EditProjectController,
  GetAvailableProjectsController,
  GetCommandGroupsController,
  GetCommandsController,
  GetCurrentProjectController,
  GetUserConfigController,
  OpenProjectController,
  RemoveCommandController,
  RemoveCommandFromCommandGroupController,
  ReorderCommandGroupsController,
  ReorderCommandsController,
  SaveUserConfigController,
  UpdateCommandGroupController,
} from "../../wailsjs/go/main/WailsControllers";
import type { domain } from "../../wailsjs/go/models.ts";
import { GetComputedPath } from "../../wailsjs/go/path/UiPathHelper";
import { BrowserOpenURL, EventsOff, EventsOn } from "../../wailsjs/runtime";

export const dataService = {
  addCommand: AddCommandController,
  duplicateCommand: DuplicateCommandController,
  editCommand: EditCommandController,
  reorderCommands: ReorderCommandsController,
  getAvailableProjects: GetAvailableProjectsController,
  editCommandGroup: UpdateCommandGroupController,
  createCommandGroup: CreateCommandGroupController,
  deleteCommandGroup: DeleteCommandGroupController,
  reorderCommandGroups: ReorderCommandGroupsController,
  removeCommandFromGroup: RemoveCommandFromCommandGroupController,
  getCommandGroups: GetCommandGroupsController,
  getCommands: GetCommandsController,
  getCurrentProject:
    GetCurrentProjectController as () => Promise<domain.Project | null>,
  getUserConfig: GetUserConfigController,
  removeCommand: RemoveCommandController,
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
