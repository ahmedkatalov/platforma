import { Link, useParams } from "react-router-dom";

import { useGetStudentCourseQuery } from "@/features/admin/api/coursesApi";
import type { LessonKind, LessonProgress } from "@/shared/types";
import { Badge, Card, EmptyState, PageHeader, Progress, Spinner } from "@/shared/ui";
import { IconBook, IconCheck, IconChevron, IconClock } from "@/shared/ui/icons";

const KIND_LABEL: Record<LessonKind, string> = {
  text: "Теория",
  quiz: "Квиз",
  terminal: "Терминал",
  code: "Код",
};

const KIND_TONE: Record<LessonKind, "default" | "accent" | "success" | "warning"> = {
  text: "default",
  quiz: "accent",
  terminal: "success",
  code: "warning",
};

export default function CoursePage() {
  const { slug = "" } = useParams();
  const { data, isLoading, isError } = useGetStudentCourseQuery(slug, { skip: !slug });

  if (isLoading) {
    return (
      <div className="grid place-items-center py-20 text-accent">
        <Spinner size={32} />
      </div>
    );
  }

  if (isError || !data) {
    return (
      <Card>
        <EmptyState
          title="Курс не найден"
          description="Возможно, он ещё не опубликован"
          icon={<IconBook size={32} />}
          action={
            <Link to="/learn/courses" className="btn btn-secondary">
              К списку курсов
            </Link>
          }
        />
      </Card>
    );
  }

  const { course, enrolled } = data;
  const modules = course.modules ?? [];

  const progressByLesson = new Map<string, LessonProgress>();
  (data.progress ?? []).forEach((item) => progressByLesson.set(item.lessonId, item));

  const allLessons = modules.flatMap((module) => module.lessons ?? []);
  const doneCount = allLessons.filter(
    (lesson) => progressByLesson.get(lesson.id)?.status === "completed",
  ).length;
  const coursePercent = allLessons.length ? (doneCount / allLessons.length) * 100 : 0;

  const totalMinutes = allLessons.reduce((sum, lesson) => sum + lesson.durationMin, 0);

  // Следующий непройденный урок — кнопка «Продолжить».
  const nextLesson = allLessons.find(
    (lesson) => progressByLesson.get(lesson.id)?.status !== "completed",
  );

  return (
    <>
      <PageHeader
        title={course.title}
        subtitle={course.subtitle}
        actions={
          <>
            <Link to="/learn/courses" className="btn btn-ghost">
              К списку курсов
            </Link>
            {enrolled && nextLesson && (
              <Link
                to={`/learn/courses/${course.slug}/lessons/${nextLesson.id}`}
                className="btn btn-primary"
              >
                {doneCount > 0 ? "Продолжить" : "Начать обучение"}
                <IconChevron size={16} />
              </Link>
            )}
          </>
        }
      />

      <div className="mb-[var(--gap)] flex flex-wrap items-center gap-2">
        {enrolled ? (
          <Badge tone="success">Курс доступен вам</Badge>
        ) : (
          <Badge tone="warning">Нужен доступ от администратора</Badge>
        )}
        {course.tags.map((tag) => (
          <Badge key={tag} tone="accent">
            {tag}
          </Badge>
        ))}
        <Badge>
          <IconClock size={12} /> ~{Math.max(1, Math.round(totalMinutes / 60))} ч
        </Badge>
      </div>

      {enrolled && allLessons.length > 0 && (
        <Card className="mb-[var(--gap)] p-[var(--pad)]">
          <div className="mb-2 flex items-center justify-between">
            <span className="text-sm font-semibold text-fg">Ваш прогресс</span>
            <span className="text-sm font-bold text-accent">
              {doneCount} из {allLessons.length} · {Math.round(coursePercent)}%
            </span>
          </div>
          <Progress value={coursePercent} />
        </Card>
      )}

      {course.description && (
        <Card className="mb-[var(--gap)] p-[var(--pad)]">
          <p className="whitespace-pre-line text-sm text-muted">{course.description}</p>
        </Card>
      )}

      <Card className="p-[var(--pad)]">
        <h2 className="mb-4 text-base font-bold text-fg">Программа курса</h2>

        {modules.length === 0 ? (
          <EmptyState title="Программа ещё готовится" icon={<IconBook size={32} />} />
        ) : (
          <div className="space-y-3">
            {modules.map((module, index) => {
              const lessons = module.lessons ?? [];
              const moduleDone = lessons.filter(
                (lesson) => progressByLesson.get(lesson.id)?.status === "completed",
              ).length;

              return (
                <div key={module.id} className="card-flat overflow-hidden">
                  <div className="flex items-center gap-3 border-b border-line p-3">
                    <div className="min-w-0 flex-1">
                      <p className="text-sm font-bold text-fg">
                        <span className="text-faint">Модуль {index + 1}. </span>
                        {module.title}
                      </p>
                      {module.summary && (
                        <p className="mt-0.5 text-xs text-muted">{module.summary}</p>
                      )}
                    </div>
                    {enrolled && lessons.length > 0 && (
                      <Badge tone={moduleDone === lessons.length ? "success" : "default"}>
                        {moduleDone}/{lessons.length}
                      </Badge>
                    )}
                  </div>

                  <ul className="divide-y divide-[var(--border)]">
                    {lessons.map((lesson) => {
                      const state = progressByLesson.get(lesson.id);
                      const done = state?.status === "completed";

                      const row = (
                        <>
                          <span
                            className={`grid h-5 w-5 shrink-0 place-items-center rounded-full border ${
                              done
                                ? "border-[var(--success)] bg-[var(--success)] text-white"
                                : "border-line text-faint"
                            }`}
                          >
                            {done && <IconCheck size={12} />}
                          </span>
                          <Badge tone={KIND_TONE[lesson.kind]}>{KIND_LABEL[lesson.kind]}</Badge>
                          <span
                            className={`min-w-0 flex-1 truncate text-sm ${
                              done ? "text-muted" : "text-fg"
                            }`}
                          >
                            {lesson.title}
                          </span>
                          {state?.bestScore != null && lesson.kind !== "text" && (
                            <span className="shrink-0 text-xs font-semibold text-accent">
                              {Math.round(state.bestScore)}%
                            </span>
                          )}
                          <span className="shrink-0 text-xs text-faint">
                            {lesson.durationMin} мин
                          </span>
                        </>
                      );

                      return (
                        <li key={lesson.id}>
                          {enrolled ? (
                            <Link
                              to={`/learn/courses/${course.slug}/lessons/${lesson.id}`}
                              className="flex items-center gap-3 px-3 py-2.5 transition-colors hover:bg-surface-hover"
                            >
                              {row}
                              <IconChevron size={16} className="shrink-0 text-faint" />
                            </Link>
                          ) : (
                            <div className="flex items-center gap-3 px-3 py-2.5 opacity-70">{row}</div>
                          )}
                        </li>
                      );
                    })}

                    {lessons.length === 0 && (
                      <li className="px-3 py-3 text-center text-xs text-faint">
                        Уроки скоро появятся
                      </li>
                    )}
                  </ul>
                </div>
              );
            })}
          </div>
        )}
      </Card>
    </>
  );
}
