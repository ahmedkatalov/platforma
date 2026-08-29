import { useEffect, useState } from "react";

import { IconClose, IconTerminal } from "@/shared/ui/icons";

// Событие beforeinstallprompt есть только в спецификации Chromium.
type InstallEvent = Event & {
  prompt: () => Promise<void>;
  userChoice: Promise<{ outcome: "accepted" | "dismissed" }>;
};

const DISMISS_KEY = "platforma.installDismissed";

// Плашка «Установить приложение». На Android/десктопе Chrome показывает
// системную установку, здесь мы даём удобную кнопку. На iPhone установки
// через событие нет — там подсказываем добавить на экран «Домой» вручную.
export default function InstallPrompt() {
  const [deferred, setDeferred] = useState<InstallEvent | null>(null);
  const [showIosHint, setShowIosHint] = useState(false);

  useEffect(() => {
    const standalone =
      window.matchMedia("(display-mode: standalone)").matches ||
      (window.navigator as { standalone?: boolean }).standalone === true;
    if (standalone || localStorage.getItem(DISMISS_KEY)) return;

    const onPrompt = (event: Event) => {
      event.preventDefault();
      setDeferred(event as InstallEvent);
    };
    window.addEventListener("beforeinstallprompt", onPrompt);

    const ua = window.navigator.userAgent;
    const isIos = /iphone|ipad|ipod/i.test(ua);
    const isSafari = /safari/i.test(ua) && !/crios|fxios|chrome/i.test(ua);
    if (isIos && isSafari) setShowIosHint(true);

    const onInstalled = () => {
      setDeferred(null);
      setShowIosHint(false);
    };
    window.addEventListener("appinstalled", onInstalled);

    return () => {
      window.removeEventListener("beforeinstallprompt", onPrompt);
      window.removeEventListener("appinstalled", onInstalled);
    };
  }, []);

  const dismiss = () => {
    localStorage.setItem(DISMISS_KEY, "1");
    setDeferred(null);
    setShowIosHint(false);
  };

  const install = async () => {
    if (!deferred) return;
    await deferred.prompt();
    await deferred.userChoice;
    setDeferred(null);
  };

  if (!deferred && !showIosHint) return null;

  return (
    <div className="fixed inset-x-4 bottom-4 z-[90] mx-auto max-w-md sm:left-auto sm:right-4">
      <div className="card flex items-start gap-3 p-4">
        <span
          className="grid h-10 w-10 shrink-0 place-items-center rounded-[var(--radius-md)] text-accent-fg"
          style={{ background: "var(--gradient)" }}
        >
          <IconTerminal size={20} />
        </span>

        <div className="min-w-0 flex-1">
          <p className="text-sm font-bold text-fg">Установить приложение</p>
          {showIosHint ? (
            <p className="mt-1 text-xs text-muted">
              Нажмите «Поделиться», затем «На экран „Домой“» — платформа откроется как приложение.
            </p>
          ) : (
            <p className="mt-1 text-xs text-muted">
              Быстрый доступ с рабочего стола, полноэкранный режим и работа даже офлайн.
            </p>
          )}

          {!showIosHint && (
            <button className="btn btn-primary mt-3 h-9" onClick={install}>
              Установить
            </button>
          )}
        </div>

        <button
          className="btn btn-ghost h-8 w-8 shrink-0 !p-0"
          onClick={dismiss}
          aria-label="Скрыть"
        >
          <IconClose size={16} />
        </button>
      </div>
    </div>
  );
}
