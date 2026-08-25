import { Fragment, useEffect, useRef } from "react";

import { getCommandGroupSectionOpenLocalStorageKey } from "@/constants/localStorage.ts";
import { useTheme } from "@/contexts/theme.tsx";
import { eventService } from "@/contracts/service.ts";
import { Event, type EventData } from "@/contracts/types.ts";
import { removeKeyFromLocalStorage } from "@/helpers/localStorage.ts";
import { XTERM_THEMES } from "@/screens/LogsScreen/components/CommandTerminal.tsx";
import {
	formatLogTimestamp,
	prependTimestamp,
} from "@/screens/LogsScreen/helpers.ts";
import { clearLogTail, recordLogTail } from "@/store/commandLogsTail.ts";
import { terminalStore } from "@/store/terminalStore.ts";
import { CommandStatus } from "@/types/CommandStatus.ts";
import { detectMissingEnvironmentPathFailure } from "@/useCases/command/detectMissingEnvironmentPathFailure.tsx";
import { recordCommandsErrors } from "@/useCases/command/recordCommandsErrors.ts";
import { updateCommandStatus } from "@/useCases/command/updateCommandStatus.ts";

export const EventListenersContainer = () => {
	const { theme } = useTheme();

	useEffect(() => {
		terminalStore.getState().setThemeAll(XTERM_THEMES[theme]);
	}, [theme]);

	const logsBuffer = useRef(new Map<string, string[]>());
	const errorBuffer = useRef<string[]>([]);

	useEffect(() => {
		const interval = setInterval(() => {
			// Process error buffer
			recordCommandsErrors(errorBuffer.current);
			errorBuffer.current = [];

			// Process logs buffer
			if (logsBuffer.current.size > 0) {
				const bufferCopy = new Map(logsBuffer.current);
				logsBuffer.current.clear();

				// Write directly to already-open terminals (bypasses React re-render cycle).
				const ts = formatLogTimestamp(new Date());
				const { terminals, bufferLogs } = terminalStore.getState();
				for (const [commandId, lines] of bufferCopy) {
					recordLogTail(commandId, lines);

					const stamped = lines.map((line) => prependTimestamp(line, ts));
					const term = terminals.get(commandId);
					if (term) {
						for (const line of stamped) {
							term.writeln(line);
						}
					} else {
						bufferLogs(commandId, stamped);
					}
				}
			}
		}, 30); // Flush every 30ms

		return () => clearInterval(interval);
	}, []);

	// Register events listeners
	useEffect(() => {
		eventService.eventsOn(
			Event.NEW_LOG_ENTRY,
			(data: EventData[Event.NEW_LOG_ENTRY]) => {
				const { id, line } = data;
				if (!logsBuffer.current.has(id)) {
					logsBuffer.current.set(id, []);
				}
				logsBuffer.current.get(id)?.push(line);
			},
		);

		eventService.eventsOn(
			Event.PROCESS_FINISHED,
			(data: EventData[Event.PROCESS_FINISHED]) => {
				updateCommandStatus(data, CommandStatus.IDLE);
				detectMissingEnvironmentPathFailure(
					data,
					logsBuffer.current.get(data) ?? [],
				);
			},
		);

		eventService.eventsOn(
			Event.PROCESS_STARTED,
			(data: EventData[Event.PROCESS_STARTED]) => {
				updateCommandStatus(data, CommandStatus.RUNNING);
				// Not-yet-flushed lines belong to the previous run; without this
				// they would land in the reset terminal and the cleared tail.
				logsBuffer.current.delete(data);
				clearLogTail(data);
				terminalStore.getState().terminals.get(data)?.reset();
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
