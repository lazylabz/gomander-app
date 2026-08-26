import { act, useEffect } from "react";
import { createRoot } from "react-dom/client";
import {
	afterEach,
	beforeEach,
	describe,
	expect,
	it,
	type Mock,
	vi,
} from "vitest";
import { z } from "zod";

import { createAutosaver, useAutosavedForm } from "@/hooks/useAutosavedForm.ts";

const DEBOUNCE_MS = 300;

describe("createAutosaver", () => {
	let values: { name: string };
	let submit: Mock<() => Promise<void>>;
	let onPendingChange: Mock<(isPending: boolean) => void>;

	const buildSut = () =>
		createAutosaver({
			getValues: () => values,
			submit,
			onPendingChange,
		});

	beforeEach(() => {
		vi.useFakeTimers();
		values = { name: "initial" };
		submit = vi.fn(async () => {});
		onPendingChange = vi.fn();
	});

	afterEach(() => {
		vi.useRealTimers();
	});

	it("Should not submit before the debounce window has elapsed", async () => {
		// Arrange
		const sut = buildSut();

		// Act
		values = { name: "changed" };
		sut.onChange();
		await vi.advanceTimersByTimeAsync(DEBOUNCE_MS - 1);

		// Assert
		expect(submit).not.toHaveBeenCalled();
	});

	it("Should submit once when several changes land inside the window", async () => {
		// Arrange
		const sut = buildSut();

		// Act
		values = { name: "c" };
		sut.onChange();
		await vi.advanceTimersByTimeAsync(100);
		values = { name: "ch" };
		sut.onChange();
		await vi.advanceTimersByTimeAsync(100);
		values = { name: "cha" };
		sut.onChange();
		await vi.advanceTimersByTimeAsync(DEBOUNCE_MS);

		// Assert
		expect(submit).toHaveBeenCalledOnce();
	});

	it("Should submit again for a change made after a save", async () => {
		// Arrange
		const sut = buildSut();
		values = { name: "first" };
		sut.onChange();
		await vi.advanceTimersByTimeAsync(DEBOUNCE_MS);

		// Act
		values = { name: "second" };
		sut.onChange();
		await vi.advanceTimersByTimeAsync(DEBOUNCE_MS);

		// Assert
		expect(submit).toHaveBeenCalledTimes(2);
	});

	it("Should stay quiet when the values are the ones already saved", async () => {
		// Arrange
		const sut = buildSut();

		// Act
		sut.onChange();
		await vi.advanceTimersByTimeAsync(DEBOUNCE_MS);

		// Assert
		expect(submit).not.toHaveBeenCalled();
		expect(onPendingChange).not.toHaveBeenCalled();
	});

	it("Should stay quiet when the change is undone inside the window", async () => {
		// Arrange
		const sut = buildSut();

		// Act
		values = { name: "changed" };
		sut.onChange();
		values = { name: "initial" };
		sut.onChange();
		await vi.advanceTimersByTimeAsync(DEBOUNCE_MS);

		// Assert
		expect(submit).not.toHaveBeenCalled();
		expect(onPendingChange).toHaveBeenLastCalledWith(false);
	});

	it("Should raise the pending flag on a change and clear it once saved", async () => {
		// Arrange
		const sut = buildSut();

		// Act
		values = { name: "changed" };
		sut.onChange();

		// Assert
		expect(onPendingChange).toHaveBeenLastCalledWith(true);

		// Act
		await vi.advanceTimersByTimeAsync(DEBOUNCE_MS);

		// Assert
		expect(onPendingChange).toHaveBeenLastCalledWith(false);
	});

	it("Should keep the pending flag raised when a change lands while saving", async () => {
		// Arrange
		let resolveSubmit: () => void = () => {};
		submit = vi.fn(
			() =>
				new Promise<void>((resolve) => {
					resolveSubmit = resolve;
				}),
		);
		const sut = buildSut();

		values = { name: "changed" };
		sut.onChange();
		await vi.advanceTimersByTimeAsync(DEBOUNCE_MS);

		// Act
		values = { name: "changed again" };
		sut.onChange();
		resolveSubmit();
		await vi.advanceTimersByTimeAsync(0);

		// Assert
		expect(onPendingChange).toHaveBeenLastCalledWith(true);
	});

	it("Should report a failed submit and clear the pending flag", async () => {
		// Arrange
		const consoleError = vi
			.spyOn(console, "error")
			.mockImplementation(() => {});
		submit = vi.fn(async () => {
			throw new Error("boom");
		});
		const sut = buildSut();

		// Act
		values = { name: "changed" };
		sut.onChange();
		await vi.advanceTimersByTimeAsync(DEBOUNCE_MS);

		// Assert
		expect(consoleError).toHaveBeenCalledOnce();
		expect(onPendingChange).toHaveBeenLastCalledWith(false);

		consoleError.mockRestore();
	});

	it("Should still hold values a failed submit dropped as unsaved", async () => {
		// Arrange a first submit that fails, so its values were never persisted
		const consoleError = vi
			.spyOn(console, "error")
			.mockImplementation(() => {});
		submit = vi.fn(async () => {
			throw new Error("boom");
		});
		const sut = buildSut();

		values = { name: "changed" };
		sut.onChange();
		await vi.advanceTimersByTimeAsync(DEBOUNCE_MS);
		submit.mockImplementation(async () => {});

		// Act - the form reports the same values again
		sut.onChange();
		await vi.advanceTimersByTimeAsync(DEBOUNCE_MS);

		// Assert
		expect(submit).toHaveBeenCalledTimes(2);

		consoleError.mockRestore();
	});

	it("Should not submit a cancelled window", async () => {
		// Arrange
		const sut = buildSut();

		// Act
		values = { name: "changed" };
		sut.onChange();
		sut.cancel();
		await vi.advanceTimersByTimeAsync(DEBOUNCE_MS);

		// Assert
		expect(submit).not.toHaveBeenCalled();
	});
});

describe("useAutosavedForm", () => {
	const schema = z.object({ name: z.string().min(1) });

	type Save = (values: { name: string }) => Promise<void>;

	beforeEach(() => {
		// react-dom needs this to accept the act() calls that drive the hook.
		Object.assign(globalThis, { IS_REACT_ACT_ENVIRONMENT: true });
		vi.useFakeTimers();
	});

	afterEach(() => {
		vi.useRealTimers();
	});

	const renderAutosavedForm = async (save: Save) => {
		let isPending = false;
		let setName: (value: string) => void = () => {};

		const Probe = () => {
			const { form, isPending: pending } = useAutosavedForm({
				schema,
				defaultValues: { name: "initial" },
				save,
			});

			isPending = pending;

			useEffect(() => {
				setName = (value) => form.setValue("name", value);
			}, [form]);

			return null;
		};

		const root = createRoot(document.createElement("div"));
		await act(async () => {
			root.render(<Probe />);
		});

		return {
			isPending: () => isPending,
			setName: async (value: string) => {
				await act(async () => {
					setName(value);
				});
			},
			wait: async (ms: number) => {
				await act(async () => {
					await vi.advanceTimersByTimeAsync(ms);
				});
			},
			unmount: () => act(() => root.unmount()),
		};
	};

	it("Should save the edited values once the debounce window closes", async () => {
		// Arrange
		const save = vi.fn<Save>(async () => {});
		const sut = await renderAutosavedForm(save);

		// Act
		await sut.setName("changed");

		// Assert
		expect(save).not.toHaveBeenCalled();
		expect(sut.isPending()).toBe(true);

		// Act
		await sut.wait(DEBOUNCE_MS);

		// Assert
		expect(save).toHaveBeenCalledOnce();
		expect(save.mock.calls[0][0]).toEqual({ name: "changed" });
		expect(sut.isPending()).toBe(false);

		await sut.unmount();
	});

	it("Should not save values the schema rejects", async () => {
		// Arrange
		const save = vi.fn<Save>(async () => {});
		const sut = await renderAutosavedForm(save);

		// Act
		await sut.setName("");
		await sut.wait(DEBOUNCE_MS);

		// Assert
		expect(save).not.toHaveBeenCalled();

		await sut.unmount();
	});
});
