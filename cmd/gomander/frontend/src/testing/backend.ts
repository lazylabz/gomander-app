import {
	createInMemoryBackend,
	type InMemoryBackend,
	type InMemoryBackendState,
} from "@/contracts/adapters/inMemory.ts";
import { setBackendServices } from "@/contracts/service.ts";

export const installInMemoryBackend = (
	initialState?: Partial<InMemoryBackendState>,
): InMemoryBackend => {
	const backend = createInMemoryBackend(initialState);
	setBackendServices(backend);
	return backend;
};
