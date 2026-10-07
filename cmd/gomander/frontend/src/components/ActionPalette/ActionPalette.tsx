import { defaultFilter } from "cmdk";
import { useState } from "react";
import { useTranslation } from "react-i18next";

import {
	type Action,
	groupActionsBySection,
} from "@/components/ActionPalette/action.ts";
import { useActions } from "@/components/ActionPalette/useActions.ts";
import {
	CommandDialog,
	CommandEmpty,
	CommandGroup,
	CommandInput,
	CommandItem,
	CommandList,
} from "@/design-system/components/ui/command.tsx";
import { useShortcut } from "@/hooks/useShortcut.ts";

// Items are keyed by Action id, which cmdk would otherwise fuzzy-match against the search too.
const matchKeywordsOnly = (
	_value: string,
	search: string,
	keywords?: string[],
) => defaultFilter(keywords?.join(" ") ?? "", search);

export const ActionPalette = () => {
	const { t } = useTranslation();
	const [open, setOpen] = useState(false);
	const actions = useActions();

	useShortcut("Mod-K", () => setOpen((wasOpen) => !wasOpen));

	const handleSelect = (action: Action) => {
		setOpen(false);
		action.perform();
	};

	return (
		<CommandDialog
			open={open}
			onOpenChange={setOpen}
			title={t("actionPalette.title")}
			description={t("actionPalette.description")}
			commandProps={{ filter: matchKeywordsOnly }}
		>
			<CommandInput placeholder={t("actionPalette.placeholder")} />
			<CommandList>
				<CommandEmpty>{t("actionPalette.empty")}</CommandEmpty>
				{groupActionsBySection(actions).map(({ section, actions }) => (
					<CommandGroup
						key={section}
						heading={t(`actionPalette.sections.${section}`)}
					>
						{actions.map((action) => (
							<CommandItem
								key={action.id}
								value={action.id}
								keywords={[action.label, ...action.keywords]}
								onSelect={() => handleSelect(action)}
							>
								<action.icon />
								<span className="truncate">{action.label}</span>
								{action.kind && (
									<span className="ml-auto" title={action.kind.label}>
										<action.kind.icon aria-hidden />
										<span className="sr-only">{action.kind.label}</span>
									</span>
								)}
							</CommandItem>
						))}
					</CommandGroup>
				))}
			</CommandList>
		</CommandDialog>
	);
};
