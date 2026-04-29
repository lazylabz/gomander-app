import { TestTubeDiagonal } from "lucide-react";
import { useTranslation } from "react-i18next";
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle,
} from "@/design-system/components/ui/card.tsx";
import { Checkbox } from "@/design-system/components/ui/checkbox.tsx";
import { Label } from "@/design-system/components/ui/label.tsx";
import { useExperimentalFeatsStore } from "@/store/experimentalFeatsStore.ts";
import { changeExperimentalFeat } from "@/useCases/experimentalFeats/changeExperimentalFeat.ts";

export const ExperimentalSettings = () => {
	const { t } = useTranslation();
	const xtermjs = useExperimentalFeatsStore(
		(state) => state.experimentalFeats.xtermjs,
	);

	const handleXtermjsChange = (checked: boolean) => {
		changeExperimentalFeat("xtermjs", checked);
	};

	return (
		<Card>
			<CardHeader>
				<CardTitle className="flex items-center space-x-2">
					<TestTubeDiagonal size={20} />
					<span>{t("settings.experimental.title")}</span>
				</CardTitle>
				<CardDescription>
					{t("settings.experimental.description")}
				</CardDescription>
			</CardHeader>
			<CardContent>
				<div className="flex items-center space-x-2">
					<Checkbox
						id="xtermjs"
						checked={xtermjs}
						onCheckedChange={handleXtermjsChange}
						className="mt-0.5"
					/>
					<Label htmlFor="xtermjs" className="flex flex-col gap-1 items-start">
						<span>{t("settings.experimental.xtermjs.label")}</span>
						<span className="text-muted-foreground text-sm font-normal">
							{t("settings.experimental.xtermjs.description")}
						</span>
					</Label>
				</div>
			</CardContent>
		</Card>
	);
};
