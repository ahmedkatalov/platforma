import { lazy, Suspense, useEffect } from "react";
import { Navigate, Route, Routes, useLocation } from "react-router-dom";

import { sessionEnded, userRefreshed } from "@/features/auth/authSlice";
// Страницы входа лёгкие и нужны сразу при первом заходе — грузим их обычным импортом.
import ForgotPasswordPage from "@/features/auth/ui/ForgotPasswordPage";
import LoginPage from "@/features/auth/ui/LoginPage";
import CertificatePage from "@/features/certificates/ui/CertificatePage";

// Остальные страницы — по требованию (route-splitting): студент на телефоне не тянет
// код админки и графиков, пока они реально не открыты. Резко уменьшает первый бандл.
const AppearancePage = lazy(() => import("@/features/admin/ui/AppearancePage"));
const AuditPage = lazy(() => import("@/features/admin/ui/AuditPage"));
const AdminCertificatesPage = lazy(() => import("@/features/admin/ui/CertificatesPage"));
const AdminCoursesPage = lazy(() => import("@/features/admin/ui/CoursesPage"));
const CourseEditorPage = lazy(() => import("@/features/admin/ui/CourseEditorPage"));
const AdminDashboardPage = lazy(() => import("@/features/admin/ui/DashboardPage"));
const RequestsPage = lazy(() => import("@/features/admin/ui/RequestsPage"));
const SettingsPage = lazy(() => import("@/features/admin/ui/SettingsPage"));
const StudentDetailPage = lazy(() => import("@/features/admin/ui/StudentDetailPage"));
const StudentsPage = lazy(() => import("@/features/admin/ui/StudentsPage"));
const LessonPage = lazy(() => import("@/features/learning/ui/LessonPage"));
const StudentCoursePage = lazy(() => import("@/features/student/ui/CoursePage"));
const CatalogPage = lazy(() => import("@/features/student/ui/CatalogPage"));
const StudentDashboardPage = lazy(() => import("@/features/student/ui/DashboardPage"));
const ProfilePage = lazy(() => import("@/features/student/ui/ProfilePage"));
const NotesPage = lazy(() => import("@/features/student/ui/NotesPage"));
const QuizzesPage = lazy(() => import("@/features/student/ui/QuizzesPage"));
// StatsPage тянет recharts (~100 КБ) — теперь грузится только на странице статистики.
const StudentStatsPage = lazy(() => import("@/features/student/ui/StatsPage"));
// Песочница тянет тяжёлый движок Linux — грузим её отдельным чанком по требованию.
const SandboxPage = lazy(() => import("@/features/learning/ui/SandboxPage"));
import { useGetMeQuery } from "@/shared/api/meApi";
import { tokenStorage } from "@/shared/api/tokenStorage";
import { Spinner } from "@/shared/ui";
import AppShell from "@/widgets/AppShell";

import { useAppDispatch, useAppSelector } from "./store";

// Загружает профиль по сохранённому токену и следит за протухшей сессией.
function useBootstrapSession() {
  const dispatch = useAppDispatch();
  const hasToken = Boolean(tokenStorage.access());
  const { data, isError, isLoading } = useGetMeQuery(undefined, { skip: !hasToken });

  useEffect(() => {
    if (data?.user) dispatch(userRefreshed(data.user));
  }, [data, dispatch]);

  useEffect(() => {
    if (isError) dispatch(sessionEnded());
  }, [isError, dispatch]);

  useEffect(() => {
    const onExpired = () => dispatch(sessionEnded());
    window.addEventListener("platforma:session-expired", onExpired);
    return () => window.removeEventListener("platforma:session-expired", onExpired);
  }, [dispatch]);

  return hasToken && isLoading;
}

function FullScreenLoader() {
  return (
    <div className="grid min-h-screen place-items-center text-accent">
      <Spinner size={36} />
    </div>
  );
}

function RequireAuth({ role }: { role?: "admin" | "student" }) {
  const user = useAppSelector((state) => state.auth.user);
  const location = useLocation();

  if (!user) return <Navigate to="/login" state={{ from: location }} replace />;
  if (role && user.role !== role) {
    return <Navigate to={user.role === "admin" ? "/admin" : "/learn"} replace />;
  }
  return <AppShell />;
}

function PublicOnly({ children }: { children: React.ReactNode }) {
  const user = useAppSelector((state) => state.auth.user);
  if (user) return <Navigate to={user.role === "admin" ? "/admin" : "/learn"} replace />;
  return <>{children}</>;
}

// Песочница-терминал доступна только тем, у кого есть курс с терминалом (или админу).
function RequireSandbox({ children }: { children: React.ReactNode }) {
  const user = useAppSelector((state) => state.auth.user);
  const { data: me, isLoading } = useGetMeQuery(undefined, { skip: !user });
  if (isLoading) return <FullScreenLoader />;
  const allowed = user?.role === "admin" || Boolean(me?.sandboxAvailable);
  if (!allowed) return <Navigate to="/learn" replace />;
  return <>{children}</>;
}

export default function App() {
  const booting = useBootstrapSession();
  const user = useAppSelector((state) => state.auth.user);

  if (booting) return <FullScreenLoader />;

  return (
    <Routes>
      <Route
        path="/login"
        element={
          <PublicOnly>
            <LoginPage />
          </PublicOnly>
        }
      />
      {/* Регистрация закрыта: аккаунты создаёт администратор. Старую ссылку уводим на вход. */}
      <Route path="/register" element={<Navigate to="/login" replace />} />
      <Route
        path="/forgot-password"
        element={
          <PublicOnly>
            <ForgotPasswordPage />
          </PublicOnly>
        }
      />

      {/* Сертификат открыт всем: ссылку можно показать работодателю */}
      <Route path="/certificates/:serial" element={<CertificatePage />} />

      <Route path="/admin" element={<RequireAuth role="admin" />}>
        <Route index element={<AdminDashboardPage />} />
        <Route path="students" element={<StudentsPage />} />
        <Route path="students/:id" element={<StudentDetailPage />} />
        <Route path="requests" element={<RequestsPage />} />
        <Route path="settings" element={<SettingsPage />} />
        <Route path="courses" element={<AdminCoursesPage />} />
        <Route path="courses/:id" element={<CourseEditorPage />} />
        <Route path="certificates" element={<AdminCertificatesPage />} />
        <Route path="appearance" element={<AppearancePage />} />
        <Route path="audit" element={<AuditPage />} />
        <Route path="profile" element={<ProfilePage />} />
      </Route>

      <Route path="/learn" element={<RequireAuth />}>
        <Route index element={<CatalogPage />} />
        <Route path="dashboard" element={<StudentDashboardPage />} />
        <Route path="courses" element={<CatalogPage />} />
        <Route path="courses/:slug" element={<StudentCoursePage />} />
        <Route path="courses/:slug/lessons/:lessonId" element={<LessonPage />} />
        <Route path="quizzes" element={<QuizzesPage />} />
        <Route
          path="sandbox"
          element={
            <RequireSandbox>
              <Suspense fallback={<FullScreenLoader />}>
                <SandboxPage />
              </Suspense>
            </RequireSandbox>
          }
        />
        <Route path="notes" element={<NotesPage />} />
        <Route path="stats" element={<StudentStatsPage />} />
        <Route path="profile" element={<ProfilePage />} />
      </Route>

      <Route
        path="*"
        element={<Navigate to={user ? (user.role === "admin" ? "/admin" : "/learn") : "/login"} replace />}
      />
    </Routes>
  );
}
