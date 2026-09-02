import { Component, type ErrorInfo, type ReactNode } from "react";

/*
 * Верхнеуровневый перехватчик ошибок. Без него любая ошибка рендера роняет всё
 * приложение в белый экран (частая причина «не грузится на телефоне»). Здесь мы
 * показываем понятный экран и даём кнопку «жёсткой» перезагрузки: снимаем service
 * worker и чистим кэши — это лечит ситуацию, когда у пользователя застряла старая
 * сломанная версия.
 */
type Props = { children: ReactNode };
type State = { error: Error | null };

async function hardReload() {
  try {
    if ("serviceWorker" in navigator) {
      const regs = await navigator.serviceWorker.getRegistrations();
      await Promise.all(regs.map((r) => r.unregister()));
    }
    if ("caches" in window) {
      const keys = await caches.keys();
      await Promise.all(keys.map((k) => caches.delete(k)));
    }
  } catch {
    // неважно — всё равно перезагружаемся
  }
  location.reload();
}

export default class ErrorBoundary extends Component<Props, State> {
  state: State = { error: null };

  static getDerivedStateFromError(error: Error): State {
    return { error };
  }

  componentDidCatch(error: Error, info: ErrorInfo) {
    // Оставляем след в консоли для диагностики на реальном устройстве.
    console.error("Сбой интерфейса:", error, info.componentStack);
  }

  render() {
    if (!this.state.error) return this.props.children;
    return (
      <div
        style={{
          minHeight: "100vh",
          display: "grid",
          placeItems: "center",
          padding: "24px",
          background: "#0b1020",
          color: "#e5e9f0",
          fontFamily: "system-ui, -apple-system, Segoe UI, Roboto, sans-serif",
        }}
      >
        <div style={{ maxWidth: 420, textAlign: "center" }}>
          <div style={{ fontSize: 40, marginBottom: 12 }}>⚠️</div>
          <h1 style={{ fontSize: 20, fontWeight: 700, margin: "0 0 8px" }}>Что-то пошло не так</h1>
          <p style={{ fontSize: 14, lineHeight: 1.6, color: "#9aa4bf", margin: "0 0 20px" }}>
            Приложение не смогло отрисовать страницу. Обычно помогает перезагрузка.
            Если не помогло — нажмите «Сбросить и перезагрузить»: это очистит устаревший кэш.
          </p>
          <div style={{ display: "flex", flexWrap: "wrap", gap: 8, justifyContent: "center" }}>
            <button
              onClick={() => location.reload()}
              style={{
                padding: "10px 18px",
                borderRadius: 10,
                border: "none",
                fontWeight: 600,
                fontSize: 14,
                cursor: "pointer",
                background: "#3b82f6",
                color: "#fff",
              }}
            >
              Перезагрузить
            </button>
            <button
              onClick={hardReload}
              style={{
                padding: "10px 18px",
                borderRadius: 10,
                border: "1px solid #2a3550",
                fontWeight: 600,
                fontSize: 14,
                cursor: "pointer",
                background: "transparent",
                color: "#e5e9f0",
              }}
            >
              Сбросить и перезагрузить
            </button>
          </div>
        </div>
      </div>
    );
  }
}
