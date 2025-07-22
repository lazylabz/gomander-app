import { CirclePlus, Terminal } from "lucide-react";
import { useState } from "react";

import { CreateCommandModal } from "@/components/modals/Command/CreateCommandModal.tsx";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu.tsx";

export const CreateMenu = () => {
  const [modalOpen, setModalOpen] = useState(false);

  const openCreateCommandModal = () => setModalOpen(true);

  return (
    <>
      <CreateCommandModal open={modalOpen} setOpen={setModalOpen} />
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
    </>
  );
};
