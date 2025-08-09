import { zodResolver } from "@hookform/resolvers/zod";
import { Route, Save } from "lucide-react";
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
import { Form } from "@/components/ui/form.tsx";
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

  const navigate = useNavigate();

  const form = useForm<FormType>({
    resolver: zodResolver(formSchema),
    values: {
      environmentPaths:
        userConfig.environmentPaths.map((p) => ({ value: p })) || [],
    },
  });

  const onSubmit = async (data: FormType) => {
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
        <Button type="submit" className="self-end">
          <Save />
          Save
        </Button>
      </form>
    </Form>
  );
};
