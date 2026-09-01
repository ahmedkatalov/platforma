import type { ReactNode } from "react";
import { Link } from "react-router-dom";

import { useTheme } from "@/shared/theme/ThemeProvider";
import { Moon, Sun, Terminal } from "lucide-react";

// Общая обёртка экранов входа: логотип, стеклянная карточка, переключатель темы.
export default function AuthLayout({
  title,
  subtitle,
  children,
  footer,
}: {
  title: string;
  subtitle: string;
  children: ReactNode;
  footer?: ReactNode;
}) {
  const { mode, toggleMode } = useTheme();

  return (
    <div className="flex min-h-screen items-center justify-center px-4 py-10">
      <button
        className="btn btn-secondary btn-icon fixed right-4 top-4"
        onClick={toggleMode}
        aria-label="Сменить тему"
      >
        {mode === "dark" ? <Sun size={18} /> : <Moon size={18} />}
      </button>

      <div className="w-full max-w-md">
        <Link to="/login" className="mb-6 flex items-center justify-center gap-3">
          <span
            className="grid h-11 w-11 place-items-center rounded-[var(--radius-md)] text-accent-fg"
            style={{ background: "var(--gradient)" }}
          >
            <Terminal size={24} />
          </span>
          <span className="text-xl font-extrabold tracking-tight">
            Okvion <span className="gradient-text">Learning</span>
          </span>
        </Link>

        <div className="card p-6 sm:p-8">
          <h1 className="text-xl font-bold text-fg">{title}</h1>
          <p className="mt-1 text-sm text-muted">{subtitle}</p>
          <div className="mt-6">{children}</div>
        </div>

        {footer && <div className="mt-4 text-center text-sm text-muted">{footer}</div>}
      </div>
    </div>
  );
}
