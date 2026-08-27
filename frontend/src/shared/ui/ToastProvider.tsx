import {
  createContext,
  useCallback,
  useContext,
  useMemo,
  useState,
  type ReactNode,
} from "react";

import { IconCheck, IconClose } from "./icons";

type ToastKind = "success" | "error" | "info";
type Toast = { id: number; kind: ToastKind; text: string };

type ToastContextValue = {
  success: (text: string) => void;
  error: (text: string) => void;
  info: (text: string) => void;
};

const ToastContext = createContext<ToastContextValue | null>(null);

const TONE: Record<ToastKind, string> = {
  success: "border-[var(--success)] text-success",
  error: "border-[var(--danger)] text-danger",
  info: "border-accent-border text-accent",
};

export function ToastProvider({ children }: { children: ReactNode }) {
  const [toasts, setToasts] = useState<Toast[]>([]);

  const push = useCallback((kind: ToastKind, text: string) => {
    const id = Date.now() + Math.random();
    setToasts((current) => [...current, { id, kind, text }]);
    window.setTimeout(() => {
      setToasts((current) => current.filter((t) => t.id !== id));
    }, 5000);
  }, []);

  const value = useMemo<ToastContextValue>(
    () => ({
      success: (text) => push("success", text),
      error: (text) => push("error", text),
      info: (text) => push("info", text),
    }),
    [push],
  );

  return (
    <ToastContext.Provider value={value}>
      {children}
      <div className="pointer-events-none fixed bottom-4 right-4 z-[100] flex w-[min(24rem,calc(100vw-2rem))] flex-col gap-2">
        {toasts.map((toast) => (
          <div
            key={toast.id}
            className={`card pointer-events-auto flex items-start gap-3 border-l-4 p-3 text-sm ${TONE[toast.kind]}`}
            role="status"
          >
            {toast.kind === "success" ? <IconCheck size={18} /> : <IconClose size={18} />}
            <p className="flex-1 text-fg">{toast.text}</p>
            <button
              className="text-faint transition-colors hover:text-fg"
              onClick={() => setToasts((c) => c.filter((t) => t.id !== toast.id))}
              aria-label="Закрыть уведомление"
            >
              <IconClose size={16} />
            </button>
          </div>
        ))}
      </div>
    </ToastContext.Provider>
  );
}

export function useToast(): ToastContextValue {
  const ctx = useContext(ToastContext);
  if (!ctx) throw new Error("useToast должен вызываться внутри ToastProvider");
  return ctx;
}
