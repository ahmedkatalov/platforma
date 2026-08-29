import { useRef, useState } from "react";

import type { LessonKind } from "@/shared/types";
import { Badge, Button, Field, Input, Select, Textarea } from "@/shared/ui";
import { IconPlus, IconTrash } from "@/shared/ui/icons";

import AssetPicker from "./AssetPicker";

// Визуальные конструкторы содержимого уроков. Работают поверх того же JSON,
// который хранится в базе, — можно в любой момент переключиться на «сырой» режим.

type AnyRecord = Record<string, unknown>;

const uid = (prefix: string) => `${prefix}${Math.random().toString(36).slice(2, 7)}`;

function Section({
  title,
  children,
  action,
}: {
  title: string;
  children: React.ReactNode;
  action?: React.ReactNode;
}) {
  return (
    <div className="card-flat p-3">
      <div className="mb-3 flex items-center justify-between gap-2">
        <p className="text-sm font-bold text-fg">{title}</p>
        {action}
      </div>
      {children}
    </div>
  );
}

// --- Теория ---

function TextEditor({ value, onChange }: { value: AnyRecord; onChange: (next: AnyRecord) => void }) {
  const [pickerOpen, setPickerOpen] = useState(false);
  const areaRef = useRef<HTMLTextAreaElement>(null);
  // Запоминаем позицию курсора: при открытии окна выбора текстовое поле теряет
  // фокус, поэтому берём место вставки в момент клика по кнопке.
  const cursorRef = useRef<number | null>(null);

  const body = String(value.body ?? "");

  const rememberCursor = () => {
    cursorRef.current = areaRef.current?.selectionStart ?? null;
  };

  const openPicker = () => {
    rememberCursor();
    setPickerOpen(true);
  };

  // Вставляем markdown-картинку туда, где стоял курсор (или в конец текста).
  const insert = (markdown: string) => {
    const at = cursorRef.current ?? body.length;
    const before = body.slice(0, at).replace(/\s+$/, "");
    const after = body.slice(at).replace(/^\s+/, "");
    const next = `${before}\n\n${markdown}\n\n${after}`.replace(/^\n+/, "");
    onChange({ ...value, body: next });
  };

  return (
    <>
      <div className="mb-2 flex items-center justify-between gap-2">
        <span className="text-sm font-semibold text-muted">Текст урока (Markdown)</span>
        <Button variant="ghost" onMouseDown={rememberCursor} onClick={openPicker}>
          Вставить картинку
        </Button>
      </div>

      <Textarea
        ref={areaRef}
        value={body}
        onChange={(e) => onChange({ ...value, body: e.target.value })}
        onSelect={rememberCursor}
        onKeyUp={rememberCursor}
        onClick={rememberCursor}
        rows={16}
        className="font-mono text-xs"
        spellCheck={false}
      />
      <p className="mt-1 text-xs text-faint">
        Поддерживаются заголовки, списки, таблицы, блоки кода и картинки.
      </p>

      <AssetPicker open={pickerOpen} onClose={() => setPickerOpen(false)} onPick={insert} />
    </>
  );
}

// --- Квиз ---

type EditorOption = { id: string; text: string; correct?: boolean };
type EditorQuestion = {
  id: string;
  text: string;
  hint?: string;
  explanation?: string;
  multiple?: boolean;
  options: EditorOption[];
};

function QuizEditor({ value, onChange }: { value: AnyRecord; onChange: (next: AnyRecord) => void }) {
  const questions = (value.questions as EditorQuestion[] | undefined) ?? [];
  const passScore = Number(value.passScore ?? 70);

  const update = (next: EditorQuestion[]) => onChange({ ...value, questions: next });

  const patchQuestion = (index: number, patch: Partial<EditorQuestion>) => {
    const next = questions.map((q, i) => (i === index ? { ...q, ...patch } : q));
    update(next);
  };

  const addQuestion = () => {
    update([
      ...questions,
      {
        id: uid("q"),
        text: "",
        options: [
          { id: uid("o"), text: "", correct: true },
          { id: uid("o"), text: "", correct: false },
        ],
      },
    ]);
  };

  return (
    <div className="space-y-4">
      <div className="grid gap-4 sm:grid-cols-2">
        <Field label="Порог прохождения, %" hint="Сколько процентов нужно набрать">
          <Input
            type="number"
            min={1}
            max={100}
            value={passScore}
            onChange={(e) => onChange({ ...value, passScore: Number(e.target.value) })}
          />
        </Field>
        <Field label="Вступление" hint="Необязательный текст перед первым вопросом">
          <Input
            value={String(value.intro ?? "")}
            onChange={(e) => onChange({ ...value, intro: e.target.value })}
          />
        </Field>
      </div>

      {questions.map((question, questionIndex) => (
        <Section
          key={question.id}
          title={`Вопрос ${questionIndex + 1}`}
          action={
            <Button
              variant="ghost"
              className="h-8 !px-2 text-danger"
              onClick={() => update(questions.filter((_, i) => i !== questionIndex))}
              title="Удалить вопрос"
            >
              <IconTrash size={16} />
            </Button>
          }
        >
          <div className="space-y-3">
            <Field label="Текст вопроса">
              <Input
                value={question.text}
                onChange={(e) => patchQuestion(questionIndex, { text: e.target.value })}
                placeholder="Какой командой посмотреть запущенные контейнеры?"
              />
            </Field>

            <label className="flex cursor-pointer items-center gap-2 text-sm text-muted">
              <input
                type="checkbox"
                checked={Boolean(question.multiple)}
                onChange={(e) => patchQuestion(questionIndex, { multiple: e.target.checked })}
                className="h-4 w-4 accent-[var(--accent)]"
              />
              Несколько правильных ответов
            </label>

            <div>
              <p className="label">Варианты ответа</p>
              <ul className="space-y-2">
                {question.options.map((option, optionIndex) => (
                  <li key={option.id} className="flex items-center gap-2">
                    <label
                      className="flex shrink-0 cursor-pointer items-center gap-1.5 text-xs text-muted"
                      title="Отметьте правильные варианты"
                    >
                      <input
                        type={question.multiple ? "checkbox" : "radio"}
                        name={`correct-${question.id}`}
                        checked={Boolean(option.correct)}
                        onChange={(e) => {
                          const options = question.options.map((o, i) => {
                            if (i === optionIndex) return { ...o, correct: e.target.checked };
                            return question.multiple ? o : { ...o, correct: false };
                          });
                          patchQuestion(questionIndex, { options });
                        }}
                        className="h-4 w-4 accent-[var(--accent)]"
                      />
                      верно
                    </label>

                    <Input
                      value={option.text}
                      onChange={(e) => {
                        const options = question.options.map((o, i) =>
                          i === optionIndex ? { ...o, text: e.target.value } : o,
                        );
                        patchQuestion(questionIndex, { options });
                      }}
                      placeholder={`Вариант ${optionIndex + 1}`}
                    />

                    <Button
                      variant="ghost"
                      className="h-8 shrink-0 !px-2 text-danger"
                      onClick={() =>
                        patchQuestion(questionIndex, {
                          options: question.options.filter((_, i) => i !== optionIndex),
                        })
                      }
                      title="Удалить вариант"
                    >
                      <IconTrash size={14} />
                    </Button>
                  </li>
                ))}
              </ul>

              <Button
                variant="ghost"
                className="mt-2"
                icon={<IconPlus size={14} />}
                onClick={() =>
                  patchQuestion(questionIndex, {
                    options: [...question.options, { id: uid("o"), text: "", correct: false }],
                  })
                }
              >
                Вариант
              </Button>
            </div>

            <Field label="Пояснение" hint="Показывается студенту после ответа">
              <Input
                value={question.explanation ?? ""}
                onChange={(e) => patchQuestion(questionIndex, { explanation: e.target.value })}
              />
            </Field>
          </div>
        </Section>
      ))}

      <Button variant="secondary" icon={<IconPlus size={16} />} onClick={addQuestion}>
        Добавить вопрос
      </Button>
    </div>
  );
}

// --- Тренажёр терминала ---

type EditorTask = {
  id: string;
  prompt: string;
  expected?: string[] | string;
  pattern?: string;
  hint?: string;
  success?: string;
};

function TerminalEditor({
  value,
  onChange,
}: {
  value: AnyRecord;
  onChange: (next: AnyRecord) => void;
}) {
  const tasks = (value.tasks as EditorTask[] | undefined) ?? [];

  const update = (next: EditorTask[]) => onChange({ ...value, tasks: next });
  const patchTask = (index: number, patch: Partial<EditorTask>) =>
    update(tasks.map((task, i) => (i === index ? { ...task, ...patch } : task)));

  const expectedText = (task: EditorTask) =>
    Array.isArray(task.expected) ? task.expected.join("\n") : (task.expected ?? "");

  return (
    <div className="space-y-4">
      <Field label="Вступление" hint="Что студент увидит перед терминалом">
        <Textarea
          value={String(value.intro ?? "")}
          onChange={(e) => onChange({ ...value, intro: e.target.value })}
          rows={2}
        />
      </Field>

      {tasks.map((task, index) => (
        <Section
          key={task.id}
          title={`Задание ${index + 1}`}
          action={
            <Button
              variant="ghost"
              className="h-8 !px-2 text-danger"
              onClick={() => update(tasks.filter((_, i) => i !== index))}
              title="Удалить задание"
            >
              <IconTrash size={16} />
            </Button>
          }
        >
          <div className="space-y-3">
            <Field label="Что нужно сделать">
              <Input
                value={task.prompt}
                onChange={(e) => patchTask(index, { prompt: e.target.value })}
                placeholder="Выведите список запущенных контейнеров"
              />
            </Field>

            <Field
              label="Допустимые команды"
              hint="По одной в строке — засчитывается любая из них"
            >
              <Textarea
                value={expectedText(task)}
                onChange={(e) =>
                  patchTask(index, {
                    expected: e.target.value.split("\n").map((line) => line.trim()).filter(Boolean),
                  })
                }
                rows={3}
                className="font-mono text-xs"
                spellCheck={false}
              />
            </Field>

            <div className="grid gap-3 sm:grid-cols-2">
              <Field label="Подсказка">
                <Input
                  value={task.hint ?? ""}
                  onChange={(e) => patchTask(index, { hint: e.target.value })}
                />
              </Field>
              <Field label="Текст при успехе">
                <Input
                  value={task.success ?? ""}
                  onChange={(e) => patchTask(index, { success: e.target.value })}
                />
              </Field>
            </div>

            <Field
              label="Или регулярное выражение"
              hint="Необязательно: подойдёт, когда вариантов команды много"
            >
              <Input
                value={task.pattern ?? ""}
                onChange={(e) => patchTask(index, { pattern: e.target.value })}
                className="font-mono text-xs"
                placeholder="^docker run .*nginx$"
              />
            </Field>
          </div>
        </Section>
      ))}

      <Button
        variant="secondary"
        icon={<IconPlus size={16} />}
        onClick={() => update([...tasks, { id: uid("t"), prompt: "", expected: [] }])}
      >
        Добавить задание
      </Button>
    </div>
  );
}

// --- Практика с кодом ---

type EditorCheck = { type: string; value: string; message?: string };

function CodeContentEditor({
  value,
  onChange,
}: {
  value: AnyRecord;
  onChange: (next: AnyRecord) => void;
}) {
  const checks = (value.checks as EditorCheck[] | undefined) ?? [];

  const patchCheck = (index: number, patch: Partial<EditorCheck>) =>
    onChange({
      ...value,
      checks: checks.map((check, i) => (i === index ? { ...check, ...patch } : check)),
    });

  return (
    <div className="space-y-4">
      <div className="grid gap-4 sm:grid-cols-2">
        <Field label="Язык" hint="Подпись над редактором">
          <Select
            value={String(value.language ?? "yaml")}
            onChange={(e) => onChange({ ...value, language: e.target.value })}
          >
            <option value="yaml">YAML</option>
            <option value="bash">Bash</option>
            <option value="dockerfile">Dockerfile</option>
            <option value="json">JSON</option>
            <option value="hcl">HCL / Terraform</option>
          </Select>
        </Field>

        <Field label="Подсказка">
          <Input
            value={String(value.hint ?? "")}
            onChange={(e) => onChange({ ...value, hint: e.target.value })}
          />
        </Field>
      </div>

      <Field label="Текст задания (Markdown)">
        <Textarea
          value={String(value.task ?? "")}
          onChange={(e) => onChange({ ...value, task: e.target.value })}
          rows={5}
        />
      </Field>

      <Field label="Заготовка кода" hint="С чего студент начинает">
        <Textarea
          value={String(value.starter ?? "")}
          onChange={(e) => onChange({ ...value, starter: e.target.value })}
          rows={8}
          className="font-mono text-xs"
          spellCheck={false}
        />
      </Field>

      <Section
        title="Проверки решения"
        action={
          <Button
            variant="ghost"
            icon={<IconPlus size={14} />}
            onClick={() =>
              onChange({ ...value, checks: [...checks, { type: "contains", value: "", message: "" }] })
            }
          >
            Проверка
          </Button>
        }
      >
        {checks.length === 0 ? (
          <p className="text-xs text-faint">
            Без проверок урок засчитается за любое непустое решение.
          </p>
        ) : (
          <ul className="space-y-2">
            {checks.map((check, index) => (
              <li key={index} className="grid gap-2 sm:grid-cols-[9rem_1fr_1fr_auto]">
                <Select
                  value={check.type}
                  onChange={(e) => patchCheck(index, { type: e.target.value })}
                >
                  <option value="contains">содержит</option>
                  <option value="notContains">не содержит</option>
                  <option value="regex">регулярное выражение</option>
                </Select>

                <Input
                  value={check.value}
                  onChange={(e) => patchCheck(index, { value: e.target.value })}
                  placeholder="8080:8080"
                  className="font-mono text-xs"
                />

                <Input
                  value={check.message ?? ""}
                  onChange={(e) => patchCheck(index, { message: e.target.value })}
                  placeholder="Текст для студента"
                />

                <Button
                  variant="ghost"
                  className="h-[var(--row-h)] !px-2 text-danger"
                  onClick={() =>
                    onChange({ ...value, checks: checks.filter((_, i) => i !== index) })
                  }
                  title="Удалить проверку"
                >
                  <IconTrash size={16} />
                </Button>
              </li>
            ))}
          </ul>
        )}
      </Section>
    </div>
  );
}

// --- Ссылки на материалы ---

type EditorResource = { title: string; url: string; note?: string };

function ResourcesEditor({
  value,
  onChange,
}: {
  value: AnyRecord;
  onChange: (next: AnyRecord) => void;
}) {
  const resources = (value.resources as EditorResource[] | undefined) ?? [];

  const patch = (index: number, next: Partial<EditorResource>) =>
    onChange({
      ...value,
      resources: resources.map((item, i) => (i === index ? { ...item, ...next } : item)),
    });

  return (
    <Section
      title="Материалы по теме"
      action={
        <Button
          variant="ghost"
          icon={<IconPlus size={14} />}
          onClick={() => onChange({ ...value, resources: [...resources, { title: "", url: "" }] })}
        >
          Ссылка
        </Button>
      }
    >
      {resources.length === 0 ? (
        <p className="text-xs text-faint">
          Добавьте официальную документацию или спецификацию — блок появится под уроком.
        </p>
      ) : (
        <ul className="space-y-2">
          {resources.map((item, index) => (
            <li key={index} className="grid gap-2 sm:grid-cols-[1fr_1fr_auto]">
              <Input
                value={item.title}
                onChange={(e) => patch(index, { title: e.target.value })}
                placeholder="Название материала"
              />
              <Input
                value={item.url}
                onChange={(e) => patch(index, { url: e.target.value })}
                placeholder="https://…"
                className="font-mono text-xs"
              />
              <Button
                variant="ghost"
                className="h-[var(--row-h)] !px-2 text-danger"
                onClick={() =>
                  onChange({ ...value, resources: resources.filter((_, i) => i !== index) })
                }
                title="Удалить ссылку"
              >
                <IconTrash size={16} />
              </Button>

              <Input
                value={item.note ?? ""}
                onChange={(e) => patch(index, { note: e.target.value })}
                placeholder="Чем полезно — короткая заметка"
                className="sm:col-span-3"
              />
            </li>
          ))}
        </ul>
      )}
    </Section>
  );
}

// --- Переключатель «визуально / JSON» ---

export default function LessonContentEditor({
  kind,
  value,
  onChange,
  raw,
  onRawChange,
  rawError,
}: {
  kind: LessonKind;
  value: AnyRecord;
  onChange: (next: AnyRecord) => void;
  raw: string;
  onRawChange: (next: string) => void;
  rawError?: string;
}) {
  const [mode, setMode] = useState<"visual" | "json">("visual");

  return (
    <div>
      <div className="mb-3 flex items-center justify-between gap-2">
        <p className="text-sm font-bold text-fg">Содержимое урока</p>
        <div className="flex gap-1">
          <button
            type="button"
            onClick={() => setMode("visual")}
            className={`rounded-[var(--radius-sm)] px-2.5 py-1 text-xs font-semibold transition-colors ${
              mode === "visual" ? "bg-accent-soft text-accent" : "text-muted hover:bg-surface-2"
            }`}
          >
            Конструктор
          </button>
          <button
            type="button"
            onClick={() => setMode("json")}
            className={`rounded-[var(--radius-sm)] px-2.5 py-1 text-xs font-semibold transition-colors ${
              mode === "json" ? "bg-accent-soft text-accent" : "text-muted hover:bg-surface-2"
            }`}
          >
            JSON
          </button>
        </div>
      </div>

      {mode === "json" ? (
        <>
          <Textarea
            value={raw}
            onChange={(e) => onRawChange(e.target.value)}
            rows={16}
            className="font-mono text-xs"
            spellCheck={false}
          />
          {rawError && <p className="mt-1 text-xs font-medium text-danger">{rawError}</p>}
        </>
      ) : kind === "quiz" ? (
        <QuizEditor value={value} onChange={onChange} />
      ) : kind === "terminal" ? (
        <TerminalEditor value={value} onChange={onChange} />
      ) : kind === "code" ? (
        <CodeContentEditor value={value} onChange={onChange} />
      ) : (
        <TextEditor value={value} onChange={onChange} />
      )}

      {mode === "visual" && (
        <div className="mt-4">
          <ResourcesEditor value={value} onChange={onChange} />
        </div>
      )}

      {mode === "visual" && kind === "quiz" && (
        <p className="mt-2 flex items-center gap-2 text-xs text-faint">
          <Badge tone="accent">важно</Badge>
          Правильные ответы не уходят студенту — проверка выполняется на сервере.
        </p>
      )}
    </div>
  );
}
