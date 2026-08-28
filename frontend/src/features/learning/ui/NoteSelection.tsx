import { useCallback, useEffect, useRef, useState, type ReactNode } from "react";

import { useCreateNoteMutation } from "@/shared/api/meApi";
import { apiErrorMessage } from "@/shared/api/baseApi";
import { useToast } from "@/shared/ui/ToastProvider";

// Обёртка вокруг содержимого урока: выделили текст — появляется кнопка
// «Сохранить в заметки». Сохранённая цитата попадает на страницу «Заметки»
// вместе с названием курса, модуля и урока.
export default function NoteSelection({
  lessonId,
  children,
}: {
  lessonId: string;
  children: ReactNode;
}) {
  const containerRef = useRef<HTMLDivElement>(null);
  const [popup, setPopup] = useState<{ x: number; y: number; text: string } | null>(null);
  const [createNote, { isLoading }] = useCreateNoteMutation();
  const toast = useToast();

  const hide = useCallback(() => setPopup(null), []);

  useEffect(() => {
    const onSelectionEnd = () => {
      // Даём браузеру закончить выделение, затем читаем его.
      window.setTimeout(() => {
        const selection = window.getSelection();
        const container = containerRef.current;

        if (!selection || selection.isCollapsed || !container) {
          setPopup(null);
          return;
        }

        const text = selection.toString().trim();
        if (text.length < 3 || text.length > 2000) {
          setPopup(null);
          return;
        }

        // Кнопку показываем только для выделений внутри урока.
        const range = selection.getRangeAt(0);
        if (!container.contains(range.commonAncestorContainer)) {
          setPopup(null);
          return;
        }

        const rect = range.getBoundingClientRect();
        const host = container.getBoundingClientRect();

        setPopup({
          x: Math.min(Math.max(rect.left + rect.width / 2 - host.left, 90), host.width - 90),
          y: rect.top - host.top,
          text,
        });
      }, 0);
    };

    document.addEventListener("mouseup", onSelectionEnd);
    document.addEventListener("touchend", onSelectionEnd);
    return () => {
      document.removeEventListener("mouseup", onSelectionEnd);
      document.removeEventListener("touchend", onSelectionEnd);
    };
  }, []);

  useEffect(() => {
    hide();
  }, [lessonId, hide]);

  const save = async () => {
    if (!popup) return;
    try {
      await createNote({ lessonId, quote: popup.text }).unwrap();
      toast.success("Сохранено в заметки");
      window.getSelection()?.removeAllRanges();
      hide();
    } catch (err) {
      toast.error(apiErrorMessage(err, "Не удалось сохранить заметку"));
    }
  };

  return (
    <div ref={containerRef} className="relative">
      {children}

      {popup && (
        <button
          onMouseDown={(e) => {
            // mousedown раньше mouseup снимет выделение — гасим его.
            e.preventDefault();
            e.stopPropagation();
          }}
          onClick={save}
          disabled={isLoading}
          className="absolute z-30 -translate-x-1/2 -translate-y-full whitespace-nowrap rounded-full border border-accent-border bg-surface-solid px-3 py-1.5 text-xs font-bold text-accent shadow-[var(--shadow-md)] transition-transform hover:scale-105"
          style={{ left: popup.x, top: popup.y - 8 }}
        >
          {isLoading ? "Сохраняю…" : "＋ В заметки"}
        </button>
      )}
    </div>
  );
}
