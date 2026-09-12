import clsx from "clsx";
import { Bookmark, BookmarkCheck } from "lucide-react";

import {
  useAddBookmarkMutation,
  useGetBookmarksQuery,
  useRemoveBookmarkMutation,
} from "@/shared/api/meApi";
import { apiErrorMessage } from "@/shared/api/baseApi";
import { useToast } from "@/shared/ui/ToastProvider";

// Кнопка «Сохранить/Сохранено» для урока (kind=lesson) или главы (kind=module).
// compact — только иконка (для заголовков модулей в списке курса).
export function BookmarkButton({
  kind,
  refId,
  compact = false,
  className,
}: {
  kind: "lesson" | "module";
  refId: string;
  compact?: boolean;
  className?: string;
}) {
  const { data: bookmarks = [] } = useGetBookmarksQuery();
  const [add, { isLoading: adding }] = useAddBookmarkMutation();
  const [remove, { isLoading: removing }] = useRemoveBookmarkMutation();
  const toast = useToast();

  const saved = bookmarks.some((b) => b.kind === kind && b.refId === refId);
  const busy = adding || removing;

  const toggle = async () => {
    try {
      if (saved) {
        await remove({ kind, refId }).unwrap();
      } else {
        await add({ kind, refId }).unwrap();
        toast.success(kind === "module" ? "Глава сохранена" : "Урок сохранён");
      }
    } catch (err) {
      toast.error(apiErrorMessage(err, "Не удалось сохранить"));
    }
  };

  const title = saved ? "Убрать из сохранённого" : "Сохранить, чтобы вернуться позже";
  const Icon = saved ? BookmarkCheck : Bookmark;

  if (compact) {
    return (
      <button
        type="button"
        onClick={(e) => {
          e.preventDefault();
          e.stopPropagation();
          void toggle();
        }}
        disabled={busy}
        aria-label={title}
        title={title}
        className={clsx(
          "btn btn-ghost btn-icon btn-sm shrink-0",
          saved && "text-accent",
          className,
        )}
      >
        <Icon size={18} />
      </button>
    );
  }

  return (
    <button
      type="button"
      onClick={() => void toggle()}
      disabled={busy}
      title={title}
      className={clsx("btn btn-sm", saved ? "btn-primary" : "btn-secondary", className)}
    >
      <Icon size={16} />
      {saved ? "Сохранено" : "Сохранить"}
    </button>
  );
}
