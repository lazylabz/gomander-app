import { zodResolver } from "@hookform/resolvers/zod";
import { Route, Save, WandSparkles } from "lucide-react";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router";

import { Button } from "@/components/ui/button.tsx";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card.tsx";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form.tsx";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { useTheme } from "@/contexts/theme.tsx";
import { EnvironmentPathsField } from "@/screens/SettingsScreen/tabs/ProjectSettings/components/EnvironmentPathsField.tsx";
import {
  formSchema,
  type FormType,
} from "@/screens/SettingsScreen/tabs/ProjectSettings/formSchema.ts";
import { EnvironmentPathsInfoDialog } from "@/screens/SettingsScreen/tabs/UserSettings/components/EnvironmentPathsInfoDialog.tsx";
import { useUserConfigurationStore } from "@/store/userConfigurationStore.ts";
import { saveUserConfig } from "@/useCases/userConfig/saveUserConfig.ts";

export const UserSettings = () => {
  const userConfig = useUserConfigurationStore((state) => state.userConfig);
  const { setRawTheme, rawTheme } = useTheme();

  const navigate = useNavigate();

  const form = useForm<FormType>({
    resolver: zodResolver(formSchema),
    values: {
      environmentPaths:
        userConfig.environmentPaths.map((p) => ({ value: p })) || [],
      theme: rawTheme || "system",
    },
  });

  const onSubmit = async (data: FormType) => {
    setRawTheme(data.theme);

    await saveUserConfig({
      lastOpenedProjectId: userConfig.lastOpenedProjectId,
      environmentPaths: data.environmentPaths.map((path) => path.value),
    });

    form.reset();

    navigate(-1);
  };

  return (
    <Form {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className="w-full h-full flex flex-col justify-between"
      >
        <div className="flex flex-col gap-2">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center">
                <Route size={20} />
                <span className="ml-2 mr-1">Environment paths</span>
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
            <CardContent>
              <FormField
                control={form.control}
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
                    <FormDescription>
                      <p className="text-xs">
                        (The system theme will adapt to your operating system's
                        theme settings)
                      </p>
                    </FormDescription>
                    <FormMessage />
                  </FormItem>
                )}
              />
            </CardContent>
          </Card>
        </div>
        <Button type="submit" className="self-end">
          <Save />
          Save
        </Button>
      </form>
    </Form>
  );
};
