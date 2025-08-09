import {
  HoverCard,
  HoverCardContent,
  HoverCardTrigger,
} from "@/components/ui/hover-card.tsx";

export const EnvironmentPathsInfoDialog = () => {
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
          You may need to set this when working with version managers like{" "}
          <code>nvm</code> and <code>pyenv</code> or when your path is modified
          by files as <code>.bashrc</code>, <code>.zshrc</code>, or{" "}
          <code>.profile</code>.
        </p>
        <p>
          <code>e.g. /path/to/.nvm/versions/node/&lt;version&gt;/bin/</code>
        </p>
      </HoverCardContent>
    </HoverCard>
  );
};
