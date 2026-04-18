import { useFormContext } from "react-hook-form";
import { useTranslation } from "react-i18next";

import { FSInput } from "@/components/inputs/FSInput.tsx";
import type { FormSchemaType } from "@/components/modals/Command/common/formSchema.ts";
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/design-system/components/ui/form.tsx";

export const CommandWorkingDirectoryField = () => {
  const { t } = useTranslation();
  const form = useFormContext<FormSchemaType>();

  return (
    <FormField
      control={form.control}
      name="workingDirectory"
      render={({ field }) => (
        <FormItem>
          <FormLabel>{t('commandForm.workingDirectoryLabel')}</FormLabel>
          <FormControl>
            <FSInput
              autoComplete="off"
              autoCorrect="off"
              autoCapitalize="off"
              {...field}
            />
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
    />
  );
};
