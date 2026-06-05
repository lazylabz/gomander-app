import { Trans, useTranslation } from "react-i18next";
import { useNavigate } from "react-router";

import { Button } from "@/design-system/components/ui/button.tsx";
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from "@/design-system/components/ui/dialog.tsx";
import { ScreenRoutes } from "@/routes.ts";
import { SettingsTab } from "@/screens/SettingsScreen/context/settingsContext.tsx";
import { useMissingEnvironmentPathStore } from "@/store/missingEnvironmentPathStore.ts";

export const MissingEnvironmentPathDialog = () => {
	const { t } = useTranslation();
	const navigate = useNavigate();

	const dialogOpen = useMissingEnvironmentPathStore(
		(state) => state.dialogOpen,
	);
	const setDialogOpen = useMissingEnvironmentPathStore(
		(state) => state.setDialogOpen,
	);

	const handleGoToSettings = () => {
		setDialogOpen(false);
		navigate(ScreenRoutes.Settings, { state: { tab: SettingsTab.User } });
	};

	return (
		<Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
			<DialogContent className="[&_code]:bg-accent [&_code]:rounded [&_code]:px-1">
				<DialogHeader>
					<DialogTitle>{t("modal.missingPath.title")}</DialogTitle>
					<DialogDescription asChild>
						<p className="text-sm text-muted-foreground">
							<Trans
								i18nKey="modal.missingPath.description"
								components={{ code: <code /> }}
							/>
						</p>
					</DialogDescription>
				</DialogHeader>

				<div className="flex flex-col gap-2 text-sm">
					<p className="font-medium">{t("modal.missingPath.howToTitle")}</p>
					<ol className="list-decimal flex flex-col gap-1 pl-5">
						<li>
							<Trans
								i18nKey="modal.missingPath.howToStep1"
								components={{
									whichcmd: (
										<code>
											which {t("modal.missingPath.commandPlaceholder")}
										</code>
									),
									wherecmd: (
										<code>
											where {t("modal.missingPath.commandPlaceholder")}
										</code>
									),
								}}
							/>
						</li>
						<li>{t("modal.missingPath.howToStep2")}</li>
					</ol>
				</div>

				<DialogFooter>
					<Button type="button" onClick={handleGoToSettings}>
						{t("modal.missingPath.goToSettings")}
					</Button>
				</DialogFooter>
			</DialogContent>
		</Dialog>
	);
};
