import { Download, ExternalLink, Heart } from "lucide-react";

import {
  Avatar,
  AvatarFallback,
  AvatarImage,
} from "@/components/ui/avatar.tsx";
import { Button } from "@/components/ui/button.tsx";
import { Dialog, DialogContent } from "@/components/ui/dialog.tsx";
import { useVersionContext } from "@/contexts/version.tsx";
import { cn } from "@/lib/utils.ts";

export const AboutModal = ({
  open,
  setOpen,
}: {
  open: boolean;
  setOpen: (open: boolean) => void;
}) => {
  const { currentVersion, newVersion } = useVersionContext();

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogContent className="sm:max-w-[628px]">
        {/* Logo and Title */}
        <div className="text-center">
          <Avatar className="size-16 mx-auto mb-4 rounded-xl">
            <AvatarImage src="/app-logo.png" />
            <AvatarFallback className="rounded-xl text-xl font-extralight bg-card-foreground text-card">
              G.
            </AvatarFallback>
          </Avatar>
          <h3 className="text-xl font-bold text-foreground">
            Gomander
            <span className="ml-2 font-normal text-sm text-muted-foreground">
              v{currentVersion}
            </span>
          </h3>
        </div>
        {/* Update Notice */}
        {newVersion && (
          <div className="bg-sky-50 dark:bg-sky-950/40 border border-b-0 border-sky-200 dark:border-sky-950 shadow-sm shadow-sky-200 dark:shadow-sky-950 rounded-lg p-4">
            <div className="flex items-center">
              <div className="flex-1">
                <p className="text-sm font-medium text-foreground mb-1">
                  Version {newVersion} is available
                </p>
                <p className="text-xs text-muted-foreground">
                  Get the latest features and improvements
                </p>
              </div>
              <Button
                onClick={() => {}}
                variant="outline"
                className="cursor-pointer"
              >
                <Download className="w-3 h-3" />
                Download Update
              </Button>
            </div>
          </div>
        )}
        {/* Description */}
        <p className="text-sm text-muted-foreground leading-relaxed">
          A simple GUI for managing and organizing your development commands.
          Built to streamline your workflow and eliminate terminal chaos.
        </p>
        {/* CTAs */}
        <div className={cn("flex gap-4", newVersion ? "flex-row" : "flex-col")}>
          {/* GitHub CTA */}
          <button
            onClick={() => {}}
            className="cursor-pointer w-full flex flex-1 items-center justify-center gap-3 px-4 py-3 bg-card hover:bg-accent border border-border shadow-sm hover:shadow-md rounded-lg transition-colors group"
          >
            <svg
              className="w-5 h-5 text-foreground"
              role="img"
              viewBox="0 0 24 24"
              xmlns="http://www.w3.org/2000/svg"
            >
              <path
                fill="currentColor"
                d="M12 .297c-6.63 0-12 5.373-12 12 0 5.303 3.438 9.8 8.205 11.385.6.113.82-.258.82-.577 0-.285-.01-1.04-.015-2.04-3.338.724-4.042-1.61-4.042-1.61C4.422 18.07 3.633 17.7 3.633 17.7c-1.087-.744.084-.729.084-.729 1.205.084 1.838 1.236 1.838 1.236 1.07 1.835 2.809 1.305 3.495.998.108-.776.417-1.305.76-1.605-2.665-.3-5.466-1.332-5.466-5.93 0-1.31.465-2.38 1.235-3.22-.135-.303-.54-1.523.105-3.176 0 0 1.005-.322 3.3 1.23.96-.267 1.98-.399 3-.405 1.02.006 2.04.138 3 .405 2.28-1.552 3.285-1.23 3.285-1.23.645 1.653.24 2.873.12 3.176.765.84 1.23 1.91 1.23 3.22 0 4.61-2.805 5.625-5.475 5.92.42.36.81 1.096.81 2.22 0 1.606-.015 2.896-.015 3.286 0 .315.21.69.825.57C20.565 22.092 24 17.592 24 12.297c0-6.627-5.373-12-12-12"
              />
            </svg>
            <div className="text-left">
              <div className="font-medium text-foreground text-sm">
                Any feedback?
              </div>
              <div className="text-xs text-muted-foreground">
                Visit our GitHub repository
              </div>
            </div>
            <ExternalLink className="size-4 text-muted-foreground group-hover:text-foreground transition-colors ml-auto" />
          </button>

          {/* LazyLabz CTA */}
          <button
            onClick={() => {}}
            className="cursor-pointer w-full flex flex-1 items-center justify-center gap-3 px-4 py-3 bg-card hover:bg-accent border border-border shadow-sm hover:shadow-md rounded-lg transition-all group"
          >
            <Heart className="size-5 text-foreground" />
            <div className="text-left">
              <div className="font-medium text-foreground text-sm">
                Meet the Team
              </div>
              <div className="text-xs text-muted-foreground">
                Learn more about LazyLabz
              </div>
            </div>
            <ExternalLink className="size-4 text-muted-foreground group-hover:text-foreground transition-colors ml-auto" />
          </button>
        </div>
      </DialogContent>
    </Dialog>
  );
};
