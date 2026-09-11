import { useEffect, useState } from "react";
import { Sparkles } from "lucide-react";

import { useGetAiSettingsQuery, useSaveAiSettingsMutation } from "@/features/admin/api/adminApi";
import { apiErrorMessage } from "@/shared/api/baseApi";
import { Badge, Button, Card, PageHeader, Spinner } from "@/shared/ui";
import { useToast } from "@/shared/ui/ToastProvider";

export default function AiSettingsPage() {
  const { data, isLoading } = useGetAiSettingsQuery();
  const [save, { isLoading: saving }] = useSaveAiSettingsMutation();
  const toast = useToast();
  const [enabled, setEnabled] = useState(false);

  useEffect(() => {
    if (data) setEnabled(data.enabled);
  }, [data]);

  const configured = Boolean(data?.configured);

  const onSave = async () => {
    try {
      await save({ enabled }).unwrap();
      toast.success("Настройки ИИ сохранены");
    } catch (err) {
      toast.error(apiErrorMessage(err, "Не удалось сохранить настройки"));
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
          <Button variant="primary" loading={saving} onClick={onSave} disabled={!configured}>
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
                Задаётся на сервере в переменной <code>GEMINI_API_KEY</code>
              </p>
            </div>
            {configured ? (
              <Badge tone="success">Подключён</Badge>
            ) : (
              <Badge tone="warning">Не задан</Badge>
            )}
          </div>

          {!configured && (
            <p className="rounded-[var(--radius-md)] bg-[var(--warning-soft)] px-3 py-2 text-sm text-warning">
              Пока ключ не задан на сервере, помощник не работает — тумблер ни на что не влияет.
              Получить бесплатный ключ: aistudio.google.com/apikey, затем добавьте
              <code> GEMINI_API_KEY</code> в окружение бэкенда и перезапустите.
            </p>
          )}

          <label className="flex cursor-pointer items-start gap-3 rounded-[var(--radius-md)] p-2 hover:bg-surface-2">
            <input
              type="checkbox"
              checked={enabled}
              onChange={(e) => setEnabled(e.target.checked)}
              disabled={!configured}
              className="mt-0.5 h-4 w-4 accent-[var(--accent)]"
            />
            <span>
              <span className="block text-sm font-semibold text-fg">Включить ИИ-помощника</span>
              <span className="block text-xs text-muted">
                Студенты, выделив текст урока, смогут открыть чат и задать вопрос по теме
              </span>
            </span>
          </label>

          <p className="rounded-[var(--radius-md)] bg-surface-2 px-3 py-2 text-xs text-muted">
            ⚠️ Выделенный студентом текст и его вопросы отправляются в Google Gemini для генерации
            ответа. Включая помощника, вы соглашаетесь с этим. Частота запросов ограничена, чтобы
            не выйти за бесплатные лимиты.
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
                : "Задайте ключ на сервере, чтобы включить помощника."}
            </p>
          )}
        </Card>
      </div>
    </>
  );
}
