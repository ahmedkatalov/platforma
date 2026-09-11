-- Прокси для исходящих запросов к Gemini. Нужен, если сервер в регионе,
-- который Google блокирует (ошибка «User location is not supported»).
-- Может содержать логин/пароль в URL, поэтому наружу (на фронтенд) не отдаётся.
ALTER TABLE platform_ai ADD COLUMN IF NOT EXISTS proxy TEXT;
