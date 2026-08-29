// Service worker платформы. Задача: сделать приложение устанавливаемым
// и дать базовый офлайн-режим, НЕ ломая обновления и API.
const CACHE = "platforma-v1";
const SHELL = ["/", "/index.html", "/manifest.webmanifest", "/pwa-192.png", "/apple-touch-icon.png"];

self.addEventListener("install", (event) => {
  event.waitUntil(
    caches.open(CACHE).then((c) => c.addAll(SHELL)).then(() => self.skipWaiting()),
  );
});

self.addEventListener("activate", (event) => {
  event.waitUntil(
    caches
      .keys()
      .then((keys) => Promise.all(keys.filter((k) => k !== CACHE).map((k) => caches.delete(k))))
      .then(() => self.clients.claim()),
  );
});

self.addEventListener("fetch", (event) => {
  const { request } = event;
  if (request.method !== "GET") return;

  const url = new URL(request.url);
  if (url.origin !== self.location.origin) return;

  // API никогда не кешируем — всегда свежие данные.
  if (url.pathname.startsWith("/api")) return;

  // Страницы: сеть в приоритете (свежий index.html после деплоя),
  // офлайн — отдаём сохранённую оболочку.
  if (request.mode === "navigate") {
    event.respondWith(fetch(request).catch(() => caches.match("/index.html")));
    return;
  }

  // Хешированные ассеты не меняются — берём из кэша, иначе качаем и сохраняем.
  if (url.pathname.startsWith("/assets/")) {
    event.respondWith(
      caches.match(request).then(
        (hit) =>
          hit ||
          fetch(request).then((resp) => {
            const copy = resp.clone();
            caches.open(CACHE).then((c) => c.put(request, copy));
            return resp;
          }),
      ),
    );
    return;
  }

  // Прочее: сеть, при офлайне — из кэша, если есть.
  event.respondWith(fetch(request).catch(() => caches.match(request)));
});
