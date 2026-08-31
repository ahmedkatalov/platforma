import { useRef, useState, type ChangeEvent, type FormEvent } from "react";
import { Link } from "react-router-dom";

import {
  useCreateCourseMutation,
  useDeleteCourseMutation,
  useGetAdminCoursesQuery,
  useImportCourseMutation,
  type CoursePayload,
} from "@/features/admin/api/coursesApi";
import { apiErrorMessage } from "@/shared/api/baseApi";
import { tokenStorage } from "@/shared/api/tokenStorage";
import type { Course, CourseLevel, CourseStatus } from "@/shared/types";
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

const STATUS_LABEL: Record<CourseStatus, string> = {
  draft: "Черновик",
  published: "Опубликован",
  archived: "В архиве",
};

const LEVEL_LABEL: Record<CourseLevel, string> = {
  beginner: "Начальный",
  intermediate: "Средний",
  advanced: "Продвинутый",
};

const EMPTY_COURSE: CoursePayload = {
  slug: "",
  title: "",
  subtitle: "",
  description: "",
  coverUrl: "",
  level: "beginner",
  tags: [],
  status: "draft",
  position: 0,
};

// Латиница из названия: «DevOps с нуля» → devops-s-nulya.
const TRANSLIT: Record<string, string> = {
  а: "a", б: "b", в: "v", г: "g", д: "d", е: "e", ё: "e", ж: "zh", з: "z", и: "i",
  й: "y", к: "k", л: "l", м: "m", н: "n", о: "o", п: "p", р: "r", с: "s", т: "t",
  у: "u", ф: "f", х: "h", ц: "c", ч: "ch", ш: "sh", щ: "sch", ъ: "", ы: "y", ь: "",
  э: "e", ю: "yu", я: "ya",
};

export function slugify(value: string): string {
  return value
    .toLowerCase()
    .split("")
    .map((char) => TRANSLIT[char] ?? char)
    .join("")
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-+|-+$/g, "");
}

export default function CoursesPage() {
  const { data: courses = [], isLoading } = useGetAdminCoursesQuery();
  const [deleteCourse] = useDeleteCourseMutation();
  const [importCourse, { isLoading: importing }] = useImportCourseMutation();
  const [createOpen, setCreateOpen] = useState(false);
  const fileRef = useRef<HTMLInputElement>(null);
  const toast = useToast();

  const remove = async (course: Course) => {
    if (!window.confirm(`Удалить курс «${course.title}» вместе со всеми модулями и уроками?`)) return;
    try {
      await deleteCourse(course.id).unwrap();
      toast.success("Курс удалён");
    } catch (err) {
      toast.error(apiErrorMessage(err));
    }
  };

  // Загрузка курса из JSON-пакета: читаем файл, отправляем на сервер.
  const onFile = async (event: ChangeEvent<HTMLInputElement>) => {
    const input = event.target;
    const file = input.files?.[0];
    if (!file) return;

    // ВАЖНО: читаем файл ДО сброса input.value — иначе ссылка на файл
    // становится недействительной и браузер бросает NotReadableError.
    let raw: string;
    try {
      raw = (await file.text()).replace(/^﻿/, "").trim();
    } catch {
      input.value = "";
      toast.error(
        "Не удалось прочитать файл. Если он лежит в облаке (iCloud, Google Drive) — " +
          "сначала скачайте его на компьютер, затем выберите снова.",
      );
      return;
    }
    input.value = ""; // теперь можно сбросить — файл уже прочитан

    if (!raw.startsWith("{")) {
      toast.error(
        `Это не похоже на пакет курса. Нужен JSON-файл вида «...course.json». Выбран: ${file.name} (${Math.round(
          file.size / 1024,
        )} КБ)`,
      );
      return;
    }

    const runImport = async (replace: boolean) => {
      const res = await importCourse({ raw, replace }).unwrap();
      toast.success(`Загружено: ${res.modules} глав, ${res.lessons} уроков (черновик)`);
    };

    try {
      await runImport(false);
    } catch (err) {
      if ((err as { status?: number }).status === 409) {
        if (window.confirm(`${apiErrorMessage(err)}\n\nЗаменить существующий курс новым из файла?`)) {
          try {
            await runImport(true);
          } catch (err2) {
            toast.error(apiErrorMessage(err2, "Не удалось заменить курс"));
          }
        }
        return;
      }
      toast.error(apiErrorMessage(err, "Не удалось загрузить курс"));
    }
  };

  // Скачивание курса файлом: авторизованный запрос + сохранение blob.
  const exportCourse = async (course: Course) => {
    try {
      const res = await fetch(`/api/admin/courses/${course.id}/export`, {
        headers: { Authorization: `Bearer ${tokenStorage.access() ?? ""}` },
      });
      if (!res.ok) throw new Error();
      const url = URL.createObjectURL(await res.blob());
      const link = document.createElement("a");
      link.href = url;
      link.download = `${course.slug}.course.json`;
      document.body.appendChild(link);
      link.click();
      link.remove();
      URL.revokeObjectURL(url);
    } catch {
      toast.error("Не удалось скачать курс");
    }
  };

  return (
    <>
      <input
        ref={fileRef}
        type="file"
        accept=".json,application/json"
        className="hidden"
        onChange={onFile}
      />
      <PageHeader
        title="Курсы"
        subtitle="Программа обучения, модули и уроки"
        actions={
          <div className="flex gap-2">
            <Button
              variant="secondary"
              icon={<Book size={18} />}
              loading={importing}
              onClick={() => fileRef.current?.click()}
            >
              Загрузить из файла
            </Button>
            <Button variant="primary" icon={<Plus size={18} />} onClick={() => setCreateOpen(true)}>
              Новый курс
            </Button>
          </div>
        }
      />

      {isLoading ? (
        <div className="grid place-items-center py-20 text-accent">
          <Spinner size={32} />
        </div>
      ) : courses.length === 0 ? (
        <Card>
          <EmptyState
            title="Курсов пока нет"
            description="Создайте курс по DevOps и наполните его модулями, квизами и практикой в терминале"
            icon={<Book size={32} />}
            action={
              <Button variant="primary" onClick={() => setCreateOpen(true)}>
                Создать курс
              </Button>
            }
          />
        </Card>
      ) : (
        <div className="grid gap-[var(--gap)] md:grid-cols-2 xl:grid-cols-3">
          {courses.map((course) => (
            <Card key={course.id} className="flex flex-col p-[var(--pad)]">
              <div className="mb-3 flex items-start justify-between gap-2">
                <div className="min-w-0">
                  <h2 className="truncate text-base font-bold text-fg">{course.title}</h2>
                  <p className="truncate text-xs text-faint">/{course.slug}</p>
                </div>
                <Badge tone={course.status === "published" ? "success" : "default"}>
                  {STATUS_LABEL[course.status]}
                </Badge>
              </div>

              {course.subtitle && (
                <p className="mb-3 line-clamp-2 text-sm text-muted">{course.subtitle}</p>
              )}

              <div className="mb-4 flex flex-wrap gap-1.5">
                <Badge>{LEVEL_LABEL[course.level]}</Badge>
                {course.tags.slice(0, 3).map((tag) => (
                  <Badge key={tag} tone="accent">
                    {tag}
                  </Badge>
                ))}
              </div>

              <div className="mb-4 grid grid-cols-3 gap-2 text-center">
                <div className="card-flat py-2">
                  <p className="text-lg font-bold text-fg">{course.modulesCount}</p>
                  <p className="text-[11px] text-faint">модулей</p>
                </div>
                <div className="card-flat py-2">
                  <p className="text-lg font-bold text-fg">{course.lessonsCount}</p>
                  <p className="text-[11px] text-faint">уроков</p>
                </div>
                <div className="card-flat py-2">
                  <p className="text-lg font-bold text-fg">{course.studentsCount}</p>
                  <p className="text-[11px] text-faint">студентов</p>
                </div>
              </div>

              <div className="mt-auto flex gap-2">
                <Link to={`/admin/courses/${course.id}`} className="btn btn-primary flex-1">
                  <Edit2 size={16} />
                  Редактировать
                </Link>
                <Link
                  to={`/learn/courses/${course.slug}`}
                  className="btn btn-secondary"
                  title="Открыть курс глазами студента"
                >
                  <Book size={16} />
                </Link>
                <Button
                  variant="ghost"
                  onClick={() => exportCourse(course)}
                  title="Скачать курс файлом (для загрузки в другом месте)"
                >
                  <ChevronRight size={16} className="rotate-90" />
                </Button>
                <Button
                  variant="ghost"
                  className="text-danger"
                  onClick={() => remove(course)}
                  title="Удалить курс"
                >
                  <Trash2 size={16} />
                </Button>
              </div>
            </Card>
          ))}
        </div>
      )}

      <CreateCourseModal open={createOpen} onClose={() => setCreateOpen(false)} />
    </>
  );
}

function CreateCourseModal({ open, onClose }: { open: boolean; onClose: () => void }) {
  const [form, setForm] = useState<CoursePayload>(EMPTY_COURSE);
  const [tags, setTags] = useState("");
  const [error, setError] = useState("");
  const [createCourse, { isLoading }] = useCreateCourseMutation();
  const toast = useToast();

  const submit = async (event: FormEvent) => {
    event.preventDefault();
    setError("");

    try {
      await createCourse({
        ...form,
        slug: form.slug || slugify(form.title),
        tags: tags
          .split(",")
          .map((tag) => tag.trim())
          .filter(Boolean),
      }).unwrap();

      toast.success("Курс создан");
      setForm(EMPTY_COURSE);
      setTags("");
      onClose();
    } catch (err) {
      setError(apiErrorMessage(err, "Не удалось создать курс"));
    }
  };

  return (
    <Modal
      open={open}
      onClose={onClose}
      title="Новый курс"
      width="36rem"
      footer={
        <>
          <Button onClick={onClose}>Отмена</Button>
          <Button variant="primary" type="submit" form="create-course" loading={isLoading}>
            Создать
          </Button>
        </>
      }
    >
      <form id="create-course" onSubmit={submit} className="space-y-4">
        <Field label="Название">
          <Input
            value={form.title}
            onChange={(e) =>
              setForm((current) => ({
                ...current,
                title: e.target.value,
                slug: current.slug || slugify(e.target.value),
              }))
            }
            placeholder="DevOps с нуля"
            required
          />
        </Field>

        <Field label="Адрес (slug)" hint="Латиница, цифры и дефисы — используется в ссылке">
          <Input
            value={form.slug}
            onChange={(e) => setForm({ ...form, slug: slugify(e.target.value) })}
            placeholder="devops-basics"
            required
          />
        </Field>

        <Field label="Подзаголовок">
          <Input
            value={form.subtitle}
            onChange={(e) => setForm({ ...form, subtitle: e.target.value })}
            placeholder="Linux, Docker, CI/CD и Kubernetes на практике"
          />
        </Field>

        <Field label="Описание">
          <Textarea
            value={form.description}
            onChange={(e) => setForm({ ...form, description: e.target.value })}
            placeholder="Чему научится студент"
          />
        </Field>

        <div className="grid gap-4 sm:grid-cols-2">
          <Field label="Уровень">
            <Select
              value={form.level}
              onChange={(e) => setForm({ ...form, level: e.target.value as CourseLevel })}
            >
              <option value="beginner">Начальный</option>
              <option value="intermediate">Средний</option>
              <option value="advanced">Продвинутый</option>
            </Select>
          </Field>

          <Field label="Статус">
            <Select
              value={form.status}
              onChange={(e) => setForm({ ...form, status: e.target.value as CourseStatus })}
            >
              <option value="draft">Черновик</option>
              <option value="published">Опубликован</option>
              <option value="archived">В архиве</option>
            </Select>
          </Field>
        </div>

        <Field label="Теги" hint="Через запятую: docker, ci-cd, linux">
          <Input value={tags} onChange={(e) => setTags(e.target.value)} placeholder="docker, linux" />
        </Field>

        {error && (
          <p className="rounded-[var(--radius-md)] bg-[var(--danger-soft)] px-3 py-2 text-sm text-danger">
            {error}
          </p>
        )}
      </form>
    </Modal>
  );
}

export { STATUS_LABEL as COURSE_STATUS_LABEL, LEVEL_LABEL as COURSE_LEVEL_LABEL };
