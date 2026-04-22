import { CirclePlus, Group, Terminal } from "lucide-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";

import { CreateCommandModal } from "@/components/modals/Command/CreateCommandModal.tsx";
import { CreateCommandGroupModal } from "@/components/modals/CommandGroup/CreateCommandGroupModal.tsx";
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuLabel,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
} from "@/design-system/components/ui/dropdown-menu.tsx";
import {
	Tooltip,
	TooltipContent,
	TooltipTrigger,
} from "@/design-system/components/ui/tooltip.tsx";
import { useCommandStore } from "@/store/commandStore.ts";

export const CreateMenu = () => {
	const { t } = useTranslation();
	const commands = useCommandStore((state) => state.commands);
	const [createCommandModalOpen, setCreateCommandModalOpen] = useState(false);
	const [createCommandGroupModalOpen, setCreateCommandGroupModalOpen] =
		useState(false);

	const openCreateCommandModal = () => setCreateCommandModalOpen(true);
	const openCreateCommandGroupModal = () =>
		setCreateCommandGroupModalOpen(true);

	const hasCommands = commands.length > 0;

	return (
		<>
			<CreateCommandModal
				open={createCommandModalOpen}
				setOpen={setCreateCommandModalOpen}
			/>
			<CreateCommandGroupModal
				open={createCommandGroupModalOpen}
				setOpen={setCreateCommandGroupModalOpen}
			/>
			<DropdownMenu>
				<DropdownMenuTrigger className="outline-0">
					<CirclePlus
						className="text-muted-foreground cursor-pointer hover:text-primary"
						size={18}
					/>
				</DropdownMenuTrigger>
				<DropdownMenuContent>
					<DropdownMenuLabel>{t("sidebar.create.title")}</DropdownMenuLabel>
					<DropdownMenuSeparator />
					<DropdownMenuItem
						onClick={openCreateCommandModal}
						className="flex flex-row items-center justify-start"
					>
						<Terminal />
						{t("sidebar.create.command")}
					</DropdownMenuItem>
					<Tooltip delayDuration={500}>
						<TooltipTrigger>
							<DropdownMenuItem
								onClick={openCreateCommandGroupModal}
								className="flex flex-row items-center justify-start"
								disabled={!hasCommands}
							>
								<Group />
								{t("sidebar.create.commandGroup")}
							</DropdownMenuItem>
						</TooltipTrigger>
						{!hasCommands && (
							<TooltipContent side="bottom">
								{t("sidebar.create.noCommandsTooltip")}
							</TooltipContent>
						)}
					</Tooltip>
				</DropdownMenuContent>
			</DropdownMenu>
		</>
	);
};
