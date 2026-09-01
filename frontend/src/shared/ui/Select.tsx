import {
  Children,
  isValidElement,
  useEffect,
  useLayoutEffect,
  useMemo,
  useRef,
  useState,
  type ReactNode,
} from "react";
import clsx from "clsx";
import { Check, ChevronDown } from "lucide-react";

export type SelectOption = { value: string; label: ReactNode; disabled?: boolean };

type SelectProps = {
  value?: string;
  onChange?: (value: string) => void;
  options?: SelectOption[];
  children?: ReactNode; // <option> элементы — как у нативного select
  placeholder?: string;
  className?: string;
  disabled?: boolean;
  id?: string;
  "aria-label"?: string;
};

// Разбираем <option>…</option> в список опций — чтобы не переписывать все места вызова.
function parseOptions(children: ReactNode): SelectOption[] {
  const out: SelectOption[] = [];
  Children.forEach(children, (child) => {
    if (isValidElement(child) && child.type === "option") {
      const p = child.props as { value?: string; children?: ReactNode; disabled?: boolean };
      out.push({ value: String(p.value ?? ""), label: p.children ?? "", disabled: p.disabled });
    }
  });
  return out;
}

// Кастомный селект в стиле платформы: нативный список браузера не поддаётся теме.
export function Select({
  value = "",
  onChange,
  options,
  children,
  placeholder = "Выберите…",
  className,
  disabled,
  id,
  "aria-label": ariaLabel,
}: SelectProps) {
  const opts = useMemo(() => options ?? parseOptions(children), [options, children]);
  const [open, setOpen] = useState(false);
  const [up, setUp] = useState(false);
  const [active, setActive] = useState(0);
  const rootRef = useRef<HTMLDivElement>(null);
  const listRef = useRef<HTMLDivElement>(null);

  const selected = opts.find((o) => o.value === value);
  const selectedIdx = Math.max(0, opts.findIndex((o) => o.value === value));

  const openMenu = () => {
    if (disabled) return;
    setActive(selectedIdx);
    // Открываться вверх, если снизу мало места.
    const r = rootRef.current?.getBoundingClientRect();
    if (r) setUp(window.innerHeight - r.bottom < 260 && r.top > window.innerHeight - r.bottom);
    setOpen(true);
  };

  const pick = (o: SelectOption) => {
    if (o.disabled) return;
    onChange?.(o.value);
    setOpen(false);
  };

  useEffect(() => {
    if (!open) return;
    const onDoc = (e: MouseEvent) => {
      if (rootRef.current && !rootRef.current.contains(e.target as Node)) setOpen(false);
    };
    document.addEventListener("mousedown", onDoc);
    return () => document.removeEventListener("mousedown", onDoc);
  }, [open]);

  // Прокручиваем к активному пункту.
  useLayoutEffect(() => {
    if (!open) return;
    const el = listRef.current?.querySelector<HTMLElement>(`[data-idx="${active}"]`);
    el?.scrollIntoView({ block: "nearest" });
  }, [open, active]);

  const onKey = (e: React.KeyboardEvent) => {
    if (disabled) return;
    if (!open) {
      if (e.key === "Enter" || e.key === " " || e.key === "ArrowDown") {
        e.preventDefault();
        openMenu();
      }
      return;
    }
    if (e.key === "Escape") {
      setOpen(false);
    } else if (e.key === "ArrowDown") {
      e.preventDefault();
      setActive((i) => Math.min(opts.length - 1, i + 1));
    } else if (e.key === "ArrowUp") {
      e.preventDefault();
      setActive((i) => Math.max(0, i - 1));
    } else if (e.key === "Enter") {
      e.preventDefault();
      if (opts[active]) pick(opts[active]);
    } else if (e.key === "Tab") {
      setOpen(false);
    }
  };

  return (
    <div ref={rootRef} className={clsx("relative", className)}>
      <button
        type="button"
        id={id}
        aria-label={ariaLabel}
        aria-haspopup="listbox"
        aria-expanded={open}
        disabled={disabled}
        onClick={() => (open ? setOpen(false) : openMenu())}
        onKeyDown={onKey}
        className="input flex cursor-pointer items-center gap-2 text-left"
      >
        <span className={clsx("flex-1 truncate", selected ? "text-fg" : "text-faint")}>
          {selected ? selected.label : placeholder}
        </span>
        <ChevronDown
          size={16}
          className={clsx("shrink-0 text-faint transition-transform duration-150", open && "rotate-180")}
        />
      </button>

      {open && (
        <div
          ref={listRef}
          role="listbox"
          className={clsx(
            "menu-surface anim-pop absolute z-40 max-h-64 w-full overflow-auto p-1",
            up ? "bottom-full mb-1.5" : "top-full mt-1.5",
          )}
        >
          {opts.length === 0 ? (
            <p className="px-2.5 py-2 text-sm text-faint">Нет вариантов</p>
          ) : (
            opts.map((o, i) => {
              const isSel = o.value === value;
              return (
                <div
                  key={o.value + i}
                  data-idx={i}
                  data-active={i === active}
                  role="option"
                  aria-selected={isSel}
                  onMouseEnter={() => setActive(i)}
                  onClick={() => pick(o)}
                  className={clsx("menu-item justify-between", o.disabled && "cursor-not-allowed opacity-40")}
                >
                  <span className="truncate">{o.label}</span>
                  {isSel && <Check size={15} className="shrink-0 text-accent" />}
                </div>
              );
            })
          )}
        </div>
      )}
    </div>
  );
}
