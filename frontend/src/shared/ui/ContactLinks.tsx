import { MessageCircle, Send } from "lucide-react";

import type { ContactSettings } from "@/shared/types";
import { telegramUrl, whatsappUrl } from "@/shared/lib/contacts";

// Кнопки «Написать в Telegram / WhatsApp». Ничего не рисует, если связь выключена.
export function ContactLinks({
  contacts,
  className,
  size = "md",
}: {
  contacts?: ContactSettings | null;
  className?: string;
  size?: "sm" | "md";
}) {
  if (!contacts?.enabled) return null;
  const tg = telegramUrl(contacts.telegram);
  const wa = whatsappUrl(contacts.whatsapp);
  if (!tg && !wa) return null;

  const cls = `btn btn-secondary ${size === "sm" ? "h-9 text-sm" : ""}`;

  return (
    <div className={className ?? "flex flex-wrap gap-2"}>
      {tg && (
        <a href={tg} target="_blank" rel="noopener noreferrer" className={cls}>
          <Send size={16} /> Telegram
        </a>
      )}
      {wa && (
        <a href={wa} target="_blank" rel="noopener noreferrer" className={cls}>
          <MessageCircle size={16} /> WhatsApp
        </a>
      )}
    </div>
  );
}
