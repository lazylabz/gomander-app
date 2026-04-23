import { useFormContext } from "react-hook-form";
import { useTranslation } from "react-i18next";

import type { FormSchemaType } from "@/components/modals/Command/common/formSchema.ts";
import {
	FormControl,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from "@/design-system/components/ui/form.tsx";
import { Input } from "@/design-system/components/ui/input.tsx";

export const CommandNameField = () => {
	const { t } = useTranslation();
	const form = useFormContext<FormSchemaType>();

	return (
		<FormField
			control={form.control}
			name="name"
			render={({ field }) => (
				<FormItem>
					<FormLabel>{t("commandForm.nameLabel")}</FormLabel>
					<FormControl>
						<Input
							autoComplete="off"
							autoCorrect="off"
							autoCapitalize="off"
							placeholder="My awesome command"
							{...field}
						/>
					</FormControl>
					<FormMessage />
				</FormItem>
			)}
		/>
	);
};
