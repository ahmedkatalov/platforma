import { useEffect, useMemo, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";

import { useAppSelector } from "@/app/store";
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
  Edit2,
  HelpCircle,
  Lock,
  Terminal,
  X,
} from "lucide-react";

import { groupThemes, themeProgress } from "@/features/learning/lib/themes";

// Иконка типа урока — понятнее короткой подписи вроде «Терм».
function kindIcon(kind: LessonKind, size = 12) {
  if (kind === "terminal") return <Terminal size={size} />;
  if (kind === "code") return <Edit2 size={size} />;
  if (kind === "quiz") return <HelpCircle size={size} />;
  return <Book size={size} />;
}

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
  const currentUser = useAppSelector((state) => state.auth.user);

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
  // Админ видит все главы (превью). Студенту глава открыта только если она в moduleAccess.
  const isAdminView = currentUser?.role === "admin";
  const moduleAccess = courseData?.moduleAccess ?? {};
  const chapterUnlocked = (moduleId: string) => isAdminView || moduleAccess[moduleId] === true;

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
      {/* Затемнение под мобильной шторкой содержания. */}
      {asideOpen && (
        <div
          className="fixed inset-0 z-40 bg-black/60 xl:hidden"
          onClick={() => setAsideOpen(false)}
          aria-hidden="true"
        />
      )}

      {/* Содержание курса: на телефоне — шторка снизу, на десктопе — колонка слева. */}
      <aside
        className={
          asideOpen
            ? "sheet-panel fixed inset-x-0 bottom-0 z-50 overflow-y-auto xl:static xl:z-auto xl:block xl:max-h-none xl:animate-none xl:overflow-visible xl:rounded-none xl:border-t-0 xl:bg-transparent xl:pb-0 xl:shadow-none"
            : "hidden xl:block"
        }
      >
        {/* Мобильная шапка шторки: «ручка», заголовок, закрытие. */}
        <div className="xl:hidden">
          <div className="sheet-grip" />
          <div className="flex items-center justify-between px-4 pb-2 pt-1">
            <p className="text-sm font-bold text-fg">Содержание курса</p>
            <button
              className="btn btn-ghost btn-icon btn-sm"
              onClick={() => setAsideOpen(false)}
              aria-label="Закрыть содержание"
            >
              <X size={18} />
            </button>
          </div>
        </div>

        <Card
          className={`relative p-3 max-xl:border-0 max-xl:bg-transparent max-xl:shadow-none max-xl:backdrop-blur-none xl:sticky xl:top-24 ${contentsCollapsed ? "xl:p-2" : "xl:p-3"}`}
        >
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
            className={`btn btn-ghost btn-icon btn-sm absolute right-2 top-2 hidden ${
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

            <nav className="max-h-[62vh] space-y-1 overflow-y-auto pr-1">
              {modules.map((module, moduleIndex) => {
                const themes = groupThemes(module);
                const expanded = openModuleId === module.id;
                const lessons = module.lessons ?? [];
                const isDone = (id: string) =>
                  progressByLesson.get(id)?.status === "completed";
                const moduleDone = lessons.filter((l) => isDone(l.id)).length;
                const moduleTotal = lessons.length;
                const activeChapter = lessons.some((l) => l.id === lesson.id);
                const chapterComplete = moduleTotal > 0 && moduleDone === moduleTotal;

                // Закрытая глава: с замком, приглушённая, ведёт на страницу курса
                // (там можно запросить доступ). Уроки не кликабельны.
                if (!chapterUnlocked(module.id)) {
                  return (
                    <section key={module.id}>
                      <Link
                        to={`/learn/courses/${data.courseSlug}`}
                        className="flex w-full items-center gap-2.5 rounded-[var(--radius-md)] px-2 py-2 opacity-75 transition-opacity hover:opacity-100"
                        title="Глава закрыта — запросите доступ на странице курса"
                      >
                        <span className="grid h-7 w-7 shrink-0 place-items-center rounded-[var(--radius-md)] bg-surface-2 text-faint">
                          <Lock size={14} />
                        </span>
                        <span className="min-w-0 flex-1">
                          <span className="block truncate text-sm font-semibold text-muted">
                            {module.title}
                          </span>
                          <span className="block text-[11px] text-faint">
                            Глава {moduleIndex + 1} · доступ закрыт
                          </span>
                        </span>
                      </Link>
                    </section>
                  );
                }

                const renderRow = (item: Lesson) => {
                  const active = item.id === lesson.id;
                  const done = isDone(item.id);
                  return (
                    <li key={item.id}>
                      <Link
                        to={`/learn/courses/${data.courseSlug}/lessons/${item.id}`}
                        className={`flex items-center gap-2.5 rounded-[var(--radius-md)] px-2 py-1.5 text-sm transition-colors ${
                          active
                            ? "bg-accent-soft font-semibold text-accent"
                            : "text-muted hover:bg-surface-2 hover:text-fg"
                        }`}
                      >
                        <span
                          className={`grid h-5 w-5 shrink-0 place-items-center rounded-full ${
                            done
                              ? "bg-[var(--success)] text-white"
                              : active
                                ? "bg-accent text-accent-fg"
                                : "bg-surface-2 text-faint"
                          }`}
                        >
                          {done ? <Check size={12} /> : kindIcon(item.kind)}
                        </span>
                        <span className="min-w-0 flex-1 truncate">{item.title}</span>
                      </Link>
                    </li>
                  );
                };

                return (
                  <section key={module.id}>
                    <button
                      type="button"
                      className={`flex w-full cursor-pointer items-center gap-2.5 rounded-[var(--radius-md)] px-2 py-2 text-left transition-colors ${
                        activeChapter ? "bg-surface-2" : "hover:bg-surface-2"
                      }`}
                      onClick={() =>
                        setOpenModuleId((current) => (current === module.id ? null : module.id))
                      }
                      aria-expanded={expanded}
                      aria-controls={`module-${module.id}`}
                    >
                      <span
                        className={`grid h-7 w-7 shrink-0 place-items-center rounded-[var(--radius-md)] text-xs font-bold ${
                          chapterComplete
                            ? "bg-[var(--success)] text-white"
                            : activeChapter
                              ? "text-accent-fg"
                              : "bg-surface-2 text-muted"
                        }`}
                        style={
                          activeChapter && !chapterComplete
                            ? { background: "var(--gradient)" }
                            : undefined
                        }
                      >
                        {chapterComplete ? <Check size={15} /> : moduleIndex + 1}
                      </span>
                      <span className="min-w-0 flex-1">
                        <span className="block truncate text-sm font-semibold text-fg">
                          {module.title}
                        </span>
                        <span className="block text-[11px] text-faint">
                          Глава {moduleIndex + 1} · {moduleDone}/{moduleTotal} уроков
                        </span>
                      </span>
                      <ChevronRight
                        size={16}
                        className={`shrink-0 text-faint transition-transform duration-200 ${
                          expanded ? "rotate-90" : ""
                        }`}
                        aria-hidden="true"
                      />
                    </button>

                    <div
                      id={`module-${module.id}`}
                      aria-hidden={!expanded}
                      className={`grid transition-[grid-template-rows] duration-200 ease-out ${
                        expanded ? "grid-rows-[1fr]" : "grid-rows-[0fr]"
                      }`}
                    >
                      <div
                        className={`overflow-hidden transition-opacity duration-200 ${
                          expanded ? "opacity-100" : "opacity-0"
                        }`}
                      >
                        <div className="ml-3.5 space-y-3 border-l border-line pb-2 pl-2 pt-1.5">
                          {themes.map((theme, ti) => {
                            const tp = themeProgress(theme, isDone);
                            return (
                              <div key={theme.key}>
                                <div className="flex items-center justify-between gap-2 px-1 pb-1">
                                  <span className="min-w-0 truncate text-[11px] font-semibold uppercase tracking-wide text-muted">
                                    Тема {ti + 1} · {theme.title}
                                  </span>
                                  <span className="shrink-0 text-[10px] font-semibold text-faint">
                                    {tp.done}/{tp.total}
                                  </span>
                                </div>
                                <ul className="space-y-0.5">
                                  {theme.pages.map((page) => renderRow(page))}
                                  {theme.quiz && renderRow(theme.quiz)}
                                </ul>
                              </div>
                            );
                          })}
                        </div>
                      </div>
                    </div>
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

        {/* Навигация между уроками: на телефоне для теории/квиза — липкая нижняя
            панель с явной «Далее»; для терминала/кода — обычный ряд, чтобы липкая
            панель не перекрывала ввод команд и кнопку проверки. На десктопе — ряд. */}
        <div
          className={`mt-[var(--gap)] flex items-center gap-3 pt-[var(--gap)] ${
            lesson.kind === "text" || lesson.kind === "quiz" ? "action-bar" : ""
          }`}
        >
          {data.prevLessonId ? (
            <Link
              to={`/learn/courses/${data.courseSlug}/lessons/${data.prevLessonId}`}
              className="btn btn-secondary btn-icon shrink-0 xl:w-auto xl:px-4"
              aria-label="Предыдущий урок"
            >
              <ChevronRight size={18} className="rotate-180" />
              <span className="hidden xl:inline">Предыдущий</span>
            </Link>
          ) : (
            <span className="hidden xl:block" />
          )}

          {data.nextLessonId ? (
            <Link
              to={`/learn/courses/${data.courseSlug}/lessons/${data.nextLessonId}`}
              className="btn btn-primary flex-1 xl:ml-auto xl:flex-none"
            >
              Далее
              <ChevronRight size={18} />
            </Link>
          ) : (
            <Link
              to={`/learn/courses/${data.courseSlug}`}
              className="btn btn-primary flex-1 xl:ml-auto xl:flex-none"
            >
              <Check size={16} />
              К программе курса
            </Link>
          )}
        </div>
      </div>
    </div>
  );
}
