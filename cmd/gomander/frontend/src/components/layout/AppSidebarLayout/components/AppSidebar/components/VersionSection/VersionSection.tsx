import { Info } from "lucide-react";

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
    <p className="text-sm text-muted-foreground flex items-center gap-2">
      {currentVersion ? `v${currentVersion}` : "..."}
      {newVersion && (
        <Tooltip>
          <TooltipTrigger>
            <Info
              className="text-orange-400 dark:text-yellow-400 cursor-pointer"
              size={16}
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
            <Info
              className="text-green-600 dark:text-green-200 cursor-pointer"
              size={16}
              onClick={openLatestReleasePage}
            />
          </TooltipTrigger>
          <TooltipContent>You are using the latest version!</TooltipContent>
        </Tooltip>
      )}
      {errorLoadingNewVersion && (
        <Tooltip>
          <TooltipTrigger>
            <Info
              className="text-red-600 dark:text-red-400 cursor-pointer"
              size={16}
              onClick={openLatestReleasePage}
            />
          </TooltipTrigger>
          <TooltipContent>Error checking for new version</TooltipContent>
        </Tooltip>
      )}
    </p>
  );
};
