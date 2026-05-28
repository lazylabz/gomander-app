import { Fragment, useEffect, useRef } from "react";

import { getCommandGroupSectionOpenLocalStorageKey } from "@/constants/localStorage.ts";
import { useTheme } from "@/contexts/theme.tsx";
import { eventService } from "@/contracts/service.ts";
import { Event, type EventData } from "@/contracts/types.ts";
import { removeKeyFromLocalStorage } from "@/helpers/localStorage.ts";
import {
	formatLogTimestamp,
	prependTimestamp,
} from "@/screens/ExperimentalLogsScreen/helpers.ts";
import { useCommandStore } from "@/store/commandStore.ts";
import { terminalStore } from "@/store/terminalStore.ts";
import { useUserConfigurationStore } from "@/store/userConfigurationStore.ts";
import { CommandStatus } from "@/types/CommandStatus.ts";
import { cleanCommandLogs } from "@/useCases/command/cleanCommandLogs.ts";
import { recordCommandsErrors } from "@/useCases/command/recordCommandsErrors.ts";
import { updateCommandStatus } from "@/useCases/command/updateCommandStatus.ts";
import { XTERM_THEMES } from "../../screens/ExperimentalLogsScreen/components/CommandTerminal.tsx";

export const EventListenersContainer = () => {
	const addLogs = useCommandStore((state) => state.addLogs);
	const userConfig = useUserConfigurationStore((state) => state.userConfig);
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
				addLogs(bufferCopy, userConfig.logLineLimit);
				logsBuffer.current.clear();

				// Write directly to already-open terminals (bypasses React re-render cycle).
				// timestamps stay xterm-only.
				const ts = formatLogTimestamp(new Date());
				const { terminals, bufferLogs } = terminalStore.getState();
				for (const [commandId, lines] of bufferCopy) {
					const stamped = lines.map((line) => prependTimestamp(line, ts));
					const term = terminals.get(commandId);
					if (term) {
						for (const line of stamped) term.writeln(line);
					} else {
						bufferLogs(commandId, stamped);
					}
				}
			}
		}, 30); // Flush every 30ms

		return () => clearInterval(interval);
	}, [addLogs, userConfig.logLineLimit]);

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
			(data: EventData[Event.PROCESS_FINISHED]) =>
				updateCommandStatus(data, CommandStatus.IDLE),
		);

		eventService.eventsOn(
			Event.PROCESS_STARTED,
			(data: EventData[Event.PROCESS_STARTED]) => {
				updateCommandStatus(data, CommandStatus.RUNNING);
				cleanCommandLogs(data);
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
