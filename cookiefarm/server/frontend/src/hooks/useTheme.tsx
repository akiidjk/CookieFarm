import {
  createContext,
  useContext,
  useEffect,
  useState,
  type ReactNode,
} from "react";

const themeStorageKey = "cookiefarm-theme-mode";

export type ThemeMode = "light" | "dark";

type ThemeContextValue = {
  mode: ThemeMode;
  setMode: (mode: ThemeMode) => void;
  toggleMode: () => void;
};

const ThemeContext = createContext<ThemeContextValue | null>(null);

function getInitialThemeMode(): ThemeMode {
  return "dark";
}

export function ThemeProvider(props: { children: ReactNode }) {
  const [mode, setMode] = useState<ThemeMode>(getInitialThemeMode);

  useEffect(() => {
    document.documentElement.setAttribute("data-theme", "kumo");
    document.documentElement.setAttribute("data-mode", "dark");
    localStorage.setItem(themeStorageKey, "dark");
  }, [mode]);

  return (
    <ThemeContext.Provider
      value={{
        mode,
        setMode,
        toggleMode: () => {
          setMode((current) => (current === "dark" ? "light" : "dark"));
        },
      }}
    >
      {props.children}
    </ThemeContext.Provider>
  );
}

export function useTheme() {
  const context = useContext(ThemeContext);
  if (!context) {
    throw new Error("useTheme must be used within ThemeProvider");
  }
  return context;
}
