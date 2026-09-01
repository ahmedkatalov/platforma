import type { ContactSettings } from "@/shared/types";

// Ссылка на Telegram из @username или полной ссылки.
export function telegramUrl(handle?: string): string | null {
  if (!handle) return null;
  const h = handle
    .trim()
    .replace(/^https?:\/\/(t\.me|telegram\.me)\//i, "")
    .replace(/^@/, "");
  return h ? `https://t.me/${h}` : null;
}

// Ссылка на WhatsApp из номера (оставляем только цифры).
export function whatsappUrl(number?: string): string | null {
  if (!number) return null;
  const digits = number.replace(/\D/g, "");
  return digits ? `https://wa.me/${digits}` : null;
}

// Есть ли включённые контакты с хотя бы одним рабочим каналом.
export function hasAnyContact(c?: ContactSettings | null): boolean {
  return Boolean(c?.enabled && (telegramUrl(c.telegram) || whatsappUrl(c.whatsapp)));
}
