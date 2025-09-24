import { FolderIcon } from "lucide-react";
import type { ChangeEvent } from "react";

import { Input } from "@/components/ui/input.tsx";
import { helpersService } from "@/contracts/service.ts";
import { cn } from "@/lib/utils.ts";

export const FSInput = (props: React.ComponentProps<typeof Input>) => {
  const handleAskForDirPath = async () => {
    const path = await helpersService.askForDirPath();
    if (path) {
      props.onChange?.({
        target: { value: path },
      } as ChangeEvent<HTMLInputElement>);
    }
  };

  return (
    <div className="w-full relative">
      <Input {...props} className={cn("pr-8", props.className)} />
      <button
        type="button"
        onClick={handleAskForDirPath}
        className="absolute right-2 top-2 z-10 bg-background"
      >
        <FolderIcon
          size={20}
          className="text-foreground opacity-75 hover:opacity-100 transition-opacity cursor-pointer"
        />
      </button>
    </div>
  );
};
