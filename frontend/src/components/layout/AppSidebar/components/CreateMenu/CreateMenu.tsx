import { CirclePlus, Terminal } from "lucide-react";

import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu.tsx";
import { Modals, useModalsContext } from "@/contexts/ModalsContext.tsx";

export const CreateMenu = () => {
  const { setOpenModal } = useModalsContext();

  const openCreateCommandModal = () => setOpenModal(Modals.CREATE)(true);

  return (
    <DropdownMenu>
      <DropdownMenuTrigger className="outline-0">
        <CirclePlus
          className="text-muted-foreground cursor-pointer hover:text-primary"
          size={18}
        />
      </DropdownMenuTrigger>
      <DropdownMenuContent>
        <DropdownMenuLabel>Create</DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuItem
          onClick={openCreateCommandModal}
          className="flex flex-row items-center justify-start"
        >
          <Terminal />
          Command
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
};
