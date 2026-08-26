import { useStore } from "zustand/react";
import { createStore } from "zustand/vanilla";

export const ALL_COMMANDS_SECTION_ID = "all-commands";

const STORAGE_KEY = "sidebar-sections-open";

type OpenSections = Record<string, boolean>;

const readStoredSections = (): OpenSections => {
	try {
		const stored = window.localStorage.getItem(STORAGE_KEY);
		return stored ? (JSON.parse(stored) as OpenSections) : {};
	} catch {
		return {};
	}
};

const storeSections = (openSections: OpenSections) => {
	try {
		window.localStorage.setItem(STORAGE_KEY, JSON.stringify(openSections));
	} catch {
		// Ignore write errors
	}
};

type SidebarSectionsStore = {
	openSections: OpenSections;
};

const sidebarSectionsStore = createStore<SidebarSectionsStore>()(() => ({
	openSections: readStoredSections(),
}));

const persist = (openSections: OpenSections) => {
	storeSections(openSections);
	sidebarSectionsStore.setState({ openSections });
};

export const isSidebarSectionOpen = (sectionId: string): boolean =>
	sidebarSectionsStore.getState().openSections[sectionId] ?? false;

export const useIsSidebarSectionOpen = (sectionId: string): boolean =>
	useStore(
		sidebarSectionsStore,
		(state) => state.openSections[sectionId] ?? false,
	);

export const setSidebarSectionOpen = (
	sectionId: string,
	open: boolean,
): void => {
	persist({
		...sidebarSectionsStore.getState().openSections,
		[sectionId]: open,
	});
};

export const forgetSidebarSection = (sectionId: string): void => {
	const { [sectionId]: _forgotten, ...remaining } =
		sidebarSectionsStore.getState().openSections;

	persist(remaining);
};
