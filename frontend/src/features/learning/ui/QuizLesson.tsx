import { useEffect, useMemo, useRef, useState } from "react";

import { useSubmitQuizMutation } from "@/features/learning/api/lessonApi";
import { apiErrorMessage } from "@/shared/api/baseApi";
import type {
  Certificate, LessonProgress, QuizContent, QuizResult } from "@/shared/types";
import { Badge, Button, Card, Progress } from "@/shared/ui";
import { IconCheck, IconClose, IconClock } from "@/shared/ui/icons";
import { useToast } from "@/shared/ui/ToastProvider";

import Markdown from "./Markdown";

type Answers = Record<string, string[]>;

// Квиз: вопросы по одному, замер времени на каждый, проверка на сервере.
export default function QuizLesson({
  lessonId,
  content,
  progress,
  onDone,
}: {
  lessonId: string;
  content: QuizContent;
  progress?: LessonProgress;
  onDone: (certificate?: Certificate | null) => void;
}) {
  const questions = content.questions ?? [];
  const passScore = content.passScore ?? 70;

  const [index, setIndex] = useState(0);
  const [answers, setAnswers] = useState<Answers>({});
  const [timings, setTimings] = useState<Record<string, number>>({});
  const [result, setResult] = useState<QuizResult | null>(null);

  const [submitQuiz, { isLoading }] = useSubmitQuizMutation();
  const toast = useToast();

  const questionStart = useRef(Date.now());
  const quizStart = useRef(Date.now());

  useEffect(() => {
    setIndex(0);
    setAnswers({});
    setTimings({});
    setResult(null);
    questionStart.current = Date.now();
    quizStart.current = Date.now();
  }, [lessonId]);

  const question = questions[index];
  const answered = Object.keys(answers).length;

  const resultByQuestion = useMemo(() => {
    const map = new Map<string, QuizResult["questions"][number]>();
    result?.questions.forEach((item) => map.set(item.questionId, item));
    return map;
  }, [result]);

  if (questions.length === 0) {
    return (
      <Card className="p-[var(--pad)]">
        <p className="text-sm text-muted">Вопросы для этого квиза ещё не добавлены.</p>
      </Card>
    );
  }

  // Фиксируем время, потраченное на текущий вопрос, и переходим дальше.
  const rememberTiming = () => {
    if (!question) return;
    const spent = Math.round((Date.now() - questionStart.current) / 1000);
    setTimings((current) => ({
      ...current,
      [question.id]: (current[question.id] ?? 0) + spent,
    }));
    questionStart.current = Date.now();
  };

  const choose = (optionId: string) => {
    if (!question) return;
    setAnswers((current) => {
      const picked = current[question.id] ?? [];
      if (question.multiple) {
        return {
          ...current,
          [question.id]: picked.includes(optionId)
            ? picked.filter((id) => id !== optionId)
            : [...picked, optionId],
        };
      }
      return { ...current, [question.id]: [optionId] };
    });
  };

  const goTo = (next: number) => {
    rememberTiming();
    setIndex(Math.max(0, Math.min(questions.length - 1, next)));
  };

  const submit = async () => {
    rememberTiming();
    const totalSeconds = Math.round((Date.now() - quizStart.current) / 1000);

    const payload = questions.map((q) => ({
      questionId: q.id,
      optionIds: answers[q.id] ?? [],
      secondsSpent: timings[q.id] ?? 0,
    }));

    try {
      const data = await submitQuiz({ id: lessonId, answers: payload, seconds: totalSeconds }).unwrap();
      setResult(data);
      if (data.passed) {
        toast.success(`Квиз пройден: ${Math.round(data.score)}%`);
        onDone(data.certificate);
      } else {
        toast.error(`Набрано ${Math.round(data.score)}% — нужно ${Math.round(data.passScore)}%`);
      }
    } catch (err) {
      toast.error(apiErrorMessage(err));
    }
  };

  const retry = () => {
    setResult(null);
    setAnswers({});
    setTimings({});
    setIndex(0);
    questionStart.current = Date.now();
    quizStart.current = Date.now();
  };

  // --- Экран результата ---
  if (result) {
    return (
      <>
        <Card className="p-[var(--pad)] sm:p-8">
          <div className="flex flex-wrap items-center justify-between gap-4">
            <div>
              <p className="text-sm text-muted">Результат</p>
              <p className="text-4xl font-extrabold text-fg">{Math.round(result.score)}%</p>
              <p className="mt-1 text-sm text-muted">
                {result.correctCount} из {result.totalCount} · порог {Math.round(result.passScore)}%
              </p>
            </div>
            <Badge tone={result.passed ? "success" : "danger"}>
              {result.passed ? "Квиз пройден" : "Нужно повторить"}
            </Badge>
          </div>

          <div className="mt-4">
            <Progress
              value={result.score}
              tone={result.passed ? "var(--success)" : "var(--danger)"}
            />
          </div>

          <div className="mt-4 flex flex-wrap gap-2">
            <Button onClick={retry}>Пройти заново</Button>
            {result.passed && (
              <Button variant="primary" onClick={() => onDone()}>
                Дальше
              </Button>
            )}
          </div>
        </Card>

        <div className="mt-[var(--gap)] space-y-[var(--gap)]">
          {questions.map((q, i) => {
            const outcome = resultByQuestion.get(q.id);
            const correctIds = new Set(outcome?.correctOptionIds ?? []);
            const chosenIds = new Set(outcome?.chosenOptionIds ?? []);

            return (
              <Card key={q.id} className="p-[var(--pad)]">
                <div className="mb-3 flex items-start gap-3">
                  <span
                    className={`grid h-7 w-7 shrink-0 place-items-center rounded-full text-xs font-bold ${
                      outcome?.correct
                        ? "bg-[var(--success-soft)] text-success"
                        : "bg-[var(--danger-soft)] text-danger"
                    }`}
                  >
                    {outcome?.correct ? <IconCheck size={16} /> : <IconClose size={16} />}
                  </span>
                  <p className="flex-1 font-semibold text-fg">
                    <span className="text-faint">{i + 1}. </span>
                    {q.text}
                  </p>
                  <span className="shrink-0 text-xs text-faint">
                    {timings[q.id] ?? 0} c
                  </span>
                </div>

                <ul className="space-y-1.5">
                  {q.options.map((option) => {
                    const isCorrect = correctIds.has(option.id);
                    const isChosen = chosenIds.has(option.id);

                    return (
                      <li
                        key={option.id}
                        className={`flex items-center gap-2 rounded-[var(--radius-md)] border px-3 py-2 text-sm ${
                          isCorrect
                            ? "border-[var(--success)] bg-[var(--success-soft)] text-fg"
                            : isChosen
                              ? "border-[var(--danger)] bg-[var(--danger-soft)] text-fg"
                              : "border-line text-muted"
                        }`}
                      >
                        <span className="flex-1">{option.text}</span>
                        {isCorrect && <Badge tone="success">верно</Badge>}
                        {isChosen && !isCorrect && <Badge tone="danger">ваш выбор</Badge>}
                      </li>
                    );
                  })}
                </ul>

                {outcome?.explanation && (
                  <div className="mt-3 rounded-[var(--radius-md)] bg-surface-2 px-3 py-2 text-sm text-muted">
                    {outcome.explanation}
                  </div>
                )}
              </Card>
            );
          })}
        </div>
      </>
    );
  }

  // --- Прохождение ---
  const picked = answers[question.id] ?? [];

  return (
    <>
      {content.intro && index === 0 && (
        <Card className="mb-[var(--gap)] p-[var(--pad)]">
          <Markdown>{content.intro}</Markdown>
        </Card>
      )}

      <Card className="p-[var(--pad)] sm:p-8">
        <div className="mb-4 flex items-center justify-between gap-3">
          <span className="text-sm font-semibold text-muted">
            Вопрос {index + 1} из {questions.length}
          </span>
          <div className="flex items-center gap-2">
            {progress?.bestScore != null && (
              <Badge tone="accent">лучший: {Math.round(progress.bestScore)}%</Badge>
            )}
            <Badge>
              <IconClock size={12} /> порог {Math.round(passScore)}%
            </Badge>
          </div>
        </div>

        <div className="mb-5">
          <Progress value={((index + 1) / questions.length) * 100} />
        </div>

        <h2 className="mb-1 text-lg font-bold text-fg">{question.text}</h2>
        {question.multiple && (
          <p className="mb-3 text-xs text-faint">Можно выбрать несколько вариантов</p>
        )}

        <ul className="mt-4 space-y-2">
          {question.options.map((option) => {
            const active = picked.includes(option.id);
            return (
              <li key={option.id}>
                <button
                  onClick={() => choose(option.id)}
                  className={`flex w-full items-center gap-3 rounded-[var(--radius-md)] border p-3 text-left text-sm transition-colors ${
                    active
                      ? "border-[var(--accent)] bg-accent-soft text-fg"
                      : "border-line text-muted hover:bg-surface-2 hover:text-fg"
                  }`}
                >
                  <span
                    className={`grid h-5 w-5 shrink-0 place-items-center border text-accent ${
                      question.multiple ? "rounded-[0.3rem]" : "rounded-full"
                    } ${active ? "border-[var(--accent)] bg-accent-soft" : "border-line"}`}
                  >
                    {active && <IconCheck size={14} />}
                  </span>
                  {option.text}
                </button>
              </li>
            );
          })}
        </ul>

        {question.hint && (
          <p className="mt-4 rounded-[var(--radius-md)] bg-surface-2 px-3 py-2 text-xs text-muted">
            Подсказка: {question.hint}
          </p>
        )}

        <div className="mt-6 flex flex-wrap items-center justify-between gap-3">
          <Button onClick={() => goTo(index - 1)} disabled={index === 0}>
            Назад
          </Button>

          <span className="text-xs text-faint">
            отвечено {answered} из {questions.length}
          </span>

          {index < questions.length - 1 ? (
            <Button variant="primary" onClick={() => goTo(index + 1)} disabled={picked.length === 0}>
              Далее
            </Button>
          ) : (
            <Button
              variant="primary"
              onClick={submit}
              loading={isLoading}
              disabled={answered < questions.length}
            >
              Завершить квиз
            </Button>
          )}
        </div>
      </Card>
    </>
  );
}
