import { useEffect, useRef, useState } from "react";

import { useCheckCodeMutation } from "@/features/learning/api/lessonApi";
import { apiErrorMessage } from "@/shared/api/baseApi";
import type {
  Certificate, CodeCheckResult, CodeContent, LessonProgress } from "@/shared/types";
import { Badge, Button, Card, Progress } from "@/shared/ui";
import { Check, ChevronRight, X } from "lucide-react";
import { useToast } from "@/shared/ui/ToastProvider";

import CodeEditor from "./CodeEditor";
import LessonResources from "./LessonResources";
import Markdown from "./Markdown";

// Практика с кодом: студент правит конфигурацию, сервер прогоняет проверки.
export default function CodeLesson({
  lessonId,
  content,
  progress,
  onDone,
}: {
  lessonId: string;
  content: CodeContent;
  progress?: LessonProgress;
  onDone: (certificate?: Certificate | null) => void;
}) {
  const storageKey = `platforma.code.${lessonId}`;

  const [code, setCode] = useState(() => {
    return localStorage.getItem(storageKey) ?? content.starter ?? "";
  });
  const [result, setResult] = useState<CodeCheckResult | null>(null);
  const [showHint, setShowHint] = useState(false);

  const [checkCode, { isLoading }] = useCheckCodeMutation();
  const toast = useToast();
  const startedAt = useRef(Date.now());

  useEffect(() => {
    setCode(localStorage.getItem(storageKey) ?? content.starter ?? "");
    setResult(null);
    setShowHint(false);
    startedAt.current = Date.now();
  }, [lessonId, content.starter, storageKey]);

  // Черновик не теряется при переходе между уроками.
  useEffect(() => {
    localStorage.setItem(storageKey, code);
  }, [code, storageKey]);

  const submit = async () => {
    const seconds = Math.round((Date.now() - startedAt.current) / 1000);
    try {
      const data = await checkCode({ id: lessonId, code, seconds }).unwrap();
      setResult(data);
      if (data.passed) {
        toast.success("Все проверки пройдены");
        // Не перелистываем сразу: показываем результат, дальше — по кнопке.
      } else {
        toast.error("Часть проверок не пройдена — посмотрите список справа");
      }
    } catch (err) {
      toast.error(apiErrorMessage(err));
    }
  };

  const passedChecks = result?.checks.filter((check) => check.ok).length ?? 0;
  const totalChecks = result?.checks.length ?? 0;

  return (
    <div className="grid gap-[var(--gap)] lg:grid-cols-5">
      <div className="space-y-[var(--gap)] lg:col-span-3">
        <Card className="p-[var(--pad)]">
          <h2 className="mb-3 text-base font-bold text-fg">Задание</h2>
          <Markdown>{content.task ?? ""}</Markdown>
        </Card>

        <Card className="p-[var(--pad)]">
          <div className="mb-3 flex items-center justify-between">
            <h2 className="text-base font-bold text-fg">Решение</h2>
            {/* Десктоп: действия в шапке. */}
            <div className="hidden gap-2 lg:flex">
              <Button onClick={() => setCode(content.starter ?? "")}>Сбросить</Button>
              <Button variant="primary" onClick={submit} loading={isLoading}>
                Проверить
              </Button>
            </div>
          </div>

          <CodeEditor value={code} onChange={setCode} language={content.language} />

          {/* Телефон: действия под редактором — логичный порядок «написал → проверил»,
              «Проверить» крупная и основная. */}
          <div className="mt-3 flex gap-2 lg:hidden">
            <Button className="flex-1" onClick={() => setCode(content.starter ?? "")}>
              Сбросить
            </Button>
            <Button className="flex-[2]" variant="primary" onClick={submit} loading={isLoading}>
              Проверить
            </Button>
          </div>
        </Card>
      </div>

      <div className="space-y-[var(--gap)] lg:col-span-2">
        <Card className="p-[var(--pad)]">
          <div className="mb-3 flex items-center justify-between">
            <h2 className="text-base font-bold text-fg">Проверки</h2>
            {result && (
              <Badge tone={result.passed ? "success" : "warning"}>
                {passedChecks} из {totalChecks}
              </Badge>
            )}
          </div>

          {!result ? (
            <p className="text-sm text-muted">
              Нажмите «Проверить» — покажем, какие требования уже выполнены.
            </p>
          ) : (
            <>
              <div className="mb-4">
                <Progress
                  value={result.score}
                  tone={result.passed ? "var(--success)" : "var(--warning)"}
                />
              </div>

              {result.passed && (
                <Button
                  variant="primary"
                  className="mb-4 w-full"
                  onClick={() => onDone(result.certificate)}
                >
                  Дальше
                  <ChevronRight size={16} />
                </Button>
              )}

              <ul className="space-y-2">
                {result.checks.map((check, i) => (
                  <li
                    key={i}
                    className={`flex items-start gap-2 rounded-[var(--radius-md)] border p-2.5 text-sm ${
                      check.ok
                        ? "border-[var(--success)] bg-[var(--success-soft)]"
                        : "border-line text-muted"
                    }`}
                  >
                    <span className={check.ok ? "text-success" : "text-faint"}>
                      {check.ok ? <Check size={16} /> : <X size={16} />}
                    </span>
                    <span className={check.ok ? "text-fg" : undefined}>{check.message}</span>
                  </li>
                ))}
              </ul>

              {result.passed && (
                <Button variant="primary" className="mt-4 w-full" onClick={() => onDone()}>
                  Дальше
                </Button>
              )}
            </>
          )}

          {content.hint && (
            <div className="mt-4 border-t border-line pt-3">
              {showHint ? (
                <p className="text-sm text-muted">
                  <span className="font-semibold text-faint">Подсказка: </span>
                  {content.hint}
                </p>
              ) : (
                <button
                  className="btn btn-ghost btn-sm -ml-2 text-accent"
                  onClick={() => setShowHint(true)}
                >
                  Показать подсказку
                </button>
              )}
            </div>
          )}

          {progress?.attempts ? (
            <p className="mt-3 text-xs text-faint">
              Попыток: {progress.attempts}
              {progress.bestScore != null && ` · лучший результат ${Math.round(progress.bestScore)}%`}
            </p>
          ) : null}
        </Card>

        <LessonResources items={content.resources} />
      </div>
    </div>
  );
}
