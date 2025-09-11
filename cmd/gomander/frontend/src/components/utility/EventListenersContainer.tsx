import { Fragment, useEffect, useRef } from "react";

import { eventService } from "@/contracts/service.ts";
import { Event, type EventData } from "@/contracts/types.ts";
import { useCommandStore } from "@/store/commandStore.ts";
import { CommandStatus } from "@/types/CommandStatus.ts";
import { updateCommandStatus } from "@/useCases/command/updateCommandStatus.ts";

export const EventListenersContainer = () => {
  const addLogs = useCommandStore((state) => state.addLogs);

  const logsBuffer = useRef(new Map<string, string[]>());

  useEffect(() => {
    const interval = setInterval(() => {
      if (logsBuffer.current.size > 0) {
        addLogs(logsBuffer.current);
        logsBuffer.current.clear();
      }
    }, 30); // Flush every 30ms

    return () => clearInterval(interval);
  }, [addLogs]);

  // Register events listeners
  useEffect(() => {
    eventService.eventsOn(
      Event.NEW_LOG_ENTRY,
      (data: EventData[Event.NEW_LOG_ENTRY]) => {
        const { id, line } = data;
        if (!logsBuffer.current.has(id)) {
          logsBuffer.current.set(id, []);
        }
        logsBuffer.current.get(id)!.push(line);
      },
    );

    eventService.eventsOn(
      Event.PROCESS_FINISHED,
      (data: EventData[Event.PROCESS_FINISHED]) =>
        updateCommandStatus(data, CommandStatus.IDLE),
    );

    eventService.eventsOn(
      Event.PROCESS_STARTED,
      (data: EventData[Event.PROCESS_STARTED]) =>
        updateCommandStatus(data, CommandStatus.RUNNING),
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
