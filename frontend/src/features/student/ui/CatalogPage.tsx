import { Link } from "react-router-dom";
import { ArrowRight, Book, Clock, Layers, Lock } from "lucide-react";

import {
  useGetStudentCoursesQuery,
  useRequestCourseAccessMutation,
} from "@/features/admin/api/coursesApi";
import { apiErrorMessage } from "@/shared/api/baseApi";
import type { Course, CourseLevel } from "@/shared/types";
import { Badge, Button, Card, EmptyState, PageHeader, Progress, Spinner } from "@/shared/ui";
import { useToast } from "@/shared/ui/ToastProvider";

const LEVEL_LABEL: Record<CourseLevel, string> = {
  beginner: "Начальный",
  intermediate: "Средний",
  advanced: "Продвинутый",
};

type CatalogItem = {
  course: Course;
  enrolled: boolean;
  completedLessons: number;
  requestStatus?: string;
};

export default function CatalogPage() {
  const { data: catalog = [], isLoading } = useGetStudentCoursesQuery();
  const [requestAccess, { isLoading: requesting }] = useRequestCourseAccessMutation();
  const toast = useToast();

  const ask = async (courseId: string) => {
    try {
      await requestAccess({ courseId }).unwrap();
      toast.success("Заявка отправлена — администратор откроет курс");
    } catch (err) {
      toast.error(apiErrorMessage(err, "Не удалось отправить заявку"));
    }
  };

  if (isLoading) {
    return (
      <div className="grid place-items-center py-20 text-accent">
        <Spinner size={32} />
      </div>
    );
  }

  const openCount = catalog.filter((c) => c.enrolled).length;

  return (
    <>
      <PageHeader
        title="Наши курсы"
        subtitle="Выбирайте курс и учитесь на практике. Закрытые курсы можно запросить у преподавателя."
        actions={
          catalog.length > 0 ? (
            <Badge tone="accent">
              {openCount} из {catalog.length} открыто
            </Badge>
          ) : undefined
        }
      />

      {catalog.length === 0 ? (
        <Card>
          <EmptyState
            title="Курсов пока нет"
            description="Опубликованные курсы появятся здесь"
            icon={<Book size={32} />}
          />
        </Card>
      ) : (
        <div className="grid gap-[var(--gap)] sm:grid-cols-2 xl:grid-cols-3">
          {catalog.map((item) => (
            <CourseCard key={item.course.id} item={item} onRequest={ask} requesting={requesting} />
          ))}
        </div>
      )}
    </>
  );
}

function CourseCard({
  item,
  onRequest,
  requesting,
}: {
  item: CatalogItem;
  onRequest: (courseId: string) => void;
  requesting: boolean;
}) {
  const { course, enrolled, completedLessons, requestStatus } = item;
  const locked = !enrolled;
  const pending = requestStatus === "pending";
  const percent = enrolled && course.lessonsCount ? (completedLessons / course.lessonsCount) * 100 : 0;
  const done = percent >= 100;

  const cover = (
    <div className="relative h-32 overflow-hidden">
      {course.coverUrl ? (
        <img
          src={course.coverUrl}
          alt=""
          className={`h-full w-full object-cover transition-transform duration-500 group-hover:scale-110 ${
            locked ? "opacity-70 grayscale-[35%]" : ""
          }`}
        />
      ) : (
        <div
          className={`h-full w-full transition-transform duration-500 group-hover:scale-110 ${
            locked ? "opacity-80" : ""
          }`}
          style={{ background: "var(--gradient)" }}
        />
      )}
      <div className="absolute inset-0 bg-gradient-to-t from-black/70 via-black/10 to-transparent" />

      <div className="absolute left-3 top-3">
        <span className="rounded-full bg-black/45 px-2 py-0.5 text-[11px] font-semibold text-white backdrop-blur-sm">
          {LEVEL_LABEL[course.level]}
        </span>
      </div>
      {locked ? (
        <span className="absolute right-3 top-3 grid h-9 w-9 place-items-center rounded-full bg-black/50 text-white backdrop-blur-sm">
          <Lock size={16} />
        </span>
      ) : done ? (
        <span className="absolute right-3 top-3">
          <Badge tone="success">Пройден</Badge>
        </span>
      ) : null}

      <h2 className="absolute inset-x-3 bottom-2 line-clamp-2 text-base font-bold text-white drop-shadow-md">
        {course.title}
      </h2>
    </div>
  );

  const body = (
    <div className="flex flex-1 flex-col gap-3 p-4">
      {course.subtitle && (
        <p className="line-clamp-2 min-h-[2.5rem] text-sm text-muted">{course.subtitle}</p>
      )}

      {course.tags.length > 0 && (
        <div className="flex flex-wrap gap-1.5">
          {course.tags.slice(0, 3).map((tag) => (
            <Badge key={tag} tone="accent">
              {tag}
            </Badge>
          ))}
        </div>
      )}

      <div className="flex items-center gap-4 text-xs text-faint">
        <span className="flex items-center gap-1">
          <Layers size={13} /> {course.modulesCount} глав
        </span>
        <span className="flex items-center gap-1">
          <Book size={13} /> {course.lessonsCount} уроков
        </span>
      </div>

      {enrolled && course.lessonsCount > 0 && (
        <div>
          <div className="mb-1 flex justify-between text-[11px] text-faint">
            <span>
              {completedLessons} из {course.lessonsCount}
            </span>
            <span className="font-bold text-accent">{Math.round(percent)}%</span>
          </div>
          <Progress value={percent} />
        </div>
      )}

      <div className="mt-auto pt-1">
        {enrolled ? (
          <span className="btn btn-primary pointer-events-none w-full">
            {done ? "Открыть курс" : completedLessons > 0 ? "Продолжить" : "Начать обучение"}
            <ArrowRight size={16} />
          </span>
        ) : pending ? (
          <span className="flex w-full items-center justify-center gap-2 rounded-[var(--radius-md)] bg-[var(--warning-soft)] px-3 py-2 text-sm font-semibold text-warning">
            <Clock size={15} /> Заявка на рассмотрении
          </span>
        ) : (
          <Button
            variant="primary"
            className="w-full"
            loading={requesting}
            onClick={() => onRequest(course.id)}
          >
            <Lock size={15} /> Запросить доступ
          </Button>
        )}
      </div>
    </div>
  );

  const cardCls =
    "group relative flex flex-col overflow-hidden rounded-[var(--radius-lg)] border border-line bg-surface shadow-[var(--shadow-md)] transition-all duration-300 hover:-translate-y-1.5 hover:border-accent-border hover:shadow-[0_22px_55px_-24px_var(--accent)]";

  return enrolled ? (
    <Link to={`/learn/courses/${course.slug}`} className={cardCls}>
      {cover}
      {body}
    </Link>
  ) : (
    <div className={cardCls}>
      {cover}
      {body}
    </div>
  );
}
