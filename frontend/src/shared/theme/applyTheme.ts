import { darken, lighten, mix, readableOn, rgba } from "./color";
import type { ThemeSettings } from "./types";

// Какая тема реально включена сейчас (system разворачивается в dark/light).
export function resolveMode(settings: ThemeSettings): "dark" | "light" {
  if (settings.mode !== "system") return settings.mode;
  if (typeof window === "undefined") return "dark";
  return window.matchMedia?.("(prefers-color-scheme: dark)").matches ? "dark" : "light";
}

// Раскладывает настройки в CSS-переменные на <html>. Всё оформление
// в приложении читает только эти переменные — поэтому смена темы мгновенная.
export function applyTheme(settings: ThemeSettings): "dark" | "light" {
  const root = document.documentElement;
  const mode = resolveMode(settings);
  const dark = mode === "dark";

  const base = dark ? settings.darkBase : settings.lightBase;
  const accent = settings.accent;
  const accent2 = settings.accent2;

  // Поверхности: в тёмной теме подсветляем базу, в светлой — затемняем.
  const step = (amount: number) => (dark ? lighten(base, amount) : darken(base, amount));

  const text = dark ? lighten(base, 0.92) : darken(base, 0.86);
  const textMuted = dark
    ? mix(text, base, 0.42 / settings.contrast)
    : mix(text, base, 0.4 / settings.contrast);

  const vars: Record<string, string> = {
    "--bg": base,
    "--bg-deep": dark ? darken(base, 0.35) : darken(base, 0.03),
    "--surface": rgba(step(0.07), dark ? 0.72 : 0.86),
    "--surface-2": rgba(step(0.11), dark ? 0.6 : 0.7),
    "--surface-solid": step(0.07),
    "--surface-hover": rgba(step(0.15), dark ? 0.75 : 0.9),
    "--border": rgba(dark ? "#ffffff" : "#0b1020", dark ? 0.1 : 0.12),
    "--border-strong": rgba(dark ? "#ffffff" : "#0b1020", dark ? 0.18 : 0.2),

    "--text": text,
    "--text-muted": textMuted,
    "--text-faint": mix(text, base, 0.62),

    "--accent": accent,
    "--accent-2": accent2,
    "--accent-soft": rgba(accent, dark ? 0.16 : 0.12),
    "--accent-border": rgba(accent, 0.42),
    "--accent-text": readableOn(accent),
    "--accent-hover": dark ? lighten(accent, 0.12) : darken(accent, 0.08),
    "--gradient": `linear-gradient(135deg, ${accent}, ${accent2})`,

    "--success": dark ? "#34d399" : "#059669",
    "--warning": dark ? "#fbbf24" : "#d97706",
    "--danger": dark ? "#f87171" : "#dc2626",
    "--info": dark ? "#60a5fa" : "#2563eb",
    "--success-soft": rgba(dark ? "#34d399" : "#059669", dark ? 0.16 : 0.12),
    "--warning-soft": rgba(dark ? "#fbbf24" : "#d97706", dark ? 0.16 : 0.12),
    "--danger-soft": rgba(dark ? "#f87171" : "#dc2626", dark ? 0.16 : 0.12),

    "--radius-sm": `${0.375 * settings.radius}rem`,
    "--radius-md": `${0.625 * settings.radius}rem`,
    "--radius-lg": `${1 * settings.radius}rem`,
    "--radius-xl": `${1.5 * settings.radius}rem`,

    "--glass-blur": `${settings.glass}px`,
    "--font-scale": `${settings.fontScale}`,
    "--gap": settings.density === "compact" ? "0.625rem" : "1rem",
    "--pad": settings.density === "compact" ? "0.75rem" : "1.25rem",
    "--row-h": settings.density === "compact" ? "2.25rem" : "2.75rem",

    "--glow-1": rgba(accent, 0.34 * settings.glow),
    "--glow-2": rgba(accent2, 0.3 * settings.glow),
    "--glow-3": rgba(dark ? lighten(accent2, 0.2) : accent, 0.18 * settings.glow),
    "--glow-opacity": `${Math.min(1, settings.glow)}`,

    "--shadow-sm": dark
      ? "0 1px 2px rgba(0,0,0,0.4)"
      : "0 1px 2px rgba(15,23,42,0.08)",
    "--shadow-md": dark
      ? "0 12px 32px rgba(0,0,0,0.45)"
      : "0 12px 32px rgba(15,23,42,0.10)",
    "--shadow-lg": dark
      ? "0 24px 64px rgba(0,0,0,0.55)"
      : "0 24px 64px rgba(15,23,42,0.14)",
  };

  for (const [key, value] of Object.entries(vars)) {
    root.style.setProperty(key, value);
  }

  root.classList.toggle("dark", dark);
  root.classList.toggle("light", !dark);
  root.classList.toggle("glow-animated", settings.glowAnimated && settings.glow > 0);
  root.style.colorScheme = mode;

  return mode;
}
