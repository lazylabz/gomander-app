import { Trans, useTranslation } from "react-i18next";

import {
  HoverCard,
  HoverCardContent,
  HoverCardTrigger,
} from "@/design-system/components/ui/hover-card.tsx";

export const EnvironmentPathsInfoDialog = () => {
  const { t } = useTranslation();

  return (
    <HoverCard openDelay={100}>
      <HoverCardTrigger className="text-xs self-center cursor-help text-muted-foreground hover:text-foreground border rounded-full size-4 flex items-center justify-center">
        ?
      </HoverCardTrigger>
      <HoverCardContent
        sideOffset={10}
        className="w-100 text-sm flex flex-col gap-2 [&>p>code]:bg-accent"
      >
        <p>
          <Trans
            i18nKey="userSettingsForm.envPathsHelpBody"
            components={{ code: <code /> }}
          />
        </p>
        <p className="text-xs">
          {t('userSettingsForm.envPathsHelpExample')}
        </p>
      </HoverCardContent>
    </HoverCard>
  );
};
