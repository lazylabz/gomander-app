import { translationsService } from "@/contracts/service";
import { fetchCommandGroups } from "@/queries/fetchCommandGroups.ts";
import { fetchCommands } from "@/queries/fetchCommands.ts";
import { fetchProject } from "@/queries/fetchProject.ts";

export const loadAllProjectData = async () => {
  await Promise.all([fetchCommands(), fetchCommandGroups()]);
  await fetchProject();

  // TODO: Delete this (testing i18n PoC)
  const languages = await translationsService.getSupportedLanguages();
  console.log("languages: ", languages);

  const translations = await translationsService.getTranslation('es');
  console.log("translations: ", translations)
  console.log("sidebar title: ", translations["sidebar.title"])
};
