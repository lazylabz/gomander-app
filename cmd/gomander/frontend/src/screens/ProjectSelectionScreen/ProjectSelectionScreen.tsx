import { ChevronDownIcon, Import, Plus } from "lucide-react";
import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate } from "react-router";

import { CreateProjectModal } from "@/components/modals/Project/CreateProjectModal.tsx";
import { DeleteProjectModal } from "@/components/modals/Project/DeleteProjectModal.tsx";
import { ImportProjectModal } from "@/components/modals/Project/ImportProjectModal.tsx";
import type { ProjectBlueprint } from "@/contracts/types.ts";
import { Button } from "@/design-system/components/ui/button.tsx";
import { ButtonGroup } from "@/design-system/components/ui/button-group.tsx";
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuTrigger,
} from "@/design-system/components/ui/dropdown-menu.tsx";
import { fetchAvailableProjects } from "@/queries/fetchAvailableProjects.ts";
import { ScreenRoutes } from "@/routes.ts";
import { ProjectCard } from "@/screens/ProjectSelectionScreen/components/ProjectCard.tsx";
import { useProjectStore } from "@/store/projectStore.ts";
import { deleteProject } from "@/useCases/project/deleteProject.ts";
import {
	type ImportSource,
	pickProjectToImport,
} from "@/useCases/project/pickProjectToImport.ts";

export const ProjectSelectionScreen = () => {
	const { t } = useTranslation();
	const [projectIdBeingDeleted, setProjectIdBeingDeleted] = useState<
		string | null
	>(null);
	const project = useProjectStore((state) => state.projectInfo);

	const navigate = useNavigate();

	const [projectBeingImported, setProjectBeingImported] =
		useState<ProjectBlueprint | null>(null);

	const [createProjectModalOpen, setCreateProjectModalOpen] = useState(false);

	const availableProjects = useProjectStore((state) => state.availableProjects);

	const openCreateProjectModal = () => {
		setCreateProjectModalOpen(true);
	};

	const handleDeleteProject = (projectId: string) => async () => {
		setProjectIdBeingDeleted(projectId);
	};

	const confirmDeleteProject = async () => {
		if (!projectIdBeingDeleted) {
			return;
		}

		if (await deleteProject(projectIdBeingDeleted)) {
			setProjectIdBeingDeleted(null);
		}
	};

	const cancelDeleteProject = () => {
		setProjectIdBeingDeleted(null);
	};

	const handleImportProject = (source: ImportSource) => async () => {
		const projectToImport = await pickProjectToImport(source);
		if (projectToImport) {
			setProjectBeingImported(projectToImport);
		}
	};

	useEffect(() => {
		fetchAvailableProjects();
	}, []);

	useEffect(() => {
		if (project) {
			navigate(ScreenRoutes.Logs);
		}
	}, [navigate, project]);

	const hasProjects = availableProjects.length > 0;

	return (
		<>
			<CreateProjectModal
				open={createProjectModalOpen}
				setOpen={setCreateProjectModalOpen}
			/>
			<ImportProjectModal
				open={!!projectBeingImported}
				onClose={() => setProjectBeingImported(null)}
				project={projectBeingImported}
			/>
			<DeleteProjectModal
				open={!!projectIdBeingDeleted}
				onConfirm={confirmDeleteProject}
				onClose={cancelDeleteProject}
			/>
			<div className="w-full h-full flex flex-col items-center justify-center gap-10">
				<h1 className="text-3xl">
					{hasProjects
						? t("projectSelection.openTitle")
						: t("projectSelection.welcome")}
				</h1>
				{hasProjects && (
					<div className="flex flex-col items-center justify-center gap-2">
						{availableProjects.map((p) => (
							<ProjectCard
								key={p.id}
								project={p}
								handleDeleteProject={handleDeleteProject(p.id)}
							/>
						))}
					</div>
				)}
				{!hasProjects && <p>{t("projectSelection.emptyState")}</p>}
				<div className="flex flex-row items-center gap-2 justify-center">
					<Button onClick={openCreateProjectModal} variant="ghost">
						<Plus /> {t("projectSelection.createButton")}
					</Button>
					<ButtonGroup>
						<Button
							variant="ghost"
							onClick={handleImportProject("exportedProject")}
						>
							<Import />
							{t("projectSelection.importButton")}
						</Button>
						<DropdownMenu>
							<DropdownMenuTrigger asChild>
								<Button
									variant="ghost"
									size="icon"
									aria-label={t("projectSelection.moreOptions")}
								>
									<ChevronDownIcon />
								</Button>
							</DropdownMenuTrigger>
							<DropdownMenuContent align="end" className="w-52">
								<DropdownMenuItem onClick={handleImportProject("packageJson")}>
									{t("projectSelection.importPackageJson")}
								</DropdownMenuItem>
							</DropdownMenuContent>
						</DropdownMenu>
					</ButtonGroup>
				</div>
			</div>
		</>
	);
};
