import { Link } from "react-router-dom";

import { useAppSelector } from "@/app/store";
import {
  useGetStudentCourseQuery,
  useGetStudentCoursesQuery,
} from "@/features/admin/api/coursesApi";
import { useGetMyStatsQuery } from "@/shared/api/meApi";
import type { LessonKind, LessonProgress } from "@/shared/types";
import { Badge, Card, EmptyState, PageHeader, Progress, Spinner, StatCard } from "@/shared/ui";
import { Book, BarChart3, Clock, Flame, Terminal, CheckCircle } from "lucide-react";

const KIND_LABEL: Record<LessonKind, string> = {
  text: "Теория",
  quiz: "Квиз",
  terminal: "Тренажёр",
  code: "Практика с кодом",
};

export default function DashboardPage() {
  const user = useAppSelector((state) => state.auth.user);
  const { data: stats, isLoading } = useGetMyStatsQuery(30);
  const { data: catalog = [] } = useGetStudentCoursesQuery();

  const myCourses = catalog.filter((item) => item.enrolled);
  const activeSlug = myCourses[0]?.course.slug ?? "";

  // Подтягиваем структуру первого курса, чтобы предложить следующий урок.
  const { data: activeCourse } = useGetStudentCourseQuery(activeSlug, { skip: !activeSlug });

  if (isLoading || !stats) {
    return (
      <div className="grid place-items-center py-20 text-accent">
        <Spinner size={32} />
      </div>
    );
  }

  const progressByLesson = new Map<string, LessonProgress>();
  (activeCourse?.progress ?? []).forEach((item) => progressByLesson.set(item.lessonId, item));

  const lessons = (activeCourse?.course.modules ?? []).flatMap((module) => module.lessons ?? []);
  const nextLesson = lessons.find(
    (lesson) => progressByLesson.get(lesson.id)?.status !== "completed",
  );

  const hour = new Date().getHours();
  const greeting = hour < 12 ? "Доброе утро" : hour < 18 ? "Добрый день" : "Добрый вечер";

  return (
    <>
      <PageHeader
        title={`${greeting}, ${user?.fullName?.split(" ")[0] || "студент"}!`}
        subtitle="Ваш прогресс по курсам DevOps"
        actions={
          <Link to="/learn/courses" className="btn btn-primary">
            <Book size={18} />
            К обучению
          </Link>
        }
      />

      <div className="grid gap-[var(--gap)] sm:grid-cols-2 xl:grid-cols-4">
        <StatCard
          label="Прогресс"
          value={`${Math.round(stats.summary.progress)}%`}
          hint={`${stats.summary.lessonsCompleted} из ${stats.summary.lessonsTotal} уроков`}
          icon={<BarChart3 size={20} />}
        />
        <StatCard
          label="Средний балл"
          value={stats.quiz.attempts > 0 ? `${Math.round(stats.quiz.averageScore)}%` : "—"}
          hint={`${stats.quiz.attempts} попыток в квизах`}
          icon={<CheckCircle size={20} />}
        />
        <StatCard
          label="Дней подряд"
          value={stats.streak}
          hint={`${stats.summary.daysVisited} дней всего`}
          icon={<Flame size={20} />}
        />
        <StatCard
          label="Время обучения"
          value={`${Math.floor(stats.summary.minutesSpent / 60)} ч`}
          hint={`${stats.summary.minutesSpent} минут`}
          icon={<Clock size={20} />}
        />
      </div>

      <Card className="mt-[var(--gap)] flex flex-wrap items-center justify-between gap-4 p-[var(--pad)]">
        <div className="flex min-w-0 items-start gap-3">
          <span
            className="grid h-11 w-11 shrink-0 place-items-center rounded-[var(--radius-md)] text-accent-fg"
            style={{ background: "var(--gradient)" }}
          >
            <Terminal size={22} />
          </span>
          <div className="min-w-0">
            <p className="text-base font-bold text-fg">Песочница — Linux-терминал</p>
            <p className="mt-0.5 text-sm text-muted">
              Настоящий Linux прямо в браузере: практикуйся свободно и выполняй задания с автопроверкой.
            </p>
          </div>
        </div>
        <Link to="/learn/sandbox" className="btn btn-primary shrink-0">
          <Terminal size={18} />
          Открыть песочницу
        </Link>
      </Card>

      {nextLesson && activeCourse && (
        <Card className="mt-[var(--gap)] flex flex-wrap items-center justify-between gap-4 p-[var(--pad)]">
          <div className="min-w-0">
            <p className="text-xs font-semibold uppercase tracking-wide text-faint">
              Продолжить обучение
            </p>
            <p className="mt-1 truncate text-base font-bold text-fg">{nextLesson.title}</p>
            <p className="mt-0.5 text-sm text-muted">
              {activeCourse.course.title} · {KIND_LABEL[nextLesson.kind]} ·{" "}
              {nextLesson.durationMin} мин
            </p>
          </div>
          <Link
            to={`/learn/courses/${activeCourse.course.slug}/lessons/${nextLesson.id}`}
            className="btn btn-primary"
          >
            <Terminal size={18} />
            Открыть урок
          </Link>
        </Card>
      )}

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
            icon={<Book size={32} />}
            action={
              <Link to="/learn/courses" className="btn btn-secondary">
                Посмотреть каталог
              </Link>
            }
          />
        ) : (
          <div className="grid gap-[var(--gap)] md:grid-cols-2">
            {myCourses.map(({ course, completedLessons }) => {
              const total = course.lessonsCount;
              const done = Math.min(completedLessons ?? 0, total);
              const percent = total > 0 ? (done / total) * 100 : 0;

              return (
                <Link
                  key={course.id}
                  to={`/learn/courses/${course.slug}`}
                  className="card-flat block p-4 transition-colors hover:bg-surface-hover"
                >
                  <div className="mb-2 flex items-start justify-between gap-2">
                    <h3 className="text-sm font-bold text-fg">{course.title}</h3>
                    {total > 0 && percent === 100 && <Badge tone="success">Пройден</Badge>}
                  </div>

                  {course.subtitle && (
                    <p className="mb-3 line-clamp-2 text-xs text-muted">{course.subtitle}</p>
                  )}

                  {total > 0 && (
                    <div className="mb-3">
                      <div className="mb-1 flex justify-between text-[11px] text-faint">
                        <span>
                          {done} из {total} уроков
                        </span>
                        <span className="font-bold text-accent">{Math.round(percent)}%</span>
                      </div>
                      <Progress value={percent} />
                    </div>
                  )}

                  <div className="flex flex-wrap gap-1.5">
                    {course.tags.slice(0, 3).map((tag) => (
                      <Badge key={tag} tone="accent">
                        {tag}
                      </Badge>
                    ))}
                  </div>
                </Link>
              );
            })}
          </div>
        )}
      </Card>
    </>
  );
}
