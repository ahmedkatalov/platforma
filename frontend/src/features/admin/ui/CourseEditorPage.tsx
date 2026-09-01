import { useEffect, useState, type FormEvent } from "react";
import { Link, useParams } from "react-router-dom";

import {
  useCreateLessonMutation,
  useCreateModuleMutation,
  useDeleteLessonMutation,
  useDeleteModuleMutation,
  useGetAdminCourseQuery,
  useUpdateCourseMutation,
  useUpdateLessonMutation,
  useUpdateModuleMutation,
  type CoursePayload,
} from "@/features/admin/api/coursesApi";
import { apiErrorMessage } from "@/shared/api/baseApi";
import type { Course, CourseLevel, CourseStatus, Lesson, LessonKind, Module } from "@/shared/types";
import {
  Badge,
  Button,
  Card,
  EmptyState,
  Field,
  Input,
  Modal,
  PageHeader,
  Select,
  Spinner,
  Textarea,
} from "@/shared/ui";
import { Book, ChevronRight, Edit2, Plus, Trash2 } from "lucide-react";
import { useToast } from "@/shared/ui/ToastProvider";

import LessonContentEditor from "./LessonContentEditor";
import { slugify } from "./CoursesPage";

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

// Заготовки содержимого — чтобы урок сразу открывался в нужном формате.
const KIND_TEMPLATE: Record<LessonKind, Record<string, unknown>> = {
  text: { body: "## Заголовок\n\nТекст урока в Markdown." },
  quiz: {
    passScore: 70,
    questions: [
      {
        id: "q1",
        text: "Какой командой посмотреть запущенные контейнеры?",
        options: [
          { id: "a", text: "docker ps", correct: true },
          { id: "b", text: "docker images", correct: false },
        ],
        explanation: "docker ps показывает запущенные контейнеры.",
      },
    ],
  },
  terminal: {
    intro: "Потренируйтесь работать с контейнерами.",
    tasks: [
      {
        id: "t1",
        prompt: "Выведите список запущенных контейнеров",
        expected: "docker ps",
        hint: "Команда начинается с docker",
      },
    ],
  },
  code: {
    language: "yaml",
    starter: "version: '3.9'\nservices:\n  app:\n    image: nginx\n",
    task: "Добавьте проброс порта 8080:80",
  },
};

export default function CourseEditorPage() {
  const { id = "" } = useParams();
  const { data: course, isLoading } = useGetAdminCourseQuery(id, { skip: !id });

  const [moduleModal, setModuleModal] = useState<Module | "new" | null>(null);
  const [lessonModal, setLessonModal] = useState<{ moduleId: string; lesson: Lesson | null } | null>(
    null,
  );

  if (isLoading || !course) {
    return (
      <div className="grid place-items-center py-20 text-accent">
        <Spinner size={32} />
      </div>
    );
  }

  return (
    <>
      <PageHeader
        title={course.title}
        subtitle={`/${course.slug} · ${course.modulesCount} модулей · ${course.lessonsCount} уроков`}
        actions={
          <>
            <Link to={`/learn/courses/${course.slug}`} className="btn btn-secondary">
              <Book size={16} />
              Открыть как студент
            </Link>
            <Link to="/admin/courses" className="btn btn-ghost">
              К списку курсов
            </Link>
          </>
        }
      />

      <div className="grid gap-[var(--gap)] lg:grid-cols-5">
        <CourseSettings course={course} />

        <Card className="p-[var(--pad)] lg:col-span-3">
          <div className="mb-4 flex items-center justify-between">
            <h2 className="text-base font-bold text-fg">Программа курса</h2>
            <Button
              variant="primary"
              icon={<Plus size={16} />}
              onClick={() => setModuleModal("new")}
            >
              Модуль
            </Button>
          </div>

          {(course.modules ?? []).length === 0 ? (
            <EmptyState
              title="Модулей пока нет"
              description="Разбейте курс на модули, а модули — на уроки"
              icon={<Book size={32} />}
            />
          ) : (
            <div className="space-y-3">
              {(course.modules ?? []).map((module, index) => (
                <ModuleBlock
                  key={module.id}
                  index={index + 1}
                  module={module}
                  courseId={course.id}
                  onEdit={() => setModuleModal(module)}
                  onAddLesson={() => setLessonModal({ moduleId: module.id, lesson: null })}
                  onEditLesson={(lesson) => setLessonModal({ moduleId: module.id, lesson })}
                />
              ))}
            </div>
          )}
        </Card>
      </div>

      <ModuleModal
        open={moduleModal !== null}
        module={moduleModal === "new" ? null : moduleModal}
        courseId={course.id}
        onClose={() => setModuleModal(null)}
      />

      <LessonModal
        open={lessonModal !== null}
        courseId={course.id}
        moduleId={lessonModal?.moduleId ?? ""}
        lesson={lessonModal?.lesson ?? null}
        onClose={() => setLessonModal(null)}
      />
    </>
  );
}

function CourseSettings({ course }: { course: Course }) {
  const [form, setForm] = useState<CoursePayload>({
    slug: course.slug,
    title: course.title,
    subtitle: course.subtitle,
    description: course.description,
    coverUrl: course.coverUrl,
    level: course.level,
    tags: course.tags,
    status: course.status,
    position: course.position,
  });
  const [tags, setTags] = useState(course.tags.join(", "));
  const [updateCourse, { isLoading }] = useUpdateCourseMutation();
  const toast = useToast();

  const save = async (event: FormEvent) => {
    event.preventDefault();
    try {
      await updateCourse({
        id: course.id,
        ...form,
        tags: tags.split(",").map((t) => t.trim()).filter(Boolean),
      }).unwrap();
      toast.success("Курс сохранён");
    } catch (err) {
      toast.error(apiErrorMessage(err));
    }
  };

  return (
    <Card className="p-[var(--pad)] lg:col-span-2">
      <h2 className="mb-4 text-base font-bold text-fg">Настройки курса</h2>

      <form onSubmit={save} className="space-y-4">
        <Field label="Название">
          <Input
            value={form.title}
            onChange={(e) => setForm({ ...form, title: e.target.value })}
            required
          />
        </Field>

        <Field label="Адрес (slug)">
          <Input
            value={form.slug}
            onChange={(e) => setForm({ ...form, slug: slugify(e.target.value) })}
            required
          />
        </Field>

        <Field label="Подзаголовок">
          <Input
            value={form.subtitle}
            onChange={(e) => setForm({ ...form, subtitle: e.target.value })}
          />
        </Field>

        <Field label="Описание">
          <Textarea
            value={form.description}
            onChange={(e) => setForm({ ...form, description: e.target.value })}
          />
        </Field>

        <div className="grid gap-4 sm:grid-cols-2">
          <Field label="Уровень">
            <Select
              value={form.level}
              onChange={(v) => setForm({ ...form, level: v as CourseLevel })}
            >
              <option value="beginner">Начальный</option>
              <option value="intermediate">Средний</option>
              <option value="advanced">Продвинутый</option>
            </Select>
          </Field>

          <Field label="Статус">
            <Select
              value={form.status}
              onChange={(v) => setForm({ ...form, status: v as CourseStatus })}
            >
              <option value="draft">Черновик</option>
              <option value="published">Опубликован</option>
              <option value="archived">В архиве</option>
            </Select>
          </Field>
        </div>

        <Field label="Теги" hint="Через запятую">
          <Input value={tags} onChange={(e) => setTags(e.target.value)} />
        </Field>

        <Button type="submit" variant="primary" className="w-full" loading={isLoading}>
          Сохранить курс
        </Button>
      </form>
    </Card>
  );
}

function ModuleBlock({
  index,
  module,
  courseId,
  onEdit,
  onAddLesson,
  onEditLesson,
}: {
  index: number;
  module: Module;
  courseId: string;
  onEdit: () => void;
  onAddLesson: () => void;
  onEditLesson: (lesson: Lesson) => void;
}) {
  const [open, setOpen] = useState(true);
  const [deleteModule] = useDeleteModuleMutation();
  const [deleteLesson] = useDeleteLessonMutation();
  const toast = useToast();

  const removeModule = async () => {
    if (!window.confirm(`Удалить модуль «${module.title}» со всеми уроками?`)) return;
    try {
      await deleteModule({ moduleId: module.id, courseId }).unwrap();
      toast.success("Модуль удалён");
    } catch (err) {
      toast.error(apiErrorMessage(err));
    }
  };

  const removeLesson = async (lesson: Lesson) => {
    if (!window.confirm(`Удалить урок «${lesson.title}»?`)) return;
    try {
      await deleteLesson({ lessonId: lesson.id, courseId }).unwrap();
      toast.success("Урок удалён");
    } catch (err) {
      toast.error(apiErrorMessage(err));
    }
  };

  return (
    <div className="card-flat overflow-hidden">
      <div className="flex items-center gap-2 p-3">
        <button
          className="btn btn-ghost h-8 w-8 !p-0"
          onClick={() => setOpen((v) => !v)}
          aria-label={open ? "Свернуть" : "Развернуть"}
        >
          <ChevronRight
            size={16}
            className={`transition-transform ${open ? "rotate-90" : ""}`}
          />
        </button>

        <div className="min-w-0 flex-1">
          <p className="truncate text-sm font-bold text-fg">
            <span className="text-faint">{index}. </span>
            {module.title}
          </p>
          {module.summary && <p className="truncate text-xs text-faint">{module.summary}</p>}
        </div>

        <Badge>{module.lessons?.length ?? 0} уроков</Badge>

        <Button variant="ghost" className="h-8 !px-2" onClick={onEdit} title="Изменить модуль">
          <Edit2 size={16} />
        </Button>
        <Button
          variant="ghost"
          className="h-8 !px-2 text-danger"
          onClick={removeModule}
          title="Удалить модуль"
        >
          <Trash2 size={16} />
        </Button>
      </div>

      {open && (
        <div className="border-t border-line px-3 py-2">
          {(module.lessons ?? []).length === 0 ? (
            <p className="py-3 text-center text-xs text-faint">В модуле пока нет уроков</p>
          ) : (
            <ul className="divide-y divide-[var(--border)]">
              {(module.lessons ?? []).map((lesson) => (
                <li key={lesson.id} className="flex items-center gap-2 py-2">
                  <Badge tone={KIND_TONE[lesson.kind]}>{KIND_LABEL[lesson.kind]}</Badge>
                  <span className="min-w-0 flex-1 truncate text-sm text-fg">{lesson.title}</span>
                  <span className="shrink-0 text-xs text-faint">{lesson.durationMin} мин</span>
                  <Button
                    variant="ghost"
                    className="h-7 !px-2"
                    onClick={() => onEditLesson(lesson)}
                    title="Изменить урок"
                  >
                    <Edit2 size={14} />
                  </Button>
                  <Button
                    variant="ghost"
                    className="h-7 !px-2 text-danger"
                    onClick={() => removeLesson(lesson)}
                    title="Удалить урок"
                  >
                    <Trash2 size={14} />
                  </Button>
                </li>
              ))}
            </ul>
          )}

          <Button variant="ghost" className="mt-1 w-full" icon={<Plus size={16} />} onClick={onAddLesson}>
            Добавить урок
          </Button>
        </div>
      )}
    </div>
  );
}

function ModuleModal({
  open,
  module,
  courseId,
  onClose,
}: {
  open: boolean;
  module: Module | null;
  courseId: string;
  onClose: () => void;
}) {
  const [title, setTitle] = useState("");
  const [summary, setSummary] = useState("");
  const [position, setPosition] = useState(0);

  const [createModule, { isLoading: creating }] = useCreateModuleMutation();
  const [updateModule, { isLoading: updating }] = useUpdateModuleMutation();
  const toast = useToast();

  useEffect(() => {
    setTitle(module?.title ?? "");
    setSummary(module?.summary ?? "");
    setPosition(module?.position ?? 0);
  }, [module, open]);

  const submit = async (event: FormEvent) => {
    event.preventDefault();
    try {
      if (module) {
        await updateModule({ moduleId: module.id, courseId, title, summary, position }).unwrap();
        toast.success("Модуль обновлён");
      } else {
        await createModule({ courseId, title, summary, position }).unwrap();
        toast.success("Модуль добавлен");
      }
      onClose();
    } catch (err) {
      toast.error(apiErrorMessage(err));
    }
  };

  return (
    <Modal
      open={open}
      onClose={onClose}
      title={module ? "Изменить модуль" : "Новый модуль"}
      footer={
        <>
          <Button onClick={onClose}>Отмена</Button>
          <Button
            variant="primary"
            type="submit"
            form="module-form"
            loading={creating || updating}
          >
            Сохранить
          </Button>
        </>
      }
    >
      <form id="module-form" onSubmit={submit} className="space-y-4">
        <Field label="Название модуля">
          <Input
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="Основы Linux"
            required
          />
        </Field>

        <Field label="Краткое описание">
          <Textarea
            value={summary}
            onChange={(e) => setSummary(e.target.value)}
            rows={3}
            placeholder="Что разбираем в модуле"
          />
        </Field>

        <Field label="Позиция" hint="0 — добавить в конец">
          <Input
            type="number"
            value={position}
            onChange={(e) => setPosition(Number(e.target.value))}
            min={0}
          />
        </Field>
      </form>
    </Modal>
  );
}

function LessonModal({
  open,
  courseId,
  moduleId,
  lesson,
  onClose,
}: {
  open: boolean;
  courseId: string;
  moduleId: string;
  lesson: Lesson | null;
  onClose: () => void;
}) {
  const [title, setTitle] = useState("");
  const [kind, setKind] = useState<LessonKind>("text");
  const [summary, setSummary] = useState("");
  const [durationMin, setDurationMin] = useState(10);
  const [position, setPosition] = useState(0);
  const [content, setContent] = useState<Record<string, unknown>>({});
  const [raw, setRaw] = useState("{}");
  const [rawError, setRawError] = useState("");
  const [error, setError] = useState("");

  const [createLesson, { isLoading: creating }] = useCreateLessonMutation();
  const [updateLesson, { isLoading: updating }] = useUpdateLessonMutation();
  const toast = useToast();

  useEffect(() => {
    if (!open) return;
    setTitle(lesson?.title ?? "");
    setKind(lesson?.kind ?? "text");
    setSummary(lesson?.summary ?? "");
    setDurationMin(lesson?.durationMin ?? 10);
    setPosition(lesson?.position ?? 0);

    const initial = (lesson?.content ?? KIND_TEMPLATE[lesson?.kind ?? "text"]) as Record<
      string,
      unknown
    >;
    setContent(initial);
    setRaw(JSON.stringify(initial, null, 2));
    setRawError("");
    setError("");
  }, [lesson, open]);

  // Содержимое хранится объектом; сырой JSON держим синхронно для режима «JSON».
  const applyContent = (next: Record<string, unknown>) => {
    setContent(next);
    setRaw(JSON.stringify(next, null, 2));
    setRawError("");
  };

  const applyRaw = (next: string) => {
    setRaw(next);
    try {
      setContent(JSON.parse(next || "{}") as Record<string, unknown>);
      setRawError("");
    } catch {
      setRawError("Некорректный JSON — исправьте, иначе изменения не сохранятся");
    }
  };

  // При смене типа подставляем шаблон, если урок ещё не наполнен.
  const changeKind = (next: LessonKind) => {
    setKind(next);

    const isTemplate = Object.values(KIND_TEMPLATE).some(
      (template) => JSON.stringify(template) === JSON.stringify(content),
    );
    if (!lesson || isTemplate || Object.keys(content).length === 0) {
      applyContent(KIND_TEMPLATE[next] as Record<string, unknown>);
    }
  };

  const submit = async (event: FormEvent) => {
    event.preventDefault();
    setError("");

    if (rawError) {
      setError("Исправьте JSON содержимого урока");
      return;
    }

    const payload = { title, kind, summary, content, durationMin, position };

    try {
      if (lesson) {
        await updateLesson({ lessonId: lesson.id, courseId, ...payload }).unwrap();
        toast.success("Урок обновлён");
      } else {
        await createLesson({ moduleId, courseId, ...payload }).unwrap();
        toast.success("Урок добавлен");
      }
      onClose();
    } catch (err) {
      setError(apiErrorMessage(err));
    }
  };

  return (
    <Modal
      open={open}
      onClose={onClose}
      title={lesson ? "Изменить урок" : "Новый урок"}
      width="42rem"
      footer={
        <>
          <Button onClick={onClose}>Отмена</Button>
          <Button variant="primary" type="submit" form="lesson-form" loading={creating || updating}>
            Сохранить
          </Button>
        </>
      }
    >
      <form id="lesson-form" onSubmit={submit} className="space-y-4">
        <Field label="Название урока">
          <Input
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="Установка Docker"
            required
          />
        </Field>

        <div className="grid gap-4 sm:grid-cols-3">
          <Field label="Тип">
            <Select value={kind} onChange={(v) => changeKind(v as LessonKind)}>
              <option value="text">Теория</option>
              <option value="quiz">Квиз</option>
              <option value="terminal">Терминал</option>
              <option value="code">Редактор кода</option>
            </Select>
          </Field>

          <Field label="Длительность, мин">
            <Input
              type="number"
              value={durationMin}
              onChange={(e) => setDurationMin(Number(e.target.value))}
              min={1}
            />
          </Field>

          <Field label="Позиция" hint="0 — в конец">
            <Input
              type="number"
              value={position}
              onChange={(e) => setPosition(Number(e.target.value))}
              min={0}
            />
          </Field>
        </div>

        <Field label="Краткое описание">
          <Input
            value={summary}
            onChange={(e) => setSummary(e.target.value)}
            placeholder="О чём урок"
          />
        </Field>

        <LessonContentEditor
          kind={kind}
          value={content}
          onChange={applyContent}
          raw={raw}
          onRawChange={applyRaw}
          rawError={rawError}
        />

        {error && (
          <p className="rounded-[var(--radius-md)] bg-[var(--danger-soft)] px-3 py-2 text-sm text-danger">
            {error}
          </p>
        )}
      </form>
    </Modal>
  );
}
