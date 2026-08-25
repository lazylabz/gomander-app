import { Fragment, useEffect, useRef } from "react";

import {
	appendCommandOutput,
	resetCommandOutput,
	setCommandOutputTheme,
} from "@/commandOutput/commandOutput.ts";
import { getCommandGroupSectionOpenLocalStorageKey } from "@/constants/localStorage.ts";
import { useTheme } from "@/contexts/theme.tsx";
import { eventService } from "@/contracts/service.ts";
import { Event, type EventData } from "@/contracts/types.ts";
import { removeKeyFromLocalStorage } from "@/helpers/localStorage.ts";
import { CommandStatus } from "@/types/CommandStatus.ts";
import { detectMissingEnvironmentPathFailure } from "@/useCases/command/detectMissingEnvironmentPathFailure.tsx";
import { recordCommandsErrors } from "@/useCases/command/recordCommandsErrors.ts";
import { updateCommandStatus } from "@/useCases/command/updateCommandStatus.ts";

// Errors are batched the same way log lines are: several can land in one tick.
const ERROR_FLUSH_INTERVAL_MS = 30;

export const EventListenersContainer = () => {
	const { theme } = useTheme();

	useEffect(() => {
		setCommandOutputTheme(theme);
	}, [theme]);

	const errorBuffer = useRef<string[]>([]);

	useEffect(() => {
		const interval = setInterval(() => {
			// Recording writes a fresh array into the store, which re-renders every
			// command row; on an empty tick that is 33 renders a second for nothing.
			if (errorBuffer.current.length === 0) {
				return;
			}
			recordCommandsErrors(errorBuffer.current);
			errorBuffer.current = [];
		}, ERROR_FLUSH_INTERVAL_MS);

		return () => clearInterval(interval);
	}, []);

	// Register events listeners
	useEffect(() => {
		eventService.eventsOn(
			Event.NEW_LOG_ENTRY,
			(data: EventData[Event.NEW_LOG_ENTRY]) =>
				appendCommandOutput(data.id, [data.line]),
		);

		eventService.eventsOn(
			Event.PROCESS_FINISHED,
			(data: EventData[Event.PROCESS_FINISHED]) => {
				updateCommandStatus(data, CommandStatus.IDLE);
				detectMissingEnvironmentPathFailure(data);
			},
		);

		eventService.eventsOn(
			Event.PROCESS_STARTED,
			(data: EventData[Event.PROCESS_STARTED]) => {
				updateCommandStatus(data, CommandStatus.RUNNING);
				resetCommandOutput(data);
			},
		);

		eventService.eventsOn(
			Event.COMMAND_GROUP_DELETED,
			(data: EventData[Event.COMMAND_GROUP_DELETED]) =>
				removeKeyFromLocalStorage(
					getCommandGroupSectionOpenLocalStorageKey(data),
				),
		);

		eventService.eventsOn(
			Event.COMMAND_ERROR_DETECTED,
			(data: EventData[Event.COMMAND_ERROR_DETECTED]) =>
				errorBuffer.current.push(data),
		);

		// Clean listeners on all events
		return () =>
			eventService.eventsOff(
				Object.values(Event)[0],
				...Object.values(Event).slice(1),
			);
	});

	return <Fragment />;
};
