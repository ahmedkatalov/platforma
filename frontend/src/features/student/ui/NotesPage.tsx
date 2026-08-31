import { useMemo, useState } from "react";
import { Link } from "react-router-dom";

import {
  useDeleteNoteMutation,
  useGetMyNotesQuery,
  useUpdateNoteMutation,
} from "@/shared/api/meApi";
import { apiErrorMessage } from "@/shared/api/baseApi";
import type { Note } from "@/shared/types";
import {
  Badge,
  Button,
  Card,
  EmptyState,
  Input,
  PageHeader,
  Spinner,
  Textarea,
} from "@/shared/ui";
import { Book, Check, Edit2, Search, Trash2 } from "lucide-react";
import { useToast } from "@/shared/ui/ToastProvider";

const dateFmt = new Intl.DateTimeFormat("ru-RU", {
  day: "2-digit",
  month: "long",
  hour: "2-digit",
  minute: "2-digit",
});

export default function NotesPage() {
  const { data: notes = [], isLoading } = useGetMyNotesQuery();
  const [search, setSearch] = useState("");

  const groups = useMemo(() => {
    const query = search.trim().toLowerCase();
    const filtered = query
      ? notes.filter(
          (note) =>
            note.quote.toLowerCase().includes(query) ||
            note.body.toLowerCase().includes(query) ||
            note.lessonTitle.toLowerCase().includes(query) ||
            note.moduleTitle.toLowerCase().includes(query),
        )
      : notes;

    // Группируем по курсу, внутри — свежие сверху (как пришли с сервера).
    const byCourse = new Map<string, Note[]>();
    for (const note of filtered) {
      const list = byCourse.get(note.courseTitle) ?? [];
      list.push(note);
      byCourse.set(note.courseTitle, list);
    }
    return [...byCourse.entries()];
  }, [notes, search]);

  if (isLoading) {
    return (
      <div className="grid place-items-center py-20 text-accent">
        <Spinner size={32} />
      </div>
    );
  }

  return (
    <>
      <PageHeader
        title="Заметки"
        subtitle="Выделите текст в любом уроке и нажмите «В заметки» — цитата сохранится сюда"
      />

      {notes.length > 0 && (
        <Card className="mb-[var(--gap)] p-[var(--pad)]">
          <div className="relative">
            <Search
              size={16}
              className="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-faint"
            />
            <Input
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              placeholder="Поиск по цитатам, комментариям и урокам"
              className="pl-9"
            />
          </div>
        </Card>
      )}

      {groups.length === 0 ? (
        <Card>
          <EmptyState
            title={notes.length === 0 ? "Заметок пока нет" : "Ничего не найдено"}
            description={
              notes.length === 0
                ? "Откройте урок, выделите важную мысль мышкой и нажмите «В заметки»"
                : "Попробуйте другой запрос"
            }
            icon={<Edit2 size={32} />}
            action={
              notes.length === 0 ? (
                <Link to="/learn/courses" className="btn btn-primary">
                  <Book size={16} />
                  К урокам
                </Link>
              ) : undefined
            }
          />
        </Card>
      ) : (
        <div className="space-y-[var(--gap)]">
          {groups.map(([courseTitle, items]) => (
            <section key={courseTitle}>
              <h2 className="mb-2 px-1 text-sm font-bold uppercase tracking-wide text-faint">
                {courseTitle} · {items.length}
              </h2>
              <div className="space-y-[var(--gap)]">
                {items.map((note) => (
                  <NoteCard key={note.id} note={note} />
                ))}
              </div>
            </section>
          ))}
        </div>
      )}
    </>
  );
}

function NoteCard({ note }: { note: Note }) {
  const [editing, setEditing] = useState(false);
  const [draft, setDraft] = useState(note.body);

  const [updateNote, { isLoading: saving }] = useUpdateNoteMutation();
  const [deleteNote] = useDeleteNoteMutation();
  const toast = useToast();

  const save = async () => {
    try {
      await updateNote({ id: note.id, body: draft }).unwrap();
      setEditing(false);
      toast.success("Комментарий сохранён");
    } catch (err) {
      toast.error(apiErrorMessage(err));
    }
  };

  const remove = async () => {
    if (!window.confirm("Удалить заметку?")) return;
    try {
      await deleteNote(note.id).unwrap();
      toast.success("Заметка удалена");
    } catch (err) {
      toast.error(apiErrorMessage(err));
    }
  };

  return (
    <Card className="p-[var(--pad)]">
      {/* Где сделана заметка: модуль → урок, ссылка ведёт прямо в него */}
      <div className="mb-3 flex flex-wrap items-center gap-2 text-xs">
        <Badge tone="accent">{note.moduleTitle}</Badge>
        <Link
          to={`/learn/courses/${note.courseSlug}/lessons/${note.lessonId}`}
          className="font-semibold text-accent hover:underline"
        >
          {note.lessonTitle}
        </Link>
        <span className="ml-auto text-faint">{dateFmt.format(new Date(note.createdAt))}</span>
      </div>

      <blockquote className="rounded-[var(--radius-md)] border-l-4 border-[var(--accent)] bg-accent-soft px-4 py-3 text-sm leading-relaxed text-fg">
        {note.quote}
      </blockquote>

      {editing ? (
        <div className="mt-3 space-y-2">
          <Textarea
            value={draft}
            onChange={(e) => setDraft(e.target.value)}
            rows={3}
            placeholder="Свой комментарий: почему это важно, что запомнить"
            autoFocus
          />
          <div className="flex gap-2">
            <Button
              variant="primary"
              className="h-8"
              icon={<Check size={14} />}
              onClick={save}
              loading={saving}
            >
              Сохранить
            </Button>
            <Button
              variant="ghost"
              className="h-8"
              onClick={() => {
                setDraft(note.body);
                setEditing(false);
              }}
            >
              Отмена
            </Button>
          </div>
        </div>
      ) : (
        <div className="mt-3 flex items-start gap-2">
          {note.body ? (
            <p className="flex-1 text-sm text-muted">{note.body}</p>
          ) : (
            <button
              className="flex-1 text-left text-xs text-faint hover:text-accent"
              onClick={() => setEditing(true)}
            >
              ＋ добавить комментарий
            </button>
          )}
          <Button
            variant="ghost"
            className="h-7 !px-2"
            onClick={() => setEditing(true)}
            title="Изменить комментарий"
          >
            <Edit2 size={14} />
          </Button>
          <Button
            variant="ghost"
            className="h-7 !px-2 text-danger"
            onClick={remove}
            title="Удалить заметку"
          >
            <Trash2 size={14} />
          </Button>
        </div>
      )}
    </Card>
  );
}
