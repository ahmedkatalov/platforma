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
  // x — центр выделения по горизонтали, y — верх выделения (относительно контейнера).
  // Попап всегда рисуем НАД выделением.
  const [popup, setPopup] = useState<{ x: number; y: number; text: string } | null>(null);
  const [aiContext, setAiContext] = useState<string | null>(null);
  const [createNote, { isLoading }] = useCreateNoteMutation();
  const { data: aiStatus } = useGetAiStatusQuery();
  const aiEnabled = Boolean(aiStatus?.enabled);
  const toast = useToast();

  const hide = useCallback(() => setPopup(null), []);

  useEffect(() => {
    let timer = 0;

    const evaluate = () => {
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

      // Кнопки показываем только для выделений внутри урока.
      const range = selection.getRangeAt(0);
      if (!container.contains(range.commonAncestorContainer)) {
        setPopup(null);
        return;
      }

      const rect = range.getBoundingClientRect();
      // Геометрия ещё не готова (бывает в момент создания выделения) — подождём.
      if (rect.width === 0 && rect.height === 0) return;
      const host = container.getBoundingClientRect();

      setPopup({
        x: Math.min(Math.max(rect.left + rect.width / 2 - host.left, 120), host.width - 120),
        y: rect.top - host.top,
        text,
      });
    };

    // selectionchange срабатывает сразу при выделении (в т.ч. первом на телефоне),
    // поэтому кнопки больше не требуют «второго касания». Дебаунс — чтобы не
    // дёргаться, пока тянут маркеры выделения.
    const schedule = () => {
      window.clearTimeout(timer);
      timer = window.setTimeout(evaluate, 180);
    };

    document.addEventListener("selectionchange", schedule);
    document.addEventListener("mouseup", evaluate); // на десктопе — без задержки
    return () => {
      window.clearTimeout(timer);
      document.removeEventListener("selectionchange", schedule);
      document.removeEventListener("mouseup", evaluate);
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

  const actions = (
    <>
      <button
        onClick={save}
        disabled={isLoading}
        className="rounded-full px-3 py-1.5 text-sm font-bold text-accent transition-transform hover:scale-105"
      >
        {isLoading ? "Сохраняю…" : "＋ В заметки"}
      </button>
      {aiEnabled && (
        <>
          <span className="h-4 w-px bg-line" />
          <button
            onClick={askAi}
            className="flex items-center gap-1 rounded-full px-3 py-1.5 text-sm font-bold text-accent transition-transform hover:scale-105"
          >
            <Sparkles size={14} /> Спросить у ИИ
          </button>
        </>
      )}
    </>
  );

  return (
    <div ref={containerRef} className="relative">
      {/* Попапы рисуем ПЕРЕД текстом и делаем невыбираемыми (select-none), иначе
          при расширении выделения оно «утекает» в них и охватывает всю страницу.
          Оба варианта — НАД выделением; на телефоне повыше, чтобы разъехаться с
          системным меню браузера. */}
      {popup && (
        <>
          {/* Десктоп/планшет */}
          <div
            onMouseDown={(e) => {
              e.preventDefault();
              e.stopPropagation();
            }}
            className="absolute z-30 hidden select-none items-center gap-1 whitespace-nowrap rounded-full border border-line bg-surface-solid p-1 shadow-[var(--shadow-md)] md:flex"
            style={{
              left: popup.x,
              top: popup.y - 8,
              transform: "translate(-50%, -100%)",
              WebkitUserSelect: "none",
              userSelect: "none",
            }}
          >
            {actions}
          </div>

          {/* Телефон — над выделением, повыше (~на 32px выше десктопного) */}
          <div
            onMouseDown={(e) => {
              e.preventDefault();
              e.stopPropagation();
            }}
            className="absolute z-40 flex select-none items-center gap-1 whitespace-nowrap rounded-full border border-line bg-surface-solid p-1 shadow-[var(--shadow-md)] md:hidden"
            style={{
              left: popup.x,
              top: popup.y - 40,
              transform: "translate(-50%, -100%)",
              WebkitUserSelect: "none",
              userSelect: "none",
            }}
          >
            {actions}
          </div>
        </>
      )}

      {children}

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
