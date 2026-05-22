import { Download, ExternalLink, Heart, Loader2, Rocket } from "lucide-react";
import { useCallback, useEffect, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { useVersionContext } from "@/contexts/version.tsx";
import { externalBrowserService, helpersService } from "@/contracts/service.ts";
import {
	Avatar,
	AvatarFallback,
	AvatarImage,
} from "@/design-system/components/ui/avatar.tsx";
import { Button } from "@/design-system/components/ui/button.tsx";
import {
	Dialog,
	DialogContent,
} from "@/design-system/components/ui/dialog.tsx";
import { cn } from "@/design-system/lib/utils.ts";
import { GithubIcon } from "@/icons/GithubIcon.tsx";

type OS = "darwin" | "linux" | "widows";

export const AboutModal = ({
	open,
	setOpen,
}: {
	open: boolean;
	setOpen: (open: boolean) => void;
}) => {
	const { t } = useTranslation();

	const {
		currentVersion,
		newVersion,
		updateStatus,
		downloadLatestRelease,
		installLatestRelease,
	} = useVersionContext();

	const [installDisclaimerModalOpen, setInstallDisclaimersModalOpen] =
		useState(false);
	const [os, setOs] = useState<OS | null>(null);

	const fetchOs = useCallback(async () => {
		const res = await helpersService.getOs();
		setOs(res as OS);
	}, []);

	useEffect(() => {
		fetchOs();
	}, [fetchOs]);

	const isDownloaded =
		updateStatus === "downloaded" || updateStatus === "installing";
	const isBusy =
		updateStatus === "downloading" || updateStatus === "installing";

	const buttonInfo = useMemo(() => {
		switch (updateStatus) {
			case "downloading":
				return {
					label: t("aboutModal.downloading"),
					icon: <Loader2 className="size-4 animate-spin" />,
				};
			case "installing":
				return {
					label: t("aboutModal.installing"),
					icon: <Loader2 className="size-4 animate-spin" />,
				};
			case "downloaded":
				return {
					label: t("aboutModal.installUpdate"),
					icon: <Rocket className="size-4" />,
				};
			default:
				return {
					label: t("aboutModal.downloadUpdate"),
					icon: <Download className="size-4" />,
				};
		}
	}, [t, updateStatus]);

	const handleGithubClick = () => {
		const url = `https://github.com/lazylabz/gomander-app`;
		externalBrowserService.browserOpenURL(url);
	};

	const handleTeamClick = () => {
		const url = `https://lazylabz.github.io/`;
		externalBrowserService.browserOpenURL(url);
	};

	const handleInstallClick = async () => {
		setOpen(false);
		setInstallDisclaimersModalOpen(true);
	};

	const handleCloseDisclaimerModal = () => {
		setInstallDisclaimersModalOpen(false);
		setOpen(true);
	};

	if (installDisclaimerModalOpen) {
		return (
			<Dialog
				open={installDisclaimerModalOpen}
				onOpenChange={handleCloseDisclaimerModal}
			>
				<DialogContent className="sm:max-w-[628px]">
					<p>test</p>
					<p>{os === "darwin" && <p>MACOS DETECTED</p>}</p>
					<Button
						onClick={installLatestRelease}
						disabled={isBusy}
						variant="outline"
						className="cursor-pointer"
					>
						{buttonInfo.icon}
						{buttonInfo.label}
					</Button>
				</DialogContent>
			</Dialog>
		);
	}

	return (
		<Dialog open={open} onOpenChange={setOpen}>
			<DialogContent className="sm:max-w-[628px]">
				{/* Logo and Title */}
				<div className="text-center">
					<Avatar className="size-16 mx-auto mb-4 rounded-xl">
						<AvatarImage src="/app-logo.png" />
						<AvatarFallback className="rounded-xl text-xl font-extralight bg-card-foreground text-card">
							G.
						</AvatarFallback>
					</Avatar>
					<h3 className="text-xl font-bold text-foreground">
						Gomander
						<span className="ml-2 font-normal text-sm text-muted-foreground">
							{t("aboutModal.version", { version: currentVersion })}
						</span>
					</h3>
				</div>
				{/* Update Notice */}
				{newVersion && (
					<div className="bg-sky-50 dark:bg-sky-950/40 border border-b-0 border-sky-200 dark:border-sky-950 shadow-sm shadow-sky-200 dark:shadow-sky-950 rounded-lg p-4">
						<div className="flex items-center">
							<div className="flex-1">
								<p className="text-sm font-medium text-foreground mb-1">
									{t("aboutModal.newVersion", { version: newVersion })}
								</p>
								<p className="text-xs text-muted-foreground">
									{t("aboutModal.newVersionSubtitle")}
								</p>
							</div>
							<Button
								onClick={
									isDownloaded ? handleInstallClick : downloadLatestRelease
								}
								disabled={isBusy}
								variant="outline"
								className="cursor-pointer"
							>
								{buttonInfo.icon}
								{buttonInfo.label}
							</Button>
						</div>
					</div>
				)}
				{/* Description */}
				<p className="text-sm text-muted-foreground leading-relaxed">
					{t("aboutModal.description")}
				</p>
				{/* CTAs */}
				<div className={cn("flex gap-4", newVersion ? "flex-row" : "flex-col")}>
					{/* GitHub CTA */}
					<button
						type="button"
						onClick={handleGithubClick}
						className="focus-visible:outline-none cursor-pointer w-full flex flex-1 items-center justify-center gap-3 px-4 py-3 bg-card hover:bg-accent border border-border shadow-sm hover:shadow-md rounded-lg transition-colors group"
					>
						<GithubIcon className="size-5 text-foreground" />
						<div className="text-left">
							<div className="font-medium text-foreground text-sm">
								{t("aboutModal.feedbackTitle")}
							</div>
							<div className="text-xs text-muted-foreground">
								{t("aboutModal.feedbackSubtitle")}
							</div>
						</div>
						<ExternalLink className="size-4 text-muted-foreground group-hover:text-foreground transition-colors ml-auto" />
					</button>

					{/* LazyLabz CTA */}
					<button
						type="button"
						onClick={handleTeamClick}
						className="focus-visible:outline-none cursor-pointer w-full flex flex-1 items-center justify-center gap-3 px-4 py-3 bg-card hover:bg-accent border border-border shadow-sm hover:shadow-md rounded-lg transition-all group"
					>
						<Heart className="size-5 text-foreground" />
						<div className="text-left">
							<div className="font-medium text-foreground text-sm">
								{t("aboutModal.teamTitle")}
							</div>
							<div className="text-xs text-muted-foreground">
								{t("aboutModal.teamSubtitle")}
							</div>
						</div>
						<ExternalLink className="size-4 text-muted-foreground group-hover:text-foreground transition-colors ml-auto" />
					</button>
				</div>
			</DialogContent>
		</Dialog>
	);
};
