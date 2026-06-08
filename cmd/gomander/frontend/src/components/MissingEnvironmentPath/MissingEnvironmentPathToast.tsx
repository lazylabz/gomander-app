import { useTranslation } from "react-i18next";
import { toast } from "sonner";

import { missingEnvironmentPathStore } from "@/store/missingEnvironmentPathStore.ts";

export const MISSING_ENV_PATH_TOAST_ID = "missing-environment-path";

export const MissingEnvironmentPathToast = () => {
	const { t } = useTranslation();

	const handleLearnMore = () => {
		missingEnvironmentPathStore.getState().setDialogOpen(true);
		toast.dismiss(MISSING_ENV_PATH_TOAST_ID);
	};

	return (
		<div className="flex flex-col gap-2 items-start">
			<p>{t("toast.missingPath.message")}</p>
			<button type="button" className="underline" onClick={handleLearnMore}>
				{t("toast.missingPath.action")}
			</button>
		</div>
	);
};
