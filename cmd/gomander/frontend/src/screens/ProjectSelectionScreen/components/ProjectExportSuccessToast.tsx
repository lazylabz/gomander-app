import { toast } from "sonner";

import { externalBrowserService } from "@/contracts/service.ts";

export const ProjectExportSuccessToast = ({
  exportFilePath,
  toastId,
}: {
  exportFilePath: string;
  toastId: string;
}) => {
  const handleOpenExportPath = async () => {
    const exportFolderPath = exportFilePath.substring(
      0,
      exportFilePath.lastIndexOf("/"),
    );

    externalBrowserService.browserOpenURL("file://" + exportFolderPath);

    toast.dismiss(toastId);
  };

  return (
    <div className="flex flex-col gap-2 items-start">
      <p>Project exported successfully</p>
      <button className="underline" onClick={handleOpenExportPath}>
        Open containing folder
      </button>
    </div>
  );
};
