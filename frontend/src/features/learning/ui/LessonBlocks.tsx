import { useState, type ReactNode } from "react";
import { ChevronDown } from "lucide-react";

/*
 * Визуальные учебные блоки, встраиваемые прямо в markdown урока (схема данных не
 * меняется — это по-прежнему обычный текст в поле body/task). Markdown.tsx
 * перехватывает ```flow / ```anatomy / ```reveal и типизированные цитаты [!TIP].
 * Всё лёгкое (без тяжёлых библиотек), адаптивное и в тёмной теме платформы.
 */

// ── Поток/схема отношений: вертикальные шаги со стрелками, ряды через → ──────────
export function FlowDiagram({ source }: { source: string }) {
  const rows = source
    .split("\n")
    .map((l) => l.trim())
    .filter((l) => l && !/^[↓|v]+$/i.test(l) && l !== "->");

  return (
    <div
      className="my-5 flex flex-col items-center gap-1.5 rounded-[var(--radius-md)] border border-line bg-surface-2/50 p-4"
      role="img"
      aria-label="Схема: последовательность шагов"
    >
      {rows.map((row, i) => {
        const cells = row.split(/\s*(?:→|->)\s*/).filter(Boolean);
        return (
          <div key={i} className="flex w-full flex-col items-center gap-1.5">
            <div className="flex w-full flex-wrap items-center justify-center gap-1.5">
              {cells.map((cell, j) => {
                const accent = cell.startsWith("*");
                const text = accent ? cell.slice(1).trim() : cell;
                return (
                  <div key={j} className="flex items-center gap-1.5">
                    {j > 0 && <span className="text-accent" aria-hidden="true">→</span>}
                    <span
                      className={`rounded-[var(--radius-sm)] border px-3 py-1.5 text-center text-[0.8125rem] font-medium ${
                        accent
                          ? "border-accent bg-accent-soft text-accent"
                          : "border-line bg-surface-solid text-fg"
                      }`}
                    >
                      {text}
                    </span>
                  </div>
                );
              })}
            </div>
            {i < rows.length - 1 && (
              <ChevronDown size={16} className="shrink-0 text-faint" aria-hidden="true" />
            )}
          </div>
        );
      })}
    </div>
  );
}

// ── Анатомия конфигурации: код + пояснение к строке по тапу ─────────────────────
export function CodeAnatomy({ source }: { source: string }) {
  const lines = source.split("\n").map((raw) => {
    const idx = raw.indexOf(" ## ");
    if (idx === -1) return { code: raw, note: "" };
    return { code: raw.slice(0, idx), note: raw.slice(idx + 4).trim() };
  });
  // Первая аннотированная строка открыта по умолчанию — блок не выглядит пустым.
  const firstAnnotated = lines.findIndex((l) => l.note);
  const [active, setActive] = useState(firstAnnotated);

  return (
    <div className="my-5 overflow-hidden rounded-[var(--radius-md)] border border-line bg-[var(--bg-deep)]">
      <div className="border-b border-line px-3 py-1.5 text-[11px] font-semibold uppercase tracking-wide text-faint">
        Разбор по строкам — нажмите строку
      </div>
      <ul className="divide-y divide-line/50">
        {lines.map((line, i) => {
          const open = active === i && Boolean(line.note);
          const clickable = Boolean(line.note);
          return (
            <li key={i}>
              <button
                type="button"
                disabled={!clickable}
                onClick={() => setActive(open ? -1 : i)}
                className={`flex w-full items-start gap-2 px-3 py-1.5 text-left font-mono text-[0.8125rem] leading-relaxed transition-colors ${
                  open ? "bg-accent-soft" : clickable ? "hover:bg-surface-2" : ""
                }`}
              >
                <code
                  className={`min-w-0 flex-1 whitespace-pre-wrap break-words ${
                    open ? "text-accent" : "text-fg"
                  }`}
                >
                  {line.code || " "}
                </code>
                {clickable && (
                  <ChevronDown
                    size={13}
                    className={`mt-1 shrink-0 text-faint transition-transform ${open ? "rotate-180" : ""}`}
                    aria-hidden="true"
                  />
                )}
              </button>
              {open && (
                <p className="border-l-2 border-accent bg-accent-soft/40 px-4 py-2 text-[0.8125rem] leading-relaxed text-fg">
                  {line.note}
                </p>
              )}
            </li>
          );
        })}
      </ul>
    </div>
  );
}

// ── «Предскажите» / скрытый ответ: сначала думаешь, потом раскрываешь ────────────
export function Reveal({ source }: { source: string }) {
  const [shown, setShown] = useState(false);
  const parts = source.split(/\n-{3,}\n/);
  const question = (parts[0] ?? "").trim();
  const answer = (parts.slice(1).join("\n") || "").trim();

  return (
    <div className="my-5 rounded-[var(--radius-md)] border border-line bg-surface-2/50 p-4">
      <p className="flex items-start gap-2 text-[0.9375rem] text-fg">
        <span aria-hidden="true">🔮</span>
        <span className="min-w-0">{question}</span>
      </p>
      {shown ? (
        <p className="mt-3 border-l-2 border-accent bg-accent-soft px-3 py-2 text-[0.9375rem] leading-relaxed text-fg">
          {answer}
        </p>
      ) : (
        <button
          type="button"
          onClick={() => setShown(true)}
          className="btn btn-secondary btn-sm mt-3"
        >
          Показать ответ
        </button>
      )}
    </div>
  );
}

// ── Типизированные выноски: [!TIP]/[!WARNING]/[!DANGER]/… ────────────────────────
const CALLOUTS: Record<string, { icon: string; label: string; cls: string }> = {
  tip: { icon: "💡", label: "Совет", cls: "border-[var(--accent)] bg-accent-soft" },
  note: { icon: "📝", label: "Заметка", cls: "border-line bg-surface-2" },
  example: { icon: "🧪", label: "Пример", cls: "border-line bg-surface-2" },
  warning: { icon: "⚠️", label: "Внимание", cls: "border-[var(--warning)] bg-[var(--warning-soft)]" },
  danger: { icon: "🛑", label: "Опасно", cls: "border-[var(--danger)] bg-[var(--danger-soft)]" },
  remember: { icon: "📌", label: "Запомнить", cls: "border-[var(--accent)] bg-accent-soft" },
  try: { icon: "🔧", label: "Попробуйте", cls: "border-[var(--accent)] bg-accent-soft" },
  debug: { icon: "🐞", label: "Отладка", cls: "border-[var(--warning)] bg-[var(--warning-soft)]" },
  scenario: { icon: "🕒", label: "Ситуация", cls: "border-[var(--accent)] bg-accent-soft" },
};

/** Распознаёт первый маркер [!TYPE] в тексте цитаты; null если это обычная цитата. */
export function parseCalloutType(text: string): keyof typeof CALLOUTS | null {
  const m = text.trimStart().match(/^\[!(\w+)\]/i);
  if (!m) return null;
  const key = m[1].toLowerCase();
  return key in CALLOUTS ? (key as keyof typeof CALLOUTS) : null;
}

export function Callout({ type, children }: { type: keyof typeof CALLOUTS; children: ReactNode }) {
  const c = CALLOUTS[type];
  return (
    <div
      className={`my-4 flex gap-3 rounded-[var(--radius-md)] border-l-4 px-4 py-3 text-sm leading-relaxed text-fg [&_p]:text-fg [&_p]:my-0 [&_code]:text-accent ${c.cls}`}
    >
      <span aria-hidden="true" className="mt-0.5 shrink-0 text-base">
        {c.icon}
      </span>
      <div className="min-w-0">
        <p className="mb-0.5 text-[11px] font-bold uppercase tracking-wide text-faint">{c.label}</p>
        {children}
      </div>
    </div>
  );
}
