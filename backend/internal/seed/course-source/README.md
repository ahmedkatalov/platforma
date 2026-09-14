# Источник курса (import-формат)

`devops-engineer.course.final.json` — курс в формате `platforma-course` v1
(тот же, что выдаёт `go run ./cmd/exportcourse`). Это канонический
импорт-артефакт.

Реально платформа грузит курс из `../content/*.json` (seed-формат, встроен
через go:embed). Файлы `content/*.json` сгенерированы из этого источника:
каждый модуль → отдельный файл `NNN.json` вида
`{title, summary, lessons:[{title,kind,summary,durationMin,content}]}`
(поля `id` опущены — их проставляет БД при сиде по title, прогресс
сохраняется). Метаданные курса (title/subtitle/description/tags) — в `../seed.go`.

Обновление курса: заменить этот файл и пересобрать `content/*.json` тем же
преобразованием, либо импортировать через админку.
