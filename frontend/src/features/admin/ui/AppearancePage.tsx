import { useState } from "react";

import { useSavePlatformThemeMutation } from "@/features/admin/api/adminApi";
import { apiErrorMessage } from "@/shared/api/baseApi";
import { useTheme } from "@/shared/theme/ThemeProvider";
import {
  ACCENT_PRESETS,
  DEFAULT_THEME,
  TONE_PRESETS,
  type ThemeMode,
  type ThemeSettings,
} from "@/shared/theme/types";
import {
  Badge,
  Button,
  Card,
  Field,
  PageHeader,
  Progress,
  Select,
} from "@/shared/ui";
import { Check, Moon, Palette, Settings, Sun } from "lucide-react";
import { useToast } from "@/shared/ui/ToastProvider";
import Logo from "@/shared/images/svg/logo.svg";

// Круглый сэмпл цвета поверх нативного color-input.
function ColorSwatch({
  value,
  onChange,
  title,
}: {
  value: string;
  onChange: (hex: string) => void;
  title: string;
}) {
  return (
    <label
      className="relative block h-11 w-11 shrink-0 cursor-pointer rounded-full border-2 border-line transition-transform hover:scale-105"
      style={{ background: value }}
      title={title}
    >
      <input
        type="color"
        value={value}
        onChange={(e) => onChange(e.target.value)}
        className="absolute inset-0 h-full w-full cursor-pointer opacity-0"
        aria-label={title}
      />
    </label>
  );
}

function Slider({
  label,
  value,
  min,
  max,
  step,
  suffix,
  onChange,
}: {
  label: string;
  value: number;
  min: number;
  max: number;
  step: number;
  suffix?: string;
  onChange: (value: number) => void;
}) {
  return (
    <div>
      <div className="mb-1.5 flex items-center justify-between">
        <span className="text-sm font-semibold text-muted">{label}</span>
        <span className="text-xs font-bold text-accent">
          {value.toFixed(step < 1 ? 2 : 0)}
          {suffix}
        </span>
      </div>
      <input
        type="range"
        min={min}
        max={max}
        step={step}
        value={value}
        onChange={(e) => onChange(Number(e.target.value))}
        className="h-2 w-full cursor-pointer appearance-none rounded-full bg-surface-2 accent-[var(--accent)]"
      />
    </div>
  );
}

export default function AppearancePage() {
  const { settings, mode, update, replace, resetToPlatform } = useTheme();
  const [savePlatformTheme, { isLoading: applying }] =
    useSavePlatformThemeMutation();
  const [savedForAll, setSavedForAll] = useState(false);
  const toast = useToast();

  const applyForEveryone = async () => {
    try {
      await savePlatformTheme(settings).unwrap();
      setSavedForAll(true);
      window.setTimeout(() => setSavedForAll(false), 3000);
      toast.success("Оформление применено для всей платформы");
    } catch (err) {
      toast.error(apiErrorMessage(err, "Не удалось сохранить оформление"));
    }
  };

  const activeAccent = ACCENT_PRESETS.find(
    (preset) =>
      preset.accent === settings.accent && preset.accent2 === settings.accent2,
  );

  const set = (patch: Partial<ThemeSettings>) => update(patch);

  return (
    <>
      <PageHeader
        title="Оформление"
        subtitle="Настройте интерфейс под себя и примените его для всех пользователей платформы"
        actions={
          <>
            <Button onClick={() => replace(DEFAULT_THEME)}>
              Сброс к стандарту
            </Button>
            <Button onClick={resetToPlatform}>Как у платформы</Button>
            <Button
              variant="primary"
              icon={savedForAll ? <Check size={18} /> : <Palette size={18} />}
              onClick={applyForEveryone}
              loading={applying}
            >
              {savedForAll ? "Применено" : "Применить для всех"}
            </Button>
          </>
        }
      />

      <div className="grid gap-[var(--gap)] xl:grid-cols-3">
        <div className="space-y-[var(--gap)] xl:col-span-2">
          {/* Режим */}
          <Card className="p-[var(--pad)]">
            <h2 className="mb-4 text-base font-bold text-fg">Тема</h2>
            <div className="grid gap-2 sm:grid-cols-3">
              {(
                [
                  { key: "dark", label: "Тёмная", icon: <Moon size={18} /> },
                  { key: "light", label: "Светлая", icon: <Sun size={18} /> },
                  {
                    key: "system",
                    label: "Как в системе",
                    icon: <Settings size={18} />,
                  },
                ] as { key: ThemeMode; label: string; icon: React.ReactNode }[]
              ).map((option) => (
                <button
                  key={option.key}
                  onClick={() => set({ mode: option.key })}
                  className={`flex items-center justify-center gap-2 rounded-[var(--radius-md)] border p-3 text-sm font-semibold transition-colors ${
                    settings.mode === option.key
                      ? "border-[var(--accent)] bg-accent-soft text-accent"
                      : "border-line text-muted hover:bg-surface-2"
                  }`}
                >
                  {option.icon}
                  {option.label}
                </button>
              ))}
            </div>
          </Card>

          {/* Акцент */}
          <Card className="p-[var(--pad)]">
            <div className="mb-4 flex items-center justify-between">
              <h2 className="text-base font-bold text-fg">Акцентные цвета</h2>
              {activeAccent && <Badge tone="accent">{activeAccent.name}</Badge>}
            </div>

            <div className="mb-4 grid grid-cols-2 gap-2 sm:grid-cols-4">
              {ACCENT_PRESETS.map((preset) => {
                const active = preset.key === activeAccent?.key;
                return (
                  <button
                    key={preset.key}
                    onClick={() =>
                      set({ accent: preset.accent, accent2: preset.accent2 })
                    }
                    className={`flex items-center gap-2 rounded-[var(--radius-md)] border p-2.5 text-xs font-semibold transition-colors ${
                      active
                        ? "border-[var(--accent)] bg-accent-soft"
                        : "border-line hover:bg-surface-2"
                    }`}
                  >
                    <span
                      className="h-6 w-6 shrink-0 rounded-full"
                      style={{
                        background: `linear-gradient(135deg, ${preset.accent}, ${preset.accent2})`,
                      }}
                    />
                    <span className="truncate text-fg">{preset.name}</span>
                  </button>
                );
              })}
            </div>

            <div className="flex flex-wrap items-center gap-4">
              <div className="flex items-center gap-3">
                <ColorSwatch
                  value={settings.accent}
                  onChange={(accent) => set({ accent })}
                  title="Основной акцент"
                />
                <div>
                  <p className="text-sm font-semibold text-fg">Основной</p>
                  <p className="font-mono text-xs text-faint">
                    {settings.accent}
                  </p>
                </div>
              </div>

              <div className="flex items-center gap-3">
                <ColorSwatch
                  value={settings.accent2}
                  onChange={(accent2) => set({ accent2 })}
                  title="Вторичный акцент"
                />
                <div>
                  <p className="text-sm font-semibold text-fg">Вторичный</p>
                  <p className="font-mono text-xs text-faint">
                    {settings.accent2}
                  </p>
                </div>
              </div>
            </div>
          </Card>

          {/* Тон поверхностей */}
          <Card className="p-[var(--pad)]">
            <h2 className="mb-4 text-base font-bold text-fg">Тон интерфейса</h2>

            <div className="mb-4 grid grid-cols-2 gap-2 sm:grid-cols-3">
              {TONE_PRESETS.map((preset) => {
                const active = settings.tone === preset.key;
                return (
                  <button
                    key={preset.key}
                    onClick={() =>
                      set({
                        tone: preset.key,
                        darkBase: preset.darkBase,
                        lightBase: preset.lightBase,
                      })
                    }
                    className={`flex items-center gap-2 rounded-[var(--radius-md)] border p-2.5 text-xs font-semibold transition-colors ${
                      active
                        ? "border-[var(--accent)] bg-accent-soft"
                        : "border-line hover:bg-surface-2"
                    }`}
                  >
                    <span
                      className="h-6 w-6 shrink-0 rounded-md border border-line"
                      style={{
                        background:
                          mode === "dark" ? preset.darkBase : preset.lightBase,
                      }}
                    />
                    <span className="truncate text-fg">{preset.name}</span>
                  </button>
                );
              })}
            </div>

            <div className="flex flex-wrap items-center gap-4">
              <div className="flex items-center gap-3">
                <ColorSwatch
                  value={settings.darkBase}
                  onChange={(darkBase) => set({ darkBase, tone: "custom" })}
                  title="База тёмной темы"
                />
                <div>
                  <p className="text-sm font-semibold text-fg">Тёмная база</p>
                  <p className="font-mono text-xs text-faint">
                    {settings.darkBase}
                  </p>
                </div>
              </div>

              <div className="flex items-center gap-3">
                <ColorSwatch
                  value={settings.lightBase}
                  onChange={(lightBase) => set({ lightBase, tone: "custom" })}
                  title="База светлой темы"
                />
                <div>
                  <p className="text-sm font-semibold text-fg">Светлая база</p>
                  <p className="font-mono text-xs text-faint">
                    {settings.lightBase}
                  </p>
                </div>
              </div>
            </div>
          </Card>

          {/* Тонкая настройка */}
          <Card className="space-y-5 p-[var(--pad)]">
            <h2 className="text-base font-bold text-fg">Форма и свет</h2>

            <Slider
              label="Скругление углов"
              value={settings.radius}
              min={0.4}
              max={1.8}
              step={0.05}
              suffix="×"
              onChange={(radius) => set({ radius })}
            />
            <Slider
              label="Размытие стекла"
              value={settings.glass}
              min={0}
              max={32}
              step={1}
              suffix=" px"
              onChange={(glass) => set({ glass })}
            />
            <Slider
              label="Яркость фоновой подсветки"
              value={settings.glow}
              min={0}
              max={2}
              step={0.05}
              suffix="×"
              onChange={(glow) => set({ glow })}
            />
            <Slider
              label="Масштаб шрифта"
              value={settings.fontScale}
              min={0.9}
              max={1.2}
              step={0.01}
              suffix="×"
              onChange={(fontScale) => set({ fontScale })}
            />
            <Slider
              label="Контраст текста"
              value={settings.contrast}
              min={0.8}
              max={1.2}
              step={0.01}
              suffix="×"
              onChange={(contrast) => set({ contrast })}
            />

            <div className="grid gap-4 sm:grid-cols-2">
              <Field label="Плотность интерфейса">
                <Select
                  value={settings.density}
                  onChange={(e) =>
                    set({ density: e.target.value as ThemeSettings["density"] })
                  }
                >
                  <option value="comfortable">Просторная</option>
                  <option value="compact">Компактная</option>
                </Select>
              </Field>

              <label className="flex cursor-pointer items-end gap-2 pb-2 text-sm text-muted">
                <input
                  type="checkbox"
                  checked={settings.glowAnimated}
                  onChange={(e) => set({ glowAnimated: e.target.checked })}
                  className="h-4 w-4 accent-[var(--accent)]"
                />
                Живая подсветка фона
              </label>
            </div>
          </Card>
        </div>

        {/* Предпросмотр */}
        <div className="xl:sticky xl:top-24 xl:h-fit">
          <Card className="p-[var(--pad)]">
            <h2 className="mb-4 text-base font-bold text-fg">Предпросмотр</h2>

            <div className="space-y-4">
              <div className="flex items-center gap-3">
                <span
                  className="grid h-11 w-11 place-items-center rounded-[var(--radius-md)] text-accent-fg"
                  style={{ background: "#fff" }}
                >
                  <img src={Logo} alt="" className="h-8 w-8" />
                </span>
                <div>
                  <p className="text-sm font-bold text-fg">Okvion Learning</p>
                  <p className="text-xs text-faint">так выглядит интерфейс</p>
                </div>
              </div>

              <div className="card-flat p-3">
                <div className="mb-2 flex items-center justify-between">
                  <span className="text-sm font-semibold text-fg">
                    Прогресс курса
                  </span>
                  <span className="text-xs font-bold text-accent">68%</span>
                </div>
                <Progress value={68} />
              </div>

              <div className="flex flex-wrap gap-2">
                <Badge tone="accent">Docker</Badge>
                <Badge tone="success">Пройден</Badge>
                <Badge tone="warning">В процессе</Badge>
                <Badge tone="danger">Ошибка</Badge>
              </div>

              <div className="flex flex-wrap gap-2">
                <Button variant="primary">Продолжить</Button>
                <Button>Отложить</Button>
                <Button variant="ghost">Позже</Button>
              </div>

              <div className="card-flat p-3 font-mono text-xs">
                <p className="text-success">$ docker ps</p>
                <p className="text-muted">CONTAINER ID IMAGE STATUS</p>
                <p className="text-fg">a1b2c3d4e5f6 nginx Up 2 minutes</p>
              </div>

              <p className="text-xs text-muted">
                Изменения применяются сразу и сохраняются лично для вас. Кнопка
                «Применить для всех» делает это оформление стандартным для
                студентов.
              </p>
            </div>
          </Card>
        </div>
      </div>
    </>
  );
}
