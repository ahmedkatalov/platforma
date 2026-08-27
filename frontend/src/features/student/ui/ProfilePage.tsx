import { useState, type FormEvent } from "react";

import { useAppDispatch, useAppSelector } from "@/app/store";
import { useChangePasswordMutation, useUpdateProfileMutation } from "@/features/auth/api/authApi";
import { userRefreshed } from "@/features/auth/authSlice";
import { apiErrorMessage } from "@/shared/api/baseApi";
import { useTheme } from "@/shared/theme/ThemeProvider";
import { ACCENT_PRESETS, DEFAULT_THEME, TONE_PRESETS, type ThemeMode } from "@/shared/theme/types";
import { Badge, Button, Card, Field, Input, PageHeader } from "@/shared/ui";
import { IconMoon, IconSettings, IconSun } from "@/shared/ui/icons";
import { useToast } from "@/shared/ui/ToastProvider";

// Профиль: данные аккаунта, смена пароля и личные настройки оформления.
export default function ProfilePage() {
  const user = useAppSelector((state) => state.auth.user);
  const dispatch = useAppDispatch();
  const toast = useToast();

  const [fullName, setFullName] = useState(user?.fullName ?? "");
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");

  const [updateProfile, { isLoading: savingProfile }] = useUpdateProfileMutation();
  const [changePassword, { isLoading: savingPassword }] = useChangePasswordMutation();

  const { settings, mode, update, replace, resetToPlatform } = useTheme();

  if (!user) return null;

  const saveProfile = async (event: FormEvent) => {
    event.preventDefault();
    try {
      const updated = await updateProfile({ fullName }).unwrap();
      dispatch(userRefreshed(updated));
      toast.success("Профиль обновлён");
    } catch (err) {
      toast.error(apiErrorMessage(err));
    }
  };

  const savePassword = async (event: FormEvent) => {
    event.preventDefault();
    try {
      await changePassword({ currentPassword, newPassword }).unwrap();
      setCurrentPassword("");
      setNewPassword("");
      toast.success("Пароль изменён");
    } catch (err) {
      toast.error(apiErrorMessage(err));
    }
  };

  return (
    <>
      <PageHeader title="Профиль" subtitle="Данные аккаунта и оформление интерфейса" />

      <div className="grid gap-[var(--gap)] lg:grid-cols-2">
        <Card className="p-[var(--pad)]">
          <h2 className="mb-4 text-base font-bold text-fg">Аккаунт</h2>

          <div className="mb-4 flex flex-wrap gap-2">
            <Badge tone={user.role === "admin" ? "accent" : "default"}>
              {user.role === "admin" ? "Администратор" : "Студент"}
            </Badge>
            {user.emailVerified && <Badge tone="success">Почта подтверждена</Badge>}
          </div>

          <form onSubmit={saveProfile} className="space-y-4">
            <Field label="Почта" hint="Изменить может только администратор">
              <Input value={user.email} disabled />
            </Field>

            <Field label="Имя и фамилия">
              <Input value={fullName} onChange={(e) => setFullName(e.target.value)} required />
            </Field>

            <Button type="submit" variant="primary" loading={savingProfile}>
              Сохранить
            </Button>
          </form>
        </Card>

        <Card className="p-[var(--pad)]">
          <h2 className="mb-4 text-base font-bold text-fg">Смена пароля</h2>

          <form onSubmit={savePassword} className="space-y-4">
            <Field label="Текущий пароль">
              <Input
                type="password"
                value={currentPassword}
                onChange={(e) => setCurrentPassword(e.target.value)}
                autoComplete="current-password"
                required
              />
            </Field>

            <Field label="Новый пароль" hint="Минимум 8 символов, буквы и цифры">
              <Input
                type="password"
                value={newPassword}
                onChange={(e) => setNewPassword(e.target.value)}
                autoComplete="new-password"
                required
              />
            </Field>

            <Button type="submit" variant="primary" loading={savingPassword}>
              Изменить пароль
            </Button>
          </form>
        </Card>
      </div>

      <Card className="mt-[var(--gap)] p-[var(--pad)]">
        <div className="mb-4 flex flex-wrap items-center justify-between gap-3">
          <h2 className="text-base font-bold text-fg">Оформление под себя</h2>
          <div className="flex gap-2">
            <Button onClick={resetToPlatform}>Как на платформе</Button>
            <Button onClick={() => replace(DEFAULT_THEME)}>Стандартное</Button>
          </div>
        </div>

        <div className="space-y-5">
          <div>
            <p className="label">Тема</p>
            <div className="grid gap-2 sm:grid-cols-3">
              {(
                [
                  { key: "dark", label: "Тёмная", icon: <IconMoon size={18} /> },
                  { key: "light", label: "Светлая", icon: <IconSun size={18} /> },
                  { key: "system", label: "Как в системе", icon: <IconSettings size={18} /> },
                ] as { key: ThemeMode; label: string; icon: React.ReactNode }[]
              ).map((option) => (
                <button
                  key={option.key}
                  onClick={() => update({ mode: option.key })}
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
          </div>

          <div>
            <p className="label">Акцент</p>
            <div className="flex flex-wrap gap-2">
              {ACCENT_PRESETS.map((preset) => {
                const active = preset.accent === settings.accent;
                return (
                  <button
                    key={preset.key}
                    onClick={() => update({ accent: preset.accent, accent2: preset.accent2 })}
                    className={`h-10 w-10 rounded-full border-2 transition-transform hover:scale-110 ${
                      active ? "border-[var(--accent)]" : "border-transparent"
                    }`}
                    style={{
                      background: `linear-gradient(135deg, ${preset.accent}, ${preset.accent2})`,
                    }}
                    title={preset.name}
                    aria-label={preset.name}
                  />
                );
              })}
            </div>
          </div>

          <div>
            <p className="label">Тон интерфейса</p>
            <div className="flex flex-wrap gap-2">
              {TONE_PRESETS.map((preset) => {
                const active = settings.tone === preset.key;
                return (
                  <button
                    key={preset.key}
                    onClick={() =>
                      update({
                        tone: preset.key,
                        darkBase: preset.darkBase,
                        lightBase: preset.lightBase,
                      })
                    }
                    className={`flex items-center gap-2 rounded-[var(--radius-md)] border px-3 py-2 text-xs font-semibold transition-colors ${
                      active ? "border-[var(--accent)] bg-accent-soft" : "border-line hover:bg-surface-2"
                    }`}
                  >
                    <span
                      className="h-5 w-5 rounded border border-line"
                      style={{ background: mode === "dark" ? preset.darkBase : preset.lightBase }}
                    />
                    <span className="text-fg">{preset.name}</span>
                  </button>
                );
              })}
            </div>
          </div>

          {user.role === "admin" && (
            <p className="text-xs text-faint">
              Полная настройка оформления и применение для всех — в разделе «Оформление».
            </p>
          )}
        </div>
      </Card>
    </>
  );
}
