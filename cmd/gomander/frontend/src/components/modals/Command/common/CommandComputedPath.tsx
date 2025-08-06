import { useEffect, useState } from "react";
import { useFormContext } from "react-hook-form";

import type { FormSchemaType } from "@/components/modals/Command/common/formSchema.ts";
import { helpersService } from "@/contracts/service.ts";
import { useProjectStore } from "@/store/projectStore.ts";

export const CommandComputedPath = () => {
  const projectBaseWorkingDirectory =
    useProjectStore((state) => state.projectInfo?.baseWorkingDirectory) || "";
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
      Will run in: <span className="font-medium">{computedPath}</span>
    </p>
  );
};
