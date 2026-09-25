import { render } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import type { ReactElement } from "react";
import { type Location, MemoryRouter, useLocation } from "react-router";

import { SidebarProvider } from "@/design-system/components/ui/sidebar.tsx";
import { ScreenRoutes } from "@/routes.ts";

type Options = {
	route?: ScreenRoutes;
	routeState?: unknown;
};

export const renderWithProviders = (
	ui: ReactElement,
	{ route = ScreenRoutes.ProjectSelection, routeState }: Options = {},
) => {
	let currentLocation: Location | undefined;

	const LocationProbe = () => {
		currentLocation = useLocation();
		return null;
	};

	const user = userEvent.setup();
	const result = render(
		<MemoryRouter initialEntries={[{ pathname: route, state: routeState }]}>
			<SidebarProvider>{ui}</SidebarProvider>
			<LocationProbe />
		</MemoryRouter>,
	);

	return {
		...result,
		user,
		location: (): Location => {
			if (!currentLocation) {
				throw new Error("The router has not rendered yet");
			}
			return currentLocation;
		},
	};
};
