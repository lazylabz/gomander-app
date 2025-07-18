import { ChevronRight, Folder, Play } from "lucide-react";

import { CreateMenu } from "@/components/layout/AppSidebar/components/CreateMenu/CreateMenu.tsx";
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible.tsx";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar.tsx";
import { useDataContext } from "@/contexts/DataContext.tsx";

export const AppSidebar = () => {
  const { commands, execCommand, setActiveCommandId, activeCommandId } =
    useDataContext();

  const handleCommandClick = (commandId: string) => async () => {
    setActiveCommandId(commandId);
    await execCommand(commandId);
  };

  const onCommandSectionClick = (commandId: string) => () => {
    setActiveCommandId(commandId);
  };

  return (
    <Sidebar collapsible="icon">
      <SidebarHeader className="flex flex-row items-center justify-between p-2">
        <h1 className="text-xl font-semibold pl-2">Gomander</h1>
        <CreateMenu />
      </SidebarHeader>
      <SidebarContent>
        <Collapsible className="group/collapsible">
          <SidebarGroup>
            <SidebarGroupLabel
              asChild
              className="group/label text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground text-sm"
            >
              <CollapsibleTrigger className="group flex items-center gap-2 p-2">
                <ChevronRight className="transition-transform group-data-[state=open]:rotate-90" />
                <Folder />
                <p>Your commands</p>
              </CollapsibleTrigger>
            </SidebarGroupLabel>
            <CollapsibleContent className="pl-4">
              <SidebarGroupContent>
                <SidebarMenu>
                  {Object.values(commands).map((command) => (
                    <SidebarMenuItem key={command.id}>
                      <SidebarMenuButton
                        asChild
                        isActive={activeCommandId === command.id}
                      >
                        <div
                          onClick={onCommandSectionClick(command.id)}
                          className="flex flex-row justify-between items-center w-full"
                        >
                          {command.name}
                          <Play
                            className="text-muted-foreground cursor-pointer hover:text-primary"
                            onClick={handleCommandClick(command.id)}
                          />
                        </div>
                      </SidebarMenuButton>
                    </SidebarMenuItem>
                  ))}
                </SidebarMenu>
              </SidebarGroupContent>
            </CollapsibleContent>
          </SidebarGroup>
        </Collapsible>
      </SidebarContent>
      <SidebarFooter />
    </Sidebar>
  );
};
