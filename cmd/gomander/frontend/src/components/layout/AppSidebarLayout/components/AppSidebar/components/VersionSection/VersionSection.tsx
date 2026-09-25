import { Info } from "lucide-react";
import { useTranslation } from "react-i18next";

import {
	Tooltip,
	TooltipContent,
	TooltipTrigger,
} from "@/design-system/components/ui/tooltip.tsx";
import { useReleaseStore } from "@/store/releaseStore.ts";

export const VersionSection = ({
	openAboutModal,
}: {
	openAboutModal: () => void;
}) => {
	const { t } = useTranslation();
	const currentRelease = useReleaseStore((state) => state.currentRelease);
	const newRelease = useReleaseStore((state) => state.newRelease);
	const checkFailed = useReleaseStore((state) => state.checkFailed);

	return (
		<Tooltip>
			<TooltipTrigger className="cursor-pointer" onClick={openAboutModal}>
				<p className="text-sm text-muted-foreground flex items-center gap-2">
					{currentRelease
						? t("sidebar.version.current", { version: currentRelease })
						: "..."}
					{newRelease && (
						<>
							<Info
								className="text-orange-400 dark:text-yellow-400 cursor-pointer"
								size={16}
								onClick={openAboutModal}
							/>
							<TooltipContent>
								<span className="font-semibold">
									{t("sidebar.version.newAvailable", { version: newRelease })}
								</span>
							</TooltipContent>
						</>
					)}
					{currentRelease && !newRelease && !checkFailed && (
						<>
							<Info
								className="text-green-600 dark:text-green-200 cursor-pointer"
								size={16}
								onClick={openAboutModal}
							/>
							<TooltipContent>{t("sidebar.version.latest")}</TooltipContent>
						</>
					)}
					{checkFailed && (
						<>
							<Info
								className="text-red-600 dark:text-red-400 cursor-pointer"
								size={16}
								onClick={openAboutModal}
							/>
							<TooltipContent>{t("sidebar.version.checkError")}</TooltipContent>
						</>
					)}
				</p>
			</TooltipTrigger>
		</Tooltip>
	);
};
