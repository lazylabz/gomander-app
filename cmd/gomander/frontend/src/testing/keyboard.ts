import type { UserEvent } from "@testing-library/user-event";
import { expect } from "vitest";

// Reaches the control the way a keyboard user does, so a control left out of
// the tab order fails here instead of being focused programmatically.
export const tabToAndPressEnter = async (
	user: UserEvent,
	control: HTMLElement,
) => {
	const focusable = document.querySelectorAll(
		"button, [href], input, select, textarea, [tabindex]",
	).length;
	for (let i = 0; i < focusable && document.activeElement !== control; i++) {
		await user.tab();
	}
	expect(control).toHaveFocus();
	await user.keyboard("{Enter}");
};
