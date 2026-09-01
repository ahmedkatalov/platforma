import { useEffect, useMemo, useRef, useState } from "react";
import clsx from "clsx";
import { Calendar, ChevronLeft, ChevronRight, X } from "lucide-react";

const WEEKDAYS = ["Пн", "Вт", "Ср", "Чт", "Пт", "Сб", "Вс"];
const monthFmt = new Intl.DateTimeFormat("ru-RU", { month: "long", year: "numeric" });
const valueFmt = new Intl.DateTimeFormat("ru-RU", { day: "2-digit", month: "2-digit", year: "numeric" });

const pad = (n: number) => String(n).padStart(2, "0");
const toValue = (d: Date) => `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}`;

function parseValue(v?: string): Date | null {
  if (!v) return null;
  const [y, m, d] = v.split("-").map(Number);
  if (!y || !m || !d) return null;
  return new Date(y, m - 1, d);
}

const sameDay = (a: Date, b: Date) =>
  a.getFullYear() === b.getFullYear() && a.getMonth() === b.getMonth() && a.getDate() === b.getDate();

// Кастомный выбор даты в стиле платформы (нативный календарь браузера не тянет тёмную тему).
export function DatePicker({
  value,
  onChange,
  placeholder = "дд.мм.гггг",
  className,
}: {
  value?: string;
  onChange: (value: string) => void;
  placeholder?: string;
  className?: string;
}) {
  const selected = parseValue(value);
  const [open, setOpen] = useState(false);
  const [view, setView] = useState<Date>(() => selected ?? new Date());
  const rootRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!open) return;
    const onDoc = (e: MouseEvent) => {
      if (rootRef.current && !rootRef.current.contains(e.target as Node)) setOpen(false);
    };
    const onKey = (e: KeyboardEvent) => e.key === "Escape" && setOpen(false);
    document.addEventListener("mousedown", onDoc);
    document.addEventListener("keydown", onKey);
    return () => {
      document.removeEventListener("mousedown", onDoc);
      document.removeEventListener("keydown", onKey);
    };
  }, [open]);

  const days = useMemo(() => {
    const first = new Date(view.getFullYear(), view.getMonth(), 1);
    const start = new Date(first);
    start.setDate(first.getDate() - ((first.getDay() + 6) % 7)); // неделя с понедельника
    return Array.from({ length: 42 }, (_, i) => {
      const d = new Date(start);
      d.setDate(start.getDate() + i);
      return d;
    });
  }, [view]);

  const today = new Date();
  const label = selected ? valueFmt.format(selected) : "";

  const toggle = () => {
    if (!open) setView(selected ?? new Date());
    setOpen((v) => !v);
  };
  const pick = (d: Date) => {
    onChange(toValue(d));
    setOpen(false);
  };
  const shiftMonth = (delta: number) =>
    setView((v) => new Date(v.getFullYear(), v.getMonth() + delta, 1));

  return (
    <div ref={rootRef} className={clsx("relative", className)}>
      <button
        type="button"
        onClick={toggle}
        className="input flex items-center gap-2 text-left"
        aria-haspopup="dialog"
        aria-expanded={open}
      >
        <Calendar size={16} className="shrink-0 text-faint" />
        <span className={clsx("flex-1 truncate", label ? "text-fg" : "text-faint")}>
          {label || placeholder}
        </span>
        {label && (
          <span
            role="button"
            tabIndex={-1}
            title="Очистить"
            className="grid h-5 w-5 shrink-0 place-items-center rounded-full text-faint hover:bg-surface-hover hover:text-danger"
            onClick={(e) => {
              e.stopPropagation();
              onChange("");
              setOpen(false);
            }}
          >
            <X size={13} />
          </span>
        )}
      </button>

      {open && (
        <div
          className="absolute left-0 top-full z-40 mt-2 w-[17.5rem] rounded-[var(--radius-lg)] border border-line bg-surface-solid p-3 shadow-2xl"
          role="dialog"
        >
          <div className="mb-2 flex items-center justify-between">
            <button
              type="button"
              className="grid h-8 w-8 place-items-center rounded-[var(--radius-md)] text-muted transition-colors hover:bg-surface-2 hover:text-fg"
              onClick={() => shiftMonth(-1)}
              aria-label="Предыдущий месяц"
            >
              <ChevronLeft size={16} />
            </button>
            <span className="text-sm font-bold capitalize text-fg">{monthFmt.format(view)}</span>
            <button
              type="button"
              className="grid h-8 w-8 place-items-center rounded-[var(--radius-md)] text-muted transition-colors hover:bg-surface-2 hover:text-fg"
              onClick={() => shiftMonth(1)}
              aria-label="Следующий месяц"
            >
              <ChevronRight size={16} />
            </button>
          </div>

          <div className="mb-1 grid grid-cols-7 text-center text-[11px] font-semibold text-faint">
            {WEEKDAYS.map((w) => (
              <span key={w} className="py-1">
                {w}
              </span>
            ))}
          </div>

          <div className="grid grid-cols-7 gap-0.5">
            {days.map((d, i) => {
              const inMonth = d.getMonth() === view.getMonth();
              const isSel = selected && sameDay(d, selected);
              const isToday = sameDay(d, today);
              return (
                <button
                  key={i}
                  type="button"
                  onClick={() => pick(d)}
                  className={clsx(
                    "grid h-8 place-items-center rounded-[var(--radius-md)] text-sm transition-colors",
                    isSel
                      ? "bg-accent font-bold text-accent-fg"
                      : isToday
                        ? "border border-accent-border text-accent hover:bg-surface-2"
                        : inMonth
                          ? "text-fg hover:bg-surface-2"
                          : "text-faint hover:bg-surface-2",
                  )}
                >
                  {d.getDate()}
                </button>
              );
            })}
          </div>

          <div className="mt-2 flex items-center justify-between border-t border-line pt-2 text-xs">
            <button
              type="button"
              className="font-semibold text-muted transition-colors hover:text-danger"
              onClick={() => {
                onChange("");
                setOpen(false);
              }}
            >
              Очистить
            </button>
            <button
              type="button"
              className="font-semibold text-accent hover:underline"
              onClick={() => pick(today)}
            >
              Сегодня
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
