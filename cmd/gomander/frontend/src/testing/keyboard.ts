import type { UserEvent } from "@testing-library/user-event";
import { expect } from "vitest";

export const pressEnterOn = async (user: UserEvent, control: HTMLElement) => {
	control.focus();
	expect(control).toHaveFocus();
	await user.keyboard("{Enter}");
};
