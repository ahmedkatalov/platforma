import { useEffect, useMemo, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";

import { useGetStudentCourseQuery } from "@/features/admin/api/coursesApi";
import {
  useGetLessonQuery,
  useStartLessonMutation,
} from "@/features/learning/api/lessonApi";
import type {
  CodeContent,
  LessonKind,
  LessonProgress,
  QuizContent,
  TerminalContent,
  TextContent,
} from "@/shared/types";
import { Badge, Button, Card, EmptyState, Progress, Spinner } from "@/shared/ui";
import {
  IconBook,
  IconCheck,
  IconChevron,
  IconClock,
  IconTerminal,
} from "@/shared/ui/icons";

import CodeLesson from "./CodeLesson";
import QuizLesson from "./QuizLesson";
import TerminalLesson from "./TerminalLesson";
import TextLesson from "./TextLesson";

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

// Плеер урока: слева содержание курса, справа сам урок нужного типа.
export default function LessonPage() {
  const { slug = "", lessonId = "" } = useParams();
  const navigate = useNavigate();

  const { data, isLoading, isError, error } = useGetLessonQuery(lessonId, { skip: !lessonId });
  const { data: courseData } = useGetStudentCourseQuery(slug, { skip: !slug });
  const [startLesson] = useStartLessonMutation();

  const [asideOpen, setAsideOpen] = useState(false);

  useEffect(() => {
    if (lessonId) void startLesson(lessonId);
    setAsideOpen(false);
    window.scrollTo({ top: 0 });
  }, [lessonId, startLesson]);

  const progressByLesson = useMemo(() => {
    const map = new Map<string, LessonProgress>();
    data?.progress.forEach((item) => map.set(item.lessonId, item));
    return map;
  }, [data]);

  const modules = courseData?.course.modules ?? [];

  const totals = useMemo(() => {
    const all = modules.flatMap((module) => module.lessons ?? []);
    const done = all.filter((lesson) => progressByLesson.get(lesson.id)?.status === "completed");
    return { total: all.length, done: done.length };
  }, [modules, progressByLesson]);

  if (isLoading) {
    return (
      <div className="grid place-items-center py-20 text-accent">
        <Spinner size={32} />
      </div>
    );
  }

  if (isError || !data) {
    const status = (error as { status?: number } | undefined)?.status;
    return (
      <Card>
        <EmptyState
          title={status === 403 ? "Курс вам ещё не открыт" : "Урок не найден"}
          description={
            status === 403
              ? "Попросите администратора назначить вам этот курс"
              : "Возможно, урок удалён или ещё не опубликован"
          }
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

  const { lesson } = data;
  const progress = progressByLesson.get(lesson.id);
  const completed = progress?.status === "completed";

  const goNext = () => {
    if (data.nextLessonId) {
      navigate(`/learn/courses/${data.courseSlug}/lessons/${data.nextLessonId}`);
    } else {
      navigate(`/learn/courses/${data.courseSlug}`);
    }
  };

  const content = lesson.content as Record<string, unknown>;

  const body = (() => {
    switch (lesson.kind) {
      case "quiz":
        return (
          <QuizLesson
            lessonId={lesson.id}
            content={content as unknown as QuizContent}
            progress={progress}
            onDone={goNext}
          />
        );
      case "terminal":
        return (
          <TerminalLesson
            lessonId={lesson.id}
            content={content as unknown as TerminalContent}
            tasks={data.tasks}
            progress={progress}
            onDone={goNext}
          />
        );
      case "code":
        return (
          <CodeLesson
            lessonId={lesson.id}
            content={content as unknown as CodeContent}
            progress={progress}
            onDone={goNext}
          />
        );
      default:
        return (
          <TextLesson
            lessonId={lesson.id}
            content={content as unknown as TextContent}
            progress={progress}
            onDone={goNext}
          />
        );
    }
  })();

  return (
    <div className="grid gap-[var(--gap)] xl:grid-cols-[18rem_1fr]">
      {/* Содержание курса */}
      <aside className={`xl:block ${asideOpen ? "block" : "hidden"}`}>
        <Card className="xl:sticky xl:top-24 p-3">
          <Link
            to={`/learn/courses/${data.courseSlug}`}
            className="mb-3 block rounded-[var(--radius-md)] p-2 transition-colors hover:bg-surface-2"
          >
            <p className="text-sm font-bold text-fg">{data.courseTitle}</p>
            <p className="mt-0.5 text-xs text-faint">
              {totals.done} из {totals.total} уроков пройдено
            </p>
            <div className="mt-2">
              <Progress value={totals.total ? (totals.done / totals.total) * 100 : 0} />
            </div>
          </Link>

          <nav className="max-h-[60vh] space-y-3 overflow-y-auto">
            {modules.map((module, moduleIndex) => (
              <div key={module.id}>
                <p className="mb-1 px-2 text-[11px] font-bold uppercase tracking-wide text-faint">
                  {moduleIndex + 1}. {module.title}
                </p>
                <ul className="space-y-0.5">
                  {(module.lessons ?? []).map((item) => {
                    const state = progressByLesson.get(item.id);
                    const active = item.id === lesson.id;

                    return (
                      <li key={item.id}>
                        <Link
                          to={`/learn/courses/${data.courseSlug}/lessons/${item.id}`}
                          className={`flex items-center gap-2 rounded-[var(--radius-md)] px-2 py-1.5 text-sm transition-colors ${
                            active
                              ? "bg-accent-soft font-semibold text-accent"
                              : "text-muted hover:bg-surface-2 hover:text-fg"
                          }`}
                        >
                          <span
                            className={`grid h-4 w-4 shrink-0 place-items-center rounded-full border text-[10px] ${
                              state?.status === "completed"
                                ? "border-[var(--success)] bg-[var(--success)] text-white"
                                : "border-line"
                            }`}
                          >
                            {state?.status === "completed" && <IconCheck size={10} />}
                          </span>
                          <span className="min-w-0 flex-1 truncate">{item.title}</span>
                          {item.kind !== "text" && (
                            <span className="shrink-0 text-[10px] text-faint">
                              {KIND_LABEL[item.kind].slice(0, 4)}
                            </span>
                          )}
                        </Link>
                      </li>
                    );
                  })}
                </ul>
              </div>
            ))}
          </nav>
        </Card>
      </aside>

      {/* Урок */}
      <div className="min-w-0">
        <div className="mb-4 flex flex-wrap items-center gap-2">
          <Button className="xl:hidden" onClick={() => setAsideOpen((v) => !v)}>
            <IconBook size={16} />
            Содержание
          </Button>
          <Badge tone={KIND_TONE[lesson.kind]}>{KIND_LABEL[lesson.kind]}</Badge>
          <Badge>
            <IconClock size={12} /> {lesson.durationMin} мин
          </Badge>
          {completed && <Badge tone="success">Пройдено</Badge>}
          <span className="text-xs text-faint">{data.moduleTitle}</span>
        </div>

        <h1 className="mb-1 text-2xl font-bold tracking-tight text-fg sm:text-3xl">
          {lesson.title}
        </h1>
        {lesson.summary && <p className="mb-6 text-sm text-muted">{lesson.summary}</p>}

        {body}

        <div className="mt-[var(--gap)] flex items-center justify-between gap-3">
          {data.prevLessonId ? (
            <Link
              to={`/learn/courses/${data.courseSlug}/lessons/${data.prevLessonId}`}
              className="btn btn-secondary"
            >
              <IconChevron size={16} className="rotate-180" />
              Предыдущий
            </Link>
          ) : (
            <span />
          )}

          {data.nextLessonId ? (
            <Link
              to={`/learn/courses/${data.courseSlug}/lessons/${data.nextLessonId}`}
              className="btn btn-secondary"
            >
              Следующий
              <IconChevron size={16} />
            </Link>
          ) : (
            <Link to={`/learn/courses/${data.courseSlug}`} className="btn btn-secondary">
              <IconTerminal size={16} />
              К программе курса
            </Link>
          )}
        </div>
      </div>
    </div>
  );
}
