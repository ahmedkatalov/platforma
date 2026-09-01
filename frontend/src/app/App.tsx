import { lazy, Suspense, useEffect } from "react";
import { Navigate, Route, Routes, useLocation } from "react-router-dom";

import { sessionEnded, userRefreshed } from "@/features/auth/authSlice";
import ForgotPasswordPage from "@/features/auth/ui/ForgotPasswordPage";
import LoginPage from "@/features/auth/ui/LoginPage";
import AppearancePage from "@/features/admin/ui/AppearancePage";
import AuditPage from "@/features/admin/ui/AuditPage";
import AdminCertificatesPage from "@/features/admin/ui/CertificatesPage";
import AdminCoursesPage from "@/features/admin/ui/CoursesPage";
import CourseEditorPage from "@/features/admin/ui/CourseEditorPage";
import AdminDashboardPage from "@/features/admin/ui/DashboardPage";
import RequestsPage from "@/features/admin/ui/RequestsPage";
import SettingsPage from "@/features/admin/ui/SettingsPage";
import StudentDetailPage from "@/features/admin/ui/StudentDetailPage";
import StudentsPage from "@/features/admin/ui/StudentsPage";
import CertificatePage from "@/features/certificates/ui/CertificatePage";
import LessonPage from "@/features/learning/ui/LessonPage";
import StudentCoursePage from "@/features/student/ui/CoursePage";
import CatalogPage from "@/features/student/ui/CatalogPage";
import StudentDashboardPage from "@/features/student/ui/DashboardPage";
import ProfilePage from "@/features/student/ui/ProfilePage";
import NotesPage from "@/features/student/ui/NotesPage";
import QuizzesPage from "@/features/student/ui/QuizzesPage";
import StudentStatsPage from "@/features/student/ui/StatsPage";
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
            <Suspense fallback={<FullScreenLoader />}>
              <SandboxPage />
            </Suspense>
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
