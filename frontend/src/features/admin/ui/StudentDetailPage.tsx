import { useState } from "react";
import { Link, useParams } from "react-router-dom";
import { Area, AreaChart, CartesianGrid, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";

import {
  useEnrollMutation,
  useGetUserQuery,
  useResetUserPasswordMutation,
  useSetDueDateMutation,
  useUnenrollMutation,
  useUpdateUserMutation,
} from "@/features/admin/api/adminApi";
import { useGetAdminCoursesQuery } from "@/features/admin/api/coursesApi";
import { ChapterAccess } from "@/features/admin/ui/ChapterAccess";
import { apiErrorMessage } from "@/shared/api/baseApi";
import { lastSeenLabel } from "@/shared/lib/time";
import type { Attempt, CreatedStudent, UserStatus } from "@/shared/types";
import {
  Badge,
  Button,
  Card,
  EmptyState,
  Field,
  Input,
  Modal,
  PageHeader,
  Progress,
  Select,
  Spinner,
  StatCard,
} from "@/shared/ui";
import { Book, BarChart3, Clock, Key, Trash2 } from "lucide-react";
import { useToast } from "@/shared/ui/ToastProvider";

const STATUS_LABEL: Record<UserStatus, string> = {
  active: "Активен",
  invited: "Приглашён",
  blocked: "Заблокирован",
};

const dayFmt = new Intl.DateTimeFormat("ru-RU", { day: "2-digit", month: "short" });
const dueFmt = new Intl.DateTimeFormat("ru-RU", { day: "2-digit", month: "2-digit", year: "2-digit" });
const stampFmt = new Intl.DateTimeFormat("ru-RU", {
  day: "2-digit",
  month: "2-digit",
  hour: "2-digit",
  minute: "2-digit",
});

const ATTEMPT_LABEL: Record<Attempt["kind"], string> = {
  quiz: "Квиз",
  terminal: "Терминал",
  code: "Код",
};

export default function StudentDetailPage() {
  const { id = "" } = useParams();
  const { data, isLoading } = useGetUserQuery(id, { skip: !id });
  const { data: courses = [] } = useGetAdminCoursesQuery();

  const [updateUser] = useUpdateUserMutation();
  const [enroll, { isLoading: enrolling }] = useEnrollMutation();
  const [unenroll] = useUnenrollMutation();
  const [resetPassword, { isLoading: resetting }] = useResetUserPasswordMutation();

  const [courseId, setCourseId] = useState("");
  const [dueDate, setDueDate] = useState("");
  const [setEnrollmentDue] = useSetDueDateMutation();
  const [newPassword, setNewPassword] = useState<CreatedStudent | null>(null);
  const toast = useToast();

  if (isLoading || !data) {
    return (
      <div className="grid place-items-center py-20 text-accent">
        <Spinner size={32} />
      </div>
    );
  }

  const { user, summary, enrollments, activity, attempts, quiz } = data;
  const enrolledIds = new Set(enrollments.map((e) => e.courseId));
  const available = courses.filter((course) => !enrolledIds.has(course.id));

  const chartData = activity.map((day) => ({
    day: dayFmt.format(new Date(day.day)),
    Минуты: Math.round(day.secondsSpent / 60),
  }));

  const changeStatus = async (status: UserStatus) => {
    try {
      await updateUser({ id: user.id, status }).unwrap();
      toast.success("Статус обновлён");
    } catch (err) {
      toast.error(apiErrorMessage(err));
    }
  };

  const assignCourse = async () => {
    if (!courseId) return;
    try {
      await enroll({ userId: user.id, courseId, dueDate: dueDate || undefined }).unwrap();
      setCourseId("");
      setDueDate("");
      toast.success(dueDate ? "Курс назначен со сроком прохождения" : "Курс назначен");
    } catch (err) {
      toast.error(apiErrorMessage(err));
    }
  };

  // Срок можно поменять в любой момент; пустое значение снимает дедлайн.
  const changeDue = async (id: string, value: string) => {
    try {
      await setEnrollmentDue({ userId: user.id, courseId: id, dueDate: value }).unwrap();
      toast.success(value ? "Срок обновлён" : "Срок снят");
    } catch (err) {
      toast.error(apiErrorMessage(err));
    }
  };

  const removeCourse = async (id: string) => {
    try {
      await unenroll({ userId: user.id, courseId: id }).unwrap();
      toast.success("Курс снят");
    } catch (err) {
      toast.error(apiErrorMessage(err));
    }
  };

  const issueNewPassword = async () => {
    if (!window.confirm("Сгенерировать новый пароль? Все активные сессии студента закроются.")) return;
    try {
      const result = await resetPassword({ id: user.id, sendMail: true }).unwrap();
      setNewPassword(result);
    } catch (err) {
      toast.error(apiErrorMessage(err));
    }
  };

  return (
    <>
      <PageHeader
        title={user.fullName || user.email}
        subtitle={user.email}
        actions={
          <>
            <Link to="/admin/students" className="btn btn-ghost">
              К списку
            </Link>
            <Button icon={<Key size={18} />} onClick={issueNewPassword} loading={resetting}>
              Новый пароль
            </Button>
            {user.status === "blocked" ? (
              <Button variant="primary" onClick={() => changeStatus("active")}>
                Разблокировать
              </Button>
            ) : (
              <Button variant="danger" onClick={() => changeStatus("blocked")}>
                Заблокировать
              </Button>
            )}
          </>
        }
      />

      <div className="mb-[var(--gap)] flex flex-wrap gap-2">
        <Badge tone={user.role === "admin" ? "accent" : "default"}>
          {user.role === "admin" ? "Администратор" : "Студент"}
        </Badge>
        <Badge
          tone={
            user.status === "active" ? "success" : user.status === "invited" ? "warning" : "danger"
          }
        >
          {STATUS_LABEL[user.status]}
        </Badge>
        {user.emailVerified && <Badge tone="success">Почта подтверждена</Badge>}
        <Badge tone={summary.online ? "success" : "default"}>
          {summary.online
            ? "● Онлайн"
            : summary.lastSeenAt
              ? `Был в сети ${lastSeenLabel(summary.lastSeenAt)}`
              : "Ещё не заходил"}
        </Badge>
      </div>

      <div className="grid gap-[var(--gap)] sm:grid-cols-2 xl:grid-cols-4">
        <StatCard
          label="Прогресс"
          value={`${Math.round(summary.progress)}%`}
          hint={`${summary.lessonsCompleted} из ${summary.lessonsTotal} уроков`}
          icon={<BarChart3 size={20} />}
        />
        <StatCard
          label="Курсов"
          value={summary.courses}
          hint="Назначено студенту"
          icon={<Book size={20} />}
        />
        <StatCard
          label="Дней посещения"
          value={summary.daysVisited}
          hint="Всего за время обучения"
          icon={<Clock size={20} />}
        />
        <StatCard
          label="Времени на платформе"
          value={`${Math.floor(summary.minutesSpent / 60)} ч`}
          hint={`${summary.minutesSpent} минут суммарно`}
          icon={<Clock size={20} />}
        />
      </div>

      <div className="mt-[var(--gap)] grid gap-[var(--gap)] lg:grid-cols-5">
        <Card className="p-[var(--pad)] lg:col-span-3">
          <h2 className="mb-4 text-base font-bold text-fg">Активность за 30 дней</h2>
          {chartData.length === 0 ? (
            <EmptyState title="Студент ещё не заходил" icon={<Clock size={32} />} />
          ) : (
            <div className="h-56">
              <ResponsiveContainer width="100%" height="100%">
                <AreaChart data={chartData}>
                  <defs>
                    <linearGradient id="activityFill" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="0%" stopColor="var(--accent)" stopOpacity={0.5} />
                      <stop offset="100%" stopColor="var(--accent)" stopOpacity={0} />
                    </linearGradient>
                  </defs>
                  <CartesianGrid strokeDasharray="3 3" stroke="var(--border)" vertical={false} />
                  <XAxis
                    dataKey="day"
                    tick={{ fill: "var(--text-muted)", fontSize: 11 }}
                    axisLine={false}
                    tickLine={false}
                  />
                  <YAxis
                    tick={{ fill: "var(--text-muted)", fontSize: 11 }}
                    axisLine={false}
                    tickLine={false}
                    allowDecimals={false}
                  />
                  <Tooltip
                    contentStyle={{
                      background: "var(--surface-solid)",
                      border: "1px solid var(--border)",
                      borderRadius: "var(--radius-md)",
                      color: "var(--text)",
                    }}
                  />
                  <Area
                    type="monotone"
                    dataKey="Минуты"
                    stroke="var(--accent)"
                    strokeWidth={2}
                    fill="url(#activityFill)"
                  />
                </AreaChart>
              </ResponsiveContainer>
            </div>
          )}
        </Card>

        <Card className="p-[var(--pad)] lg:col-span-2">
          <h2 className="mb-4 text-base font-bold text-fg">Курсы студента</h2>

          <div className="mb-4 space-y-2">
            <Field label="Курс">
              <Select value={courseId} onChange={(e) => setCourseId(e.target.value)}>
                <option value="">Выберите курс…</option>
                {available.map((course) => (
                  <option key={course.id} value={course.id}>
                    {course.title}
                  </option>
                ))}
              </Select>
            </Field>

            <div className="flex items-end gap-2">
              <Field label="Срок прохождения" hint="Необязательно — напомним студенту письмом">
                <Input
                  type="date"
                  value={dueDate}
                  onChange={(e) => setDueDate(e.target.value)}
                />
              </Field>
              <Button variant="primary" onClick={assignCourse} disabled={!courseId} loading={enrolling}>
                Назначить
              </Button>
            </div>
          </div>

          {enrollments.length === 0 ? (
            <p className="py-6 text-center text-sm text-muted">Курсы ещё не назначены</p>
          ) : (
            <ul className="space-y-2">
              {enrollments.map((enrollment) => (
                <li key={enrollment.id} className="card-flat p-3">
                  <div className="flex items-center gap-3">
                    <div className="min-w-0 flex-1">
                      <p className="truncate text-sm font-semibold text-fg">
                        {enrollment.courseTitle}
                      </p>
                      <p className="text-xs text-faint">{enrollment.courseSlug}</p>
                    </div>
                    {enrollment.dueDate && (
                      <Badge tone={new Date(enrollment.dueDate) < new Date() ? "danger" : "warning"}>
                        до {dueFmt.format(new Date(enrollment.dueDate))}
                      </Badge>
                    )}
                    <Button
                      variant="ghost"
                      className="h-8 !px-2 text-danger"
                      onClick={() => removeCourse(enrollment.courseId)}
                      title="Снять курс"
                    >
                      <Trash2 size={16} />
                    </Button>
                  </div>

                  <div className="mt-2 flex items-center gap-2">
                    <span className="text-xs text-faint">Срок:</span>
                    <Input
                      type="date"
                      value={enrollment.dueDate ? enrollment.dueDate.slice(0, 10) : ""}
                      onChange={(e) => changeDue(enrollment.courseId, e.target.value)}
                      className="h-8 text-xs"
                    />
                  </div>

                  <ChapterAccess userId={id} courseId={enrollment.courseId} />
                </li>
              ))}
            </ul>
          )}

          <div className="mt-4">
            <p className="mb-1.5 text-xs font-semibold text-faint">Общий прогресс</p>
            <Progress value={summary.progress} />
          </div>
        </Card>
      </div>

      <div className="mt-[var(--gap)] grid gap-[var(--gap)] lg:grid-cols-5">
        <Card className="p-[var(--pad)] lg:col-span-2">
          <h2 className="mb-4 text-base font-bold text-fg">Квизы</h2>

          {quiz.attempts === 0 ? (
            <p className="py-4 text-center text-sm text-muted">Студент ещё не проходил квизы</p>
          ) : (
            <div className="space-y-3 text-sm">
              <div className="card-flat p-3">
                <div className="mb-1 flex items-center justify-between">
                  <span className="text-muted">Средний балл</span>
                  <span className="font-bold text-accent">{Math.round(quiz.averageScore)}%</span>
                </div>
                <Progress value={quiz.averageScore} />
              </div>

              <div className="grid grid-cols-2 gap-2">
                <div className="card-flat p-3 text-center">
                  <p className="text-lg font-bold text-fg">
                    {quiz.passed}/{quiz.attempts}
                  </p>
                  <p className="text-[11px] text-faint">успешных попыток</p>
                </div>
                <div className="card-flat p-3 text-center">
                  <p className="text-lg font-bold text-fg">{Math.round(quiz.accuracy)}%</p>
                  <p className="text-[11px] text-faint">верных ответов</p>
                </div>
                <div className="card-flat p-3 text-center">
                  <p className="text-lg font-bold text-fg">
                    {quiz.avgSecondsPerQuestion.toFixed(1)} c
                  </p>
                  <p className="text-[11px] text-faint">на вопрос</p>
                </div>
                <div className="card-flat p-3 text-center">
                  <p className="text-lg font-bold text-fg">{Math.round(quiz.bestScore)}%</p>
                  <p className="text-[11px] text-faint">лучший результат</p>
                </div>
              </div>
            </div>
          )}
        </Card>

        <Card className="overflow-hidden lg:col-span-3">
          <div className="border-b border-line px-[var(--pad)] py-3">
            <h2 className="text-base font-bold text-fg">Последние попытки</h2>
          </div>

          {attempts.length === 0 ? (
            <EmptyState title="Попыток пока нет" icon={<BarChart3 size={32} />} />
          ) : (
            <div className="max-h-80 overflow-y-auto">
              <table className="w-full text-sm">
                <thead className="sticky top-0 bg-surface-solid">
                  <tr className="border-b border-line text-left text-xs uppercase tracking-wide text-faint">
                    <th className="px-4 py-2.5 font-semibold">Когда</th>
                    <th className="px-4 py-2.5 font-semibold">Урок</th>
                    <th className="px-4 py-2.5 font-semibold">Тип</th>
                    <th className="px-4 py-2.5 font-semibold">Балл</th>
                  </tr>
                </thead>
                <tbody>
                  {attempts.map((attempt) => (
                    <tr key={attempt.id} className="border-b border-line/60 last:border-0">
                      <td className="whitespace-nowrap px-4 py-2.5 text-muted">
                        {stampFmt.format(new Date(attempt.createdAt))}
                      </td>
                      <td className="px-4 py-2.5 text-fg">{attempt.lessonTitle}</td>
                      <td className="px-4 py-2.5">
                        <Badge tone={attempt.kind === "quiz" ? "accent" : "default"}>
                          {ATTEMPT_LABEL[attempt.kind]}
                        </Badge>
                      </td>
                      <td className="px-4 py-2.5">
                        <span
                          className={`font-bold ${attempt.passed ? "text-success" : "text-warning"}`}
                        >
                          {Math.round(attempt.score)}%
                        </span>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </Card>
      </div>

      <Modal
        open={Boolean(newPassword)}
        onClose={() => setNewPassword(null)}
        title="Новый пароль выдан"
        footer={
          <Button variant="primary" onClick={() => setNewPassword(null)}>
            Готово
          </Button>
        }
      >
        {newPassword && (
          <div className="space-y-3 text-sm">
            <div className="card-flat p-3 font-mono">
              <p>
                <span className="text-faint">Логин: </span>
                <span className="font-bold text-fg">{newPassword.user.email}</span>
              </p>
              <p>
                <span className="text-faint">Пароль: </span>
                <span className="font-bold text-accent">{newPassword.tempPassword}</span>
              </p>
            </div>
            {newPassword.mailSent ? (
              <p className="rounded-[var(--radius-md)] bg-[var(--success-soft)] px-3 py-2 text-success">
                Пароль отправлен студенту на почту
              </p>
            ) : (
              <p className="rounded-[var(--radius-md)] bg-[var(--warning-soft)] px-3 py-2 text-warning">
                {newPassword.mailError || "Письмо не отправлено — передайте пароль вручную"}
              </p>
            )}
          </div>
        )}
      </Modal>
    </>
  );
}
