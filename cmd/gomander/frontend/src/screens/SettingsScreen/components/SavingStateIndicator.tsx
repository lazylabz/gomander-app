import { CloudAlert, CloudCheck } from "lucide-react";
import { useTranslation } from "react-i18next";

import {
	Tooltip,
	TooltipContent,
	TooltipTrigger,
} from "@/design-system/components/ui/tooltip.tsx";
import { useSettingsContext } from "@/screens/SettingsScreen/context/settingsContext.tsx";

export const SavingStateIndicator = () => {
	const { t } = useTranslation();
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
				{hasPendingChanges
					? t("settings.saving.inProgress")
					: t("settings.saving.done")}
			</TooltipContent>
		</Tooltip>
	);
};
