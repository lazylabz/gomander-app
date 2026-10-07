import type { LucideIcon } from "lucide-react";

export const actionSections = ["commands"] as const;

export type ActionSection = (typeof actionSections)[number];

export type ActionKind = {
	icon: LucideIcon;
	label: string;
};

export type Action = {
	id: string;
	section: ActionSection;
	label: string;
	keywords: string[];
	icon: LucideIcon;
	kind?: ActionKind;
	perform: () => unknown;
};

export const groupActionsBySection = (
	actions: Action[],
): { section: ActionSection; actions: Action[] }[] =>
	actionSections
		.map((section) => ({
			section,
			actions: actions.filter((action) => action.section === section),
		}))
		.filter(({ actions }) => actions.length > 0);
