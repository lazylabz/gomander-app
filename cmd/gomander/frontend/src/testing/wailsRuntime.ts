// Vitest resolves @wailsio/runtime here: the real one boots on import, warns,
// connects to the page origin and leaves a timer that outlives the test
// environment. Tests reach the backend through installInMemoryBackend() instead.
const unavailable = (): never => {
	throw new Error(
		"The Wails runtime is not available in tests: install the in-memory backend",
	);
};

export const Call = { ByID: unavailable, ByName: unavailable };
export const Events = { On: unavailable, Off: unavailable };
export const Browser = { OpenURL: unavailable };
export const Create = { Events: {} };
export class CancellablePromise<T> extends Promise<T> {}
