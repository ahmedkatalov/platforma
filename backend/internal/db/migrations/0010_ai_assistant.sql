-- Настройки ИИ-помощника (вкл/выкл администратором). Синглтон-строка по
-- образцу platform_contacts. Сам ключ Gemini хранится в переменной окружения,
-- в базе его нет — тут только флаг доступности фичи.
CREATE TABLE IF NOT EXISTS platform_ai (
    id         BOOLEAN PRIMARY KEY DEFAULT TRUE CHECK (id),
    enabled    BOOLEAN NOT NULL DEFAULT TRUE,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
