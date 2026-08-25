import { afterEach, describe, expect, it, vi } from "vitest";

import { wailsBackend } from "@/contracts/adapters/wails.ts";
import {
	dataService,
	eventService,
	resetBackendServices,
} from "@/contracts/service.ts";
import { Event, type EventData } from "@/contracts/types.ts";
import { installInMemoryBackend } from "@/testing/backend.ts";
import { CommandBuilder } from "@/testing/builders/command.ts";

describe("the backend seam", () => {
	afterEach(() => {
		resetBackendServices();
	});

	it("Should route calls to the substituted adapter", async () => {
		// Arrange
		const backend = installInMemoryBackend();
		const command = new CommandBuilder().withId("cmd-1").build();

		// Act
		await dataService.addCommand(command);

		// Assert
		expect(backend.state.commands).toEqual([command]);
		expect(await dataService.getCommands()).toEqual([command]);
	});

	it("Should route calls back to the Wails adapter once reset", () => {
		// Arrange
		installInMemoryBackend();

		// Act
		resetBackendServices();

		// Assert
		expect(dataService).toBe(wailsBackend.data);
	});

	it("Should deliver an emitted log entry to a listener registered through the seam", () => {
		// Arrange
		const backend = installInMemoryBackend();
		const received: EventData[Event.NEW_LOG_ENTRY][] = [];
		eventService.eventsOn(Event.NEW_LOG_ENTRY, (data) => received.push(data));

		// Act
		backend.emit(Event.NEW_LOG_ENTRY, { id: "cmd-1", line: "hello" });

		// Assert
		expect(received).toEqual([{ id: "cmd-1", line: "hello" }]);
	});

	it.each([
		Event.PROCESS_STARTED,
		Event.PROCESS_FINISHED,
		Event.COMMAND_ERROR_DETECTED,
		Event.COMMAND_GROUP_DELETED,
	])("Should deliver an emitted %s on demand", (event) => {
		// Arrange
		const backend = installInMemoryBackend();
		const listener = vi.fn();
		eventService.eventsOn(event, listener);

		// Act
		backend.emit(event, "cmd-1");

		// Assert
		expect(listener).toHaveBeenCalledExactlyOnceWith("cmd-1");
	});

	it("Should stop delivering to a listener that unsubscribed", () => {
		// Arrange
		const backend = installInMemoryBackend();
		const listener = vi.fn();
		const unsubscribe = eventService.eventsOn(Event.PROCESS_STARTED, listener);

		// Act
		unsubscribe();
		backend.emit(Event.PROCESS_STARTED, "cmd-1");

		// Assert
		expect(listener).not.toHaveBeenCalled();
	});

	it("Should stop delivering the events turned off", () => {
		// Arrange
		const backend = installInMemoryBackend();
		const startedListener = vi.fn();
		const finishedListener = vi.fn();
		eventService.eventsOn(Event.PROCESS_STARTED, startedListener);
		eventService.eventsOn(Event.PROCESS_FINISHED, finishedListener);

		// Act
		eventService.eventsOff(Event.PROCESS_STARTED, Event.PROCESS_FINISHED);
		backend.emit(Event.PROCESS_STARTED, "cmd-1");
		backend.emit(Event.PROCESS_FINISHED, "cmd-1");

		// Assert
		expect(startedListener).not.toHaveBeenCalled();
		expect(finishedListener).not.toHaveBeenCalled();
	});
});
