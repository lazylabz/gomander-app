import { act, useEffect } from "react";
import { createRoot } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { useLoadFailure } from "@/hooks/useLoadFailure.ts";

describe("useLoadFailure", () => {
	let sut: ReturnType<typeof useLoadFailure>;
	let container: HTMLDivElement;

	const Probe = () => {
		const state = useLoadFailure();
		useEffect(() => {
			sut = state;
		});
		return null;
	};

	beforeEach(() => {
		vi.spyOn(console, "error").mockImplementation(() => {});
		Object.assign(globalThis, { IS_REACT_ACT_ENVIRONMENT: true });
		container = document.createElement("div");
		act(() => {
			createRoot(container).render(<Probe />);
		});
	});

	afterEach(() => {
		vi.restoreAllMocks();
	});

	it("Should not report a failure before anything was attempted", () => {
		// Assert
		expect(sut.hasFailed).toBe(false);
	});

	it("Should not report a failure when the load succeeds", async () => {
		// Act
		await act(async () => {
			await sut.attempt(async () => {});
		});

		// Assert
		expect(sut.hasFailed).toBe(false);
	});

	it("Should report a failure when the load rejects, without rejecting itself", async () => {
		// Act
		await act(async () => {
			await sut.attempt(async () => {
				throw new Error("boom");
			});
		});

		// Assert
		expect(sut.hasFailed).toBe(true);
	});

	it("Should run the failed load again on retry and clear the failure once it succeeds", async () => {
		// Arrange
		const load = vi
			.fn<() => Promise<void>>()
			.mockRejectedValueOnce(new Error("boom"))
			.mockResolvedValueOnce(undefined);
		await act(async () => {
			await sut.attempt(load);
		});

		// Act
		await act(async () => {
			await sut.retry();
		});

		// Assert
		expect(load).toHaveBeenCalledTimes(2);
		expect(sut.hasFailed).toBe(false);
	});

	it("Should keep the failure when the retry fails again", async () => {
		// Arrange
		const load = vi
			.fn<() => Promise<void>>()
			.mockRejectedValue(new Error("boom"));
		await act(async () => {
			await sut.attempt(load);
		});

		// Act
		await act(async () => {
			await sut.retry();
		});

		// Assert
		expect(load).toHaveBeenCalledTimes(2);
		expect(sut.hasFailed).toBe(true);
	});

	it("Should retry the latest failed load, not an earlier one", async () => {
		// Arrange
		const first = vi
			.fn<() => Promise<void>>()
			.mockRejectedValue(new Error("a"));
		const second = vi
			.fn<() => Promise<void>>()
			.mockRejectedValue(new Error("b"));
		await act(async () => {
			await sut.attempt(first);
		});
		await act(async () => {
			await sut.attempt(second);
		});

		// Act
		await act(async () => {
			await sut.retry();
		});

		// Assert
		expect(first).toHaveBeenCalledTimes(1);
		expect(second).toHaveBeenCalledTimes(2);
	});

	it("Should be pending only while a load is running", async () => {
		// Arrange
		let finish: () => void = () => {};
		const load = () =>
			new Promise<void>((resolve) => {
				finish = resolve;
			});

		// Act
		let attempting: Promise<void> = Promise.resolve();
		act(() => {
			attempting = sut.attempt(load);
		});

		// Assert
		expect(sut.isPending).toBe(true);
		await act(async () => {
			finish();
			await attempting;
		});
		expect(sut.isPending).toBe(false);
	});

	it("Should keep reporting the failure while a retry is pending", async () => {
		// Arrange
		await act(async () => {
			await sut.attempt(async () => {
				throw new Error("boom");
			});
		});
		let finish: () => void = () => {};
		const retryLoad = () =>
			new Promise<void>((resolve) => {
				finish = resolve;
			});

		// Act
		let retrying: Promise<void> = Promise.resolve();
		act(() => {
			retrying = sut.attempt(retryLoad);
		});

		// Assert
		expect(sut.hasFailed).toBe(true);
		await act(async () => {
			finish();
			await retrying;
		});
		expect(sut.hasFailed).toBe(false);
	});

	describe("when loads overlap", () => {
		const deferred = () => {
			let resolve: () => void = () => {};
			let reject: () => void = () => {};
			const promise = new Promise<void>((res, rej) => {
				resolve = res;
				reject = () => rej(new Error("boom"));
			});
			return { load: () => promise, resolve, reject };
		};

		it("Should keep the failure of the newer load when an older one succeeds afterwards", async () => {
			// Arrange
			const older = deferred();
			const newer = deferred();
			let olderAttempt: Promise<void> = Promise.resolve();
			let newerAttempt: Promise<void> = Promise.resolve();
			act(() => {
				olderAttempt = sut.attempt(older.load);
				newerAttempt = sut.attempt(newer.load);
			});

			// Act
			await act(async () => {
				newer.reject();
				await newerAttempt;
				older.resolve();
				await olderAttempt;
			});

			// Assert
			expect(sut.hasFailed).toBe(true);
			expect(sut.isPending).toBe(false);
		});

		it("Should ignore the failure of an older load once a newer one succeeded", async () => {
			// Arrange
			const older = deferred();
			const newer = deferred();
			let olderAttempt: Promise<void> = Promise.resolve();
			let newerAttempt: Promise<void> = Promise.resolve();
			act(() => {
				olderAttempt = sut.attempt(older.load);
				newerAttempt = sut.attempt(newer.load);
			});

			// Act
			await act(async () => {
				newer.resolve();
				await newerAttempt;
				older.reject();
				await olderAttempt;
			});

			// Assert
			expect(sut.hasFailed).toBe(false);
		});

		it("Should stay pending until the newest load finishes", async () => {
			// Arrange
			const older = deferred();
			const newer = deferred();
			let olderAttempt: Promise<void> = Promise.resolve();
			let newerAttempt: Promise<void> = Promise.resolve();
			act(() => {
				olderAttempt = sut.attempt(older.load);
				newerAttempt = sut.attempt(newer.load);
			});

			// Act
			await act(async () => {
				older.resolve();
				await olderAttempt;
			});

			// Assert
			expect(sut.isPending).toBe(true);
			await act(async () => {
				newer.resolve();
				await newerAttempt;
			});
			expect(sut.isPending).toBe(false);
		});
	});
});
