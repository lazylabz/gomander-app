import { ArrowLeft } from "lucide-react";
import { useLocation, useNavigate } from "react-router";

import {
  Tabs,
  TabsContent,
  TabsList,
  TabsTrigger,
} from "@/components/ui/tabs.tsx";

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
    <div className="bg-background p-8">
      <ArrowLeft onClick={handleBack} />
      <h1 className="text-3xl mb-4">Settings</h1>
      <Tabs defaultValue={defaultTab} className="w-[400px]">
        <TabsList>
          <TabsTrigger value={SettingsTab.User}>User</TabsTrigger>
          <TabsTrigger value={SettingsTab.Project}>Project</TabsTrigger>
        </TabsList>
        <TabsContent value={SettingsTab.User}>
          Make changes to your account here.
        </TabsContent>
        <TabsContent value={SettingsTab.Project}>
          Change your password here.
        </TabsContent>
      </Tabs>
    </div>
  );
};
