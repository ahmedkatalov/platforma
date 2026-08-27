-- Оформление платформы: одна строка (id = TRUE) — общие настройки от админа.
CREATE TABLE IF NOT EXISTS platform_theme (
    id          BOOLEAN PRIMARY KEY DEFAULT TRUE CHECK (id),
    settings    JSONB NOT NULL DEFAULT '{}'::jsonb,
    updated_by  UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Персональные настройки интерфейса пользователя (перекрывают общие).
CREATE TABLE IF NOT EXISTS user_preferences (
    user_id     UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    theme       JSONB NOT NULL DEFAULT '{}'::jsonb,
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
