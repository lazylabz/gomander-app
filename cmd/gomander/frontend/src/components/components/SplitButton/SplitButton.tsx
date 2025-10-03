import { ChevronDown } from "lucide-react";
import { Children, type ComponentProps, type ReactNode } from "react";

import { Button, buttonVariants } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { cn } from "@/lib/utils.ts";

export function SplitButton({
  children,
  handleDefaultAction,
  extraActions,
  variant = "default",
  size = "default",
  className,
}: {
  children: ReactNode;
  handleDefaultAction?: () => void;
  extraActions?: {
    label: string;
    handleClick: () => void;
  }[];
  variant?: ComponentProps<typeof Button>["variant"];
  size?: ComponentProps<typeof Button>["size"];
  className?: string;
}) {
  return (
    <div
      className={cn(
        "flex flex-row items-center",
        buttonVariants({ variant, size, className }),
      )}
    >
      <div
        role="button"
        className="cursor-pointer flex flex-row items-center gap-2"
        onClick={handleDefaultAction}
      >
        {children}
      </div>
      {!!extraActions?.length && (
        <DropdownMenu>
          <DropdownMenuTrigger className="m-0 p-0">
            <ChevronDown className="h-4 w-4" />
          </DropdownMenuTrigger>
          <DropdownMenuContent sideOffset={12} align="end">
            {Children.toArray(
              extraActions.map(({ label, handleClick }) => (
                <DropdownMenuItem onClick={handleClick}>
                  {label}
                </DropdownMenuItem>
              )),
            )}
          </DropdownMenuContent>
        </DropdownMenu>
      )}
    </div>
  );
}
