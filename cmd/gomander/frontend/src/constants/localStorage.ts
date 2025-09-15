export const COMMAND_GROUP_SECTION_OPEN = `command-group-section-open-`;
export const ALL_COMMANDS_SECTION_OPEN = `all-commands-section-open`;

export const getCommandGroupSectionOpenLocalStorageKey = (groupId: string) =>
  `${COMMAND_GROUP_SECTION_OPEN}${groupId}`;
