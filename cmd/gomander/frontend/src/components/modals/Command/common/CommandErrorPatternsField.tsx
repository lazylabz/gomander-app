import { useFormContext } from "react-hook-form";
import { useTranslation } from "react-i18next";

import type { FormSchemaType } from "@/components/modals/Command/common/formSchema.ts";
import {
	FormControl,
	FormDescription,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from "@/design-system/components/ui/form.tsx";
import { Textarea } from "@/design-system/components/ui/textarea.tsx";

export const CommandErrorPatternsField = () => {
	const { t } = useTranslation();
	const form = useFormContext<FormSchemaType>();

	return (
		<FormField
			control={form.control}
			name="errorPatterns"
			render={({ field }) => (
				<FormItem>
					<FormLabel>{t("commandForm.errorPatternsLabel")}</FormLabel>
					<FormDescription className="text-xs">
						{t("commandForm.errorPatternsDescription")}
					</FormDescription>
					<FormControl>
						<Textarea
							autoComplete="off"
							autoCorrect="off"
							autoCapitalize="off"
							placeholder="[nodemon] app crashed"
							{...field}
						/>
					</FormControl>
					<FormMessage />
				</FormItem>
			)}
		/>
	);
};
