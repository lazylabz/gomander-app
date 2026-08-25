import { wailsBackend } from "@/contracts/adapters/wails.ts";
import type {
	BackendServices,
	DataService,
	EventService,
	ExternalBrowserService,
	HelpersService,
	TranslationsService,
} from "@/contracts/ports.ts";

// Module-level swap rather than injecting the services into every use case: the
// exports below are live ESM bindings, so `setBackendServices` reaches every call
// site without changing any of their signatures. Only tests should call it.
export let dataService: DataService = wailsBackend.data;
export let helpersService: HelpersService = wailsBackend.helpers;
export let eventService: EventService = wailsBackend.event;
export let externalBrowserService: ExternalBrowserService =
	wailsBackend.externalBrowser;
export let translationsService: TranslationsService = wailsBackend.translations;

export const setBackendServices = (backend: BackendServices): void => {
	dataService = backend.data;
	helpersService = backend.helpers;
	eventService = backend.event;
	externalBrowserService = backend.externalBrowser;
	translationsService = backend.translations;
};

export const resetBackendServices = (): void =>
	setBackendServices(wailsBackend);
