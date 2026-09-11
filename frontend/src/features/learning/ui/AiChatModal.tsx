import { useEffect, useRef, useState, type KeyboardEvent } from "react";
import clsx from "clsx";
import { Info, Maximize2, Minimize2, Send, Sparkles, X } from "lucide-react";

import { useAskAiMutation } from "@/shared/api/meApi";
import { apiErrorMessage } from "@/shared/api/baseApi";
import type { AiMessage } from "@/shared/types";
import { Spinner } from "@/shared/ui";
import { useToast } from "@/shared/ui/ToastProvider";
import Markdown from "@/features/learning/ui/Markdown";

const NOTICE_KEY = "platforma.aiNoticeHidden";
const EXPANDED_KEY = "platforma.aiChatExpanded";

function readFlag(key: string): boolean {
  try {
    return localStorage.getItem(key) === "1";
  } catch {
    return false;
  }
}

function writeFlag(key: string, on: boolean) {
  try {
    if (on) localStorage.setItem(key, "1");
    else localStorage.removeItem(key);
  } catch {
    /* приватный режим — просто не сохраняем */
  }
}

export default function AiChatModal({
  open,
  onClose,
  context,
  lessonId,
  perMinute = 6,
  perHour = 60,
}: {
  open: boolean;
  onClose: () => void;
  context: string;
  lessonId?: string;
  perMinute?: number;
  perHour?: number;
}) {
  const [messages, setMessages] = useState<AiMessage[]>([]);
  const [input, setInput] = useState("");
  const [noticeHidden, setNoticeHidden] = useState(false);
  const [expanded, setExpanded] = useState(false);
  const [askAi, { isLoading }] = useAskAiMutation();
  const toast = useToast();
  const scrollRef = useRef<HTMLDivElement>(null);

  // Настройки студента (память между сессиями).
  useEffect(() => {
    setNoticeHidden(readFlag(NOTICE_KEY));
    setExpanded(readFlag(EXPANDED_KEY));
  }, []);

  // Новое выделение (или закрытие) — начинаем диалог заново.
  useEffect(() => {
    if (open) {
      setMessages([]);
      setInput("");
    }
  }, [open, context]);

  // Esc закрывает, фон не скроллится, пока чат открыт.
  useEffect(() => {
    if (!open) return;
    const onKey = (e: globalThis.KeyboardEvent) => e.key === "Escape" && onClose();
    document.addEventListener("keydown", onKey);
    document.body.style.overflow = "hidden";
    return () => {
      document.removeEventListener("keydown", onKey);
      document.body.style.overflow = "";
    };
  }, [open, onClose]);

  useEffect(() => {
    scrollRef.current?.scrollTo({ top: scrollRef.current.scrollHeight, behavior: "smooth" });
  }, [messages, isLoading]);

  if (!open) return null;

  const hideNotice = () => {
    setNoticeHidden(true);
    writeFlag(NOTICE_KEY, true);
  };

  const toggleExpanded = () => {
    setExpanded((v) => {
      writeFlag(EXPANDED_KEY, !v);
      return !v;
    });
  };

  const send = async () => {
    const question = input.trim();
    if (!question || isLoading) return;

    const next: AiMessage[] = [...messages, { role: "user", text: question }];
    setMessages(next);
    setInput("");
    try {
      const res = await askAi({ lessonId, context, messages: next }).unwrap();
      setMessages([...next, { role: "assistant", text: res.answer }]);
    } catch (err) {
      // Оставляем вопрос студента в истории — можно повторить.
      toast.error(apiErrorMessage(err, "ИИ не смог ответить, попробуйте ещё раз"));
    }
  };

  const onKeyDown = (e: KeyboardEvent<HTMLTextAreaElement>) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      void send();
    }
  };

  return (
    <div
      className="fixed inset-0 z-50 flex flex-col md:items-center md:justify-center md:p-4"
      role="dialog"
      aria-modal="true"
      aria-label="Спросить у ИИ"
    >
      {/* Затемнение — только на планшете/десктопе (на телефоне чат во весь экран). */}
      <div
        className="absolute inset-0 hidden bg-black/65 backdrop-blur-[2px] md:block"
        onClick={onClose}
        aria-hidden="true"
      />

      {/* Панель: телефон — во весь экран; десктоп — карточка, размер зависит от «развернуть». */}
      <div
        className={clsx(
          "relative z-10 flex h-full w-full flex-col overflow-hidden bg-surface-solid",
          "pt-[var(--safe-top)] pb-[var(--safe-bottom)] md:p-0",
          "md:rounded-[var(--radius-lg)] md:border md:border-line md:shadow-[var(--shadow-lg)]",
          expanded
            ? "md:h-[92vh] md:w-[min(94vw,72rem)]"
            : "md:h-[85vh] md:max-h-[44rem] md:w-[min(92vw,46rem)]",
        )}
      >
        {/* Шапка */}
        <div className="flex items-center justify-between gap-3 border-b border-line px-4 py-3 sm:px-5">
          <h2 className="flex items-center gap-2 text-base font-bold text-fg">
            <Sparkles size={18} className="text-accent" /> Спросить у ИИ
          </h2>
          <div className="flex items-center gap-1">
            <button
              className="btn btn-ghost btn-icon btn-sm hidden md:inline-flex"
              onClick={toggleExpanded}
              aria-label={expanded ? "Свернуть окно" : "Развернуть окно"}
              title={expanded ? "Свернуть" : "Развернуть"}
            >
              {expanded ? <Minimize2 size={18} /> : <Maximize2 size={18} />}
            </button>
            <button
              className="btn btn-ghost btn-icon btn-sm"
              onClick={onClose}
              aria-label="Закрыть"
            >
              <X size={18} />
            </button>
          </div>
        </div>

        {/* Тело */}
        <div className="flex-1 overflow-y-auto px-4 py-4 sm:px-5">
          {/* Памятку студент может закрыть — больше не появится. */}
          {!noticeHidden && (
            <div className="mb-4 flex items-start gap-2 rounded-[var(--radius-md)] bg-surface-2 px-3 py-2 text-xs text-muted">
              <Info size={15} className="mt-0.5 shrink-0 text-accent" />
              <span className="flex-1">
                На каждого студента — до {perMinute} вопросов в минуту и {perHour} в час. Если ИИ
                сейчас загружен, он может ответить не сразу — просто попробуйте чуть позже.
              </span>
              <button
                className="shrink-0 rounded p-0.5 text-faint hover:text-fg"
                onClick={hideNotice}
                aria-label="Скрыть подсказку"
                title="Больше не показывать"
              >
                <X size={14} />
              </button>
            </div>
          )}

          {/* Контекст: выделенный фрагмент урока. */}
          {context && (
            <div className="mb-4 rounded-[var(--radius-md)] border-l-2 border-accent bg-surface-2 px-3 py-2">
              <p className="mb-1 text-[11px] font-semibold uppercase tracking-wide text-faint">
                Ваш фрагмент
              </p>
              <p className="line-clamp-4 text-xs text-muted">{context}</p>
            </div>
          )}

          {messages.length === 0 && !isLoading && (
            <div className="grid place-items-center gap-2 py-8 text-center">
              <Sparkles size={28} className="text-accent" />
              <p className="text-sm text-muted">
                Спросите что угодно по выделенному фрагменту: «объясни проще», «зачем это»,
                «приведи пример».
              </p>
            </div>
          )}

          <div ref={scrollRef} className="space-y-3">
            {messages.map((m, i) =>
              m.role === "user" ? (
                <div key={i} className="flex justify-end">
                  <div className="max-w-[85%] rounded-[var(--radius-md)] bg-accent-soft px-3 py-2 text-sm text-fg">
                    {m.text}
                  </div>
                </div>
              ) : (
                <div key={i} className="flex justify-start">
                  <div className="prose-sm max-w-[90%] rounded-[var(--radius-md)] border border-line bg-surface px-3 py-2 text-sm text-fg">
                    <Markdown>{m.text}</Markdown>
                  </div>
                </div>
              ),
            )}
            {isLoading && (
              <div className="flex items-center gap-2 px-1 text-sm text-muted">
                <Spinner size={14} /> ИИ думает…
              </div>
            )}
          </div>
        </div>

        {/* Поле ввода */}
        <div className="border-t border-line px-4 py-3 sm:px-5">
          <div className="flex items-end gap-2">
            <textarea
              value={input}
              onChange={(e) => setInput(e.target.value)}
              onKeyDown={onKeyDown}
              rows={2}
              placeholder="Задайте вопрос по этому фрагменту…"
              className="min-h-[2.75rem] flex-1 resize-none rounded-[var(--radius-md)] border border-line bg-surface-2 px-3 py-2 text-sm text-fg outline-none focus:border-accent-border"
            />
            <button
              className="btn btn-primary btn-icon shrink-0"
              onClick={send}
              disabled={isLoading || !input.trim()}
              aria-label="Отправить"
            >
              {isLoading ? <Spinner size={16} /> : <Send size={16} />}
            </button>
          </div>
          <p className="mt-1.5 text-[11px] text-faint">
            Ответы генерирует ИИ (Google Gemini) — они могут быть неточными. Enter — отправить,
            Shift+Enter — новая строка.
          </p>
        </div>
      </div>
    </div>
  );
}
