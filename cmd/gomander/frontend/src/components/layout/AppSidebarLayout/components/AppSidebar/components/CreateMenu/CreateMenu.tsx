import { CirclePlus, Group, Terminal } from "lucide-react";
import { useState } from "react";

import { CreateCommandModal } from "@/components/modals/Command/CreateCommandModal.tsx";
import { CreateCommandGroupModal } from "@/components/modals/CommandGroup/CreateCommandGroupModal.tsx";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu.tsx";

export const CreateMenu = () => {
  const [createCommandModalOpen, setCreateCommandModalOpen] = useState(false);
  const [createCommandGroupModalOpen, setCreateCommandGroupModalOpen] =
    useState(false);

  const openCreateCommandModal = () => setCreateCommandModalOpen(true);
  const openCreateCommandGroupModal = () =>
    setCreateCommandGroupModalOpen(true);

  return (
    <>
      <CreateCommandModal
        open={createCommandModalOpen}
        setOpen={setCreateCommandModalOpen}
      />
      <CreateCommandGroupModal
        open={createCommandGroupModalOpen}
        setOpen={setCreateCommandGroupModalOpen}
      />
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
          <DropdownMenuItem
            onClick={openCreateCommandGroupModal}
            className="flex flex-row items-center justify-start"
          >
            <Group />
            Command Group
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </>
  );
};
