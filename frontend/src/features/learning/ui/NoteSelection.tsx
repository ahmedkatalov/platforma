import { useCallback, useEffect, useRef, useState, type ReactNode } from "react";
import { Sparkles, X } from "lucide-react";

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

      {/* Десктоп/планшет — плавающий попап у выделения. */}
      {popup && (
        <div
          onMouseDown={(e) => {
            // mousedown раньше mouseup снимет выделение — гасим его.
            e.preventDefault();
            e.stopPropagation();
          }}
          className="absolute z-30 hidden -translate-x-1/2 -translate-y-full items-center gap-1 whitespace-nowrap rounded-full border border-line bg-surface-solid p-1 shadow-[var(--shadow-md)] md:flex"
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

      {/* Телефон — фиксированная панель снизу: не конфликтует с системным меню
          выделения (Копировать/Найти), которое браузер рисует у самого текста. */}
      {popup && (
        <div className="fixed inset-x-0 bottom-0 z-40 flex items-center gap-2 border-t border-line bg-surface-solid px-3 pt-2.5 pb-[calc(0.625rem+var(--safe-bottom))] shadow-[0_-10px_30px_-15px_rgba(0,0,0,0.6)] md:hidden">
          <button
            onClick={save}
            disabled={isLoading}
            className="flex-1 rounded-[var(--radius-md)] bg-accent-soft px-3 py-2.5 text-sm font-bold text-accent"
          >
            {isLoading ? "Сохраняю…" : "＋ В заметки"}
          </button>
          {aiEnabled && (
            <button
              onClick={askAi}
              className="flex flex-1 items-center justify-center gap-1.5 rounded-[var(--radius-md)] px-3 py-2.5 text-sm font-bold text-accent-fg"
              style={{ background: "var(--gradient)" }}
            >
              <Sparkles size={15} /> Спросить у ИИ
            </button>
          )}
          <button
            onClick={() => {
              window.getSelection()?.removeAllRanges();
              hide();
            }}
            aria-label="Закрыть"
            className="shrink-0 rounded-[var(--radius-md)] p-2 text-faint hover:text-fg"
          >
            <X size={18} />
          </button>
        </div>
      )}

      {aiEnabled && (
        <AiChatModal
          open={aiContext !== null}
          onClose={() => setAiContext(null)}
          context={aiContext ?? ""}
          lessonId={lessonId}
          perMinute={aiStatus?.perMinute ?? 6}
          perHour={aiStatus?.perHour ?? 60}
        />
      )}
    </div>
  );
}
