import { Info } from "lucide-react";

import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/components/ui/tooltip.tsx";
import { useVersionContext } from "@/contexts/version.tsx";

export const VersionSection = ({
  openAboutModal,
}: {
  openAboutModal: () => void;
}) => {
  const { newVersion, currentVersion, errorLoadingNewVersion } =
    useVersionContext();

  return (
    <Tooltip>
      <TooltipTrigger className="cursor-pointer" onClick={openAboutModal}>
        <p className="text-sm text-muted-foreground flex items-center gap-2">
          {currentVersion ? `v${currentVersion}` : "..."}
          {newVersion && (
            <>
              <Info
                className="text-orange-400 dark:text-yellow-400 cursor-pointer"
                size={16}
                onClick={openAboutModal}
              />
              <TooltipContent>
                <span className="font-semibold">
                  New Version v{newVersion} Available!
                </span>
              </TooltipContent>
            </>
          )}
          {currentVersion && !newVersion && !errorLoadingNewVersion && (
            <>
              <Info
                className="text-green-600 dark:text-green-200 cursor-pointer"
                size={16}
                onClick={openAboutModal}
              />
              <TooltipContent>You are using the latest version!</TooltipContent>
            </>
          )}
          {errorLoadingNewVersion && (
            <>
              <Info
                className="text-red-600 dark:text-red-400 cursor-pointer"
                size={16}
                onClick={openAboutModal}
              />
              <TooltipContent>Error checking for new version</TooltipContent>
            </>
          )}
        </p>
      </TooltipTrigger>
    </Tooltip>
  );
};
