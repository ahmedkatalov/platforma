import { Link } from "react-router-dom";
import {
  Bar,
  BarChart,
  CartesianGrid,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";

import {
  useGetAuditQuery,
  useGetOverviewQuery,
  useGetStudentsProgressQuery,
} from "@/features/admin/api/adminApi";
import { useGetAdminCoursesQuery } from "@/features/admin/api/coursesApi";
import { Badge, Card, EmptyState, PageHeader, Progress, Spinner, StatCard } from "@/shared/ui";
import { Book, BarChart3, Users } from "lucide-react";

const dateFmt = new Intl.DateTimeFormat("ru-RU", { day: "2-digit", month: "short", hour: "2-digit", minute: "2-digit" });

export default function DashboardPage() {
  const { data: overview, isLoading } = useGetOverviewQuery();
  const { data: progress = [] } = useGetStudentsProgressQuery(10);
  const { data: courses = [] } = useGetAdminCoursesQuery();
  const { data: audit = [] } = useGetAuditQuery(8);

  if (isLoading || !overview) {
    return (
      <div className="grid place-items-center py-20 text-accent">
        <Spinner size={32} />
      </div>
    );
  }

  const chartData = courses.slice(0, 8).map((course) => ({
    name: course.title.length > 16 ? `${course.title.slice(0, 16)}…` : course.title,
    Студентов: course.studentsCount,
    Уроков: course.lessonsCount,
  }));

  return (
    <>
      <PageHeader
        title="Обзор платформы"
        subtitle="Ключевые показатели по студентам, курсам и активности"
      />

      <div className="grid gap-[var(--gap)] sm:grid-cols-2 xl:grid-cols-4">
        <StatCard
          label="Студентов"
          value={overview.students}
          hint={`${overview.activeStudents} активных · ${overview.invitedStudents} приглашены`}
          icon={<Users size={20} />}
        />
        <StatCard
          label="Курсов"
          value={overview.courses}
          hint={`${overview.publishedCourses} опубликовано · ${overview.lessons} уроков`}
          icon={<Book size={20} />}
        />
        <StatCard
          label="Записей на курсы"
          value={overview.enrollments}
          hint="Всего назначений студентам"
          icon={<BarChart3 size={20} />}
        />
        <StatCard
          label="Активность"
          value={overview.activeToday}
          hint={`сегодня · ${overview.activeWeek} за неделю`}
          icon={<BarChart3 size={20} />}
        />
      </div>

      <div className="mt-[var(--gap)] grid gap-[var(--gap)] lg:grid-cols-5">
        <Card className="p-[var(--pad)] lg:col-span-3">
          <h2 className="mb-4 text-base font-bold text-fg">Курсы: студенты и объём</h2>
          {chartData.length === 0 ? (
            <EmptyState
              title="Курсов пока нет"
              description="Создайте первый курс, чтобы увидеть статистику"
              icon={<Book size={32} />}
            />
          ) : (
            <div className="h-64">
              <ResponsiveContainer width="100%" height="100%">
                <BarChart data={chartData}>
                  <CartesianGrid strokeDasharray="3 3" stroke="var(--border)" vertical={false} />
                  <XAxis
                    dataKey="name"
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
                    cursor={{ fill: "var(--accent-soft)" }}
                    contentStyle={{
                      background: "var(--surface-solid)",
                      border: "1px solid var(--border)",
                      borderRadius: "var(--radius-md)",
                      color: "var(--text)",
                    }}
                  />
                  <Bar dataKey="Студентов" fill="var(--accent)" radius={[6, 6, 0, 0]} />
                  <Bar dataKey="Уроков" fill="var(--accent-2)" radius={[6, 6, 0, 0]} />
                </BarChart>
              </ResponsiveContainer>
            </div>
          )}
        </Card>

        <Card className="p-[var(--pad)] lg:col-span-2">
          <div className="mb-4 flex items-center justify-between">
            <h2 className="text-base font-bold text-fg">Успеваемость</h2>
            <Link to="/admin/students" className="text-xs font-semibold text-accent hover:underline">
              Все студенты
            </Link>
          </div>

          {progress.length === 0 ? (
            <EmptyState title="Студентов пока нет" icon={<Users size={32} />} />
          ) : (
            <ul className="space-y-3">
              {progress.slice(0, 6).map((student) => (
                <li key={student.userId}>
                  <Link
                    to={`/admin/students/${student.userId}`}
                    className="block rounded-[var(--radius-md)] p-2 transition-colors hover:bg-surface-2"
                  >
                    <div className="mb-1.5 flex items-center justify-between gap-2">
                      <span className="truncate text-sm font-semibold text-fg">
                        {student.fullName || student.email}
                      </span>
                      <span className="shrink-0 text-xs font-bold text-accent">
                        {Math.round(student.progress)}%
                      </span>
                    </div>
                    <Progress value={student.progress} />
                    <p className="mt-1 text-[11px] text-faint">
                      {student.lessonsCompleted} из {student.lessonsTotal} уроков ·{" "}
                      {student.daysVisited} дней посещения
                    </p>
                  </Link>
                </li>
              ))}
            </ul>
          )}
        </Card>
      </div>

      <Card className="mt-[var(--gap)] p-[var(--pad)]">
        <div className="mb-4 flex items-center justify-between">
          <h2 className="text-base font-bold text-fg">Последние действия</h2>
          <Link to="/admin/audit" className="text-xs font-semibold text-accent hover:underline">
            Весь журнал
          </Link>
        </div>

        {audit.length === 0 ? (
          <p className="py-6 text-center text-sm text-muted">Пока ничего не происходило</p>
        ) : (
          <ul className="divide-y divide-[var(--border)]">
            {audit.map((entry) => (
              <li key={entry.id} className="flex items-center gap-3 py-2.5 text-sm">
                <Badge tone="accent">{entry.action}</Badge>
                <span className="truncate text-muted">{entry.actorName || "система"}</span>
                <span className="ml-auto shrink-0 text-xs text-faint">
                  {dateFmt.format(new Date(entry.createdAt))}
                </span>
              </li>
            ))}
          </ul>
        )}
      </Card>
    </>
  );
}
