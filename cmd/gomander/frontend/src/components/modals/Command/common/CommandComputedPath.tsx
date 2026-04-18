import { useEffect, useState } from "react";
import { useFormContext } from "react-hook-form";
import { useTranslation } from "react-i18next";

import type { FormSchemaType } from "@/components/modals/Command/common/formSchema.ts";
import { helpersService } from "@/contracts/service.ts";
import { useProjectStore } from "@/store/projectStore.ts";

export const CommandComputedPath = () => {
  const { t } = useTranslation();
  const projectBaseWorkingDirectory =
    useProjectStore((state) => state.projectInfo?.workingDirectory) || "";
  const [computedPath, setComputedPath] = useState(
    projectBaseWorkingDirectory || "",
  );

  const form = useFormContext<FormSchemaType>();

  const workingDirectoryWatcher = form.watch("workingDirectory");
  useEffect(() => {
    helpersService
      .getComputedPath(projectBaseWorkingDirectory, workingDirectoryWatcher)
      .then((cp) => {
        setComputedPath(cp);
      });
  }, [projectBaseWorkingDirectory, workingDirectoryWatcher]);

  return (
    <p className="-mt-4 text-xs text-muted-foreground">
      {t('commandForm.computedPath', { path: computedPath })}
    </p>
  );
};
