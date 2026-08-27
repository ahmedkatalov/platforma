// Настройки оформления платформы. Админ задаёт значения по умолчанию для всех,
// пользователь может переопределить их для себя.

export type ThemeMode = "dark" | "light" | "system";

export type ThemeSettings = {
  mode: ThemeMode;
  accent: string; // основной акцент (hex)
  accent2: string; // вторичный акцент — градиенты, подсветка (hex)
  darkBase: string; // базовый цвет тёмной темы: из него выводятся поверхности
  lightBase: string; // базовый цвет светлой темы
  tone: string; // ключ пресета тона или "custom"
  radius: number; // множитель скруглений: 0.4 .. 1.8
  fontScale: number; // масштаб шрифта: 0.9 .. 1.2
  density: "compact" | "comfortable"; // плотность интерфейса
  glass: number; // размытие «стекла», px: 0 .. 32
  glow: number; // яркость фоновой подсветки: 0 .. 2
  glowAnimated: boolean; // медленное движение подсветки
  contrast: number; // контраст текста: 0.8 .. 1.2
};

export type AccentPreset = {
  key: string;
  name: string;
  accent: string;
  accent2: string;
};

export type TonePreset = {
  key: string;
  name: string;
  darkBase: string;
  lightBase: string;
};

// Готовые пары акцентов.
export const ACCENT_PRESETS: AccentPreset[] = [
  { key: "ocean", name: "Океан", accent: "#3b82f6", accent2: "#7c3aed" },
  { key: "aurora", name: "Аврора", accent: "#22d3ee", accent2: "#6366f1" },
  { key: "emerald", name: "Изумруд", accent: "#10b981", accent2: "#06b6d4" },
  { key: "sunset", name: "Закат", accent: "#f97316", accent2: "#ef4444" },
  { key: "violet", name: "Аметист", accent: "#8b5cf6", accent2: "#ec4899" },
  { key: "rose", name: "Роза", accent: "#f43f5e", accent2: "#8b5cf6" },
  { key: "amber", name: "Янтарь", accent: "#f59e0b", accent2: "#f97316" },
  { key: "graphite", name: "Графит", accent: "#64748b", accent2: "#94a3b8" },
];

// Базовые тона поверхностей.
export const TONE_PRESETS: TonePreset[] = [
  { key: "navy", name: "Ночь", darkBase: "#0b1020", lightBase: "#f6f8fc" },
  { key: "graphite", name: "Графит", darkBase: "#111114", lightBase: "#f5f5f6" },
  { key: "black", name: "Чёрная", darkBase: "#000000", lightBase: "#ffffff" },
  { key: "forest", name: "Хвоя", darkBase: "#0a1512", lightBase: "#f3f8f5" },
  { key: "plum", name: "Слива", darkBase: "#140b1c", lightBase: "#f9f5fc" },
  { key: "soft", name: "Мягкая", darkBase: "#1a2030", lightBase: "#eef1f7" },
];

export const DEFAULT_THEME: ThemeSettings = {
  mode: "dark",
  accent: "#3b82f6",
  accent2: "#7c3aed",
  darkBase: "#0b1020",
  lightBase: "#f6f8fc",
  tone: "navy",
  radius: 1,
  fontScale: 1,
  density: "comfortable",
  glass: 18,
  glow: 1,
  glowAnimated: true,
  contrast: 1,
};

const HEX_RE = /^#([0-9a-fA-F]{3}|[0-9a-fA-F]{6})$/;

const clamp = (v: number, min: number, max: number) =>
  Math.min(max, Math.max(min, Number.isFinite(v) ? v : min));

// Приводит произвольный объект (из localStorage или с сервера) к валидным настройкам.
export function normalizeTheme(raw: unknown): ThemeSettings {
  const base = { ...DEFAULT_THEME };
  if (!raw || typeof raw !== "object") return base;

  const src = raw as Partial<ThemeSettings>;

  if (src.mode === "dark" || src.mode === "light" || src.mode === "system") base.mode = src.mode;
  if (typeof src.accent === "string" && HEX_RE.test(src.accent)) base.accent = src.accent;
  if (typeof src.accent2 === "string" && HEX_RE.test(src.accent2)) base.accent2 = src.accent2;
  if (typeof src.darkBase === "string" && HEX_RE.test(src.darkBase)) base.darkBase = src.darkBase;
  if (typeof src.lightBase === "string" && HEX_RE.test(src.lightBase)) base.lightBase = src.lightBase;
  if (typeof src.tone === "string") base.tone = src.tone;
  if (typeof src.radius === "number") base.radius = clamp(src.radius, 0.4, 1.8);
  if (typeof src.fontScale === "number") base.fontScale = clamp(src.fontScale, 0.9, 1.2);
  if (src.density === "compact" || src.density === "comfortable") base.density = src.density;
  if (typeof src.glass === "number") base.glass = clamp(src.glass, 0, 32);
  if (typeof src.glow === "number") base.glow = clamp(src.glow, 0, 2);
  if (typeof src.glowAnimated === "boolean") base.glowAnimated = src.glowAnimated;
  if (typeof src.contrast === "number") base.contrast = clamp(src.contrast, 0.8, 1.2);

  return base;
}
