import { CircleAlert, CircleCheck } from "lucide-react";

import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip.tsx";
import { useVersionContext } from "@/contexts/version.tsx";

export const VersionSection = () => {
  const {
    newVersion,
    currentVersion,
    errorLoadingNewVersion,
    openLatestReleasePage,
  } = useVersionContext();

  return (
    <p className="text-sm text-muted-foreground flex items-center gap-1">
      {currentVersion ? `v${currentVersion}` : "..."}
      {newVersion && (
        <Tooltip>
          <TooltipTrigger>
            <CircleAlert
              className="text-orange-400 dark:text-yellow-400 cursor-pointer"
              size={15}
              onClick={openLatestReleasePage}
            />
          </TooltipTrigger>
          <TooltipContent>
            <span className="font-semibold">
              New Version v{newVersion} Available!
            </span>
          </TooltipContent>
        </Tooltip>
      )}
      {currentVersion && !newVersion && !errorLoadingNewVersion && (
        <Tooltip>
          <TooltipTrigger>
            <CircleCheck
              className="text-green-600 dark:text-green-200 cursor-pointer"
              size={15}
              onClick={openLatestReleasePage}
            />
          </TooltipTrigger>
          <TooltipContent>You are using the latest version!</TooltipContent>
        </Tooltip>
      )}
      {errorLoadingNewVersion && (
        <Tooltip>
          <TooltipTrigger>
            <CircleAlert
              className="text-red-600 dark:text-red-400 cursor-pointer"
              size={15}
            />
          </TooltipTrigger>
          <TooltipContent>Error checking for new version</TooltipContent>
        </Tooltip>
      )}
    </p>
  );
};
