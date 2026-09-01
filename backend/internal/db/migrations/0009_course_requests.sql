-- Заявки студентов на доступ к КУРСУ (запись на курс) — из общей витрины курсов.
CREATE TABLE IF NOT EXISTS course_requests (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    course_id  UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    status     TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    note       TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    decided_at TIMESTAMPTZ,
    decided_by UUID REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS course_requests_status_idx ON course_requests (status, created_at);

-- Не более одной висящей заявки на пару (студент, курс).
CREATE UNIQUE INDEX IF NOT EXISTS course_requests_pending_uq
    ON course_requests (user_id, course_id) WHERE status = 'pending';
