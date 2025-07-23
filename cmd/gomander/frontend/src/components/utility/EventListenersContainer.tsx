import { Fragment, useEffect } from "react";
import { toast } from "sonner";

import { CommandStatus, useDataContext } from "@/contexts/DataContext.tsx";
import { useUserConfigDataContext } from "@/contexts/UserConfigDataContext.tsx";
import { Event, type EventData } from "@/types/contracts.ts";

import { EventsOff, EventsOn } from "../../../wailsjs/runtime";

export const EventListenersContainer = () => {
  const { fetchCommands, fetchCommandGroups, setCommandsStatus, setLogs } =
    useDataContext();
  const { fetchUserConfig } = useUserConfigDataContext();

  // Register events listeners
  useEffect(() => {
    EventsOn(Event.GET_COMMANDS, () => {
      fetchCommands();
    });

    EventsOn(Event.GET_COMMAND_GROUPS, () => {
      fetchCommandGroups();
    });

    EventsOn(Event.NEW_LOG_ENTRY, (data: EventData[Event.NEW_LOG_ENTRY]) => {
      const { id, line } = data;

      setLogs((prevLogs) => {
        const newLogs = { ...prevLogs };
        if (!newLogs[id]) {
          newLogs[id] = [];
        }
        newLogs[id].push(line);
        return newLogs;
      });
    });

    EventsOn(
      Event.ERROR_NOTIFICATION,
      (data: EventData[Event.ERROR_NOTIFICATION]) => {
        toast.error("Error", {
          description: data,
        });
      },
    );

    EventsOn(
      Event.SUCCESS_NOTIFICATION,
      (data: EventData[Event.SUCCESS_NOTIFICATION]) => {
        toast.success(data);
      },
    );

    EventsOn(
      Event.PROCESS_FINISHED,
      (data: EventData[Event.PROCESS_FINISHED]) => {
        setCommandsStatus((prevStatus) => ({
          ...prevStatus,
          [data]: CommandStatus.IDLE, // Reset status to IDLE when process finishes
        }));
      },
    );

    EventsOn(Event.GET_USER_CONFIG, () => {
      fetchUserConfig();
    });

    // Clean listeners on all events
    return () =>
      EventsOff(Object.values(Event)[0], ...Object.values(Event).slice(1));
  });

  return <Fragment />;
};
