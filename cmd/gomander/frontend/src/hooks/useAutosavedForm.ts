import { zodResolver } from "@hookform/resolvers/zod";
import { useEffect, useRef, useState } from "react";
import { type DefaultValues, type FieldValues, useForm } from "react-hook-form";
import type { ZodType } from "zod";

const AUTOSAVE_DEBOUNCE_MS = 300;

type AutosaverOptions = {
	getValues: () => unknown;
	submit: () => Promise<void>;
	onPendingChange: (isPending: boolean) => void;
};

export const createAutosaver = ({
	getValues,
	submit,
	onPendingChange,
}: AutosaverOptions) => {
	// Both sides of every comparison are reads of the same schema-shaped values, so a
	// JSON snapshot is a stable structural compare.
	const snapshot = () => JSON.stringify(getValues());

	let lastSaved = snapshot();
	let timeout: ReturnType<typeof setTimeout> | undefined;

	const flush = async () => {
		const saving = snapshot();

		if (saving === lastSaved) {
			onPendingChange(false);
			return;
		}

		// Moved forward before the submit so that the form's own submit events do not read
		// as fresh changes, and back again if the save never landed.
		const previous = lastSaved;
		lastSaved = saving;

		try {
			await submit();
		} catch (error) {
			// A save reports its own failures; reaching here means one broke that contract.
			lastSaved = previous;
			console.error("Autosave failed:", error);
		}

		onPendingChange(snapshot() !== saving);
	};

	return {
		onChange: () => {
			if (snapshot() === lastSaved) {
				return;
			}

			onPendingChange(true);
			clearTimeout(timeout);
			timeout = setTimeout(flush, AUTOSAVE_DEBOUNCE_MS);
		},
		cancel: () => clearTimeout(timeout),
	};
};

export const useAutosavedForm = <T extends FieldValues>({
	schema,
	defaultValues,
	save,
}: {
	schema: ZodType<T, T>;
	defaultValues: DefaultValues<T>;
	save: (values: T) => Promise<void>;
}) => {
	const [isPending, setIsPending] = useState(false);

	const form = useForm<T>({
		resolver: zodResolver(schema),
		defaultValues,
	});

	// Read through a ref so that a call site passing an inline save does not rebuild the
	// autosaver - and cancel the window it was waiting on - on every render.
	const saveRef = useRef(save);
	useEffect(() => {
		saveRef.current = save;
	}, [save]);

	useEffect(() => {
		const autosaver = createAutosaver({
			getValues: () => form.getValues(),
			submit: () => form.handleSubmit((values) => saveRef.current(values))(),
			onPendingChange: setIsPending,
		});

		const unsubscribe = form.subscribe({
			formState: { values: true },
			callback: autosaver.onChange,
		});

		return () => {
			autosaver.cancel();
			unsubscribe();
		};
	}, [form]);

	return { form, isPending };
};
