import { useState, type ReactNode } from "react";
import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";

import { IconCheck } from "@/shared/ui/icons";

// Якорь для оглавления: «Права доступа» → prava-dostupa-подобный стабильный id.
export function headingId(text: string): string {
  return text
    .toLowerCase()
    .replace(/[^a-zа-яё0-9\s-]/gi, "")
    .trim()
    .replace(/\s+/g, "-")
    .slice(0, 60);
}

function textOf(children: ReactNode): string {
  if (typeof children === "string") return children;
  if (Array.isArray(children)) return children.map(textOf).join("");
  if (children && typeof children === "object" && "props" in children) {
    return textOf((children as { props: { children?: ReactNode } }).props.children);
  }
  return "";
}

// Блок кода с кнопкой копирования — чтобы команды не перепечатывали руками.
function CodeBlock({ children }: { children: ReactNode }) {
  const [copied, setCopied] = useState(false);

  const copy = () => {
    void navigator.clipboard.writeText(textOf(children).trimEnd()).then(() => {
      setCopied(true);
      window.setTimeout(() => setCopied(false), 1600);
    });
  };

  return (
    <div className="group relative">
      <pre className="overflow-x-auto rounded-[var(--radius-md)] border border-line bg-[var(--bg-deep)] p-4">
        {children}
      </pre>
      <button
        type="button"
        onClick={copy}
        className={`absolute right-2 top-2 rounded-[var(--radius-sm)] border border-line px-2 py-1 text-[11px] font-semibold transition-opacity ${
          copied
            ? "bg-[var(--success-soft)] text-success opacity-100"
            : "bg-surface-solid text-muted opacity-0 hover:text-fg group-hover:opacity-100"
        }`}
        aria-label="Скопировать код"
      >
        {copied ? (
          <span className="flex items-center gap-1">
            <IconCheck size={11} /> скопировано
          </span>
        ) : (
          "копировать"
        )}
      </button>
    </div>
  );
}

// Разметка уроков: markdown, оформленный под текущую тему платформы.
export default function Markdown({ children }: { children: string }) {
  return (
    <div className="space-y-4 text-[0.9375rem] leading-[1.75] text-fg">
      <ReactMarkdown
        remarkPlugins={[remarkGfm]}
        components={{
          h1: ({ children }) => (
            <h1 className="mt-6 text-2xl font-bold text-fg first:mt-0">{children}</h1>
          ),
          h2: ({ children }) => (
            <h2
              id={headingId(textOf(children))}
              className="mt-8 scroll-mt-24 border-b border-line pb-2 text-xl font-bold text-fg first:mt-0"
            >
              {children}
            </h2>
          ),
          h3: ({ children }) => <h3 className="mt-5 text-lg font-bold text-fg">{children}</h3>,
          p: ({ children }) => <p className="text-muted">{children}</p>,
          ul: ({ children }) => (
            <ul className="list-disc space-y-2 pl-5 text-muted marker:text-accent">{children}</ul>
          ),
          ol: ({ children }) => (
            <ol className="list-decimal space-y-2 pl-5 text-muted marker:font-bold marker:text-accent">
              {children}
            </ol>
          ),
          li: ({ children }) => <li className="pl-1 leading-[1.7]">{children}</li>,
          strong: ({ children }) => <strong className="font-bold text-fg">{children}</strong>,
          img: ({ src, alt }) => (
            <img
              src={typeof src === "string" ? src : undefined}
              alt={alt ?? ""}
              loading="lazy"
              className="mx-auto max-w-full rounded-[var(--radius-md)] border border-line"
            />
          ),
          a: ({ children, href }) => (
            <a
              href={href}
              target="_blank"
              rel="noreferrer"
              className="font-medium text-accent hover:underline"
            >
              {children}
            </a>
          ),
          blockquote: ({ children }) => (
            <blockquote className="flex gap-3 rounded-[var(--radius-md)] border-l-4 border-[var(--accent)] bg-accent-soft px-4 py-3 text-sm leading-relaxed text-fg [&_p]:text-fg">
              <span aria-hidden="true" className="mt-0.5 shrink-0 text-base">
                💡
              </span>
              <span className="min-w-0">{children}</span>
            </blockquote>
          ),
          code: ({ className, children }) => {
            // Блок кода — моноширинный с прокруткой, инлайн — подсветка в тексте.
            const isBlock = Boolean(className?.startsWith("language-"));
            if (!isBlock) {
              return (
                <code className="rounded bg-surface-2 px-1.5 py-0.5 font-mono text-[0.85em] text-accent">
                  {children}
                </code>
              );
            }
            return (
              <code className="block overflow-x-auto whitespace-pre font-mono text-[0.8125rem] leading-relaxed text-fg">
                {children}
              </code>
            );
          },
          pre: ({ children }) => <CodeBlock>{children}</CodeBlock>,
          table: ({ children }) => (
            <div className="overflow-x-auto rounded-[var(--radius-md)] border border-line">
              <table className="w-full text-sm">{children}</table>
            </div>
          ),
          thead: ({ children }) => <thead className="bg-surface-2">{children}</thead>,
          th: ({ children }) => (
            <th className="border-b border-line px-3 py-2 text-left text-xs font-bold uppercase tracking-wide text-faint">
              {children}
            </th>
          ),
          td: ({ children }) => (
            <td className="border-b border-line/60 px-3 py-2 align-top text-muted">{children}</td>
          ),
          tr: ({ children }) => <tr className="even:bg-surface-2/40">{children}</tr>,
          hr: () => <hr className="border-line" />,
        }}
      >
        {children}
      </ReactMarkdown>
    </div>
  );
}
