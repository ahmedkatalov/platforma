-- Ключ Gemini и модель теперь можно задать из админки, а не только через
-- переменные окружения. Хранятся на сервере (в этой строке-синглтоне) и
-- НИКОГДА не отдаются на фронтенд — студенты их не видят. Если пусто —
-- используется значение из окружения (GEMINI_API_KEY / GEMINI_MODEL).
ALTER TABLE platform_ai ADD COLUMN IF NOT EXISTS api_key TEXT;
ALTER TABLE platform_ai ADD COLUMN IF NOT EXISTS model   TEXT;
