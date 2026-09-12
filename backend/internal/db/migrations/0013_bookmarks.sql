-- Закладки студента: сохранённые уроки (темы) и модули (главы), чтобы вернуться
-- к ним позже и перечитать. ref_id указывает на lessons.id или modules.id
-- в зависимости от kind; чистим осиротевшие закладки при удалении цели триггером
-- не нужно — при выборке join просто их не вернёт, а на удаление курса каскада нет,
-- поэтому подчищаем по ссылке в приложении при показе.
CREATE TABLE IF NOT EXISTS bookmarks (
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    kind       TEXT NOT NULL CHECK (kind IN ('lesson', 'module')),
    ref_id     UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, kind, ref_id)
);

CREATE INDEX IF NOT EXISTS bookmarks_user_idx ON bookmarks (user_id, created_at DESC);
