import { useState, type FormEvent } from "react";
import { Link } from "react-router-dom";

import {
  useCreateCourseMutation,
  useDeleteCourseMutation,
  useGetAdminCoursesQuery,
  type CoursePayload,
} from "@/features/admin/api/coursesApi";
import { apiErrorMessage } from "@/shared/api/baseApi";
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
import { IconBook, IconEdit, IconPlus, IconTrash } from "@/shared/ui/icons";
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
  const [createOpen, setCreateOpen] = useState(false);
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

  return (
    <>
      <PageHeader
        title="Курсы"
        subtitle="Программа обучения, модули и уроки"
        actions={
          <Button variant="primary" icon={<IconPlus size={18} />} onClick={() => setCreateOpen(true)}>
            Новый курс
          </Button>
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
            icon={<IconBook size={32} />}
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
                  <IconEdit size={16} />
                  Редактировать
                </Link>
                <Link
                  to={`/learn/courses/${course.slug}`}
                  className="btn btn-secondary"
                  title="Открыть курс глазами студента"
                >
                  <IconBook size={16} />
                </Link>
                <Button
                  variant="ghost"
                  className="text-danger"
                  onClick={() => remove(course)}
                  title="Удалить курс"
                >
                  <IconTrash size={16} />
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
