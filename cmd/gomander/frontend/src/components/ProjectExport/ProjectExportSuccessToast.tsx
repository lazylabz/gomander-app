import { useTranslation } from "react-i18next";
import { toast } from "sonner";

import { helpersService } from "@/contracts/service.ts";

export const ProjectExportSuccessToast = ({
	exportFilePath,
	toastId,
}: {
	exportFilePath: string;
	toastId: string;
}) => {
	const { t } = useTranslation();

	const handleOpenExportPath = async () => {
		await helpersService.openFileFolder(exportFilePath);

		toast.dismiss(toastId);
	};

	return (
		<div className="flex flex-col gap-2 items-start">
			<p>{t("toast.project.exportSuccess")}</p>
			<button
				type="button"
				className="underline"
				onClick={handleOpenExportPath}
			>
				{t("toast.project.openFolderAction")}
			</button>
		</div>
	);
};
