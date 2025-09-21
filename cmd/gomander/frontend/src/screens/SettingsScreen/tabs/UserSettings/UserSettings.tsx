import { Route, Save, WandSparkles } from "lucide-react";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";

import { Button } from "@/design-system/components/ui/button.tsx";
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
import { translationsService } from "@/contracts/service";
import { useSettingsContext } from "@/screens/SettingsScreen/contexts/settingsContext.tsx";
import { EnvironmentPathsField } from "@/screens/SettingsScreen/tabs/ProjectSettings/components/EnvironmentPathsField.tsx";
import { EnvironmentPathsInfoDialog } from "@/screens/SettingsScreen/tabs/UserSettings/components/EnvironmentPathsInfoDialog.tsx";

export const UserSettings = () => {
  const { t } = useTranslation();

  const { settingsForm, saveSettings, hasUnsavedChanges } =
    useSettingsContext();
  const [supportedLanguages, setSupportedLanguages] = useState<string[]>([]);

  useEffect(() => {
    const loadSupportedLanguages = async () => {
      try {
        const languages = await translationsService.getSupportedLanguages();
        setSupportedLanguages(languages);
      } catch (error) {
        console.error("Failed to load supported languages:", error);
      }
    };

    loadSupportedLanguages();
  }, []);

  return (
    <Form {...settingsForm}>
      <form
        onSubmit={settingsForm.handleSubmit(saveSettings)}
        className="w-full h-full flex flex-col justify-between"
      >
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
                control={settingsForm.control}
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
                            <SelectItem key={language} value={language}>
                              {language === "en"
                                ? "English"
                                : language === "es"
                                  ? "Español"
                                  : language}
                            </SelectItem>
                          ))}
                        </SelectContent>
                      </Select>
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={settingsForm.control}
                name="theme"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>Theme</FormLabel>
                    <FormControl>
                      <Select
                        onValueChange={field.onChange}
                        defaultValue={field.value}
                      >
                        <FormControl>
                          <SelectTrigger>
                            <SelectValue placeholder="Select your preferred theme" />
                          </SelectTrigger>
                        </FormControl>
                        <SelectContent>
                          <SelectItem value="system">System theme</SelectItem>
                          <SelectItem value="light">Light theme</SelectItem>
                          <SelectItem value="dark">Dark theme</SelectItem>
                        </SelectContent>
                      </Select>
                    </FormControl>
                    <FormDescription className="text-xs">
                      (The system theme will adapt to your operating system's
                      theme settings)
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={settingsForm.control}
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
                        onChange={(e) =>
                          field.onChange(Number(e.target.value))
                        }
                      />
                    </FormControl>
                    <FormDescription className="text-xs">
                      Maximum number of log lines to keep per command
                      (1-5000). The recommended value is 100. Bigger values
                      may impact performance.
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </CardContent>
          </Card>
        </div>
        <Button
          type="submit"
          className="self-end cursor-pointer"
          disabled={!hasUnsavedChanges}
        >
          <Save />
          {t("actions.save")}
        </Button>
      </form>
    </Form>
  );
};
