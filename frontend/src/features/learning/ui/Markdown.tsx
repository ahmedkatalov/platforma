import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";

// Разметка уроков: markdown, оформленный под текущую тему платформы.
export default function Markdown({ children }: { children: string }) {
  return (
    <div className="space-y-4 text-[0.9375rem] leading-relaxed text-fg">
      <ReactMarkdown
        remarkPlugins={[remarkGfm]}
        components={{
          h1: ({ children }) => (
            <h1 className="mt-6 text-2xl font-bold text-fg first:mt-0">{children}</h1>
          ),
          h2: ({ children }) => (
            <h2 className="mt-6 border-b border-line pb-2 text-xl font-bold text-fg first:mt-0">
              {children}
            </h2>
          ),
          h3: ({ children }) => <h3 className="mt-5 text-lg font-bold text-fg">{children}</h3>,
          p: ({ children }) => <p className="text-muted">{children}</p>,
          ul: ({ children }) => (
            <ul className="list-disc space-y-1.5 pl-5 text-muted marker:text-accent">{children}</ul>
          ),
          ol: ({ children }) => (
            <ol className="list-decimal space-y-1.5 pl-5 text-muted marker:text-accent marker:font-bold">
              {children}
            </ol>
          ),
          li: ({ children }) => <li className="pl-1">{children}</li>,
          strong: ({ children }) => <strong className="font-bold text-fg">{children}</strong>,
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
            <blockquote className="rounded-[var(--radius-md)] border-l-4 border-[var(--accent)] bg-accent-soft px-4 py-3 text-sm text-fg">
              {children}
            </blockquote>
          ),
          code: ({ className, children }) => {
            // Блок кода — с рамкой и горизонтальной прокруткой, инлайн — просто подсветка.
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
          pre: ({ children }) => (
            <pre className="overflow-x-auto rounded-[var(--radius-md)] border border-line bg-[var(--bg-deep)] p-4">
              {children}
            </pre>
          ),
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
            <td className="border-b border-line/60 px-3 py-2 text-muted">{children}</td>
          ),
          hr: () => <hr className="border-line" />,
        }}
      >
        {children}
      </ReactMarkdown>
    </div>
  );
}
