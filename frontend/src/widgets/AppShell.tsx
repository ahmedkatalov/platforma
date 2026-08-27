import clsx from "clsx";
import { useEffect, useState } from "react";
import { NavLink, Outlet, useLocation, useNavigate } from "react-router-dom";

import { useAppDispatch, useAppSelector } from "@/app/store";
import { useLogoutMutation } from "@/features/auth/api/authApi";
import { sessionEnded } from "@/features/auth/authSlice";
import { useTrackActivityMutation } from "@/shared/api/meApi";
import { tokenStorage } from "@/shared/api/tokenStorage";
import { useTheme } from "@/shared/theme/ThemeProvider";
import {
  IconBook,
  IconChart,
  IconClose,
  IconDashboard,
  IconLogout,
  IconMoon,
  IconPalette,
  IconSettings,
  IconShield,
  IconSun,
  IconTerminal,
  IconUsers,
} from "@/shared/ui/icons";

type NavItem = { to: string; label: string; icon: React.ReactNode; end?: boolean };

const ADMIN_NAV: NavItem[] = [
  { to: "/admin", label: "Обзор", icon: <IconDashboard size={18} />, end: true },
  { to: "/admin/students", label: "Студенты", icon: <IconUsers size={18} /> },
  { to: "/admin/courses", label: "Курсы", icon: <IconBook size={18} /> },
  { to: "/admin/appearance", label: "Оформление", icon: <IconPalette size={18} /> },
  { to: "/admin/audit", label: "Журнал", icon: <IconShield size={18} /> },
];

const STUDENT_NAV: NavItem[] = [
  { to: "/learn", label: "Главная", icon: <IconDashboard size={18} />, end: true },
  { to: "/learn/courses", label: "Мои курсы", icon: <IconBook size={18} /> },
  { to: "/learn/stats", label: "Статистика", icon: <IconChart size={18} /> },
];

// Каждые 60 секунд отправляем время, проведённое на платформе.
function useActivityTracker() {
  const [track] = useTrackActivityMutation();

  useEffect(() => {
    if (!tokenStorage.access()) return;

    void track({ seconds: 0 });
    const timer = window.setInterval(() => {
      if (document.visibilityState === "visible") void track({ seconds: 60 });
    }, 60_000);

    return () => window.clearInterval(timer);
  }, [track]);
}

export default function AppShell() {
  const user = useAppSelector((state) => state.auth.user);
  const dispatch = useAppDispatch();
  const navigate = useNavigate();
  const location = useLocation();
  const { mode, toggleMode } = useTheme();
  const [logout] = useLogoutMutation();
  const [menuOpen, setMenuOpen] = useState(false);

  useActivityTracker();

  useEffect(() => setMenuOpen(false), [location.pathname]);

  if (!user) return null;

  const isAdmin = user.role === "admin";
  const nav = isAdmin ? ADMIN_NAV : STUDENT_NAV;

  const handleLogout = async () => {
    const refreshToken = tokenStorage.refresh();
    if (refreshToken) {
      try {
        await logout({ refreshToken }).unwrap();
      } catch {
        // Сессию всё равно закрываем локально.
      }
    }
    dispatch(sessionEnded());
    navigate("/login", { replace: true });
  };

  const initials = (user.fullName || user.email)
    .split(" ")
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0]?.toUpperCase())
    .join("");

  return (
    <div className="min-h-screen">
      {/* Боковая навигация */}
      <aside
        className={clsx(
          "fixed inset-y-0 left-0 z-40 flex w-64 flex-col border-r border-line bg-surface backdrop-blur-[var(--glass-blur)] transition-transform duration-200 lg:translate-x-0",
          menuOpen ? "translate-x-0" : "-translate-x-full",
        )}
      >
        <div className="flex h-16 items-center gap-2 border-b border-line px-5">
          <span className="grid h-9 w-9 place-items-center rounded-[var(--radius-md)] text-accent-fg" style={{ background: "var(--gradient)" }}>
            <IconTerminal size={20} />
          </span>
          <div className="min-w-0">
            <p className="truncate text-sm font-bold text-fg">DevOps Platform</p>
            <p className="truncate text-[11px] text-faint">
              {isAdmin ? "Панель администратора" : "Личный кабинет"}
            </p>
          </div>
          <button
            className="btn btn-ghost ml-auto h-8 w-8 !p-0 lg:hidden"
            onClick={() => setMenuOpen(false)}
            aria-label="Закрыть меню"
          >
            <IconClose size={18} />
          </button>
        </div>

        <nav className="flex-1 space-y-1 overflow-y-auto p-3">
          {nav.map((item) => (
            <NavLink
              key={item.to}
              to={item.to}
              end={item.end}
              className={({ isActive }) =>
                clsx(
                  "flex items-center gap-3 rounded-[var(--radius-md)] px-3 py-2.5 text-sm font-medium transition-colors",
                  isActive
                    ? "bg-accent-soft text-accent"
                    : "text-muted hover:bg-surface-2 hover:text-fg",
                )
              }
            >
              {item.icon}
              {item.label}
            </NavLink>
          ))}
        </nav>

        <div className="border-t border-line p-3">
          <NavLink
            to={isAdmin ? "/admin/profile" : "/learn/profile"}
            className={({ isActive }) =>
              clsx(
                "flex items-center gap-3 rounded-[var(--radius-md)] p-2 transition-colors",
                isActive ? "bg-accent-soft" : "hover:bg-surface-2",
              )
            }
          >
            <span
              className="grid h-9 w-9 shrink-0 place-items-center rounded-full text-xs font-bold text-accent-fg"
              style={{ background: "var(--gradient)" }}
            >
              {initials}
            </span>
            <span className="min-w-0 flex-1">
              <span className="block truncate text-sm font-semibold text-fg">
                {user.fullName || user.email}
              </span>
              <span className="block truncate text-[11px] text-faint">{user.email}</span>
            </span>
            <IconSettings size={16} className="text-faint" />
          </NavLink>

          <button
            className="btn btn-ghost mt-1 w-full justify-start"
            onClick={handleLogout}
          >
            <IconLogout size={18} />
            Выйти
          </button>
        </div>
      </aside>

      {menuOpen && (
        <div
          className="fixed inset-0 z-30 bg-black/50 lg:hidden"
          onClick={() => setMenuOpen(false)}
          aria-hidden="true"
        />
      )}

      {/* Контент */}
      <div className="lg:pl-64">
        <header className="sticky top-0 z-20 flex h-16 items-center gap-3 border-b border-line bg-surface px-4 backdrop-blur-[var(--glass-blur)] sm:px-6">
          <button
            className="btn btn-ghost h-9 w-9 !p-0 lg:hidden"
            onClick={() => setMenuOpen(true)}
            aria-label="Открыть меню"
          >
            <span className="flex flex-col gap-1">
              <span className="block h-0.5 w-4 rounded bg-current" />
              <span className="block h-0.5 w-4 rounded bg-current" />
              <span className="block h-0.5 w-4 rounded bg-current" />
            </span>
          </button>

          <div className="ml-auto flex items-center gap-2">
            <button
              className="btn btn-secondary h-9 w-9 !p-0"
              onClick={toggleMode}
              aria-label={mode === "dark" ? "Включить светлую тему" : "Включить тёмную тему"}
              title={mode === "dark" ? "Светлая тема" : "Тёмная тема"}
            >
              {mode === "dark" ? <IconSun size={18} /> : <IconMoon size={18} />}
            </button>
          </div>
        </header>

        <main className="mx-auto w-full max-w-7xl px-4 py-6 sm:px-6 sm:py-8">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
