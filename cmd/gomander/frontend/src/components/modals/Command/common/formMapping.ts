import type { FormSchemaType } from "@/components/modals/Command/common/formSchema.ts";
import type { Command } from "@/contracts/types.ts";

type CommandBase = Pick<Command, "id" | "projectId" | "position">;

export const emptyFormValues = (): FormSchemaType => ({
	name: "",
	command: "",
	workingDirectory: "",
	link: "",
	errorPatterns: "",
});

export const toFormValues = (command: Command | null): FormSchemaType =>
	command
		? {
				name: command.name,
				command: command.command,
				workingDirectory: command.workingDirectory,
				link: command.link,
				errorPatterns: command.errorPatterns.join("\n"),
			}
		: emptyFormValues();

export const toCommand = (
	values: FormSchemaType,
	base: CommandBase,
): Command => ({
	id: base.id,
	projectId: base.projectId,
	position: base.position,
	name: values.name,
	command: values.command,
	workingDirectory: values.workingDirectory,
	link: values.link,
	errorPatterns: values.errorPatterns
		.split("\n")
		.filter((pattern) => pattern.trim() !== ""),
});
