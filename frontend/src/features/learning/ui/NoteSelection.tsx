import { useCallback, useEffect, useRef, useState, type ReactNode } from "react";
import { Sparkles } from "lucide-react";

import { useCreateNoteMutation, useGetAiStatusQuery } from "@/shared/api/meApi";
import { apiErrorMessage } from "@/shared/api/baseApi";
import { useToast } from "@/shared/ui/ToastProvider";
import AiChatModal from "@/features/learning/ui/AiChatModal";

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
  const [aiContext, setAiContext] = useState<string | null>(null);
  const [createNote, { isLoading }] = useCreateNoteMutation();
  const { data: aiStatus } = useGetAiStatusQuery();
  const aiEnabled = Boolean(aiStatus?.enabled);
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

  const askAi = () => {
    if (!popup) return;
    setAiContext(popup.text);
    window.getSelection()?.removeAllRanges();
    hide();
  };

  return (
    <div ref={containerRef} className="relative">
      {children}

      {popup && (
        <div
          onMouseDown={(e) => {
            // mousedown раньше mouseup снимет выделение — гасим его.
            e.preventDefault();
            e.stopPropagation();
          }}
          className="absolute z-30 flex -translate-x-1/2 -translate-y-full items-center gap-1 whitespace-nowrap rounded-full border border-line bg-surface-solid p-1 shadow-[var(--shadow-md)]"
          style={{ left: popup.x, top: popup.y - 8 }}
        >
          <button
            onClick={save}
            disabled={isLoading}
            className="rounded-full px-3 py-1 text-xs font-bold text-accent transition-transform hover:scale-105"
          >
            {isLoading ? "Сохраняю…" : "＋ В заметки"}
          </button>
          {aiEnabled && (
            <>
              <span className="h-4 w-px bg-line" />
              <button
                onClick={askAi}
                className="flex items-center gap-1 rounded-full px-3 py-1 text-xs font-bold text-accent transition-transform hover:scale-105"
              >
                <Sparkles size={13} /> Спросить у ИИ
              </button>
            </>
          )}
        </div>
      )}

      {aiEnabled && (
        <AiChatModal
          open={aiContext !== null}
          onClose={() => setAiContext(null)}
          context={aiContext ?? ""}
          lessonId={lessonId}
        />
      )}
    </div>
  );
}
