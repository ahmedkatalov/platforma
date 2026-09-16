import { useMemo, type KeyboardEvent } from "react";

import { highlightToHtml, resolveLang } from "@/features/learning/lib/highlight";

// Редактор кода с подсветкой синтаксиса (Prism) и темой в стиле VS Code Dark+.
// Намеренно без Monaco: подсвеченный слой лежит ПОД прозрачным textarea, оба
// имеют одинаковые метрики и скроллятся вместе одним внешним контейнером.
export default function CodeEditor({
  value,
  onChange,
  language,
  minRows = 14,
  readOnly = false,
}: {
  value: string;
  onChange: (value: string) => void;
  language?: string;
  minRows?: number;
  readOnly?: boolean;
}) {
  const lines = value.split("\n");
  const rows = Math.max(minRows, lines.length);
  const html = useMemo(() => highlightToHtml(value, language), [value, language]);

  const handleKeyDown = (event: KeyboardEvent<HTMLTextAreaElement>) => {
    const textarea = event.currentTarget;
    const { selectionStart, selectionEnd } = textarea;

    if (event.key === "Tab") {
      event.preventDefault();
      const next = `${value.slice(0, selectionStart)}  ${value.slice(selectionEnd)}`;
      onChange(next);
      requestAnimationFrame(() => {
        textarea.selectionStart = textarea.selectionEnd = selectionStart + 2;
      });
      return;
    }

    if (event.key === "Enter") {
      // Сохраняем отступ текущей строки (важно для YAML и вложенного кода).
      const before = value.slice(0, selectionStart);
      const lineStart = before.lastIndexOf("\n") + 1;
      const indent = before.slice(lineStart).match(/^[ \t]*/)?.[0] ?? "";
      if (!indent) return;

      event.preventDefault();
      const next = `${before}\n${indent}${value.slice(selectionEnd)}`;
      onChange(next);
      requestAnimationFrame(() => {
        const position = selectionStart + 1 + indent.length;
        textarea.selectionStart = textarea.selectionEnd = position;
      });
    }
  };

  return (
    <div className="overflow-hidden rounded-[var(--radius-md)] border border-line bg-[var(--code-bg)]">
      <div className="flex items-center justify-between border-b border-[var(--code-line)] px-3 py-1.5">
        <span className="font-mono text-[11px] uppercase tracking-wide text-[var(--code-gutter)]">
          {resolveLang(language)}
        </span>
        <span className="text-[11px] text-[var(--code-gutter)]">{lines.length} строк</span>
      </div>

      {/* Единый скролл-контейнер: и номера, и код скроллятся вместе. */}
      <div className="max-h-[26rem] overflow-auto font-mono text-[13px] leading-6 lg:max-h-[32rem]">
        <div className="flex min-w-max">
          {/* Номера строк — липкие слева, чтобы оставаться при горизонтальной прокрутке. */}
          <div
            className="sticky left-0 z-[1] shrink-0 select-none border-r border-[var(--code-line)] bg-[var(--code-bg)] px-2 py-3 text-right text-[var(--code-gutter)]"
            aria-hidden="true"
          >
            {Array.from({ length: rows }, (_, i) => (
              <div key={i}>{i + 1}</div>
            ))}
          </div>

          {/* Подсвеченный слой и прозрачный textarea лежат в одной ячейке грида. */}
          <div className="relative grid">
            <pre
              className="prism-code pointer-events-none col-start-1 row-start-1 m-0 whitespace-pre px-3 py-3"
              aria-hidden="true"
            >
              <code
                className={`language-${resolveLang(language)}`}
                dangerouslySetInnerHTML={{ __html: html + "\n" }}
              />
            </pre>
            <textarea
              value={value}
              onChange={(e) => onChange(e.target.value)}
              onKeyDown={handleKeyDown}
              readOnly={readOnly}
              rows={rows}
              wrap="off"
              spellCheck={false}
              autoCapitalize="off"
              autoComplete="off"
              autoCorrect="off"
              className="col-start-1 row-start-1 w-full resize-none overflow-hidden whitespace-pre border-0 bg-transparent px-3 py-3 leading-6 outline-none"
              style={{ color: "transparent", caretColor: "var(--code-caret)" }}
              aria-label="Редактор кода"
            />
          </div>
        </div>
      </div>
    </div>
  );
}
