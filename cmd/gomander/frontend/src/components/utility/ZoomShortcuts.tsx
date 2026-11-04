import { Fragment, useEffect } from "react";

import { helpersService } from "@/contracts/service.ts";
import { zoomShortcutOf } from "@/helpers/zoomShortcut.ts";
import { zoom } from "@/useCases/zoom/zoom.ts";

// macOS zooms natively, before the page sees the press; elsewhere Wails cannot
// name the keys, so the page catches them.
export const ZoomShortcuts = () => {
	useEffect(() => {
		let cancelled = false;

		const handleKeyDown = (event: KeyboardEvent) => {
			const direction = zoomShortcutOf(event);
			if (!direction) {
				return;
			}
			event.preventDefault();
			void zoom(direction);
		};

		const listenOutsideMacOS = async () => {
			const os = await helpersService.getOs();
			if (cancelled || os === "darwin") {
				return;
			}
			window.addEventListener("keydown", handleKeyDown);
		};

		void listenOutsideMacOS();

		return () => {
			cancelled = true;
			window.removeEventListener("keydown", handleKeyDown);
		};
	}, []);

	return <Fragment />;
};
