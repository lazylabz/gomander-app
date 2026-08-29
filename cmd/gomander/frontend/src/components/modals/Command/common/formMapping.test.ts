import { describe, expect, it } from "vitest";

import {
	emptyFormValues,
	toCommand,
	toFormValues,
} from "@/components/modals/Command/common/formMapping.ts";
import { CommandBuilder } from "@/testing/builders/command.ts";

describe("emptyFormValues", () => {
	it("Should hand back a blank form", () => {
		// Act
		const result = emptyFormValues();

		// Assert
		expect(result).toEqual({
			name: "",
			command: "",
			workingDirectory: "",
			link: "",
			errorPatterns: "",
		});
	});
});

describe("toFormValues", () => {
	it("Should map a command's editable fields onto the form", () => {
		// Arrange
		const command = new CommandBuilder()
			.withName("Dev server")
			.withCommand("pnpm dev")
			.withWorkingDirectory("apps/web")
			.withLink("http://localhost:3000")
			.build();

		// Act
		const result = toFormValues(command);

		// Assert
		expect(result).toEqual({
			name: "Dev server",
			command: "pnpm dev",
			workingDirectory: "apps/web",
			link: "http://localhost:3000",
			errorPatterns: "",
		});
	});

	it("Should write one error pattern per line", () => {
		// Arrange
		const command = new CommandBuilder()
			.withErrorPatterns("ERROR", "FATAL")
			.build();

		// Act
		const result = toFormValues(command);

		// Assert
		expect(result.errorPatterns).toBe("ERROR\nFATAL");
	});

	it("Should hand back a blank form when there is no command", () => {
		// Act
		const result = toFormValues(null);

		// Assert
		expect(result).toEqual(emptyFormValues());
	});
});

describe("toCommand", () => {
	const base = { id: "cmd-1", projectId: "project-1", position: 3 };

	const values = {
		name: "Dev server",
		command: "pnpm dev",
		workingDirectory: "apps/web",
		link: "http://localhost:3000",
		errorPatterns: "",
	};

	it("Should keep the base fields the form does not own", () => {
		// Act
		const result = toCommand(values, base);

		// Assert
		expect(result).toMatchObject(base);
	});

	it("Should read one error pattern per line", () => {
		// Act
		const result = toCommand(
			{ ...values, errorPatterns: "ERROR\nFATAL" },
			base,
		);

		// Assert
		expect(result.errorPatterns).toEqual(["ERROR", "FATAL"]);
	});

	it("Should drop the lines that hold no pattern", () => {
		// Act
		const result = toCommand(
			{ ...values, errorPatterns: "ERROR\n\n   \nFATAL\n" },
			base,
		);

		// Assert
		expect(result.errorPatterns).toEqual(["ERROR", "FATAL"]);
	});

	it("Should read no pattern out of an empty field", () => {
		// Act
		const result = toCommand({ ...values, errorPatterns: "" }, base);

		// Assert
		expect(result.errorPatterns).toEqual([]);
	});

	it("Should hand back the command it was given after a round trip", () => {
		// Arrange
		const command = new CommandBuilder()
			.withErrorPatterns("ERROR", "FATAL")
			.build();

		// Act
		const result = toCommand(toFormValues(command), command);

		// Assert
		expect(result).toEqual(command);
	});
});
