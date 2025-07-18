import { createContext, useContext, useState } from "react";

import { CreateCommandModal } from "@/components/modals/CreateCommandModal.tsx";

export enum Modals {
  CREATE,
}

export type ModalsContextType = {
  setOpenModal: (modal: Modals) => (value: boolean) => void;
};

const modalsContext = createContext<ModalsContextType>({} as ModalsContextType);

export const ModalsContextProvider = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  const [openModals, setOpenModals] = useState<Record<Modals, boolean>>({
    [Modals.CREATE]: false,
  });

  const setOpenModal = (modal: Modals) => (value: boolean) => {
    setOpenModals((prev) => ({ ...prev, [modal]: value }));
  };

  const value: ModalsContextType = {
    setOpenModal,
  };

  return (
    <modalsContext.Provider value={value}>
      <CreateCommandModal
        open={openModals[Modals.CREATE]}
        setOpen={setOpenModal(Modals.CREATE)}
      />
      {children}
    </modalsContext.Provider>
  );
};

export const useModalsContext = () => {
  return useContext(modalsContext);
};
