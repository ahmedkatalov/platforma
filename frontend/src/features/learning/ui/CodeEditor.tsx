import { useLayoutEffect, useRef, type KeyboardEvent } from "react";

// Лёгкий редактор кода: номера строк, Tab-отступы и автоотступ новой строки.
// Намеренно без Monaco — он тянет мегабайты и внешние воркеры.
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
  const textareaRef = useRef<HTMLTextAreaElement>(null);
  const gutterRef = useRef<HTMLDivElement>(null);

  const lines = value.split("\n");
  const rows = Math.max(minRows, lines.length + 1);

  // Номера строк прокручиваются вместе с текстом.
  useLayoutEffect(() => {
    const textarea = textareaRef.current;
    const gutter = gutterRef.current;
    if (!textarea || !gutter) return;

    const sync = () => {
      gutter.scrollTop = textarea.scrollTop;
    };
    textarea.addEventListener("scroll", sync);
    return () => textarea.removeEventListener("scroll", sync);
  }, []);

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
      // Сохраняем отступ текущей строки (важно для YAML).
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
    <div className="overflow-hidden rounded-[var(--radius-md)] border border-line bg-[var(--bg-deep)]">
      <div className="flex items-center justify-between border-b border-line px-3 py-1.5">
        <span className="font-mono text-[11px] uppercase tracking-wide text-faint">
          {language ?? "text"}
        </span>
        <span className="text-[11px] text-faint">{lines.length} строк</span>
      </div>

      <div className="flex max-h-[26rem] font-mono text-[13px] leading-6 lg:max-h-[32rem]">
        <div
          ref={gutterRef}
          className="select-none overflow-hidden border-r border-line bg-surface-2 px-2 py-3 text-right text-faint"
          aria-hidden="true"
        >
          {Array.from({ length: rows }, (_, i) => (
            <div key={i}>{i + 1}</div>
          ))}
        </div>

        {/* wrap="off" — код скроллится вбок, а не переносится: отступы и номера
            строк не сбиваются (важно для YAML на узком экране). */}
        <textarea
          ref={textareaRef}
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
          className="flex-1 resize-none overflow-x-auto whitespace-pre bg-transparent px-3 py-3 leading-6 text-fg outline-none"
          aria-label="Редактор кода"
        />
      </div>
    </div>
  );
}
