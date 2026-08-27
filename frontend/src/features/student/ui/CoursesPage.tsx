import { Link } from "react-router-dom";

import { useGetStudentCoursesQuery } from "@/features/admin/api/coursesApi";
import type { CourseLevel } from "@/shared/types";
import { Badge, Card, EmptyState, PageHeader, Spinner } from "@/shared/ui";
import { IconBook } from "@/shared/ui/icons";

const LEVEL_LABEL: Record<CourseLevel, string> = {
  beginner: "Начальный",
  intermediate: "Средний",
  advanced: "Продвинутый",
};

export default function CoursesPage() {
  const { data: catalog = [], isLoading } = useGetStudentCoursesQuery();

  if (isLoading) {
    return (
      <div className="grid place-items-center py-20 text-accent">
        <Spinner size={32} />
      </div>
    );
  }

  return (
    <>
      <PageHeader title="Курсы" subtitle="Программа обучения DevOps" />

      {catalog.length === 0 ? (
        <Card>
          <EmptyState
            title="Курсов пока нет"
            description="Опубликованные курсы появятся здесь"
            icon={<IconBook size={32} />}
          />
        </Card>
      ) : (
        <div className="grid gap-[var(--gap)] md:grid-cols-2 xl:grid-cols-3">
          {catalog.map(({ course, enrolled }) => (
            <Link
              key={course.id}
              to={`/learn/courses/${course.slug}`}
              className="card flex flex-col p-[var(--pad)] transition-transform hover:-translate-y-0.5"
            >
              <div className="mb-3 flex items-start justify-between gap-2">
                <h2 className="text-base font-bold text-fg">{course.title}</h2>
                {enrolled ? <Badge tone="success">Доступен</Badge> : <Badge>Каталог</Badge>}
              </div>

              {course.subtitle && (
                <p className="mb-3 line-clamp-3 flex-1 text-sm text-muted">{course.subtitle}</p>
              )}

              <div className="mb-3 flex flex-wrap gap-1.5">
                <Badge>{LEVEL_LABEL[course.level]}</Badge>
                {course.tags.slice(0, 3).map((tag) => (
                  <Badge key={tag} tone="accent">
                    {tag}
                  </Badge>
                ))}
              </div>

              <p className="mt-auto text-xs text-faint">
                {course.modulesCount} модулей · {course.lessonsCount} уроков
              </p>
            </Link>
          ))}
        </div>
      )}
    </>
  );
}
