// Регистрация service worker — только в проде. В dev-режиме Vite отдаёт модули
// напрямую, и SW только мешал бы горячей перезагрузке.
export function registerServiceWorker() {
  if (!import.meta.env.PROD) return;
  if (!("serviceWorker" in navigator)) return;

  window.addEventListener("load", () => {
    navigator.serviceWorker.register("/sw.js").catch(() => {
      // Регистрация не критична: без неё сайт работает как обычно.
    });
  });
}
