import { Link, useParams } from "react-router-dom";

import { useGetStudentCourseQuery } from "@/features/admin/api/coursesApi";
import type { LessonKind } from "@/shared/types";
import { Badge, Card, EmptyState, PageHeader, Spinner } from "@/shared/ui";
import { IconBook, IconChevron, IconClock, IconTerminal } from "@/shared/ui/icons";

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
  const totalMinutes = modules.reduce(
    (sum, module) => sum + (module.lessons ?? []).reduce((s, l) => s + l.durationMin, 0),
    0,
  );

  return (
    <>
      <PageHeader
        title={course.title}
        subtitle={course.subtitle}
        actions={
          <Link to="/learn/courses" className="btn btn-ghost">
            К списку курсов
          </Link>
        }
      />

      <div className="mb-[var(--gap)] flex flex-wrap items-center gap-2">
        {enrolled ? <Badge tone="success">Курс доступен вам</Badge> : <Badge tone="warning">Нужен доступ от администратора</Badge>}
        {course.tags.map((tag) => (
          <Badge key={tag} tone="accent">
            {tag}
          </Badge>
        ))}
        <Badge>
          <IconClock size={12} /> ~{Math.round(totalMinutes / 60)} ч
        </Badge>
      </div>

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
            {modules.map((module, index) => (
              <div key={module.id} className="card-flat overflow-hidden">
                <div className="border-b border-line p-3">
                  <p className="text-sm font-bold text-fg">
                    <span className="text-faint">Модуль {index + 1}. </span>
                    {module.title}
                  </p>
                  {module.summary && <p className="mt-0.5 text-xs text-muted">{module.summary}</p>}
                </div>

                <ul className="divide-y divide-[var(--border)]">
                  {(module.lessons ?? []).map((lesson) => (
                    <li key={lesson.id} className="flex items-center gap-3 px-3 py-2.5">
                      <Badge tone={KIND_TONE[lesson.kind]}>{KIND_LABEL[lesson.kind]}</Badge>
                      <span className="min-w-0 flex-1 truncate text-sm text-fg">{lesson.title}</span>
                      <span className="shrink-0 text-xs text-faint">{lesson.durationMin} мин</span>
                      <IconChevron size={16} className="shrink-0 text-faint" />
                    </li>
                  ))}
                  {(module.lessons ?? []).length === 0 && (
                    <li className="px-3 py-3 text-center text-xs text-faint">Уроки скоро появятся</li>
                  )}
                </ul>
              </div>
            ))}
          </div>
        )}
      </Card>

      <Card className="mt-[var(--gap)] flex items-center gap-3 p-[var(--pad)] text-sm text-muted">
        <IconTerminal size={20} className="shrink-0 text-accent" />
        Прохождение уроков — квизы, тренажёр терминала и редактор кода — подключим на следующем этапе.
      </Card>
    </>
  );
}
