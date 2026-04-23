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

export const CommandLinkField = () => {
	const { t } = useTranslation();
	const form = useFormContext<FormSchemaType>();

	return (
		<FormField
			control={form.control}
			name="link"
			render={({ field }) => (
				<FormItem>
					<FormLabel>{t("commandForm.linkLabel")}</FormLabel>
					<FormControl>
						<Input
							autoComplete="off"
							autoCorrect="off"
							autoCapitalize="off"
							placeholder="http://localhost:3000"
							{...field}
						/>
					</FormControl>
					<FormMessage />
				</FormItem>
			)}
		/>
	);
};
