import { Link } from "react-router-dom";
import { Bookmark, BookOpen, Layers, Trash2 } from "lucide-react";

import {
  useGetBookmarksQuery,
  useRemoveBookmarkMutation,
} from "@/shared/api/meApi";
import { apiErrorMessage } from "@/shared/api/baseApi";
import type { Bookmark as BookmarkT, LessonKind } from "@/shared/types";
import { Badge, Card, EmptyState, PageHeader, Spinner } from "@/shared/ui";
import { useToast } from "@/shared/ui/ToastProvider";

const KIND_LABEL: Record<LessonKind, string> = {
  text: "Теория",
  quiz: "Квиз",
  terminal: "Тренажёр",
  code: "Практика",
};

function linkTo(b: BookmarkT): string {
  return b.kind === "lesson"
    ? `/learn/courses/${b.courseSlug}/lessons/${b.refId}`
    : `/learn/courses/${b.courseSlug}`;
}

export default function SavedPage() {
  const { data: bookmarks = [], isLoading } = useGetBookmarksQuery();
  const [remove] = useRemoveBookmarkMutation();
  const toast = useToast();

  const modules = bookmarks.filter((b) => b.kind === "module");
  const lessons = bookmarks.filter((b) => b.kind === "lesson");

  const drop = async (b: BookmarkT) => {
    try {
      await remove({ kind: b.kind, refId: b.refId }).unwrap();
    } catch (err) {
      toast.error(apiErrorMessage(err, "Не удалось убрать закладку"));
    }
  };

  if (isLoading) {
    return (
      <div className="grid place-items-center py-20 text-accent">
        <Spinner size={32} />
      </div>
    );
  }

  const row = (b: BookmarkT) => (
    <div
      key={`${b.kind}:${b.refId}`}
      className="card-flat flex items-center gap-3 p-3 transition-colors hover:bg-surface-hover"
    >
      <span className="grid h-10 w-10 shrink-0 place-items-center rounded-[var(--radius-md)] bg-accent-soft text-accent">
        {b.kind === "module" ? <Layers size={18} /> : <BookOpen size={18} />}
      </span>
      <Link to={linkTo(b)} className="min-w-0 flex-1">
        <div className="mb-0.5 flex items-center gap-2">
          <span className="truncate text-sm font-bold text-fg hover:text-accent">{b.title}</span>
          {b.kind === "lesson" && b.lessonKind && (
            <Badge tone="default">{KIND_LABEL[b.lessonKind]}</Badge>
          )}
        </div>
        <p className="truncate text-xs text-faint">
          {b.courseTitle}
          {b.kind === "lesson" && b.moduleTitle ? ` · ${b.moduleTitle}` : ""}
        </p>
      </Link>
      <button
        type="button"
        onClick={() => void drop(b)}
        aria-label="Убрать из сохранённого"
        title="Убрать из сохранённого"
        className="btn btn-ghost btn-icon btn-sm shrink-0 text-faint hover:text-danger"
      >
        <Trash2 size={16} />
      </button>
    </div>
  );

  return (
    <>
      <PageHeader
        title="Сохранённое"
        subtitle="Главы и уроки, которые вы отложили, чтобы вернуться и закрепить"
      />

      {bookmarks.length === 0 ? (
        <Card className="p-[var(--pad)]">
          <EmptyState
            title="Пока ничего не сохранено"
            description="Открывая урок или главу, нажмите «Сохранить» — и они появятся здесь, чтобы легко вернуться и перечитать."
            icon={<Bookmark size={32} />}
          />
        </Card>
      ) : (
        <div className="space-y-[var(--gap)]">
          {modules.length > 0 && (
            <Card className="p-[var(--pad)]">
              <h2 className="mb-3 flex items-center gap-2 text-base font-bold text-fg">
                <Layers size={18} className="text-accent" /> Главы
                <span className="text-sm font-normal text-faint">· {modules.length}</span>
              </h2>
              <div className="space-y-2">{modules.map(row)}</div>
            </Card>
          )}

          {lessons.length > 0 && (
            <Card className="p-[var(--pad)]">
              <h2 className="mb-3 flex items-center gap-2 text-base font-bold text-fg">
                <BookOpen size={18} className="text-accent" /> Уроки
                <span className="text-sm font-normal text-faint">· {lessons.length}</span>
              </h2>
              <div className="space-y-2">{lessons.map(row)}</div>
            </Card>
          )}
        </div>
      )}
    </>
  );
}
