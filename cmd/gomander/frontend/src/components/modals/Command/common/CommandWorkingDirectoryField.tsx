import { PopoverTrigger } from "@radix-ui/react-popover";
import { FolderCode } from "lucide-react";
import { useMemo, useState } from "react";
import { useFormContext } from "react-hook-form";

import type { FormSchemaType } from "@/components/modals/Command/common/formSchema.ts";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command.tsx";
import {
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form.tsx";
import { Input } from "@/components/ui/input.tsx";
import { Popover, PopoverContent } from "@/components/ui/popover";
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import { useDataContext } from "@/contexts/DataContext.tsx";

export const CommandWorkingDirectoryField = () => {
  const form = useFormContext<FormSchemaType>();
  const { commands } = useDataContext();

  const [open, setOpen] = useState(false);

  const workingDirSuggestions = useMemo(
    () =>
      Array.from(
        new Set(
          Object.values(commands)
            .map((wd) => wd.workingDirectory.trim())
            .filter(Boolean),
        ),
      ).map((workDir) => ({
        id:
          Object.values(commands).find(
            (cmd) => cmd.workingDirectory.trim() === workDir,
          )?.id || "",
        workDir,
      })),
    [commands],
  );

  return (
    <FormField
      control={form.control}
      name="workingDirectory"
      render={({ field }) => (
        <FormItem>
          <FormLabel>Working Directory</FormLabel>
          <FormControl>
            <div className="w-full flex items-center gap-2">
              <Input
                className="flex-1"
                placeholder={"/Users/hackerman/Code"}
                {...field}
              />
              <Popover open={open} onOpenChange={setOpen}>
                <PopoverTrigger>
                  <Tooltip>
                    <TooltipTrigger type="button" className="flex items-center">
                      <FolderCode
                        role="button"
                        size={20}
                        className="text-muted-foreground cursor-pointer hover:text-primary"
                      />
                    </TooltipTrigger>
                    <TooltipContent>
                      Existing working directories
                    </TooltipContent>
                  </Tooltip>
                </PopoverTrigger>
                <PopoverContent
                  className="p-0 w-[580px]"
                  side="bottom"
                  align="end"
                >
                  <Command>
                    <CommandInput placeholder="Search recent..." />
                    <CommandList>
                      <CommandEmpty>No results found.</CommandEmpty>
                      <CommandGroup>
                        {workingDirSuggestions.map((wds) => (
                          <CommandItem
                            key={wds.id}
                            value={wds.workDir}
                            onSelect={() => {
                              field.onChange(wds.workDir);
                              setOpen(false);
                            }}
                          >
                            {wds.workDir}
                          </CommandItem>
                        ))}
                      </CommandGroup>
                    </CommandList>
                  </Command>
                </PopoverContent>
              </Popover>
            </div>
          </FormControl>
          <FormMessage />
        </FormItem>
      )}
    />
  );
};
