import { RotateCw, TriangleAlert } from "lucide-react";
import { useTranslation } from "react-i18next";

import { Button } from "@/design-system/components/ui/button.tsx";

export const LoadFailure = ({
	onRetry,
	isRetrying,
}: {
	onRetry: () => void;
	isRetrying: boolean;
}) => {
	const { t } = useTranslation();

	return (
		<div
			role="alert"
			className="flex flex-col items-center justify-center gap-4 text-center"
		>
			<TriangleAlert className="text-red-500" />
			<p>{t("common.loadFailed")}</p>
			<Button variant="outline" onClick={onRetry} disabled={isRetrying}>
				<RotateCw /> {t("common.retry")}
			</Button>
		</div>
	);
};
