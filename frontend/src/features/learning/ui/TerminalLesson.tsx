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
  const [hintShown, setHintShown] = useState<Record<string, boolean>>({});

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
        usedHint: Boolean(hintShown[current.id]),
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
        {/* Терминал */}
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
            className="h-[26rem] overflow-y-auto bg-[var(--bg-deep)] p-4 font-mono text-[13px] leading-relaxed"
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

            <div className="flex items-center gap-2">
              <span className="shrink-0 text-accent">{prompt(shell)}</span>
              <input
                ref={inputRef}
                value={input}
                onChange={(e) => setInput(e.target.value)}
                onKeyDown={onKeyDown}
                className="flex-1 bg-transparent font-mono text-[13px] text-fg outline-none"
                autoComplete="off"
                autoCapitalize="off"
                spellCheck={false}
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

                  {active && task.hint && (
                    <div className="mt-2">
                      {hintShown[task.id] ? (
                        <p className="rounded-[var(--radius-sm)] bg-surface-2 px-2.5 py-1.5 text-xs text-muted">
                          Подсказка: {task.hint}
                        </p>
                      ) : (
                        <button
                          className="text-xs font-semibold text-accent hover:underline"
                          onClick={() => setHintShown((c) => ({ ...c, [task.id]: true }))}
                        >
                          Показать подсказку
                        </button>
                      )}
                    </div>
                  )}
                </li>
              );
            })}
          </ol>

          {allDone && (
            <div className="mt-4 rounded-[var(--radius-md)] bg-[var(--success-soft)] p-3 text-sm text-success">
              Все задания выполнены. Урок засчитан.
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
