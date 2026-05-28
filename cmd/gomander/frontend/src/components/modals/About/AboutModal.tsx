import {
	BookOpen,
	Check,
	Copy,
	Download,
	ExternalLink,
	Heart,
	Loader2,
	Rocket,
	ShieldAlert,
} from "lucide-react";
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
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from "@/design-system/components/ui/dialog.tsx";
import { cn } from "@/design-system/lib/utils.ts";
import { GithubIcon } from "@/icons/GithubIcon.tsx";

type OS = "darwin" | "linux" | "windows";

const MACOS_QUARANTINE_COMMAND =
	"sudo xattr -d com.apple.quarantine /Applications/gomander.app";

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
	const [copied, setCopied] = useState(false);

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

	const handleReleaseNotesClick = () => {
		if (!newVersion) return;
		const url = `https://github.com/lazylabz/gomander-app/releases/tag/v${newVersion}`;
		externalBrowserService.browserOpenURL(url);
	};

	const handleMacosReadMoreClick = () => {
		const url = `https://github.com/lazylabz/gomander-app#macos-users---important-notice`;
		externalBrowserService.browserOpenURL(url);
	};

	const handleCopyCommand = async () => {
		await navigator.clipboard.writeText(MACOS_QUARANTINE_COMMAND);
		setCopied(true);
		setTimeout(() => setCopied(false), 2000);
	};

	const handleInstallClick = () => {
		setOpen(false);
		setInstallDisclaimersModalOpen(true);
	};

	const handleDisclaimerOpenChange = (nextOpen: boolean) => {
		if (nextOpen) return;
		setInstallDisclaimersModalOpen(false);
		setOpen(true);
	};

	return (
		<>
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
					<div
						className={cn("flex gap-4", newVersion ? "flex-row" : "flex-col")}
					>
						{/* GitHub CTA */}
						<button
							type="button"
							onClick={handleGithubClick}
							className="outline-none focus-visible:ring-[3px] focus-visible:ring-neutral-950/50 dark:focus-visible:ring-neutral-300/50 cursor-pointer w-full flex flex-1 items-center justify-center gap-3 px-4 py-3 bg-card hover:bg-accent border border-border shadow-sm hover:shadow-md rounded-lg transition-colors group"
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
							className="outline-none focus-visible:ring-[3px] focus-visible:ring-neutral-950/50 dark:focus-visible:ring-neutral-300/50 cursor-pointer w-full flex flex-1 items-center justify-center gap-3 px-4 py-3 bg-card hover:bg-accent border border-border shadow-sm hover:shadow-md rounded-lg transition-all group"
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

			<Dialog
				open={installDisclaimerModalOpen}
				onOpenChange={handleDisclaimerOpenChange}
			>
				<DialogContent className="sm:max-w-[628px]">
					<DialogHeader>
						<DialogTitle>{t("aboutModal.installDisclaimer.title")}</DialogTitle>
						<DialogDescription>
							{t("aboutModal.installDisclaimer.description")}
						</DialogDescription>
					</DialogHeader>

					<button
						type="button"
						onClick={handleReleaseNotesClick}
						className="outline-none focus-visible:ring-[3px] focus-visible:ring-neutral-950/50 dark:focus-visible:ring-neutral-300/50 cursor-pointer w-full flex items-center gap-3 px-4 py-3 bg-card hover:bg-accent border border-border shadow-sm hover:shadow-md rounded-lg transition-colors group"
					>
						<BookOpen className="size-5 text-foreground" />
						<div className="text-left">
							<div className="font-medium text-foreground text-sm">
								{t("aboutModal.installDisclaimer.releaseNotesTitle")}
							</div>
							<div className="text-xs text-muted-foreground">
								{t("aboutModal.installDisclaimer.releaseNotesSubtitle", {
									version: newVersion ?? "",
								})}
							</div>
						</div>
						<ExternalLink className="size-4 text-muted-foreground group-hover:text-foreground transition-colors ml-auto" />
					</button>

					{os === "darwin" && (
						<div className="border border-amber-200 dark:border-amber-900 bg-amber-50 dark:bg-amber-950/40 rounded-lg p-4">
							<div className="flex items-start gap-3">
								<ShieldAlert className="size-5 text-amber-600 dark:text-amber-400 mt-0.5 shrink-0" />
								<div className="flex-1 space-y-2 min-w-0">
									<p className="text-sm font-medium text-foreground">
										{t("aboutModal.installDisclaimer.macosTitle")}
									</p>
									<p className="text-xs text-muted-foreground">
										{t("aboutModal.installDisclaimer.macosDescription")}
									</p>
									<div className="flex items-center gap-2 bg-background border border-border rounded-md px-2 py-1.5">
										<code className="flex-1 text-xs font-mono break-all">
											{MACOS_QUARANTINE_COMMAND}
										</code>
										<Button
											type="button"
											variant="ghost"
											size="icon-sm"
											onClick={handleCopyCommand}
											aria-label={t("aboutModal.installDisclaimer.copyCommand")}
										>
											{copied ? (
												<Check className="size-4" />
											) : (
												<Copy className="size-4" />
											)}
										</Button>
									</div>
									<button
										type="button"
										onClick={handleMacosReadMoreClick}
										className="inline-flex items-center gap-1 text-xs text-muted-foreground hover:text-foreground underline underline-offset-2 cursor-pointer rounded-sm outline-none focus-visible:ring-[3px] focus-visible:ring-neutral-950/50 dark:focus-visible:ring-neutral-300/50"
									>
										{t("aboutModal.installDisclaimer.readMore")}
										<ExternalLink className="size-3" />
									</button>
								</div>
							</div>
						</div>
					)}

					<DialogFooter>
						<Button
							variant="outline"
							onClick={() => handleDisclaimerOpenChange(false)}
							disabled={isBusy}
						>
							{t("common.cancel")}
						</Button>
						<Button onClick={installLatestRelease} disabled={isBusy}>
							{buttonInfo.icon}
							{buttonInfo.label}
						</Button>
					</DialogFooter>
				</DialogContent>
			</Dialog>
		</>
	);
};
