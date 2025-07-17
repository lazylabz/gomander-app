import { useEffect, useState } from "react";
import { EventsOn } from "../wailsjs/runtime";
import { GetCommands } from "../wailsjs/go/main/LogServer";
import type { Command } from "./types/contracts.ts";
import { Event } from "./types/contracts.ts";

import { Button } from "@/components/ui/button.tsx";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog.tsx";

function App() {
  const [, setCommands] = useState<Record<string, Command>>({});
  const [createModalOpen, setCreateModalOpen] = useState(false);

  const refreshCommands = async () => {
    const commandsData = await GetCommands();

    setCommands(commandsData);
  };

  const openCreateModal = () => {
    setCreateModalOpen(true);
  };

  useEffect(() => {
    EventsOn(Event.GET_COMMANDS, () => {
      refreshCommands();
    });
  });

  return (
    <main className="w-full h-full bg-white">
      <Dialog open={createModalOpen} onOpenChange={setCreateModalOpen}>
        <form>
          <DialogContent className="sm:max-w-[425px]">
            <DialogHeader>
              <DialogTitle>Add new command</DialogTitle>
              <DialogDescription>
                Make changes to your profile here.
              </DialogDescription>
            </DialogHeader>
            <DialogFooter>
              <DialogClose asChild>
                <Button variant="outline">Cancel</Button>
              </DialogClose>
              <Button type="submit">Create</Button>
            </DialogFooter>
          </DialogContent>
        </form>
      </Dialog>
      <Button
        onClick={openCreateModal}
        variant="outline"
        className="fixed top-2 right-2 cursor-pointer"
      >
        Add new command
      </Button>
    </main>
  );
}

export default App;
