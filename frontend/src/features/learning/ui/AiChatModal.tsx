import { useEffect, useRef, useState, type KeyboardEvent } from "react";
import { Info, Send, Sparkles } from "lucide-react";

import { useAskAiMutation } from "@/shared/api/meApi";
import { apiErrorMessage } from "@/shared/api/baseApi";
import type { AiMessage } from "@/shared/types";
import { Modal, Spinner } from "@/shared/ui";
import { useToast } from "@/shared/ui/ToastProvider";
import Markdown from "@/features/learning/ui/Markdown";

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
  const [askAi, { isLoading }] = useAskAiMutation();
  const toast = useToast();
  const scrollRef = useRef<HTMLDivElement>(null);

  // Новое выделение (или закрытие) — начинаем диалог заново.
  useEffect(() => {
    if (open) {
      setMessages([]);
      setInput("");
    }
  }, [open, context]);

  useEffect(() => {
    scrollRef.current?.scrollTo({ top: scrollRef.current.scrollHeight, behavior: "smooth" });
  }, [messages, isLoading]);

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
    <Modal
      open={open}
      onClose={onClose}
      title="Спросить у ИИ"
      width="42rem"
      footer={
        <div className="w-full">
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
      }
    >
      {/* Дружелюбная памятка: сколько можно спрашивать и про возможную задержку. */}
      <div className="mb-4 flex items-start gap-2 rounded-[var(--radius-md)] bg-surface-2 px-3 py-2 text-xs text-muted">
        <Info size={15} className="mt-0.5 shrink-0 text-accent" />
        <span>
          На каждого студента — до {perMinute} вопросов в минуту и {perHour} в час. Если ИИ сейчас
          загружен, он может ответить не сразу — просто попробуйте чуть позже.
        </span>
      </div>

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

      <div ref={scrollRef} className="max-h-[46vh] space-y-3 overflow-y-auto">
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
    </Modal>
  );
}
