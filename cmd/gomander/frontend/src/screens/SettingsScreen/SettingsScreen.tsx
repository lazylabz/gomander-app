import { ArrowLeft, PanelsTopLeft, Settings, User } from "lucide-react";

import { Button } from "@/design-system/components/ui/button";
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from "@/design-system/components/ui/tabs.tsx";
import {
  SettingsTab,
  useSettingsContext,
} from "@/screens/SettingsScreen/contexts/settingsContext.tsx";
import { ProjectSettings } from "@/screens/SettingsScreen/tabs/ProjectSettings/ProjectSettings.tsx";
import { UserSettings } from "@/screens/SettingsScreen/tabs/UserSettings/UserSettings.tsx";

export const SettingsScreen = () => {
  const { initialTab, closeSettings } = useSettingsContext();

  return (
    <div className="bg-background p-6 flex flex-col h-full">
      <div className="flex items-center space-x-4 mb-6">
        <Button
          variant="ghost"
          size="sm"
          className="p-2 cursor-pointer"
          onClick={closeSettings}
        >
          <ArrowLeft className="h-4 w-4" />
        </Button>
        <div className="flex items-center space-x-2">
          <Settings className="h-6 w-6" />
          <h1 className="text-2xl font-bold">Settings</h1>
        </div>
      </div>
      <Tabs defaultValue={initialTab} className="w-full flex-1">
        <TabsList className="grid w-full grid-cols-2">
          <TabsTrigger
            value={SettingsTab.User}
            className="flex items-center space-x-2 cursor-pointer"
          >
            <User className="h-4 w-4" />
            <span>User Settings</span>
          </TabsTrigger>
          <TabsTrigger
            value={SettingsTab.Project}
            className="flex items-center space-x-2 cursor-pointer"
          >
            <PanelsTopLeft className="h-4 w-4" />
            <span>Project Settings</span>
          </TabsTrigger>
        </TabsList>
        <TabsContent value={SettingsTab.User}>
          <UserSettings />
        </TabsContent>
        <TabsContent value={SettingsTab.Project}>
          <ProjectSettings />
        </TabsContent>
      </Tabs>
    </div>
  );
};
