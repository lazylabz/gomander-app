import { Play } from "lucide-react";
import { describe, expect, it } from "vitest";

import {
	type Action,
	groupActionsBySection,
} from "@/components/ActionPalette/action.ts";

const anAction = (id: string): Action => ({
	id,
	section: "commands",
	label: id,
	keywords: [],
	icon: Play,
	perform: () => {},
});

describe("groupActionsBySection", () => {
	const sut = groupActionsBySection;

	it("Should keep the actions of a section in the order they came in", () => {
		// Arrange
		const actions = [anAction("b"), anAction("a")];

		// Act
		const sections = sut(actions);

		// Assert
		expect(sections).toEqual([{ section: "commands", actions }]);
	});

	it("Should leave out a section with no actions", () => {
		// Act
		const sections = sut([]);

		// Assert
		expect(sections).toEqual([]);
	});
});
