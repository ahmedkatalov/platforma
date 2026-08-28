-- Попытки прохождения уроков: квизы, терминал, редактор кода.
CREATE TABLE IF NOT EXISTS lesson_attempts (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id          UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    lesson_id        UUID NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
    kind             TEXT NOT NULL CHECK (kind IN ('quiz', 'terminal', 'code')),
    score            NUMERIC(5,2) NOT NULL DEFAULT 0,
    correct_count    INT NOT NULL DEFAULT 0,
    total_count      INT NOT NULL DEFAULT 0,
    passed           BOOLEAN NOT NULL DEFAULT FALSE,
    duration_seconds INT NOT NULL DEFAULT 0,
    details          JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS lesson_attempts_user_idx ON lesson_attempts (user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS lesson_attempts_lesson_idx ON lesson_attempts (lesson_id);

-- Прогресс по отдельным заданиям внутри урока (шаги терминала, проверки кода).
CREATE TABLE IF NOT EXISTS task_progress (
    user_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    lesson_id    UUID NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
    task_id      TEXT NOT NULL,
    attempts     INT NOT NULL DEFAULT 0,
    hints_used   INT NOT NULL DEFAULT 0,
    completed_at TIMESTAMPTZ,
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, lesson_id, task_id)
);

-- Ответы на отдельные вопросы квиза — нужны для статистики скорости и слабых тем.
CREATE TABLE IF NOT EXISTS quiz_answers (
    id               BIGSERIAL PRIMARY KEY,
    attempt_id       UUID NOT NULL REFERENCES lesson_attempts(id) ON DELETE CASCADE,
    user_id          UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    lesson_id        UUID NOT NULL REFERENCES lessons(id) ON DELETE CASCADE,
    question_id      TEXT NOT NULL,
    correct          BOOLEAN NOT NULL,
    seconds_spent    INT NOT NULL DEFAULT 0,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS quiz_answers_user_idx ON quiz_answers (user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS quiz_answers_question_idx ON quiz_answers (lesson_id, question_id);

-- Лучший результат урока держим прямо в прогрессе.
ALTER TABLE lesson_progress ADD COLUMN IF NOT EXISTS attempts INT NOT NULL DEFAULT 0;
ALTER TABLE lesson_progress ADD COLUMN IF NOT EXISTS best_score NUMERIC(5,2);
ALTER TABLE lesson_progress ADD COLUMN IF NOT EXISTS started_at TIMESTAMPTZ NOT NULL DEFAULT now();
