import { useEffect, useRef, useState } from "react";

import { useCompleteLessonMutation } from "@/features/learning/api/lessonApi";
import { apiErrorMessage } from "@/shared/api/baseApi";
import type {
  Certificate, LessonProgress, TextContent } from "@/shared/types";
import { Button, Card } from "@/shared/ui";
import { Check } from "lucide-react";
import { useToast } from "@/shared/ui/ToastProvider";

import LessonResources from "./LessonResources";
import Markdown, { headingId } from "./Markdown";

// Урок-теория: markdown плюс кнопка «Прочитал».
export default function TextLesson({
  lessonId,
  content,
  progress,
  onDone,
}: {
  lessonId: string;
  content: TextContent;
  progress?: LessonProgress;
  onDone: (certificate?: Certificate | null) => void;
}) {
  const [complete, { isLoading }] = useCompleteLessonMutation();
  const toast = useToast();
  const startedAt = useRef(Date.now());
  const [done, setDone] = useState(progress?.status === "completed");

  useEffect(() => {
    setDone(progress?.status === "completed");
    startedAt.current = Date.now();
  }, [lessonId, progress?.status]);

  const markDone = async () => {
    const seconds = Math.round((Date.now() - startedAt.current) / 1000);
    try {
      const result = await complete({ id: lessonId, seconds }).unwrap();
      setDone(true);
      toast.success("Урок отмечен пройденным");
      onDone(result.certificate);
    } catch (err) {
      toast.error(apiErrorMessage(err));
    }
  };

  const body = content.body ?? "_Материал урока пока не заполнен._";

  // Оглавление по заголовкам второго уровня — для быстрого перехода к разделу.
  const headings = body
    .split("\n")
    .filter((line) => line.startsWith("## "))
    .map((line) => line.slice(3).trim());

  return (
    <>
      {/* На телефоне теория читается прямо на фоне (без «стеклянной» карточки) —
          легче и спокойнее для длинного чтения; на sm+ остаётся карточкой. */}
      <Card className="p-[var(--pad)] max-sm:border-0 max-sm:bg-transparent max-sm:p-0 max-sm:shadow-none max-sm:backdrop-blur-none sm:p-8">
        <div className="mx-auto">
          {headings.length >= 4 && (
            <nav className="mb-6 rounded-[var(--radius-md)] border border-line bg-surface-2 p-4">
              <p className="mb-2 text-xs font-bold uppercase tracking-wide text-faint">
                В этом уроке
              </p>
              <ol className="space-y-1 text-sm">
                {headings.map((title, index) => (
                  <li key={title}>
                    <a
                      href={`#${headingId(title)}`}
                      className="text-muted transition-colors hover:text-accent"
                    >
                      <span className="mr-1.5 font-bold text-faint">{index + 1}.</span>
                      {title}
                    </a>
                  </li>
                ))}
              </ol>
            </nav>
          )}

          <Markdown>{body}</Markdown>
        </div>
      </Card>

      <LessonResources items={content.resources} />

      <Card className="mt-[var(--gap)] flex flex-wrap items-center justify-between gap-3 p-[var(--pad)]">
        <p className="text-sm text-muted">
          {done
            ? "Урок пройден. Можно двигаться дальше."
            : "Дочитали до конца? Отметьте урок пройденным."}
        </p>
        <Button
          variant={done ? "secondary" : "primary"}
          icon={done ? <Check size={18} /> : undefined}
          onClick={markDone}
          loading={isLoading}
          disabled={done}
        >
          {done ? "Пройдено" : "Отметить пройденным"}
        </Button>
      </Card>
    </>
  );
}
