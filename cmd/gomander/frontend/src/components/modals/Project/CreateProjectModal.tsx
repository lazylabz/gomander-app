import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";

import { BaseWorkingDirectoryField } from "@/components/modals/Project/common/BaseWorkingDirectoryField.tsx";
import {
	type FormSchemaType,
	formSchema,
} from "@/components/modals/Project/common/createSchema.ts";
import { ProjectNameField } from "@/components/modals/Project/common/ProjectNameField.tsx";
import { Button } from "@/design-system/components/ui/button.tsx";
import {
	Dialog,
	DialogHeader,
	DialogTitle,
} from "@/design-system/components/ui/dialog";
import { DialogContent } from "@/design-system/components/ui/dialog.tsx";
import { Form } from "@/design-system/components/ui/form.tsx";
import { createProject } from "@/useCases/project/createProject.ts";

export const CreateProjectModal = ({
	open,
	setOpen,
}: {
	open: boolean;
	setOpen: (open: boolean) => void;
}) => {
	const { t } = useTranslation();
	const form = useForm<FormSchemaType>({
		resolver: zodResolver(formSchema),
		defaultValues: {
			name: "",
			baseWorkingDirectory: "",
		},
	});

	const handleOpenChange = (open: boolean) => {
		setOpen(open);
		if (!open) {
			form.reset();
		}
	};

	const onSubmit = async (values: FormSchemaType) => {
		const created = await createProject(
			values.name,
			values.baseWorkingDirectory,
		);
		if (!created) {
			return;
		}

		setOpen(false);
		form.reset();
	};

	return (
		<Dialog open={open} onOpenChange={handleOpenChange}>
			<DialogContent>
				<DialogHeader>
					<DialogTitle>{t("modal.createProject.title")}</DialogTitle>
					<Form {...form}>
						<form
							onSubmit={form.handleSubmit(onSubmit)}
							className="w-full mt-2 flex flex-col gap-4"
						>
							<ProjectNameField<FormSchemaType> />
							<BaseWorkingDirectoryField<FormSchemaType> />
							<Button className="self-end" type="submit">
								{t("common.save")}
							</Button>
						</form>
					</Form>
				</DialogHeader>
			</DialogContent>
		</Dialog>
	);
};
