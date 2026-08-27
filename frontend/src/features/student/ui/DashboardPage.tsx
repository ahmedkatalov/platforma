import { Link } from "react-router-dom";

import { useAppSelector } from "@/app/store";
import { useGetStudentCoursesQuery } from "@/features/admin/api/coursesApi";
import { useGetMeQuery, useGetMyStatsQuery } from "@/shared/api/meApi";
import { Badge, Button, Card, EmptyState, PageHeader, Spinner, StatCard } from "@/shared/ui";
import { IconBook, IconChart, IconClock, IconFlame } from "@/shared/ui/icons";

export default function DashboardPage() {
  const user = useAppSelector((state) => state.auth.user);
  const { data: me } = useGetMeQuery();
  const { data: stats, isLoading } = useGetMyStatsQuery(30);
  const { data: catalog = [] } = useGetStudentCoursesQuery();

  if (isLoading || !stats) {
    return (
      <div className="grid place-items-center py-20 text-accent">
        <Spinner size={32} />
      </div>
    );
  }

  const enrollments = me?.enrollments ?? [];
  const myCourses = catalog.filter((item) => item.enrolled);

  const greeting = new Date().getHours() < 12 ? "Доброе утро" : new Date().getHours() < 18 ? "Добрый день" : "Добрый вечер";

  return (
    <>
      <PageHeader
        title={`${greeting}, ${user?.fullName?.split(" ")[0] || "студент"}!`}
        subtitle="Ваш прогресс по курсам DevOps"
        actions={
          <Link to="/learn/courses" className="btn btn-primary">
            <IconBook size={18} />
            К обучению
          </Link>
        }
      />

      <div className="grid gap-[var(--gap)] sm:grid-cols-2 xl:grid-cols-4">
        <StatCard
          label="Прогресс"
          value={`${Math.round(stats.summary.progress)}%`}
          hint={`${stats.summary.lessonsCompleted} из ${stats.summary.lessonsTotal} уроков`}
          icon={<IconChart size={20} />}
        />
        <StatCard
          label="Курсов"
          value={stats.summary.courses}
          hint="Назначено вам"
          icon={<IconBook size={20} />}
        />
        <StatCard
          label="Дней подряд"
          value={stats.streak}
          hint={`${stats.summary.daysVisited} дней всего`}
          icon={<IconFlame size={20} />}
        />
        <StatCard
          label="Время обучения"
          value={`${Math.floor(stats.summary.minutesSpent / 60)} ч`}
          hint={`${stats.summary.minutesSpent} минут`}
          icon={<IconClock size={20} />}
        />
      </div>

      <Card className="mt-[var(--gap)] p-[var(--pad)]">
        <div className="mb-4 flex items-center justify-between">
          <h2 className="text-base font-bold text-fg">Мои курсы</h2>
          <Link to="/learn/courses" className="text-xs font-semibold text-accent hover:underline">
            Все курсы
          </Link>
        </div>

        {myCourses.length === 0 ? (
          <EmptyState
            title="Курсы ещё не назначены"
            description="Как только администратор откроет вам курс, он появится здесь"
            icon={<IconBook size={32} />}
            action={
              <Link to="/learn/courses" className="btn btn-secondary">
                Посмотреть каталог
              </Link>
            }
          />
        ) : (
          <div className="grid gap-[var(--gap)] md:grid-cols-2">
            {myCourses.map(({ course }) => {
              const enrollment = enrollments.find((e) => e.courseId === course.id);
              return (
                <Link
                  key={course.id}
                  to={`/learn/courses/${course.slug}`}
                  className="card-flat block p-4 transition-colors hover:bg-surface-hover"
                >
                  <div className="mb-2 flex items-start justify-between gap-2">
                    <h3 className="text-sm font-bold text-fg">{course.title}</h3>
                    {enrollment?.status === "completed" && <Badge tone="success">Пройден</Badge>}
                  </div>
                  {course.subtitle && (
                    <p className="mb-3 line-clamp-2 text-xs text-muted">{course.subtitle}</p>
                  )}
                  <div className="mb-3 flex flex-wrap gap-1.5">
                    {course.tags.slice(0, 3).map((tag) => (
                      <Badge key={tag} tone="accent">
                        {tag}
                      </Badge>
                    ))}
                  </div>
                  <p className="text-xs text-faint">
                    {course.modulesCount} модулей · {course.lessonsCount} уроков
                  </p>
                </Link>
              );
            })}
          </div>
        )}
      </Card>

      <Card className="mt-[var(--gap)] flex flex-wrap items-center justify-between gap-4 p-[var(--pad)]">
        <div>
          <h2 className="text-base font-bold text-fg">Практика в терминале</h2>
          <p className="mt-1 text-sm text-muted">
            Тренажёр команд и редактор кода откроются внутри уроков курса.
          </p>
        </div>
        <Button variant="secondary" disabled>
          Скоро
        </Button>
      </Card>
    </>
  );
}
