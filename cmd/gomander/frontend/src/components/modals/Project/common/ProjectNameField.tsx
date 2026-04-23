import { type Path, useFormContext } from "react-hook-form";
import { useTranslation } from "react-i18next";

import {
	FormControl,
	FormField,
	FormItem,
	FormLabel,
	FormMessage,
} from "@/design-system/components/ui/form.tsx";
import { Input } from "@/design-system/components/ui/input.tsx";

export const ProjectNameField = <T extends { name: string }>() => {
	const { t } = useTranslation();
	const form = useFormContext<T>();

	const name = "name" satisfies keyof T as Path<T>;

	return (
		<FormField
			control={form.control}
			name={name}
			render={({ field }) => (
				<FormItem>
					<FormLabel>{t("projectForm.nameLabel")}</FormLabel>
					<FormControl>
						<Input
							autoComplete="off"
							autoCorrect="off"
							autoCapitalize="off"
							placeholder="My awesome project"
							{...field}
						/>
					</FormControl>
					<FormMessage />
				</FormItem>
			)}
		/>
	);
};
