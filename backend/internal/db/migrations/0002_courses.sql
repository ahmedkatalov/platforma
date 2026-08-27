-- Курсы, модули и уроки.
CREATE TABLE IF NOT EXISTS courses (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slug          TEXT NOT NULL UNIQUE,
    title         TEXT NOT NULL,
    subtitle      TEXT NOT NULL DEFAULT '',
    description   TEXT NOT NULL DEFAULT '',
    cover_url     TEXT NOT NULL DEFAULT '',
    level         TEXT NOT NULL DEFAULT 'beginner' CHECK (level IN ('beginner', 'intermediate', 'advanced')),
    tags          TEXT[] NOT NULL DEFAULT '{}',
    status        TEXT NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'archived')),
    position      INT NOT NULL DEFAULT 0,
    created_by    UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS modules (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id   UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    title       TEXT NOT NULL,
    summary     TEXT NOT NULL DEFAULT '',
    position    INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS modules_course_idx ON modules (course_id, position);

-- Урок: теория, квиз, эмуляция терминала или редактор кода.
CREATE TABLE IF NOT EXISTS lessons (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    module_id      UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    title          TEXT NOT NULL,
    kind           TEXT NOT NULL DEFAULT 'text' CHECK (kind IN ('text', 'quiz', 'terminal', 'code')),
    summary        TEXT NOT NULL DEFAULT '',
    content        JSONB NOT NULL DEFAULT '{}'::jsonb,
    duration_min   INT NOT NULL DEFAULT 10,
    position       INT NOT NULL DEFAULT 0,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS lessons_module_idx ON lessons (module_id, position);

-- Запись студента на курс.
CREATE TABLE IF NOT EXISTS enrollments (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    course_id    UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    status       TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'completed', 'paused')),
    assigned_by  UUID REFERENCES users(id) ON DELETE SET NULL,
    started_at   TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, course_id)
);

CREATE INDEX IF NOT EXISTS enrollments_course_idx ON enrollments (course_id);

-- Прогресс по уроку.
CREATE TABLE IF NOT EXISTS lesson_progress (
    user_id       UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    lesson_id     UUID NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
    status        TEXT NOT NULL DEFAULT 'in_progress' CHECK (status IN ('in_progress', 'completed')),
    score         NUMERIC(5,2),
    seconds_spent INT NOT NULL DEFAULT 0,
    completed_at  TIMESTAMPTZ,
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, lesson_id)
);
