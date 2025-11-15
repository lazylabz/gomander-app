import { zodResolver } from "@hookform/resolvers/zod";
import { Route, WandSparkles } from "lucide-react";
import { useEffect, useState } from "react";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { getI18n } from "react-i18next";
import { toast } from "sonner";
import { z } from "zod";

import { type Theme, useTheme } from "@/contexts/theme.tsx";
import { translationsService } from "@/contracts/service.ts";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/design-system/components/ui/card.tsx";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/design-system/components/ui/form.tsx";
import { Input } from "@/design-system/components/ui/input.tsx";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/design-system/components/ui/select";
import { parseError } from "@/helpers/errorHelpers.ts";
import { fetchUserConfig } from "@/queries/fetchUserConfig.ts";
import { EnvironmentPathsField } from "@/screens/SettingsScreen/tabs/ProjectSettings/components/EnvironmentPathsField.tsx";
import { EnvironmentPathsInfoDialog } from "@/screens/SettingsScreen/tabs/UserSettings/components/EnvironmentPathsInfoDialog.tsx";
import {
  userConfigurationStore,
  useUserConfigurationStore,
} from "@/store/userConfigurationStore.ts";
import { saveUserConfig } from "@/useCases/userConfig/saveUserConfig.ts";

const formSchema = z.object({
  environmentPaths: z.array(
    z.object({
      id: z.uuid(),
      path: z.string().min(1, "Path cannot be empty"),
    }),
  ),
  locale: z.string(),
  logLineLimit: z
    .number()
    .int()
    .min(1, "Must be at least 1")
    .max(5000, "Must be at most 5000"),
});

type FormSchemaType = z.infer<typeof formSchema>;

type SupportedLanguage = {
  value: string;
  label: string;
};

const languageValueToLabelMap: Record<string, string> = {
  en: "English",
  es: "Español",
};

const changeLanguage = async (lang: string) => {
  const i18n = getI18n();

  if (i18n.language === lang) {
    return;
  }

  if (!i18n.hasResourceBundle(lang, "translation")) {
    const translations = await translationsService.getTranslation(lang);
    i18n.addResourceBundle(lang, "translation", translations);
  }

  await i18n.changeLanguage(lang);
};

const handleSave = async (formData: FormSchemaType) => {
  const { userConfig } = userConfigurationStore.getState();

  try {
    await changeLanguage(formData.locale);
    await saveUserConfig({
      lastOpenedProjectId: userConfig.lastOpenedProjectId,
      environmentPaths: formData.environmentPaths,
      logLineLimit: formData.logLineLimit,
      locale: formData.locale,
    });
    toast.success("User settings saved successfully");
  } catch (e) {
    toast.error(parseError(e, "Failed to save user settings"));
  }

  await fetchUserConfig();
};

export const UserSettings = () => {
  const { i18n } = useTranslation();

  const userConfig = useUserConfigurationStore((state) => state.userConfig);

  const { rawTheme, setRawTheme } = useTheme();
  const [supportedLanguages, setSupportedLanguages] = useState<
    SupportedLanguage[]
  >([]);

  useEffect(() => {
    const loadSupportedLanguages = async () => {
      try {
        const languages = await translationsService.getSupportedLanguages();
        const languageOptions = languages.map(
          (lang): SupportedLanguage => ({
            value: lang,
            label: languageValueToLabelMap[lang] || lang,
          }),
        );
        setSupportedLanguages(languageOptions);
      } catch (error) {
        console.error("Failed to load supported languages:", error);
      }
    };

    loadSupportedLanguages();
  }, []);

  const form = useForm<FormSchemaType>({
    resolver: zodResolver(formSchema),
    values: {
      environmentPaths: userConfig.environmentPaths,
      logLineLimit: userConfig.logLineLimit,
      locale: i18n.language,
    },
  });

  const formWatcher = form.watch();

  useEffect(() => {
    if (!form.formState.isDirty) {
      return;
    }
    const timeout = setTimeout(async () => {
      await form.handleSubmit(handleSave)();
    }, 300);

    return () => clearTimeout(timeout);
  }, [form, formWatcher]);

  return (
    <Form {...form}>
      <form className="w-full h-full flex flex-col justify-between">
        <div className="flex flex-col gap-2">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center">
                <Route size={20} />
                <span className="ml-2 mr-2">Environment paths</span>
                <EnvironmentPathsInfoDialog />
              </CardTitle>
              <CardDescription>
                These paths will be used to resolve commands and executables.
              </CardDescription>
            </CardHeader>
            <CardContent>
              <EnvironmentPathsField />
            </CardContent>
          </Card>
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center space-x-2">
                <WandSparkles size={20} />
                <span>Preferences</span>
              </CardTitle>
              <CardDescription>Make gomander your own!</CardDescription>
            </CardHeader>
            <CardContent className="space-y-3">
              <FormField
                control={form.control}
                name="locale"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Language</FormLabel>
                    <FormControl>
                      <Select
                        onValueChange={field.onChange}
                        defaultValue={field.value}
                      >
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue placeholder="Select your preferred language" />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          {supportedLanguages.map((language) => (
                            <SelectItem
                              key={language.value}
                              value={language.value}
                            >
                              {language.label}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormItem>
                <FormLabel>Theme</FormLabel>
                <Select
                  onValueChange={(value) => {
                    setRawTheme(value as Theme);
                  }}
                  value={rawTheme}
                >
                  <SelectTrigger>
                    <SelectValue placeholder="Select your preferred theme" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="system">System theme</SelectItem>
                    <SelectItem value="light">Light theme</SelectItem>
                    <SelectItem value="dark">Dark theme</SelectItem>
                  </SelectContent>
                </Select>
                <FormDescription className="text-xs">
                  (The system theme will adapt to your operating system's theme
                  settings)
                </FormDescription>
              </FormItem>
              <FormField
                control={form.control}
                name="logLineLimit"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Log line limit</FormLabel>
                    <FormControl>
                      <Input
                        type="number"
                        min={1}
                        max={5000}
                        {...field}
                        onChange={(e) => field.onChange(Number(e.target.value))}
                      />
                    </FormControl>
                    <FormDescription className="text-xs">
                      Maximum number of log lines to keep per command (1-5000).
                      The recommended value is 100. Bigger values may impact
                      performance.
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </CardContent>
          </Card>
        </div>
      </form>
    </Form>
  );
};
