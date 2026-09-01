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
  LayoutGrid,
  BarChart3,
  ChevronLeft,
  ChevronRight,
  Check,
  X,
  LogOut,
  Edit2,
  Moon,
  Palette,
  Settings,
  Shield,
  Sun,
  Users,
  Book,
  TerminalSquare,
  Inbox,
  MessageSquare,
} from "lucide-react";
import Logo from "@/shared/images/svg/logo.svg";

type NavItem = { to: string; label: string; icon: React.ReactNode; end?: boolean };

const SIDEBAR_COLLAPSED_KEY = "platforma.sidebarCollapsed";

const ADMIN_NAV: NavItem[] = [
  { to: "/admin", label: "Обзор", icon: <LayoutGrid size={18} />, end: true },
  { to: "/admin/students", label: "Студенты", icon: <Users size={18} /> },
  { to: "/admin/requests", label: "Заявки", icon: <Inbox size={18} /> },
  { to: "/admin/courses", label: "Курсы", icon: <Book size={18} /> },
  { to: "/admin/certificates", label: "Сертификаты", icon: <Shield size={18} /> },
  { to: "/admin/settings", label: "Связь", icon: <MessageSquare size={18} /> },
  { to: "/admin/appearance", label: "Оформление", icon: <Palette size={18} /> },
  { to: "/admin/audit", label: "Журнал", icon: <Settings size={18} /> },
];

const STUDENT_NAV: NavItem[] = [
  { to: "/learn", label: "Курсы", icon: <Book size={18} />, end: true },
  { to: "/learn/dashboard", label: "Обзор", icon: <LayoutGrid size={18} /> },
  { to: "/learn/sandbox", label: "Песочница", icon: <TerminalSquare size={18} /> },
  { to: "/learn/quizzes", label: "Квизы", icon: <Check size={18} /> },
  { to: "/learn/notes", label: "Заметки", icon: <Edit2 size={18} /> },
  { to: "/learn/stats", label: "Статистика", icon: <BarChart3 size={18} /> },
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
  const [sidebarCollapsed, setSidebarCollapsed] = useState(
    () => localStorage.getItem(SIDEBAR_COLLAPSED_KEY) === "true",
  );

  useActivityTracker();

  useEffect(() => setMenuOpen(false), [location.pathname]);
  useEffect(() => {
    localStorage.setItem(SIDEBAR_COLLAPSED_KEY, String(sidebarCollapsed));
  }, [sidebarCollapsed]);

  if (!user) return null;

  const isAdmin = user.role === "admin";
  // Навигацию выбираем по разделу, а не только по роли: админ может открыть
  // курс в разделе /learn и пройти его ровно как студент.
  const inStudentArea = location.pathname.startsWith("/learn");
  const nav =
    isAdmin && inStudentArea
      ? [
          { to: "/admin/courses", label: "← В админку", icon: <ChevronRight size={18} className="rotate-180" />, end: false },
          ...STUDENT_NAV,
        ]
      : inStudentArea
        ? STUDENT_NAV
        : ADMIN_NAV;

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
          "fixed inset-y-0 left-0 z-40 flex w-64 flex-col border-r border-line bg-surface backdrop-blur-[var(--glass-blur)] transition-[width,transform] duration-200 lg:translate-x-0",
          sidebarCollapsed ? "lg:w-20" : "lg:w-64",
          menuOpen ? "translate-x-0" : "-translate-x-full",
        )}
      >
        <div className={clsx("flex h-16 items-center gap-2 border-b border-line px-5", sidebarCollapsed && "lg:justify-center lg:px-3")}>
          <span className="grid h-11 w-11 place-items-center rounded-[var(--radius-md)] text-accent-fg" style={{ background: "#fff" }}>
            <img src={Logo} alt="" className="h-8 w-8" />
          </span>
          <div className={clsx("min-w-0", sidebarCollapsed && "lg:hidden")}>
            <p className="truncate text-sm font-bold text-fg">Okvion Learning</p>
            <p className="truncate text-[11px] text-faint">
              {inStudentArea
                ? isAdmin
                  ? "Просмотр как студент"
                  : "Личный кабинет"
                : "Панель администратора"}
            </p>
          </div>
          <button
            className="btn btn-ghost ml-auto h-8 w-8 !p-0 lg:hidden"
            onClick={() => setMenuOpen(false)}
            aria-label="Закрыть меню"
          >
            <X size={18} />
          </button>
        </div>

        <nav className={clsx("flex-1 space-y-1 overflow-y-auto p-3", sidebarCollapsed ? "lg:p-2" : "lg:p-3")}>
          {nav.map((item) => (
            <NavLink
              key={item.to}
              to={item.to}
              end={item.end}
              className={({ isActive }) =>
                clsx(
                  "flex items-center gap-3 rounded-[var(--radius-md)] px-3 py-2.5 text-sm font-medium transition-colors",
                  sidebarCollapsed && "lg:justify-center lg:px-0",
                  isActive
                    ? "bg-accent-soft text-accent"
                    : "text-muted hover:bg-surface-2 hover:text-fg",
                )
              }
            >
              <span className="shrink-0">{item.icon}</span>
              <span className={clsx(sidebarCollapsed && "lg:sr-only")}>{item.label}</span>
            </NavLink>
          ))}
        </nav>

        <div className={clsx("border-t border-line p-3", sidebarCollapsed ? "lg:p-2" : "lg:p-3")}>
          <NavLink
            to={isAdmin ? "/admin/profile" : "/learn/profile"}
            className={({ isActive }) =>
              clsx(
                "flex items-center gap-3 rounded-[var(--radius-md)] p-2 transition-colors",
                sidebarCollapsed && "lg:justify-center",
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
            <span className={clsx("min-w-0 flex-1", sidebarCollapsed && "lg:hidden")}>
              <span className="block truncate text-sm font-semibold text-fg">
                {user.fullName || user.email}
              </span>
              <span className="block truncate text-[11px] text-faint">{user.email}</span>
            </span>
            <Settings size={16} className={clsx("text-faint", sidebarCollapsed && "lg:hidden")} />
          </NavLink>

          <button
            className={clsx("btn btn-ghost mt-1 w-full justify-start", sidebarCollapsed && "lg:justify-center lg:px-0")}
            onClick={handleLogout}
            title={sidebarCollapsed ? "Выйти" : undefined}
          >
            <LogOut size={18} />
            <span className={clsx(sidebarCollapsed && "lg:sr-only")}>Выйти</span>
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
      <div className={clsx("transition-[padding] duration-200", sidebarCollapsed ? "lg:pl-20" : "lg:pl-64")}>
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

          <button
            className="btn btn-ghost hidden h-9 w-9 !p-0 lg:inline-flex"
            onClick={() => setSidebarCollapsed((value) => !value)}
            aria-label={sidebarCollapsed ? "Развернуть меню" : "Свернуть меню"}
            title={sidebarCollapsed ? "Развернуть меню" : "Свернуть меню"}
          >
            {sidebarCollapsed ? <ChevronRight size={18} /> : <ChevronLeft size={18} />}
          </button>

          <div className="ml-auto flex items-center gap-2">
            <button
              className="btn btn-secondary h-9 w-9 !p-0"
              onClick={toggleMode}
              aria-label={mode === "dark" ? "Включить светлую тему" : "Включить тёмную тему"}
              title={mode === "dark" ? "Светлая тема" : "Тёмная тема"}
            >
              {mode === "dark" ? <Sun size={18} /> : <Moon size={18} />}
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
