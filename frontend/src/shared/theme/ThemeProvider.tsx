import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useRef,
  useState,
  type ReactNode,
} from "react";

import {
  useGetPreferencesQuery,
  useGetPublicThemeQuery,
  useResetPreferencesMutation,
  useSavePreferencesMutation,
} from "@/shared/api/meApi";
import { useAppSelector } from "@/app/store";

import { applyTheme, resolveMode } from "./applyTheme";
import { DEFAULT_THEME, normalizeTheme, type ThemeSettings } from "./types";

const LOCAL_KEY = "platforma.theme";

function loadLocal(): ThemeSettings | null {
  try {
    const raw = localStorage.getItem(LOCAL_KEY);
    return raw ? normalizeTheme(JSON.parse(raw)) : null;
  } catch {
    return null;
  }
}

type ThemeContextValue = {
  settings: ThemeSettings;
  mode: "dark" | "light";
  /** Частичное изменение настроек — сразу применяется и сохраняется. */
  update: (patch: Partial<ThemeSettings>) => void;
  /** Полная замена настроек (например, применение пресета). */
  replace: (next: ThemeSettings) => void;
  /** Вернуться к оформлению, заданному администратором. */
  resetToPlatform: () => void;
  /** Быстрое переключение светлая/тёмная. */
  toggleMode: () => void;
  platformTheme: ThemeSettings | null;
};

const ThemeContext = createContext<ThemeContextValue | null>(null);

export function ThemeProvider({ children }: { children: ReactNode }) {
  const user = useAppSelector((state) => state.auth.user);

  const [settings, setSettings] = useState<ThemeSettings>(() => loadLocal() ?? DEFAULT_THEME);
  const [mode, setMode] = useState<"dark" | "light">(() => resolveMode(settings));
  const hasLocal = useRef(loadLocal() !== null);

  const { data: publicTheme } = useGetPublicThemeQuery();
  const { data: preferences } = useGetPreferencesQuery(undefined, { skip: !user });
  const [savePreferences] = useSavePreferencesMutation();
  const [resetPreferences] = useResetPreferencesMutation();

  const platformTheme = useMemo(
    () => (publicTheme?.settings ? normalizeTheme(publicTheme.settings) : null),
    [publicTheme],
  );

  // Оформление платформы — база для тех, кто ничего не настраивал под себя.
  useEffect(() => {
    if (platformTheme && !hasLocal.current) {
      setSettings(platformTheme);
    }
  }, [platformTheme]);

  // Личные настройки пользователя всегда важнее общих.
  useEffect(() => {
    if (preferences?.theme) {
      const personal = normalizeTheme(preferences.theme);
      hasLocal.current = true;
      setSettings(personal);
      localStorage.setItem(LOCAL_KEY, JSON.stringify(personal));
    }
  }, [preferences]);

  useEffect(() => {
    setMode(applyTheme(settings));
  }, [settings]);

  // Следим за системной темой, когда выбран режим «как в системе».
  useEffect(() => {
    if (settings.mode !== "system") return;
    const media = window.matchMedia("(prefers-color-scheme: dark)");
    const onChange = () => setMode(applyTheme(settings));
    media.addEventListener("change", onChange);
    return () => media.removeEventListener("change", onChange);
  }, [settings]);

  const persist = useCallback(
    (next: ThemeSettings) => {
      hasLocal.current = true;
      localStorage.setItem(LOCAL_KEY, JSON.stringify(next));
      if (user) void savePreferences(next);
    },
    [savePreferences, user],
  );

  const replace = useCallback(
    (next: ThemeSettings) => {
      const normalized = normalizeTheme(next);
      setSettings(normalized);
      persist(normalized);
    },
    [persist],
  );

  const update = useCallback(
    (patch: Partial<ThemeSettings>) => {
      setSettings((current) => {
        const next = normalizeTheme({ ...current, ...patch });
        persist(next);
        return next;
      });
    },
    [persist],
  );

  const resetToPlatform = useCallback(() => {
    const next = platformTheme ?? DEFAULT_THEME;
    hasLocal.current = false;
    localStorage.removeItem(LOCAL_KEY);
    setSettings(next);
    if (user) void resetPreferences();
  }, [platformTheme, resetPreferences, user]);

  const toggleMode = useCallback(() => {
    update({ mode: resolveMode(settings) === "dark" ? "light" : "dark" });
  }, [settings, update]);

  const value = useMemo<ThemeContextValue>(
    () => ({ settings, mode, update, replace, resetToPlatform, toggleMode, platformTheme }),
    [settings, mode, update, replace, resetToPlatform, toggleMode, platformTheme],
  );

  return <ThemeContext.Provider value={value}>{children}</ThemeContext.Provider>;
}

export function useTheme(): ThemeContextValue {
  const ctx = useContext(ThemeContext);
  if (!ctx) throw new Error("useTheme должен вызываться внутри ThemeProvider");
  return ctx;
}
