import { CircleAlert, CircleCheck } from "lucide-react";
import { useEffect, useState } from "react";
import { toast } from "sonner";

import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip.tsx";
import { externalBrowserService, helpersService } from "@/contracts/service.ts";

export const VersionSection = () => {
  const [currentRelease, setCurrentRelease] = useState<string>("");
  const [newRelease, setNewRelease] = useState<string | null>(null);
  const [errorLoadingNewRelease, setErrorLoadingNewRelease] =
    useState<Error | null>(null);

  const fetchCurrentRelease = async () => {
    const release = await helpersService.getCurrentRelease();
    setCurrentRelease(release);
  };

  const checkNewRelease = async () => {
    setErrorLoadingNewRelease(null);
    try {
      const release = await helpersService.isThereANewRelease();
      if (release) {
        setNewRelease(release);
      }
    } catch (err) {
      setErrorLoadingNewRelease(err as Error);
      console.error("Error checking for new releases:", err);
      toast.error(`Error checking for new releases`);
    }
  };

  useEffect(() => {
    fetchCurrentRelease();
    checkNewRelease();
  }, []);

  const openNewReleasePage = () => {
    if (!newRelease) return;

    const url = `https://github.com/lazylabz/gomander-app/releases/tag/v${newRelease}`;

    externalBrowserService.browserOpenURL(url);
  };

  return (
    <p className="text-sm text-muted-foreground flex items-center gap-1">
      {currentRelease ? `${currentRelease}` : "..."}
      {newRelease && (
        <Tooltip>
          <TooltipTrigger>
            <CircleAlert
              className="text-orange-400 dark:text-yellow-400 cursor-pointer"
              size={15}
              onClick={openNewReleasePage}
            />
          </TooltipTrigger>
          <TooltipContent>
            <span className="font-semibold">
              New Release v{newRelease} Available!
            </span>
          </TooltipContent>
        </Tooltip>
      )}
      {currentRelease && !newRelease && !errorLoadingNewRelease && (
        <Tooltip>
          <TooltipTrigger>
            <CircleCheck
              className="text-green-600 dark:text-green-200 cursor-pointer"
              size={15}
            />
          </TooltipTrigger>
          <TooltipContent>You are using the latest version!</TooltipContent>
        </Tooltip>
      )}
      {errorLoadingNewRelease && (
        <Tooltip>
          <TooltipTrigger>
            <CircleAlert
              className="text-red-600 dark:text-red-400 cursor-pointer"
              size={15}
            />
          </TooltipTrigger>
          <TooltipContent>Error checking for new releases</TooltipContent>
        </Tooltip>
      )}
    </p>
  );
};
