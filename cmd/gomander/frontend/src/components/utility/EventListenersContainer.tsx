import { Fragment, useEffect } from "react";
import { toast } from "sonner";

import { eventService } from "@/contracts/service.ts";
import { Event, type EventData } from "@/contracts/types.ts";
import { fetchCommandGroups } from "@/queries/fetchCommandGroups.ts";
import { fetchCommands } from "@/queries/fetchCommands.ts";
import { fetchUserConfig } from "@/queries/fetchUserConfig.ts";
import { useCommandStore } from "@/store/commandStore.ts";
import { CommandStatus } from "@/types/CommandStatus.ts";
import { updateCommandStatus } from "@/useCases/command/updateCommandStatus.ts";

export const EventListenersContainer = () => {
  const commandsLogs = useCommandStore((state) => state.commandsLogs);
  const setCommandsLogs = useCommandStore((state) => state.setCommandsLogs);

  // Register events listeners
  useEffect(() => {
    eventService.eventsOn(Event.GET_COMMANDS, () => {
      fetchCommands();
    });

    eventService.eventsOn(Event.GET_COMMAND_GROUPS, () => {
      fetchCommandGroups();
    });

    eventService.eventsOn(
      Event.NEW_LOG_ENTRY,
      (data: EventData[Event.NEW_LOG_ENTRY]) => {
        const { id, line } = data;

        const newLogs = { ...commandsLogs };

        if (!newLogs[id]) {
          newLogs[id] = [];
        }
        newLogs[id].push(line);

        if (newLogs[id].length > 100) {
          newLogs[id] = newLogs[id].slice(-100); // Keep only the last 100 lines
        }

        setCommandsLogs(newLogs);
      },
    );

    eventService.eventsOn(
      Event.ERROR_NOTIFICATION,
      (data: EventData[Event.ERROR_NOTIFICATION]) => {
        toast.error("Error", {
          description: data,
        });
      },
    );

    eventService.eventsOn(
      Event.SUCCESS_NOTIFICATION,
      (data: EventData[Event.SUCCESS_NOTIFICATION]) => {
        toast.success(data);
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

    eventService.eventsOn(Event.GET_USER_CONFIG, () => {
      fetchUserConfig();
    });

    // Clean listeners on all events
    return () =>
      eventService.eventsOff(
        Object.values(Event)[0],
        ...Object.values(Event).slice(1),
      );
  });

  return <Fragment />;
};
