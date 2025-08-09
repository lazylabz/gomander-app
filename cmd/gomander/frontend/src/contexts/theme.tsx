import { createContext, useContext, useEffect, useState } from "react";

export const availableThemes = ["dark", "light", "system"] as const;

export type Theme = (typeof availableThemes)[number];

type ThemeProviderProps = {
  children: React.ReactNode;
  defaultTheme?: Theme;
  storageKey?: string;
};

type ThemeProviderState = {
  rawTheme: Theme;
  setRawTheme: (theme: Theme) => void;
  theme: "light" | "dark";
};

const initialState: ThemeProviderState = {
  rawTheme: "system",
  setRawTheme: () => null,
  theme: "light",
};

const ThemeProviderContext = createContext<ThemeProviderState>(initialState);

export function ThemeProvider({
  children,
  defaultTheme = "system",
  storageKey = "vite-ui-theme",
  ...props
}: ThemeProviderProps) {
  const [rawTheme, setRawTheme] = useState<Theme>(
    () => (localStorage.getItem(storageKey) as Theme) || defaultTheme,
  );

  useEffect(() => {
    const root = window.document.documentElement;

    root.classList.remove("light", "dark");

    if (rawTheme === "system") {
      const systemTheme = window.matchMedia("(prefers-color-scheme: dark)")
        .matches
        ? "dark"
        : "light";

      root.classList.add(systemTheme);
      return;
    }

    root.classList.add(rawTheme);
  }, [rawTheme]);

  const theme =
    rawTheme !== "system"
      ? rawTheme
      : window.matchMedia("(prefers-color-scheme: dark)").matches
        ? "dark"
        : "light";

  const value = {
    theme,
    rawTheme,
    setRawTheme: (theme: Theme) => {
      localStorage.setItem(storageKey, theme);
      setRawTheme(theme);
    },
  };

  return (
    <ThemeProviderContext.Provider {...props} value={value}>
      {children}
    </ThemeProviderContext.Provider>
  );
}

export const useTheme = () => {
  const context = useContext(ThemeProviderContext);

  if (context === undefined)
    throw new Error("useTheme must be used within a ThemeProvider");

  return context;
};
