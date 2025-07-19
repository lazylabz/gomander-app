import { ChevronRight, Folder } from "lucide-react";

import { CommandMenuItem } from "@/components/layout/AppSidebar/components/CommandMenuItem/CommandMenuItem.tsx";
import {
  Collapsible,
  CollapsibleContent,
  CollapsibleTrigger,
} from "@/components/ui/collapsible.tsx";
import {
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuItem,
} from "@/components/ui/sidebar.tsx";
import { useDataContext } from "@/contexts/DataContext.tsx";

export const MyCommandsSection = () => {
  const { commands } = useDataContext();

  return (
    <Collapsible defaultOpen className="group/collapsible">
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
                  <CommandMenuItem command={command} />
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </CollapsibleContent>
      </SidebarGroup>
    </Collapsible>
  );
};
