import { useEffect, useMemo, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";

import { useGetStudentCourseQuery } from "@/features/admin/api/coursesApi";
import {
  useGetLessonQuery,
  useStartLessonMutation,
} from "@/features/learning/api/lessonApi";
import type {
  Certificate,
  CodeContent,
  Lesson,
  LessonKind,
  LessonProgress,
  QuizContent,
  TerminalContent,
  TextContent,
} from "@/shared/types";
import { Badge, Button, Card, EmptyState, Modal, Progress, Spinner } from "@/shared/ui";
import {
  Book,
  Check,
  ChevronLeft,
  ChevronRight,
  Clock,
  Terminal,
} from "lucide-react";

import { groupThemes } from "@/features/learning/lib/themes";

import CodeLesson from "./CodeLesson";
import NoteSelection from "./NoteSelection";
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
  const [certificate, setCertificate] = useState<Certificate | null>(null);
  const [openModuleId, setOpenModuleId] = useState<string | null>(null);
  const [contentsCollapsed, setContentsCollapsed] = useState(false);

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

  // При переходе к уроку из другой главы раскрываем именно её.
  useEffect(() => {
    const activeLessonId = data?.lesson.id;
    if (!activeLessonId) return;

    const activeModule = modules.find((module) =>
      module.lessons?.some((item) => item.id === activeLessonId),
    );
    if (activeModule) setOpenModuleId(activeModule.id);
  }, [data?.lesson.id, modules]);

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
          icon={<Book size={32} />}
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

  // Если урок закрыл курс целиком — сначала показываем сертификат.
  const goNext = (issued?: Certificate | null) => {
    if (issued) {
      setCertificate(issued);
      return;
    }
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
    <div
      className={`grid gap-[var(--gap)] transition-[grid-template-columns] duration-200 ${
        contentsCollapsed ? "xl:grid-cols-[4.5rem_1fr]" : "xl:grid-cols-[18rem_1fr]"
      }`}
    >
      {/* Содержание курса */}
      <aside className={`xl:block ${asideOpen ? "block" : "hidden"}`}>

        <Card className={`relative xl:sticky xl:top-24 p-3 ${contentsCollapsed ? "xl:p-2" : "xl:p-3"}`}>
          {contentsCollapsed && (
            <button
              type="button"
              className="hidden h-16 w-full flex cursor-pointer items-center justify-center gap-1 rounded-[var(--radius-md)] bg-accent-soft text-accent transition-colors hover:bg-surface-2 xl:flex"
              onClick={() => setContentsCollapsed(false)}
              aria-label="Развернуть содержание"
              title="Содержание курса"
            >
              <Book size={18} />
              <ChevronRight size={14} />
            </button>
          )}
          <button
            className={`btn btn-ghost absolute right-2 top-2 hidden h-8 w-8 p-0! ${
              contentsCollapsed ? "xl:hidden" : "xl:inline-flex"
            }`}
            onClick={() => setContentsCollapsed((value) => !value)}
            aria-label={contentsCollapsed ? "Развернуть содержание" : "Свернуть содержание"}
            title={contentsCollapsed ? "Развернуть содержание" : "Свернуть содержание"}
          >
            {contentsCollapsed ? <ChevronRight size={18} /> : <ChevronLeft size={18} />}
          </button>

          <div className={contentsCollapsed ? "xl:hidden" : undefined}>
            <Link
              to={`/learn/courses/${data.courseSlug}`}
              className="mb-3 block rounded-[var(--radius-md)] p-2 pr-10 transition-colors hover:bg-surface-2"
            >
              <p className="text-sm font-bold text-fg">{data.courseTitle}</p>
              <p className="mt-0.5 text-xs text-faint">
                {totals.done} из {totals.total} уроков пройдено
              </p>
              <div className="mt-2">
                <Progress value={totals.total ? (totals.done / totals.total) * 100 : 0} />
              </div>
            </Link>

            <nav className="max-h-[60vh] space-y-4 overflow-y-auto">
            {modules.map((module, moduleIndex) => {
              const themes = groupThemes(module);
              const expanded = openModuleId === module.id;

              const renderRow = (item: Lesson, isQuiz: boolean) => {
                const state = progressByLesson.get(item.id);
                const active = item.id === lesson.id;
                const done = state?.status === "completed";
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
                          done
                            ? "border-[var(--success)] bg-[var(--success)] text-white"
                            : isQuiz
                              ? "border-accent-border text-accent"
                              : "border-line"
                        }`}
                      >
                        {done && <Check size={10} />}
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
              };

              return (
                <section key={module.id} className="rounded-[var(--radius-md)] border border-transparent hover:border-line">
                  <button
                    type="button"
                    className="flex w-full cursor-pointer items-center gap-2 rounded-[var(--radius-md)] px-2 py-2 text-left transition-colors hover:bg-surface-2"
                    onClick={() => setOpenModuleId((current) => (current === module.id ? null : module.id))}
                    aria-expanded={expanded}
                    aria-controls={`module-${module.id}`}
                  >
                    <span className="min-w-0 flex-1 text-[11px] font-bold uppercase tracking-wide text-faint">
                      Глава {moduleIndex + 1}. {module.title}
                    </span>
                    <ChevronRight
                      size={16}
                      className={`shrink-0 text-faint transition-transform duration-200 ${expanded ? "rotate-90" : ""}`}
                      aria-hidden="true"
                    />
                  </button>

                  {expanded && (
                    <div id={`module-${module.id}`} className="space-y-2 px-1 pb-1.5">
                      {themes.map((theme, ti) => (
                        <div key={theme.key} className="rounded-[var(--radius-md)] bg-surface-2/40 p-1.5">
                          <p className="truncate px-1.5 pb-1 text-[16px] font-semibold text-fg">
                            <span className="text-faint">{ti + 1}. </span>
                            {theme.title}
                          </p>
                          <ul className="space-y-0.5">
                            {theme.pages.map((page) => renderRow(page, false))}
                            {theme.quiz && (
                              <>
                                <li className="px-2 pt-0.5 text-[10px] font-bold uppercase tracking-wide text-accent">
                                  Проверка темы
                                </li>
                                {renderRow(theme.quiz, true)}
                              </>
                            )}
                          </ul>
                        </div>
                      ))}
                    </div>
                  )}
                </section>
              );
            })}
            </nav>
          </div>
        </Card>
      </aside>

      {/* Урок */}
      <Modal
        open={Boolean(certificate)}
        onClose={() => setCertificate(null)}
        title="Курс пройден!"
        footer={
          <>
            <Button onClick={() => setCertificate(null)}>Закрыть</Button>
            {certificate && (
              <a
                href={`/certificates/${certificate.serial}`}
                target="_blank"
                rel="noreferrer"
                className="btn btn-primary"
              >
                Открыть сертификат
              </a>
            )}
          </>
        }
      >
        {certificate && (
          <div className="space-y-4 text-center">
            <span
              className="mx-auto grid h-16 w-16 place-items-center rounded-full text-accent-fg"
              style={{ background: "var(--gradient)" }}
            >
              <Check size={32} />
            </span>

            <div>
              <p className="text-lg font-bold text-fg">{certificate.courseTitle}</p>
              <p className="mt-1 text-sm text-muted">
                Пройдено {certificate.lessonsCompleted} из {certificate.lessonsTotal} уроков ·
                средний балл {Math.round(certificate.score)}%
              </p>
            </div>

            <div className="card-flat p-3">
              <p className="text-xs text-faint">Номер сертификата</p>
              <p className="font-mono text-base font-bold text-accent">{certificate.serial}</p>
            </div>

            <p className="text-xs text-muted">
              Ссылку на сертификат можно отправить работодателю — подлинность проверяется по номеру.
            </p>
          </div>
        )}
      </Modal>

      <div className="min-w-0">
        <div className="mb-4 flex flex-wrap items-center gap-2">
          <Button className="xl:hidden" onClick={() => setAsideOpen((v) => !v)}>
            <Book size={16} />
            Содержание
          </Button>
          <Badge tone={KIND_TONE[lesson.kind]}>{KIND_LABEL[lesson.kind]}</Badge>
          <Badge>
            <Clock size={12} /> {lesson.durationMin} мин
          </Badge>
          {completed && <Badge tone="success">Пройдено</Badge>}
          <span className="text-xs text-faint">{data.moduleTitle}</span>
        </div>

        <h1 className="mb-1 text-2xl font-bold tracking-tight text-fg sm:text-3xl">
          {lesson.title}
        </h1>
        {lesson.summary && <p className="mb-6 text-sm text-muted">{lesson.summary}</p>}

        <NoteSelection lessonId={lesson.id}>{body}</NoteSelection>

        <div className="mt-[var(--gap)] flex items-center justify-between gap-3">
          {data.prevLessonId ? (
            <Link
              to={`/learn/courses/${data.courseSlug}/lessons/${data.prevLessonId}`}
              className="btn btn-secondary"
            >
              <ChevronRight size={16} className="rotate-180" />
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
              <ChevronRight size={16} />
            </Link>
          ) : (
            <Link to={`/learn/courses/${data.courseSlug}`} className="btn btn-secondary">
              <Terminal size={16} />
              К программе курса
            </Link>
          )}
        </div>
      </div>
    </div>
  );
}
