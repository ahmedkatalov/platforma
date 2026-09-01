import { useEffect, useState } from "react";
import { MessageSquare } from "lucide-react";

import { useGetContactsQuery, useSaveContactsMutation } from "@/features/admin/api/adminApi";
import { apiErrorMessage } from "@/shared/api/baseApi";
import type { ContactSettings } from "@/shared/types";
import { Button, Card, Field, Input, PageHeader, Spinner } from "@/shared/ui";
import { ContactLinks } from "@/shared/ui/ContactLinks";
import { useToast } from "@/shared/ui/ToastProvider";

const EMPTY: ContactSettings = {
  enabled: false,
  showOnLogin: false,
  telegram: "",
  whatsapp: "",
  note: "",
};

function Toggle({
  checked,
  onChange,
  label,
  hint,
}: {
  checked: boolean;
  onChange: (v: boolean) => void;
  label: string;
  hint?: string;
}) {
  return (
    <label className="flex cursor-pointer items-start gap-3 rounded-[var(--radius-md)] p-2 hover:bg-surface-2">
      <input
        type="checkbox"
        checked={checked}
        onChange={(e) => onChange(e.target.checked)}
        className="mt-0.5 h-4 w-4 accent-[var(--accent)]"
      />
      <span>
        <span className="block text-sm font-semibold text-fg">{label}</span>
        {hint && <span className="block text-xs text-muted">{hint}</span>}
      </span>
    </label>
  );
}

export default function SettingsPage() {
  const { data, isLoading } = useGetContactsQuery();
  const [save, { isLoading: saving }] = useSaveContactsMutation();
  const toast = useToast();
  const [form, setForm] = useState<ContactSettings>(EMPTY);

  useEffect(() => {
    if (data?.settings) setForm({ ...EMPTY, ...data.settings });
  }, [data]);

  const set = (patch: Partial<ContactSettings>) => setForm((f) => ({ ...f, ...patch }));

  const onSave = async () => {
    try {
      await save(form).unwrap();
      toast.success("Контакты сохранены");
    } catch (err) {
      toast.error(apiErrorMessage(err, "Не удалось сохранить контакты"));
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
        title="Связь со студентами"
        subtitle="Кнопки «Написать в Telegram/WhatsApp» — вы сами включаете и выключаете их"
        actions={
          <Button variant="primary" loading={saving} onClick={onSave}>
            Сохранить
          </Button>
        }
      />

      <div className="grid gap-[var(--gap)] lg:grid-cols-[1fr_18rem]">
        <Card className="space-y-4 p-[var(--pad)]">
          <div className="space-y-1">
            <Toggle
              checked={Boolean(form.enabled)}
              onChange={(v) => set({ enabled: v })}
              label="Включить связь"
              hint="Показывать кнопки в кабинете студента"
            />
            <Toggle
              checked={Boolean(form.showOnLogin)}
              onChange={(v) => set({ showOnLogin: v })}
              label="Показывать на странице входа"
              hint="Чтобы можно было написать вам ещё до входа"
            />
          </div>

          <Field label="Telegram" hint="Ник в формате @username или ссылка t.me/…">
            <Input
              value={form.telegram ?? ""}
              onChange={(e) => set({ telegram: e.target.value })}
              placeholder="@okvion"
            />
          </Field>

          <Field label="WhatsApp" hint="Номер в международном формате, только цифры">
            <Input
              value={form.whatsapp ?? ""}
              onChange={(e) => set({ whatsapp: e.target.value })}
              placeholder="79991234567"
            />
          </Field>

          <Field label="Подпись (необязательно)" hint="Короткий текст рядом с кнопками">
            <Input
              value={form.note ?? ""}
              onChange={(e) => set({ note: e.target.value })}
              placeholder="Пишите по любым вопросам"
            />
          </Field>
        </Card>

        <Card className="p-[var(--pad)]">
          <h2 className="mb-2 flex items-center gap-2 text-sm font-bold text-fg">
            <MessageSquare size={16} className="text-accent" /> Как увидят студенты
          </h2>
          {form.note && <p className="mb-3 text-sm text-muted">{form.note}</p>}
          {form.enabled ? (
            <ContactLinks contacts={form} className="flex flex-col gap-2" />
          ) : (
            <p className="text-sm text-faint">Связь выключена — кнопки скрыты.</p>
          )}
          {form.enabled && !form.telegram && !form.whatsapp && (
            <p className="text-sm text-faint">Заполните Telegram или WhatsApp.</p>
          )}
        </Card>
      </div>
    </>
  );
}
