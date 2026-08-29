import { Link, useParams } from "react-router-dom";

import { useGetStudentCourseQuery } from "@/features/admin/api/coursesApi";
import { groupThemes, themeProgress, type Theme } from "@/features/learning/lib/themes";
import { useGetMeQuery } from "@/shared/api/meApi";
import type { Lesson, LessonKind, LessonProgress } from "@/shared/types";
import { Badge, Card, EmptyState, PageHeader, Progress, Spinner } from "@/shared/ui";
import {
  IconBook,
  IconCheck,
  IconChevron,
  IconClock,
  IconEdit,
  IconTerminal,
} from "@/shared/ui/icons";

const dueFmt = new Intl.DateTimeFormat("ru-RU", { day: "2-digit", month: "long" });

function pageIcon(kind: LessonKind) {
  if (kind === "terminal") return <IconTerminal size={15} />;
  if (kind === "code") return <IconEdit size={15} />;
  if (kind === "quiz") return <IconCheck size={15} />;
  return <IconBook size={15} />;
}

export default function CoursePage() {
  const { slug = "" } = useParams();
  const { data, isLoading, isError } = useGetStudentCourseQuery(slug, { skip: !slug });
  const { data: me } = useGetMeQuery();

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

  const enrollment = me?.enrollments.find((item) => item.courseId === course.id);
  const deadline = enrollment?.dueDate
    ? (() => {
        const due = new Date(enrollment.dueDate);
        const days = Math.ceil((due.getTime() - Date.now()) / 86_400_000);
        return { label: dueFmt.format(due), overdue: days < 0, soon: days >= 0 && days <= 3 };
      })()
    : null;

  const progressByLesson = new Map<string, LessonProgress>();
  (data.progress ?? []).forEach((item) => progressByLesson.set(item.lessonId, item));
  const isDone = (id: string) => progressByLesson.get(id)?.status === "completed";

  const allLessons = modules.flatMap((module) => module.lessons ?? []);
  const doneCount = allLessons.filter((l) => isDone(l.id)).length;
  const coursePercent = allLessons.length ? (doneCount / allLessons.length) * 100 : 0;
  const totalMinutes = allLessons.reduce((sum, l) => sum + l.durationMin, 0);
  const nextLesson = allLessons.find((l) => !isDone(l.id));

  const lessonHref = (lesson: Lesson) => `/learn/courses/${course.slug}/lessons/${lesson.id}`;

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
              <Link to={lessonHref(nextLesson)} className="btn btn-primary">
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
        {deadline && (
          <Badge tone={deadline.overdue ? "danger" : deadline.soon ? "warning" : "default"}>
            <IconClock size={12} />
            {deadline.overdue ? `срок истёк ${deadline.label}` : `сдать до ${deadline.label}`}
          </Badge>
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

      {modules.length === 0 ? (
        <Card>
          <EmptyState title="Программа ещё готовится" icon={<IconBook size={32} />} />
        </Card>
      ) : (
        <div className="space-y-[var(--gap)]">
          {modules.map((module, index) => {
            const themes = groupThemes(module);
            const lessons = module.lessons ?? [];
            const moduleDone = lessons.filter((l) => isDone(l.id)).length;
            const modulePercent = lessons.length ? (moduleDone / lessons.length) * 100 : 0;

            return (
              <Card key={module.id} className="overflow-hidden p-0">
                {/* Шапка главы */}
                <div className="flex items-center gap-4 border-b border-line p-[var(--pad)]">
                  <span
                    className="grid h-11 w-11 shrink-0 place-items-center rounded-[var(--radius-md)] text-lg font-extrabold text-accent-fg"
                    style={{ background: "var(--gradient)" }}
                  >
                    {index + 1}
                  </span>
                  <div className="min-w-0 flex-1">
                    <p className="text-[11px] font-bold uppercase tracking-wide text-faint">
                      Глава {index + 1}
                    </p>
                    <h2 className="truncate text-base font-bold text-fg sm:text-lg">
                      {module.title}
                    </h2>
                    {module.summary && (
                      <p className="mt-0.5 line-clamp-1 text-xs text-muted">{module.summary}</p>
                    )}
                  </div>
                  <div className="hidden w-40 shrink-0 sm:block">
                    <div className="mb-1 flex justify-between text-[11px] text-faint">
                      <span>{themes.length} тем</span>
                      {enrolled && (
                        <span className="font-bold text-accent">{Math.round(modulePercent)}%</span>
                      )}
                    </div>
                    {enrolled && <Progress value={modulePercent} />}
                  </div>
                </div>

                {/* Темы главы */}
                <div className="space-y-[var(--gap)] p-[var(--pad)]">
                  {themes.map((theme, ti) => (
                    <ThemeCard
                      key={theme.key}
                      theme={theme}
                      index={ti + 1}
                      enrolled={enrolled}
                      isDone={isDone}
                      progressByLesson={progressByLesson}
                      lessonHref={lessonHref}
                    />
                  ))}
                </div>
              </Card>
            );
          })}
        </div>
      )}
    </>
  );
}

function ThemeCard({
  theme,
  index,
  enrolled,
  isDone,
  progressByLesson,
  lessonHref,
}: {
  theme: Theme;
  index: number;
  enrolled: boolean;
  isDone: (id: string) => boolean;
  progressByLesson: Map<string, LessonProgress>;
  lessonHref: (lesson: Lesson) => string;
}) {
  const prog = themeProgress(theme, isDone);
  const complete = prog.total > 0 && prog.done === prog.total;

  const Row = ({ lesson, isQuiz }: { lesson: Lesson; isQuiz: boolean }) => {
    const state = progressByLesson.get(lesson.id);
    const done = state?.status === "completed";
    const inner = (
      <>
        <span
          className={`grid h-6 w-6 shrink-0 place-items-center rounded-full border text-[11px] ${
            done
              ? "border-[var(--success)] bg-[var(--success)] text-white"
              : isQuiz
                ? "border-accent-border bg-accent-soft text-accent"
                : "border-line text-faint"
          }`}
        >
          {done ? <IconCheck size={13} /> : pageIcon(lesson.kind)}
        </span>
        <span className={`min-w-0 flex-1 truncate text-sm ${done ? "text-muted" : "text-fg"}`}>
          {lesson.title}
        </span>
        {state?.bestScore != null && lesson.kind !== "text" && (
          <span className="shrink-0 text-xs font-semibold text-accent">
            {Math.round(state.bestScore)}%
          </span>
        )}
        <span className="hidden shrink-0 text-xs text-faint sm:inline">{lesson.durationMin} мин</span>
        {enrolled && <IconChevron size={15} className="shrink-0 text-faint" />}
      </>
    );

    const cls = `flex items-center gap-3 rounded-[var(--radius-md)] px-2.5 py-2 ${
      isQuiz ? "bg-accent-soft/40" : ""
    } ${enrolled ? "transition-colors hover:bg-surface-hover" : "opacity-70"}`;

    return enrolled ? (
      <Link to={lessonHref(lesson)} className={cls}>
        {inner}
      </Link>
    ) : (
      <div className={cls}>{inner}</div>
    );
  };

  return (
    <div className="card-flat overflow-hidden">
      <div className="flex items-center gap-3 border-b border-line px-3 py-2.5">
        <span className="grid h-6 w-6 shrink-0 place-items-center rounded-full bg-surface-2 text-xs font-bold text-muted">
          {index}
        </span>
        <p className="min-w-0 flex-1 truncate text-sm font-bold text-fg">Тема: {theme.title}</p>
        {enrolled && (
          <Badge tone={complete ? "success" : "default"}>
            {prog.done}/{prog.total}
          </Badge>
        )}
      </div>

      <div className="space-y-0.5 p-2">
        {theme.pages.map((page) => (
          <Row key={page.id} lesson={page} isQuiz={false} />
        ))}

        {theme.quiz && (
          <div className="mt-1 border-t border-line pt-1">
            <p className="px-2.5 pb-1 pt-1 text-[11px] font-bold uppercase tracking-wide text-accent">
              Проверка темы
            </p>
            <Row lesson={theme.quiz} isQuiz />
          </div>
        )}
      </div>
    </div>
  );
}
