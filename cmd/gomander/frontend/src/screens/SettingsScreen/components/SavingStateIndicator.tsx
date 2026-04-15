import { CloudAlert, CloudCheck } from "lucide-react";

import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from "@/design-system/components/ui/tooltip.tsx";
import { useSettingsContext } from "@/screens/SettingsScreen/context/settingsContext.tsx";

export const SavingStateIndicator = () => {
  const { hasPendingChanges } = useSettingsContext();

  return (
    <Tooltip>
      <TooltipTrigger>
        {hasPendingChanges ? (
          <CloudAlert className="h-6 w-6 mt-1 text-yellow-400 dark:text-yellow-700" />
        ) : (
          <CloudCheck className="h-6 w-6 mt-0.5 text-green-400 dark:text-green-700" />
        )}
      </TooltipTrigger>
      <TooltipContent>
        {hasPendingChanges ? "Saving..." : "All changes saved"}
      </TooltipContent>
    </Tooltip>
  );
};
