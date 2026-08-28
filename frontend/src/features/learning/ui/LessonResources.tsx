import type { LessonResource } from "@/shared/types";
import { Card } from "@/shared/ui";

// Ссылки на первоисточники под уроком: официальная документация,
// спецификации и материалы, к которым стоит вернуться после курса.
export default function LessonResources({ items }: { items?: LessonResource[] }) {
  if (!items || items.length === 0) return null;

  return (
    <Card className="mt-[var(--gap)] p-[var(--pad)]">
      <h2 className="mb-1 text-base font-bold text-fg">Материалы по теме</h2>
      <p className="mb-4 text-xs text-faint">
        Первоисточники: документация обновляется чаще любого курса — сверяйтесь с ней.
      </p>

      <ul className="space-y-2">
        {items.map((item) => (
          <li key={item.url}>
            <a
              href={item.url}
              target="_blank"
              rel="noreferrer noopener"
              className="card-flat flex items-start gap-3 p-3 transition-colors hover:bg-surface-hover"
            >
              <span className="mt-0.5 shrink-0 text-accent" aria-hidden="true">
                <svg width={16} height={16} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth={1.75}>
                  <path d="M10 13a5 5 0 0 0 7.5.5l3-3a5 5 0 0 0-7-7l-1.5 1.5" strokeLinecap="round" />
                  <path d="M14 11a5 5 0 0 0-7.5-.5l-3 3a5 5 0 0 0 7 7l1.5-1.5" strokeLinecap="round" />
                </svg>
              </span>

              <span className="min-w-0 flex-1">
                <span className="block text-sm font-semibold text-fg">{item.title}</span>
                {item.note && <span className="mt-0.5 block text-xs text-muted">{item.note}</span>}
                <span className="mt-1 block truncate font-mono text-[11px] text-faint">
                  {item.url.replace(/^https?:\/\//, "")}
                </span>
              </span>
            </a>
          </li>
        ))}
      </ul>
    </Card>
  );
}
