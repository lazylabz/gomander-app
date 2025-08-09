import { ArrowLeft, PanelsTopLeft, Settings, User } from "lucide-react";
import { useLocation, useNavigate } from "react-router";

import { Button } from "@/components/ui/button";
import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from "@/components/ui/tabs.tsx";
import { ProjectSettings } from "@/screens/SettingsScreen/tabs/ProjectSettings/ProjectSettings.tsx";
import { UserSettings } from "@/screens/SettingsScreen/tabs/UserSettings.tsx";

export enum SettingsTab {
  User = "user",
  Project = "project",
}

export const SettingsScreen = () => {
  const navigate = useNavigate();
  const { state } = useLocation();

  const handleBack = () => {
    navigate(-1);
  };

  const defaultTab = state?.tab || SettingsTab.User;

  return (
    <div className="bg-background p-6 flex flex-col h-full">
      <div className="flex items-center space-x-4 mb-6">
        <Button variant="ghost" size="sm" className="p-2" onClick={handleBack}>
          <ArrowLeft className="h-4 w-4" />
        </Button>
        <div className="flex items-center space-x-2">
          <Settings className="h-6 w-6" />
          <h1 className="text-2xl font-bold">Settings</h1>
        </div>
      </div>
      <Tabs defaultValue={defaultTab} className="w-full flex-1">
        <TabsList className="grid w-full grid-cols-2">
          <TabsTrigger
            value={SettingsTab.User}
            className="flex items-center space-x-2"
          >
            <User className="h-4 w-4" />
            <span>User Settings</span>
          </TabsTrigger>
          <TabsTrigger
            value={SettingsTab.Project}
            className="flex items-center space-x-2"
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
