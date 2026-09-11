import { useEffect, useState } from "react";
import { Sparkles } from "lucide-react";

import { useGetAiSettingsQuery, useSaveAiSettingsMutation } from "@/features/admin/api/adminApi";
import { apiErrorMessage } from "@/shared/api/baseApi";
import { Badge, Button, Card, Field, Input, PageHeader, Spinner } from "@/shared/ui";
import { useToast } from "@/shared/ui/ToastProvider";

const DEFAULT_MODEL = "gemini-2.5-flash";

export default function AiSettingsPage() {
  const { data, isLoading } = useGetAiSettingsQuery();
  const [save, { isLoading: saving }] = useSaveAiSettingsMutation();
  const toast = useToast();

  const [enabled, setEnabled] = useState(false);
  const [model, setModel] = useState("");
  const [apiKey, setApiKey] = useState(""); // локальное поле — на сервер уходит только при вводе

  useEffect(() => {
    if (data) {
      setEnabled(data.enabled);
      setModel(data.model || "");
    }
  }, [data]);

  const configured = Boolean(data?.configured);
  const hasPanelKey = Boolean(data?.hasKey);
  const source = data?.source ?? "";

  const onSave = async () => {
    try {
      const payload: { enabled: boolean; model?: string; apiKey?: string } = {
        enabled,
        model: model.trim(),
      };
      // Ключ отправляем только если админ его ввёл — иначе оставляем прежний.
      if (apiKey.trim()) payload.apiKey = apiKey.trim();
      await save(payload).unwrap();
      setApiKey("");
      toast.success("Настройки ИИ сохранены");
    } catch (err) {
      toast.error(apiErrorMessage(err, "Не удалось сохранить настройки"));
    }
  };

  const onClearKey = async () => {
    if (!window.confirm("Удалить сохранённый ключ Gemini из панели?")) return;
    try {
      await save({ enabled, model: model.trim(), clearKey: true }).unwrap();
      setApiKey("");
      toast.success("Ключ удалён");
    } catch (err) {
      toast.error(apiErrorMessage(err, "Не удалось удалить ключ"));
    }
  };

  if (isLoading) {
    return (
      <div className="grid place-items-center py-20 text-accent">
        <Spinner size={32} />
      </div>
    );
  }

  return (
    <>
      <PageHeader
        title="ИИ-помощник"
        subtitle="Кнопка «Спросить у ИИ» в уроках — на базе Google Gemini"
        actions={
          <Button variant="primary" loading={saving} onClick={onSave}>
            Сохранить
          </Button>
        }
      />

      <div className="grid gap-[var(--gap)] lg:grid-cols-[1fr_18rem]">
        <Card className="space-y-4 p-[var(--pad)]">
          {/* Статус ключа */}
          <div className="flex items-center justify-between gap-3 rounded-[var(--radius-md)] bg-surface-2 p-3">
            <div>
              <p className="text-sm font-semibold text-fg">Ключ Gemini</p>
              <p className="text-xs text-muted">
                {source === "panel"
                  ? "Задан в этой панели (хранится на сервере)"
                  : source === "env"
                    ? "Взят из переменной окружения GEMINI_API_KEY"
                    : "Пока не задан"}
              </p>
            </div>
            {configured ? (
              <Badge tone="success">Подключён</Badge>
            ) : (
              <Badge tone="warning">Не задан</Badge>
            )}
          </div>

          <Field
            label="API-ключ Gemini"
            hint={
              hasPanelKey
                ? "Ключ сохранён. Оставьте поле пустым, чтобы не менять его; введите новый — чтобы заменить."
                : "Вставьте ключ из aistudio.google.com/apikey. Он сохранится на сервере и не попадёт к студентам."
            }
          >
            <Input
              type="password"
              value={apiKey}
              onChange={(e) => setApiKey(e.target.value)}
              placeholder={hasPanelKey ? "•••••••••• (сохранён)" : "AIza…"}
              autoComplete="new-password"
            />
          </Field>

          <Field label="Модель" hint="Например gemini-2.5-flash или gemini-2.0-flash. Пусто — по умолчанию.">
            <Input
              value={model}
              onChange={(e) => setModel(e.target.value)}
              placeholder={DEFAULT_MODEL}
              autoComplete="off"
            />
          </Field>

          {hasPanelKey && (
            <Button variant="ghost" className="text-danger" onClick={onClearKey} disabled={saving}>
              Удалить сохранённый ключ
            </Button>
          )}

          <label className="flex cursor-pointer items-start gap-3 rounded-[var(--radius-md)] p-2 hover:bg-surface-2">
            <input
              type="checkbox"
              checked={enabled}
              onChange={(e) => setEnabled(e.target.checked)}
              className="mt-0.5 h-4 w-4 accent-[var(--accent)]"
            />
            <span>
              <span className="block text-sm font-semibold text-fg">Включить ИИ-помощника</span>
              <span className="block text-xs text-muted">
                Студенты, выделив текст урока, смогут открыть чат и задать вопрос по теме
                {!configured && " (сначала задайте ключ выше)"}
              </span>
            </span>
          </label>

          <p className="rounded-[var(--radius-md)] bg-surface-2 px-3 py-2 text-xs text-muted">
            ⚠️ Выделенный студентом текст и его вопросы отправляются в Google Gemini для генерации
            ответа. Включая помощника, вы соглашаетесь с этим. Частота запросов ограничена, чтобы
            не выйти за бесплатные лимиты. Ключ хранится на сервере и студентам не передаётся.
          </p>
        </Card>

        <Card className="p-[var(--pad)]">
          <h2 className="mb-2 flex items-center gap-2 text-sm font-bold text-fg">
            <Sparkles size={16} className="text-accent" /> Что увидят студенты
          </h2>
          {configured && enabled ? (
            <p className="text-sm text-muted">
              При выделении текста в уроке рядом с «＋ В заметки» появится «Спросить у ИИ» —
              откроется чат с этим фрагментом как контекстом.
            </p>
          ) : (
            <p className="text-sm text-faint">
              {configured
                ? "Помощник выключен — кнопка скрыта."
                : "Задайте ключ Gemini, чтобы включить помощника."}
            </p>
          )}
          <p className="mt-3 text-xs text-faint">
            Бесплатный ключ: aistudio.google.com/apikey
          </p>
        </Card>
      </div>
    </>
  );
}
