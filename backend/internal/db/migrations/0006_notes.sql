-- Личные заметки студента: выделенная цитата из урока плюс свой комментарий.
CREATE TABLE IF NOT EXISTS notes (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    lesson_id  UUID NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
    quote      TEXT NOT NULL,
    body       TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS notes_user_idx ON notes (user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS notes_lesson_idx ON notes (user_id, lesson_id);
