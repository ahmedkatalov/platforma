-- Сертификат выдаётся один раз на пару студент+курс при прохождении всех уроков.
CREATE TABLE IF NOT EXISTS certificates (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    serial            TEXT NOT NULL UNIQUE,
    user_id           UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    course_id         UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    holder_name       TEXT NOT NULL,
    course_title      TEXT NOT NULL,
    score             NUMERIC(5,2) NOT NULL DEFAULT 0,
    lessons_total     INT NOT NULL DEFAULT 0,
    lessons_completed INT NOT NULL DEFAULT 0,
    revoked_at        TIMESTAMPTZ,
    issued_at         TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, course_id)
);

CREATE INDEX IF NOT EXISTS certificates_user_idx ON certificates (user_id, issued_at DESC);

-- Сроки прохождения курса и служебные отметки о напоминаниях.
ALTER TABLE enrollments ADD COLUMN IF NOT EXISTS due_date DATE;
ALTER TABLE enrollments ADD COLUMN IF NOT EXISTS deadline_notified_at TIMESTAMPTZ;
ALTER TABLE enrollments ADD COLUMN IF NOT EXISTS idle_notified_at TIMESTAMPTZ;

-- Файлы, загруженные администратором для уроков (картинки, схемы).
CREATE TABLE IF NOT EXISTS assets (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    filename     TEXT NOT NULL UNIQUE,
    original     TEXT NOT NULL DEFAULT '',
    mime         TEXT NOT NULL,
    size_bytes   BIGINT NOT NULL DEFAULT 0,
    uploaded_by  UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS assets_created_idx ON assets (created_at DESC);
