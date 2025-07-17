import {main} from "../../wailsjs/go/models";

// Types
export type Command = main.Command;


// Enums
export type Event = main.Event;
export const Event = main.Event;

export type EventData = {
    [Event.GET_COMMANDS]: null;
}