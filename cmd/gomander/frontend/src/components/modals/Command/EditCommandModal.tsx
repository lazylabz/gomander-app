import { zodResolver } from "@hookform/resolvers/zod";
import { Terminal } from "lucide-react";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";

import { CommandCommandField } from "@/components/modals/Command/common/CommandCommandField.tsx";
import { CommandComputedPath } from "@/components/modals/Command/common/CommandComputedPath.tsx";
import { CommandErrorPatternsField } from "@/components/modals/Command/common/CommandErrorPatternsField.tsx";
import { CommandLinkField } from "@/components/modals/Command/common/CommandLinkField.tsx";
import { CommandNameField } from "@/components/modals/Command/common/CommandNameField.tsx";
import { CommandWorkingDirectoryField } from "@/components/modals/Command/common/CommandWorkingDirectoryField.tsx";
import {
	toCommand,
	toFormValues,
} from "@/components/modals/Command/common/formMapping.ts";
import {
	type FormSchemaType,
	formSchema,
} from "@/components/modals/Command/common/formSchema.ts";
import type { Command } from "@/contracts/types.ts";
import {
	Accordion,
	AccordionContent,
	AccordionItem,
	AccordionTrigger,
} from "@/design-system/components/ui/accordion.tsx";
import { Button } from "@/design-system/components/ui/button.tsx";
import {
	Dialog,
	DialogClose,
	DialogContent,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from "@/design-system/components/ui/dialog.tsx";
import { Form } from "@/design-system/components/ui/form.tsx";
import { editCommand } from "@/useCases/command/editCommand.ts";

export const EditCommandModal = ({
	command,
	open,
	setOpen,
}: {
	command: Command | null;
	open: boolean;
	setOpen: (open: boolean) => void;
}) => {
	const { t } = useTranslation();
	const form = useForm<FormSchemaType>({
		resolver: zodResolver(formSchema),
		values: toFormValues(command),
	});

	const onSubmit = async (values: FormSchemaType) => {
		if (!command) {
			return;
		}

		const edited = await editCommand(toCommand(values, command));
		if (!edited) {
			return;
		}

		setOpen(false);
		form.reset();
	};

	const onOpenChange = (open: boolean) => {
		setOpen(open);
		if (!open) {
			form.reset();
		}
	};

	return (
		<Dialog open={open} onOpenChange={onOpenChange}>
			<DialogContent className="sm:max-w-[628px]">
				<Form {...form}>
					<form onSubmit={form.handleSubmit(onSubmit)} className="w-full">
						<DialogHeader className="flex flex-row items-center gap-6">
							<Terminal />
							<DialogTitle>{t("modal.editCommand.title")}</DialogTitle>
						</DialogHeader>
						<div className="my-4 space-y-2">
							<div className="space-y-6">
								<CommandNameField />
								<CommandCommandField />
								<CommandWorkingDirectoryField />
								<CommandComputedPath />
							</div>
							<Accordion type="single" collapsible>
								<AccordionItem value="1">
									<AccordionTrigger>{t("common.advanced")}</AccordionTrigger>
									<AccordionContent>
										<div className="space-y-6">
											<CommandLinkField />
											<CommandErrorPatternsField />
										</div>
									</AccordionContent>
								</AccordionItem>
							</Accordion>
						</div>
						<DialogFooter>
							<DialogClose asChild>
								<Button type="button" variant="outline">
									{t("common.cancel")}
								</Button>
							</DialogClose>
							<Button type="submit">{t("common.save")}</Button>
						</DialogFooter>
					</form>
				</Form>
			</DialogContent>
		</Dialog>
	);
};
