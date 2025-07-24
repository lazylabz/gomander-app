import {
  AddCommand,
  EditCommand,
  GetCommandGroups,
  GetCommands,
  GetUserConfig,
  RemoveCommand,
  RunCommand,
  SaveCommandGroups,
  SaveUserConfig,
  StopCommand,
} from "../../wailsjs/go/app/App";
import { EventsOff, EventsOn } from "../../wailsjs/runtime";

export const dataService = {
  addCommand: AddCommand,
  editCommand: EditCommand,
  getCommandGroups: GetCommandGroups,
  getCommands: GetCommands,
  getUserConfig: GetUserConfig,
  removeCommand: RemoveCommand,
  runCommand: RunCommand,
  saveCommandGroups: SaveCommandGroups,
  saveUserConfig: SaveUserConfig,
  stopCommand: StopCommand,
};

export const eventService = {
  eventsOn: EventsOn,
  eventsOff: EventsOff,
};
