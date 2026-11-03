import * as events from "../../bindings/events/index.ts";
import type * as domain from "../../bindings/gomander/cmd/gomander/transport/domain/models.ts";
import type * as localization from "../../bindings/gomander/internal/localization/models.ts";

// v3 types Go slices as `T[] | null`; the app has always treated them as arrays,
// as v2 typed them, so the contract keeps that promise instead of every caller.
export type WithoutNullSlices<T> =
	NonNullable<T> extends (infer Item)[]
		? WithoutNullSlices<Item>[]
		: T extends object
			? { [K in keyof T]: WithoutNullSlices<T[K]> }
			: T;

// Types
export type Command = WithoutNullSlices<domain.Command>;
export type UserConfig = WithoutNullSlices<domain.Config>;
export type CommandGroup = WithoutNullSlices<domain.CommandGroup>;
export type Project = WithoutNullSlices<domain.Project>;
export type EnvironmentPath = WithoutNullSlices<domain.EnvironmentPath>;
export type ProjectBlueprint = WithoutNullSlices<domain.ProjectBlueprint>;
export type CommandBlueprint = WithoutNullSlices<domain.CommandBlueprint>;
export type CommandGroupBlueprint =
	WithoutNullSlices<domain.CommandGroupBlueprint>;
export type Localization = WithoutNullSlices<localization.Localization>;

// Enums
export enum Event {
	NEW_LOG_ENTRY = "new_log_entry",
	PROCESS_FINISHED = "process_finished",
	PROCESS_STARTED = "process_started",
	COMMAND_GROUP_DELETED = "command_group_deleted",
	COMMAND_ERROR_DETECTED = "command_error_detected",
}

// A string enum can't take its values from the generated bindings, so this
// fails to compile instead when an Event drifts from the name Go registers.
({
	[Event.NEW_LOG_ENTRY]: events.NewLogEntry,
	[Event.PROCESS_FINISHED]: events.ProcessFinished,
	[Event.PROCESS_STARTED]: events.ProcessStarted,
	[Event.COMMAND_GROUP_DELETED]: events.CommandGroupDeleted,
	[Event.COMMAND_ERROR_DETECTED]: events.CommandErrorDetected,
}) satisfies { [E in Event]: `${E}` };

// Which of a command's two kinds of line arrived: the command's own text, sent
// once as it starts, or a line the process printed. How each one looks is
// decided in commandOutput, not by the backend that sends it.
export enum LogEntryKind {
	COMMAND = "command",
	OUTPUT = "output",
}

export type EventData = {
	[Event.NEW_LOG_ENTRY]: {
		id: string;
		line: string;
		kind: LogEntryKind;
	};
	[Event.PROCESS_FINISHED]: string;
	[Event.PROCESS_STARTED]: string;
	[Event.COMMAND_GROUP_DELETED]: string;
	[Event.COMMAND_ERROR_DETECTED]: string;
};
