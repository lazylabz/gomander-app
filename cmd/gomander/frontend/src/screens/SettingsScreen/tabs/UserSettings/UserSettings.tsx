import { Route, Save, WandSparkles } from "lucide-react";

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
import { useSettingsContext } from "@/screens/SettingsScreen/contexts/settingsContext.tsx";
import { EnvironmentPathsField } from "@/screens/SettingsScreen/tabs/ProjectSettings/components/EnvironmentPathsField.tsx";
import { EnvironmentPathsInfoDialog } from "@/screens/SettingsScreen/tabs/UserSettings/components/EnvironmentPathsInfoDialog.tsx";

export const UserSettings = () => {
  const { settingsForm, saveSettings, hasUnsavedChanges } =
    useSettingsContext();

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
        <Button
          type="submit"
          className="self-end"
          disabled={!hasUnsavedChanges}
        >
          <Save />
          Save
        </Button>
      </form>
    </Form>
  );
};
