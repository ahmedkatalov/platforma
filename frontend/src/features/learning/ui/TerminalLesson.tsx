import { useEffect, useMemo, useRef, useState } from "react";

import { useCheckTerminalMutation } from "@/features/learning/api/lessonApi";
import {
  createShell,
  execute,
  prompt,
  type ShellState,
} from "@/features/learning/lib/shell";
import { apiErrorMessage } from "@/shared/api/baseApi";
import type {
  Certificate, LessonProgress, TaskState, TerminalContent } from "@/shared/types";
import { Badge, Button, Card, Progress } from "@/shared/ui";
import { Check, Terminal } from "lucide-react";
import { useToast } from "@/shared/ui/ToastProvider";

import LessonResources from "./LessonResources";
import Markdown from "./Markdown";

type Line = { kind: "input" | "output" | "note"; text: string };

// Рендер подсказки: `код` в обратных кавычках — моноширинной плашкой.
function renderHint(text: string) {
  return text.split(/`([^`]+)`/g).map((part, i) =>
    i % 2 === 1 ? (
      <code key={i} className="rounded bg-surface-solid px-1 font-mono text-accent">
        {part}
      </code>
    ) : (
      part
    ),
  );
}

// Тренажёр терминала: команды выполняет учебный эмулятор, а правильность
// решения задания подтверждает сервер.
export default function TerminalLesson({
  lessonId,
  content,
  tasks,
  progress,
  onDone,
}: {
  lessonId: string;
  content: TerminalContent;
  tasks: TaskState[];
  progress?: LessonProgress;
  onDone: (certificate?: Certificate | null) => void;
}) {
  const taskList = content.tasks ?? [];

  const [shell, setShell] = useState<ShellState>(() => createShell());
  const [lines, setLines] = useState<Line[]>([
    { kind: "note", text: "Учебный терминал. Наберите help, чтобы увидеть список команд." },
  ]);
  const [input, setInput] = useState("");
  const [history, setHistory] = useState<string[]>([]);
  const [historyIndex, setHistoryIndex] = useState(-1);
  const [solved, setSolved] = useState<Set<string>>(
    () => new Set(tasks.filter((t) => t.completedAt).map((t) => t.taskId)),
  );
  const [hintLevel, setHintLevel] = useState<Record<string, number>>({});

  const [checkTerminal] = useCheckTerminalMutation();
  const toast = useToast();

  const screenRef = useRef<HTMLDivElement>(null);
  const inputRef = useRef<HTMLInputElement>(null);
  const startedAt = useRef(Date.now());

  useEffect(() => {
    setShell(createShell());
    setLines([{ kind: "note", text: "Учебный терминал. Наберите help, чтобы увидеть список команд." }]);
    setInput("");
    setHistory([]);
    setSolved(new Set(tasks.filter((t) => t.completedAt).map((t) => t.taskId)));
    startedAt.current = Date.now();
    // tasks меняются вместе с уроком — этого достаточно для сброса
  }, [lessonId]); // eslint-disable-line react-hooks/exhaustive-deps

  useEffect(() => {
    screenRef.current?.scrollTo({ top: screenRef.current.scrollHeight });
  }, [lines]);

  // Текущее задание — первое нерешённое.
  const current = useMemo(
    () => taskList.find((task) => !solved.has(task.id)) ?? null,
    [taskList, solved],
  );

  const allDone = taskList.length > 0 && solved.size >= taskList.length;

  const push = (...items: Line[]) => setLines((current) => [...current, ...items]);

  const run = async (raw: string) => {
    const command = raw.trim();
    push({ kind: "input", text: `${prompt(shell)} ${command}` });

    if (command) {
      setHistory((items) => [...items, command]);
      setHistoryIndex(-1);
    }

    const result = execute(shell, command);
    setShell(result.state);

    if (result.clear) {
      setLines([]);
    } else if (result.output) {
      push({ kind: "output", text: result.output });
    }

    if (!command || !current) return;

    // Проверяем команду на сервере — эталон хранится только там.
    try {
      const seconds = Math.round((Date.now() - startedAt.current) / 1000);
      const check = await checkTerminal({
        id: lessonId,
        taskId: current.id,
        command,
        usedHint: Boolean(hintLevel[current.id]),
        seconds,
      }).unwrap();

      if (check.solved) {
        setSolved((items) => new Set([...items, current.id]));
        push({ kind: "note", text: `✓ ${check.message}` });

        if (check.lessonComplete) {
          toast.success("Все задания выполнены!");
          onDone(check.certificate);
        }
      }
    } catch (err) {
      toast.error(apiErrorMessage(err));
    }
  };

  const onKeyDown = (event: React.KeyboardEvent<HTMLInputElement>) => {
    if (event.key === "Enter") {
      void run(input);
      setInput("");
      return;
    }

    // Стрелки — история команд, как в настоящем шелле.
    if (event.key === "ArrowUp") {
      event.preventDefault();
      if (history.length === 0) return;
      const next = historyIndex < 0 ? history.length - 1 : Math.max(0, historyIndex - 1);
      setHistoryIndex(next);
      setInput(history[next]);
    }
    if (event.key === "ArrowDown") {
      event.preventDefault();
      if (historyIndex < 0) return;
      const next = historyIndex + 1;
      if (next >= history.length) {
        setHistoryIndex(-1);
        setInput("");
      } else {
        setHistoryIndex(next);
        setInput(history[next]);
      }
    }
    if (event.key === "l" && event.ctrlKey) {
      event.preventDefault();
      setLines([]);
    }
  };

  return (
    <>
      {content.intro && (
        <Card className="mb-[var(--gap)] p-[var(--pad)]">
          <Markdown>{content.intro}</Markdown>
        </Card>
      )}

      <div className="grid gap-[var(--gap)] lg:grid-cols-5">
        {/* Мобильный фокус: текущее задание видно прямо над терминалом,
            чтобы не листать вниз к списку заданий во время ввода команд. */}
        {current && !allDone && (
          <div className="rounded-[var(--radius-md)] border border-accent-border bg-accent-soft p-3 lg:hidden">
            <div className="mb-1 flex items-center justify-between gap-2">
              <span className="text-[11px] font-bold uppercase tracking-wide text-accent">
                Задание {Math.min(solved.size + 1, taskList.length)} из {taskList.length}
              </span>
              <span className="w-20 shrink-0">
                <Progress value={taskList.length ? (solved.size / taskList.length) * 100 : 0} />
              </span>
            </div>
            <p className="text-sm font-medium text-fg">{current.prompt}</p>
            {current.predict && (
              <p className="mt-1.5 text-xs text-muted">
                <span aria-hidden>🔮 </span>
                {renderHint(current.predict)}
              </p>
            )}
          </div>
        )}
        {allDone && (
          <div className="flex items-center gap-2 rounded-[var(--radius-md)] bg-[var(--success-soft)] p-3 text-sm font-medium text-success lg:hidden">
            <Check size={16} /> Все задания выполнены. Урок засчитан.
          </div>
        )}

        {/* Терминал — ввод прямо в строке приглашения, как в настоящем терминале. */}
        <Card className="overflow-hidden lg:col-span-3">
          <div className="flex items-center gap-2 border-b border-line px-4 py-2.5">
            <span className="flex gap-1.5">
              <span className="h-3 w-3 rounded-full bg-[var(--danger)]" />
              <span className="h-3 w-3 rounded-full bg-[var(--warning)]" />
              <span className="h-3 w-3 rounded-full bg-[var(--success)]" />
            </span>
            <span className="ml-2 flex items-center gap-1.5 text-xs font-semibold text-muted">
              <Terminal size={14} />
              {content.shell ?? "student@devops"}
            </span>
          </div>

          <div
            ref={screenRef}
            className="h-[16rem] overflow-y-auto overscroll-contain bg-[var(--bg-deep)] p-4 font-mono text-[13px] leading-relaxed sm:h-[20rem] lg:h-[26rem]"
            onClick={() => inputRef.current?.focus()}
          >
            {lines.map((line, i) => (
              <pre
                key={i}
                className={`whitespace-pre-wrap break-words ${
                  line.kind === "input"
                    ? "text-accent"
                    : line.kind === "note"
                      ? "text-success"
                      : "text-fg"
                }`}
              >
                {line.text}
              </pre>
            ))}

            {/* Приглашение и курсор ввода — прямо в потоке вывода. */}
            <div className="flex items-center gap-2">
              <span className="shrink-0 text-accent">{prompt(shell)}</span>
              <input
                ref={inputRef}
                value={input}
                onChange={(e) => setInput(e.target.value)}
                onKeyDown={onKeyDown}
                className="min-w-0 flex-1 bg-transparent font-mono text-[13px] text-fg outline-none"
                autoComplete="off"
                autoCapitalize="off"
                autoCorrect="off"
                spellCheck={false}
                enterKeyHint="send"
                aria-label="Командная строка"
              />
            </div>
          </div>
        </Card>

        {/* Задания */}
        <Card className="p-[var(--pad)] lg:col-span-2">
          <div className="mb-4 flex items-center justify-between">
            <h2 className="text-base font-bold text-fg">Задания</h2>
            <Badge tone={allDone ? "success" : "accent"}>
              {solved.size} из {taskList.length}
            </Badge>
          </div>

          <div className="mb-4">
            <Progress
              value={taskList.length ? (solved.size / taskList.length) * 100 : 0}
              tone={allDone ? "var(--success)" : undefined}
            />
          </div>

          <ol className="space-y-2">
            {taskList.map((task, i) => {
              const done = solved.has(task.id);
              const active = current?.id === task.id;

              return (
                <li
                  key={task.id}
                  className={`rounded-[var(--radius-md)] border p-3 transition-colors ${
                    active
                      ? "border-[var(--accent)] bg-accent-soft"
                      : done
                        ? "border-line bg-[var(--success-soft)]"
                        : "border-line opacity-70"
                  }`}
                >
                  <div className="flex items-start gap-2">
                    <span
                      className={`grid h-5 w-5 shrink-0 place-items-center rounded-full text-[11px] font-bold ${
                        done ? "bg-[var(--success)] text-white" : "bg-surface-2 text-muted"
                      }`}
                    >
                      {done ? <Check size={12} /> : i + 1}
                    </span>
                    <p className={`flex-1 text-sm ${done ? "text-muted line-through" : "text-fg"}`}>
                      {task.prompt}
                    </p>
                  </div>

                  {active && task.predict && !done && (
                    <div className="mt-2 flex items-start gap-2 rounded-[var(--radius-sm)] border border-accent-border bg-accent-soft px-2.5 py-1.5 text-xs text-muted">
                      <span aria-hidden>🔮</span>
                      <span>
                        <b className="text-accent">Прежде чем вводить:</b> {renderHint(task.predict)}
                      </span>
                    </div>
                  )}

                  {active &&
                    (() => {
                      const levels =
                        task.hints && task.hints.length
                          ? task.hints
                          : task.hint
                            ? [task.hint]
                            : [];
                      if (!levels.length) return null;
                      const shown = hintLevel[task.id] ?? 0;
                      const isLast = (i: number) => i === levels.length - 1 && levels.length > 1;
                      return (
                        <div className="mt-2 space-y-1.5">
                          {levels.slice(0, shown).map((h, hi) => (
                            <p
                              key={hi}
                              className="rounded-[var(--radius-sm)] bg-surface-2 px-2.5 py-1.5 text-xs text-muted"
                            >
                              <span className="font-semibold text-faint">
                                {isLast(hi) ? "Команда" : `Подсказка ${hi + 1}`}:{" "}
                              </span>
                              {renderHint(h)}
                            </p>
                          ))}
                          {shown < levels.length && (
                            <button
                              className="text-xs font-semibold text-accent hover:underline"
                              onClick={() =>
                                setHintLevel((c) => ({ ...c, [task.id]: shown + 1 }))
                              }
                            >
                              {shown === 0
                                ? "Показать подсказку"
                                : isLast(shown)
                                  ? "Показать команду"
                                  : "Ещё подсказку"}
                              {levels.length > 1 && (
                                <span className="ml-1 text-faint">
                                  ({shown}/{levels.length})
                                </span>
                              )}
                            </button>
                          )}
                        </div>
                      );
                    })()}
                </li>
              );
            })}
          </ol>

          {allDone && (
            <div className="mt-4 rounded-[var(--radius-md)] bg-[var(--success-soft)] p-3 text-sm text-success">
              Все задания выполнены. Урок засчитан.
            </div>
          )}

          {content.challenge && (
            <div className="mt-4 rounded-[var(--radius-md)] border border-line bg-surface-2 p-3">
              <p className="mb-1 flex items-center gap-1.5 text-sm font-bold text-fg">
                <span aria-hidden>🔧</span> Измените и попробуйте сами
              </p>
              <div className="text-sm text-muted [&_code]:text-accent">
                <Markdown>{content.challenge}</Markdown>
              </div>
            </div>
          )}

          {content.debug && (
            <div className="mt-3 rounded-[var(--radius-md)] border border-line bg-surface-2 p-3">
              <p className="mb-1 flex items-center gap-1.5 text-sm font-bold text-fg">
                <span aria-hidden>🐞</span> Если что-то сломалось
              </p>
              <div className="text-sm text-muted [&_code]:text-accent">
                <Markdown>{content.debug}</Markdown>
              </div>
            </div>
          )}

          {progress?.attempts ? (
            <p className="mt-3 text-xs text-faint">Попыток: {progress.attempts}</p>
          ) : null}

          <div className="mt-4 flex gap-2">
            <Button
              onClick={() => {
                setShell(createShell());
                setLines([{ kind: "note", text: "Терминал перезапущен." }]);
              }}
            >
              Сбросить терминал
            </Button>
            {allDone && (
              <Button variant="primary" onClick={() => onDone()}>
                Дальше
              </Button>
            )}
          </div>
        </Card>
      </div>

      <LessonResources items={content.resources} />
    </>
  );
}
