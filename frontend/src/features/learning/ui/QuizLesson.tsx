import { useEffect, useMemo, useRef, useState } from "react";

import { useSubmitQuizMutation } from "@/features/learning/api/lessonApi";
import { apiErrorMessage } from "@/shared/api/baseApi";
import type {
  Certificate,
  LessonProgress,
  QuizContent,
  QuizItem,
  QuizQuestion,
  QuizResult,
} from "@/shared/types";
import { Badge, Button, Card, Input, Progress } from "@/shared/ui";
import { Check, ChevronRight, X, Clock } from "lucide-react";
import { useToast } from "@/shared/ui/ToastProvider";

import LessonResources from "./LessonResources";
import Markdown from "./Markdown";

// Для choice и order храним массив id, для blank — строку.
type Answers = Record<string, string[]>;
type Blanks = Record<string, string>;

const kindOf = (q: QuizQuestion) => q.type ?? "choice";

// Квиз с тремя типами вопросов: варианты, порядок шагов и вписать ответ.
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
  const [blanks, setBlanks] = useState<Blanks>({});
  const [timings, setTimings] = useState<Record<string, number>>({});
  const [result, setResult] = useState<QuizResult | null>(null);

  const [submitQuiz, { isLoading }] = useSubmitQuizMutation();
  const toast = useToast();

  const questionStart = useRef(Date.now());
  const quizStart = useRef(Date.now());

  const reset = () => {
    setIndex(0);
    setAnswers({});
    setBlanks({});
    setTimings({});
    setResult(null);
    questionStart.current = Date.now();
    quizStart.current = Date.now();
  };

  useEffect(reset, [lessonId]);

  const question = questions[index];

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

  // Переставляемая последовательность: order — это items, match — правые части.
  const seqItems = (q: QuizQuestion): QuizItem[] =>
    kindOf(q) === "match" ? (q.rights ?? []) : (q.items ?? []);

  // Текущий порядок (сохранённый) или исходный (как пришёл с сервера).
  const orderOf = (q: QuizQuestion): string[] =>
    answers[q.id] ?? seqItems(q).map((it) => it.id);

  const isAnswered = (q: QuizQuestion): boolean => {
    switch (kindOf(q)) {
      case "blank":
        return (blanks[q.id] ?? "").trim() !== "";
      case "order":
      case "match":
        return true; // всегда есть какой-то порядок
      default:
        return (answers[q.id] ?? []).length > 0;
    }
  };

  const answered = questions.filter(isAnswered).length;

  const rememberTiming = () => {
    if (!question) return;
    const spent = Math.round((Date.now() - questionStart.current) / 1000);
    setTimings((current) => ({ ...current, [question.id]: (current[question.id] ?? 0) + spent }));
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

  const moveItem = (dir: -1 | 1, itemId: string) => {
    if (!question) return;
    const order = [...orderOf(question)];
    const from = order.indexOf(itemId);
    const to = from + dir;
    if (to < 0 || to >= order.length) return;
    [order[from], order[to]] = [order[to], order[from]];
    setAnswers((current) => ({ ...current, [question.id]: order }));
  };

  const goTo = (next: number) => {
    rememberTiming();
    setIndex(Math.max(0, Math.min(questions.length - 1, next)));
  };

  const submit = async () => {
    rememberTiming();
    const totalSeconds = Math.round((Date.now() - quizStart.current) / 1000);

    const payload = questions.map((q) => {
      const base = { questionId: q.id, secondsSpent: timings[q.id] ?? 0 };
      switch (kindOf(q)) {
        case "order":
        case "match":
          return { ...base, order: orderOf(q) };
        case "blank":
          return { ...base, text: blanks[q.id] ?? "" };
        default:
          return { ...base, optionIds: answers[q.id] ?? [] };
      }
    });

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
            <Progress value={result.score} tone={result.passed ? "var(--success)" : "var(--danger)"} />
          </div>

          <div className="mt-4 flex flex-wrap gap-2">
            <Button onClick={reset}>Пройти заново</Button>
            {result.passed && (
              <Button variant="primary" onClick={() => onDone()}>
                Дальше
              </Button>
            )}
          </div>
        </Card>

        <LessonResources items={content.resources} />

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
                    {outcome?.correct ? <Check size={16} /> : <X size={16} />}
                  </span>
                  <p className="flex-1 font-semibold text-fg">
                    <span className="text-faint">{i + 1}. </span>
                    {q.text}
                    {q.review && (
                      <span className="ml-2 align-middle">
                        <Badge tone="warning">повторение</Badge>
                      </span>
                    )}
                  </p>
                  <span className="shrink-0 text-xs text-faint">{timings[q.id] ?? 0} c</span>
                </div>

                {/* Разбор ответа зависит от типа вопроса */}
                {kindOf(q) === "choice" && (
                  <ul className="space-y-1.5">
                    {(q.options ?? []).map((option) => {
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
                )}

                {(kindOf(q) === "order" || kindOf(q) === "blank" || kindOf(q) === "match") &&
                  outcome?.correctText && (
                  <div
                    className={`rounded-[var(--radius-md)] border px-3 py-2 text-sm ${
                      outcome.correct
                        ? "border-[var(--success)] bg-[var(--success-soft)] text-fg"
                        : "border-line text-muted"
                    }`}
                  >
                    <span className="text-faint">Правильно: </span>
                    <span className="font-medium text-fg">{outcome.correctText}</span>
                  </div>
                )}

                {outcome?.explanation && (
                  <div
                    className={`mt-3 rounded-[var(--radius-md)] border-l-2 px-3 py-2 text-sm ${
                      outcome.correct
                        ? "border-[var(--success)] bg-[var(--success-soft)] text-fg"
                        : "border-[var(--danger)] bg-surface-2 text-muted"
                    }`}
                  >
                    <span className="font-semibold text-fg">
                      {outcome.correct ? "Верно. " : "Разбор. "}
                    </span>
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
  const kind = kindOf(question);

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
              <Clock size={12} /> порог {Math.round(passScore)}%
            </Badge>
          </div>
        </div>

        <div className="mb-5">
          <Progress value={((index + 1) / questions.length) * 100} />
        </div>

        {question.review && (
          <div className="mb-2">
            <Badge tone="warning">повторение пройденного</Badge>
          </div>
        )}

        <h2 className="mb-1 text-lg font-bold text-fg">{question.text}</h2>

        {kind === "choice" && question.multiple && (
          <p className="mb-3 text-xs text-faint">Можно выбрать несколько вариантов</p>
        )}
        {kind === "order" && (
          <p className="mb-3 text-xs text-faint">Расставьте шаги по порядку кнопками ↑↓</p>
        )}
        {kind === "blank" && (
          <p className="mb-3 text-xs text-faint">Впишите ответ (команду или слово)</p>
        )}
        {kind === "match" && (
          <p className="mb-3 text-xs text-faint">
            Двигайте правые части кнопками ↑↓, чтобы каждая встала напротив своей левой
          </p>
        )}

        {/* Варианты */}
        {kind === "choice" && (
          <ul className="mt-4 space-y-2">
            {(question.options ?? []).map((option) => {
              const active = (answers[question.id] ?? []).includes(option.id);
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
                      {active && <Check size={14} />}
                    </span>
                    {option.text}
                  </button>
                </li>
              );
            })}
          </ul>
        )}

        {/* Порядок шагов */}
        {kind === "order" && (
          <ol className="mt-4 space-y-2">
            {orderOf(question).map((itemId, pos) => {
              const item = (question.items ?? []).find((it) => it.id === itemId);
              if (!item) return null;
              return (
                <li
                  key={itemId}
                  className="flex items-center gap-3 rounded-[var(--radius-md)] border border-line bg-surface-2 p-3 text-sm"
                >
                  <span className="grid h-6 w-6 shrink-0 place-items-center rounded-full bg-accent-soft text-xs font-bold text-accent">
                    {pos + 1}
                  </span>
                  <span className="min-w-0 flex-1 text-fg">{item.text}</span>
                  <span className="flex shrink-0 gap-1">
                    <button
                      className="btn btn-ghost h-7 w-7 !p-0"
                      onClick={() => moveItem(-1, itemId)}
                      disabled={pos === 0}
                      aria-label="Выше"
                    >
                      <ChevronRight size={14} className="-rotate-90" />
                    </button>
                    <button
                      className="btn btn-ghost h-7 w-7 !p-0"
                      onClick={() => moveItem(1, itemId)}
                      disabled={pos === orderOf(question).length - 1}
                      aria-label="Ниже"
                    >
                      <ChevronRight size={14} className="rotate-90" />
                    </button>
                  </span>
                </li>
              );
            })}
          </ol>
        )}

        {/* Сопоставление: левые части фиксированы, правые переставляются */}
        {kind === "match" && (
          <ol className="mt-4 space-y-2">
            {orderOf(question).map((rightId, pos) => {
              const left = (question.lefts ?? [])[pos];
              const right = (question.rights ?? []).find((r) => r.id === rightId);
              if (!right) return null;
              return (
                <li key={rightId} className="flex items-stretch gap-2">
                  <div className="flex min-w-0 flex-1 items-center gap-2 rounded-[var(--radius-md)] border border-line bg-surface-2 p-3 text-sm">
                    <span className="grid h-6 w-6 shrink-0 place-items-center rounded-full bg-accent-soft text-xs font-bold text-accent">
                      {pos + 1}
                    </span>
                    <span className="min-w-0 flex-1 font-medium text-fg">{left?.text}</span>
                  </div>
                  <div className="flex min-w-0 flex-1 items-center gap-2 rounded-[var(--radius-md)] border border-line bg-surface p-3 text-sm">
                    <span className="min-w-0 flex-1 text-muted">{right.text}</span>
                    <span className="flex shrink-0 gap-1">
                      <button
                        className="btn btn-ghost h-7 w-7 !p-0"
                        onClick={() => moveItem(-1, rightId)}
                        disabled={pos === 0}
                        aria-label="Выше"
                      >
                        <ChevronRight size={14} className="-rotate-90" />
                      </button>
                      <button
                        className="btn btn-ghost h-7 w-7 !p-0"
                        onClick={() => moveItem(1, rightId)}
                        disabled={pos === orderOf(question).length - 1}
                        aria-label="Ниже"
                      >
                        <ChevronRight size={14} className="rotate-90" />
                      </button>
                    </span>
                  </div>
                </li>
              );
            })}
          </ol>
        )}

        {/* Вписать ответ */}
        {kind === "blank" && (
          <div className="mt-4 max-w-md">
            <Input
              value={blanks[question.id] ?? ""}
              onChange={(e) => setBlanks((c) => ({ ...c, [question.id]: e.target.value }))}
              placeholder="Введите ответ…"
              className="font-mono"
              autoComplete="off"
              spellCheck={false}
            />
          </div>
        )}

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
            <Button variant="primary" onClick={() => goTo(index + 1)} disabled={!isAnswered(question)}>
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
