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
  ExportProjectController,
  GetAvailableProjectsController,
  GetCommandGroupsController,
  GetCommandsController,
  GetCurrentProjectController,
  GetProjectToImportController,
  GetUserConfigController,
  ImportProjectController,
  OpenProjectController,
  RemoveCommandController,
  RemoveCommandFromCommandGroupController,
  ReorderCommandGroupsController,
  ReorderCommandsController,
  RunCommandController,
  RunCommandGroupController,
  SaveUserConfigController,
  StopCommandController,
  StopCommandGroupController,
  UpdateCommandGroupController,
} from "../../wailsjs/go/main/WailsControllers";
import type { domain } from "../../wailsjs/go/models.ts";
import { GetComputedPath } from "../../wailsjs/go/path/UiPathHelper";
import {
  GetCurrentRelease,
  IsThereANewRelease,
} from "../../wailsjs/go/releases/ReleaseHelper";
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
  runCommandGroup: RunCommandGroupController,
  stopCommandGroup: StopCommandGroupController,
  getCommandGroups: GetCommandGroupsController,
  getCommands: GetCommandsController,
  getCurrentProject:
    GetCurrentProjectController as () => Promise<domain.Project | null>,
  getUserConfig: GetUserConfigController,
  removeCommand: RemoveCommandController,
  runCommand: RunCommandController,
  saveUserConfig: SaveUserConfigController,
  createProject: CreateProjectController,
  stopCommand: StopCommandController,
  openProject: OpenProjectController,
  closeProject: CloseProjectController,
  deleteProject: DeleteProjectController,
  exportProject: ExportProjectController,
  importProject: ImportProjectController,
  getProjectToImport: GetProjectToImportController,
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
