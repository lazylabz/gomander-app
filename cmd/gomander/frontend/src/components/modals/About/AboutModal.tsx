import { Download, ExternalLink, Heart } from "lucide-react";

import { useVersionContext } from "@/contexts/version.tsx";
import { externalBrowserService } from "@/contracts/service.ts";
import {
  Avatar,
  AvatarFallback,
  AvatarImage,
} from "@/design-system/components/ui/avatar.tsx";
import { Button } from "@/design-system/components/ui/button.tsx";
import {
  Dialog,
  DialogContent,
} from "@/design-system/components/ui/dialog.tsx";
import { cn } from "@/design-system/lib/utils.ts";
import { GithubIcon } from "@/icons/GithubIcon.tsx";

export const AboutModal = ({
  open,
  setOpen,
}: {
  open: boolean;
  setOpen: (open: boolean) => void;
}) => {
  const { currentVersion, newVersion, openLatestReleasePage } =
    useVersionContext();

  const handleGithubClick = () => {
    const url = `https://github.com/lazylabz/gomander-app`;
    externalBrowserService.browserOpenURL(url);
  };

  const handleTeamClick = () => {
    const url = `https://lazylabz.github.io/`;
    externalBrowserService.browserOpenURL(url);
  };

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
                onClick={openLatestReleasePage}
                variant="outline"
                className="cursor-pointer"
              >
                <Download className="size-4" />
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
            onClick={handleGithubClick}
            className="focus-visible:outline-none cursor-pointer w-full flex flex-1 items-center justify-center gap-3 px-4 py-3 bg-card hover:bg-accent border border-border shadow-sm hover:shadow-md rounded-lg transition-colors group"
          >
            <GithubIcon className="size-5 text-foreground" />
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
            onClick={handleTeamClick}
            className="focus-visible:outline-none cursor-pointer w-full flex flex-1 items-center justify-center gap-3 px-4 py-3 bg-card hover:bg-accent border border-border shadow-sm hover:shadow-md rounded-lg transition-all group"
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
