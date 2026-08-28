import { useEffect, useRef, useState } from "react";

import { useCompleteLessonMutation } from "@/features/learning/api/lessonApi";
import { apiErrorMessage } from "@/shared/api/baseApi";
import type {
  Certificate, LessonProgress, TextContent } from "@/shared/types";
import { Button, Card } from "@/shared/ui";
import { IconCheck } from "@/shared/ui/icons";
import { useToast } from "@/shared/ui/ToastProvider";

import Markdown from "./Markdown";

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

  return (
    <>
      <Card className="p-[var(--pad)] sm:p-8">
        <Markdown>{content.body ?? "_Материал урока пока не заполнен._"}</Markdown>
      </Card>

      <Card className="mt-[var(--gap)] flex flex-wrap items-center justify-between gap-3 p-[var(--pad)]">
        <p className="text-sm text-muted">
          {done
            ? "Урок пройден. Можно двигаться дальше."
            : "Дочитали до конца? Отметьте урок пройденным."}
        </p>
        <Button
          variant={done ? "secondary" : "primary"}
          icon={done ? <IconCheck size={18} /> : undefined}
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
