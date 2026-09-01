-- Пошаговый доступ к главам: студент видит только открытые ему главы (модули).
CREATE TABLE IF NOT EXISTS module_access (
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    module_id  UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    granted_by UUID REFERENCES users(id) ON DELETE SET NULL,
    granted_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, module_id)
);
CREATE INDEX IF NOT EXISTS module_access_module_idx ON module_access (module_id);

-- Заявки студентов на открытие следующей главы.
CREATE TABLE IF NOT EXISTS access_requests (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    module_id  UUID NOT NULL REFERENCES modules(id) ON DELETE CASCADE,
    status     TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected')),
    note       TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    decided_at TIMESTAMPTZ,
    decided_by UUID REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX IF NOT EXISTS access_requests_status_idx ON access_requests (status, created_at);

-- Не более одной висящей заявки на пару (студент, глава).
CREATE UNIQUE INDEX IF NOT EXISTS access_requests_pending_uq
    ON access_requests (user_id, module_id) WHERE status = 'pending';
