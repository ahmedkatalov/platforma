-- Контакты для связи (Telegram/WhatsApp), настраиваются администратором.
-- Синглтон-строка по образцу platform_theme.
CREATE TABLE IF NOT EXISTS platform_contacts (
    id         BOOLEAN PRIMARY KEY DEFAULT TRUE CHECK (id),
    settings   JSONB NOT NULL DEFAULT '{}'::jsonb,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
