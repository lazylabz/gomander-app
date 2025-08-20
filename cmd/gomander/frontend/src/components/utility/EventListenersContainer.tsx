import { Fragment, useEffect } from "react";

import { eventService } from "@/contracts/service.ts";
import { Event, type EventData } from "@/contracts/types.ts";
import { fetchUserConfig } from "@/queries/fetchUserConfig.ts";
import { useCommandStore } from "@/store/commandStore.ts";
import { CommandStatus } from "@/types/CommandStatus.ts";
import { updateCommandStatus } from "@/useCases/command/updateCommandStatus.ts";

export const EventListenersContainer = () => {
  const addLog = useCommandStore((state) => state.addLog);

  // Register events listeners
  useEffect(() => {
    eventService.eventsOn(Event.GET_USER_CONFIG, () => {
      fetchUserConfig();
    });

    eventService.eventsOn(
      Event.NEW_LOG_ENTRY,
      (data: EventData[Event.NEW_LOG_ENTRY]) => {
        const { id, line } = data;
        addLog(id, line);
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
