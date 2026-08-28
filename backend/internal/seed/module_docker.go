package seed

func moduleDocker() ModuleSeed {
	return ModuleSeed{
		Title:   "Docker и контейнеры",
		Summary: "Образы, контейнеры, Dockerfile и docker compose",
		Lessons: []LessonSeed{
			{
				Title:       "Образы и контейнеры",
				Kind:        "text",
				Summary:     "Чем контейнер отличается от виртуальной машины и как собрать образ",
				DurationMin: 18,
				Content: map[string]any{
					"body": "## Зачем контейнеры\n\n" +
						"Классическая беда: «на моей машине работает». Контейнер решает её тем, что " +
						"приложение едет вместе со своим окружением — библиотеками, зависимостями, настройками.\n\n" +
						"В отличие от виртуальной машины контейнер не тащит собственное ядро: он использует " +
						"ядро хоста и изолируется механизмами namespaces и cgroups. Поэтому он стартует за " +
						"секунды и весит десятки мегабайт, а не гигабайты.\n\n" +
						"## Образ и контейнер\n\n" +
						"- **Образ (image)** — неизменяемый слепок файловой системы приложения.\n" +
						"- **Контейнер** — запущенный экземпляр образа со своим слоем записи.\n\n" +
						"Из одного образа можно запустить сколько угодно контейнеров.\n\n" +
						"## Основные команды\n\n" +
						"```bash\n" +
						"docker ps                    # запущенные контейнеры\n" +
						"docker ps -a                 # включая остановленные\n" +
						"docker images                # локальные образы\n" +
						"docker pull nginx:1.27       # скачать образ\n" +
						"docker run -d -p 8080:80 nginx  # запустить в фоне с пробросом порта\n" +
						"docker logs -f my-app        # смотреть логи контейнера\n" +
						"docker exec -it my-app sh    # зайти внутрь работающего контейнера\n" +
						"docker stop my-app           # остановить\n" +
						"docker rm my-app             # удалить контейнер\n" +
						"```\n\n" +
						"## Dockerfile\n\n" +
						"```dockerfile\n" +
						"FROM golang:1.25-alpine AS build\n" +
						"WORKDIR /src\n" +
						"COPY go.mod go.sum ./\n" +
						"RUN go mod download\n" +
						"COPY . .\n" +
						"RUN go build -o /app ./cmd/api\n\n" +
						"FROM alpine:3.20\n" +
						"COPY --from=build /app /app\n" +
						"EXPOSE 8080\n" +
						"CMD [\"/app\"]\n" +
						"```\n\n" +
						"Это **многоэтапная сборка**: тяжёлый образ с компилятором остаётся на этапе build, " +
						"а в финальный образ попадает только бинарник. Так образ уменьшается в десятки раз.\n\n" +
						"## Правила, которые экономят время\n\n" +
						"1. Копируйте файлы зависимостей раньше исходников — слой с `go mod download` " +
						"или `npm ci` будет переиспользоваться из кэша.\n" +
						"2. Фиксируйте версии базовых образов: `nginx:1.27`, а не `nginx:latest`.\n" +
						"3. Не кладите секреты в образ — передавайте их переменными окружения.\n" +
						"4. Данные храните в томах (`volumes`), а не внутри контейнера." +
						"\n\n## Что важно знать сегодня\n\n" +
						"- Compose давно живёт как плагин: команда пишется `docker compose`, без дефиса, " +
						"а строка `version:` в начале файла больше не нужна.\n" +
						"- Сборкой занимается BuildKit: он кэширует слои параллельно и умеет монтировать кэш " +
						"зависимостей — `RUN --mount=type=cache`.\n" +
						"- В Kubernetes образы запускает containerd, а не Docker: собранный образ от этого " +
						"не меняется, ведь формат описан стандартом OCI.\n" +
						"- Для минимального финального образа берут `alpine` или distroless-образы: " +
						"меньше пакетов — меньше поверхность атаки и короче отчёт сканера.\n" +
						"- Multi-stage сборка — норма: компилятор и исходники остаются на этапе сборки, " +
						"в финальный образ едет только бинарник.\n\n" +
						"```dockerfile\n" +
						"# кэш зависимостей переживает пересборку\n" +
						"RUN --mount=type=cache,target=/go/pkg/mod go mod download\n" +
						"```",
					"resources": []map[string]any{
						{"title": "Документация Docker", "url": "https://docs.docker.com/", "note": "Основной справочник по CLI, сборке и Compose"},
						{"title": "Как писать Dockerfile", "url": "https://docs.docker.com/build/building/best-practices/", "note": "Официальные рекомендации: слои, кэш, размер образа"},
						{"title": "Спецификации OCI", "url": "https://opencontainers.org/", "note": "Стандарты образа и среды выполнения — на них держится совместимость"},
						{"title": "Docker Compose: описание файла", "url": "https://docs.docker.com/reference/compose-file/", "note": "Актуальный формат Compose без версии сверху файла"},
					},
				},
			},
			{
				Title:       "Тренажёр: работа с контейнерами",
				Kind:        "terminal",
				Summary:     "Запуск, логи, вход внутрь контейнера",
				DurationMin: 20,
				Content: map[string]any{
					"resources": []map[string]any{
						{
							"title": "Справочник команд Docker CLI",
							"url":   "https://docs.docker.com/reference/cli/docker/",
							"note":  "все подкоманды с флагами и примерами",
						},
					},
					"intro": "Учебный хост с установленным Docker. Выполните задания по очереди.",
					"shell": "student@devops",
					"tasks": []map[string]any{
						{
							"id":       "d1",
							"prompt":   "Выведите список запущенных контейнеров",
							"expected": []string{"docker ps"},
							"hint":     "Как ps в Linux, только для Docker",
							"success":  "docker ps показывает работающие контейнеры.",
						},
						{
							"id":       "d2",
							"prompt":   "Покажите все контейнеры, включая остановленные",
							"expected": []string{"docker ps -a", "docker ps --all"},
							"hint":     "Нужен флаг -a",
							"success":  "Флаг -a добавляет остановленные контейнеры.",
						},
						{
							"id":       "d3",
							"prompt":   "Скачайте образ nginx версии 1.27",
							"expected": []string{"docker pull nginx:1.27"},
							"hint":     "docker pull образ:тег",
							"success":  "Версию всегда лучше указывать явно.",
						},
						{
							"id":       "d4",
							"prompt":   "Запустите nginx в фоне, пробросив порт 8080 хоста на 80 контейнера",
							"expected": []string{"docker run -d -p 8080:80 nginx", "docker run -p 8080:80 -d nginx"},
							"hint":     "-d — фон, -p хост:контейнер",
							"success":  "Теперь сайт доступен на localhost:8080.",
						},
						{
							"id":       "d5",
							"prompt":   "Посмотрите логи контейнера web в режиме слежения",
							"expected": []string{"docker logs -f web", "docker logs --follow web"},
							"hint":     "docker logs с флагом -f",
							"success":  "Логи контейнера читаются так же, как tail -f.",
						},
						{
							"id":       "d6",
							"prompt":   "Зайдите внутрь контейнера web интерактивной оболочкой sh",
							"expected": []string{"docker exec -it web sh", "docker exec -ti web sh"},
							"hint":     "docker exec -it имя sh",
							"success":  "Так отлаживают контейнер изнутри.",
						},
						{
							"id":       "d7",
							"prompt":   "Остановите контейнер web",
							"expected": []string{"docker stop web"},
							"hint":     "Команда из одного слова и имя контейнера",
							"success":  "Контейнер остановлен корректно.",
						},
					},
				},
			},
			{
				Title:       "Практика: docker-compose.yml",
				Kind:        "code",
				Summary:     "Опишите приложение и базу данных одним файлом",
				DurationMin: 25,
				Content: map[string]any{
					"resources": []map[string]any{
						{
							"title": "Compose file reference",
							"url":   "https://docs.docker.com/reference/compose-file/",
							"note":  "актуальная спецификация формата: services, volumes, healthcheck, depends_on",
						},
						{
							"title": "Compose в разработке: тома и hot reload",
							"url":   "https://docs.docker.com/compose/how-tos/file-watch/",
							"note":  "как не пересобирать образ на каждое изменение кода",
						},
					},
					"language": "yaml",
					"task": "Допишите docker-compose.yml так, чтобы:\n\n" +
						"1. сервис `app` пробрасывал порт `8080:8080`;\n" +
						"2. `app` зависел от `db` через `depends_on`;\n" +
						"3. у сервиса `db` был том `pgdata` для хранения данных;\n" +
						"4. пароль базы задавался переменной окружения `POSTGRES_PASSWORD`.",
					"starter": "services:\n" +
						"  app:\n" +
						"    build: .\n" +
						"    environment:\n" +
						"      DATABASE_URL: postgres://app:secret@db:5432/app\n" +
						"\n" +
						"  db:\n" +
						"    image: postgres:16\n" +
						"    environment:\n" +
						"      POSTGRES_USER: app\n" +
						"      POSTGRES_DB: app\n",
					"hint": "Тома объявляются дважды: внутри сервиса и в корневой секции volumes.",
					"checks": []map[string]any{
						{"type": "contains", "value": "8080:8080", "message": "Порт 8080 проброшен наружу"},
						{"type": "contains", "value": "depends_on", "message": "app дожидается запуска базы"},
						{"type": "contains", "value": "POSTGRES_PASSWORD", "message": "Пароль базы задан переменной окружения"},
						{"type": "regex", "value": "(?s)volumes:.*pgdata", "message": "Объявлен том pgdata"},
						{"type": "notContains", "value": "latest", "message": "Версии образов зафиксированы, без latest"},
					},
					"solution": "services:\n" +
						"  app:\n" +
						"    build: .\n" +
						"    ports:\n" +
						"      - \"8080:8080\"\n" +
						"    depends_on:\n" +
						"      - db\n" +
						"    environment:\n" +
						"      DATABASE_URL: postgres://app:secret@db:5432/app\n" +
						"\n" +
						"  db:\n" +
						"    image: postgres:16\n" +
						"    environment:\n" +
						"      POSTGRES_USER: app\n" +
						"      POSTGRES_DB: app\n" +
						"      POSTGRES_PASSWORD: secret\n" +
						"    volumes:\n" +
						"      - pgdata:/var/lib/postgresql/data\n" +
						"\n" +
						"volumes:\n" +
						"  pgdata:\n",
				},
			},
			{
				Title:       "Проверка: Docker",
				Kind:        "quiz",
				Summary:     "Образы, слои, тома и сеть",
				DurationMin: 10,
				Content: map[string]any{
					"resources": []map[string]any{
						{
							"title": "Рекомендации по написанию Dockerfile",
							"url":   "https://docs.docker.com/build/building/best-practices/",
							"note":  "слои, кэш, многоэтапная сборка, размер образа",
						},
						{
							"title": "Open Container Initiative — спецификации",
							"url":   "https://opencontainers.org/",
							"note":  "стандарт образов и рантайма: почему образ Docker запускается в Kubernetes",
						},
					},
					"passScore": 70,
					"questions": []map[string]any{
						{
							"id":   "q1",
							"text": "Чем контейнер принципиально отличается от виртуальной машины?",
							"options": []map[string]any{
								{"id": "a", "text": "Контейнер использует ядро хоста, а не своё", "correct": true},
								{"id": "b", "text": "Контейнер всегда работает быстрее любой программы", "correct": false},
								{"id": "c", "text": "В контейнере нельзя запускать базы данных", "correct": false},
							},
							"explanation": "Изоляция достигается namespaces и cgroups, отдельного ядра нет — отсюда лёгкость и скорость старта.",
						},
						{
							"id":   "q2",
							"text": "Что означает -p 8080:80 в команде docker run?",
							"options": []map[string]any{
								{"id": "a", "text": "Порт 8080 хоста ведёт на порт 80 контейнера", "correct": true},
								{"id": "b", "text": "Порт 80 хоста ведёт на порт 8080 контейнера", "correct": false},
								{"id": "c", "text": "Контейнер получит 8080 МБ памяти", "correct": false},
							},
							"explanation": "Формат всегда хост:контейнер.",
						},
						{
							"id":       "q3",
							"text":     "Зачем нужна многоэтапная сборка (multi-stage build)?",
							"multiple": true,
							"options": []map[string]any{
								{"id": "a", "text": "Чтобы не тащить компилятор в финальный образ", "correct": true},
								{"id": "b", "text": "Чтобы уменьшить размер образа", "correct": true},
								{"id": "c", "text": "Чтобы контейнер запускался без Docker", "correct": false},
							},
							"explanation": "Собираем в одном образе, а в финальный копируем только результат сборки.",
						},
						{
							"id":   "q4",
							"text": "Где правильно хранить данные базы, работающей в контейнере?",
							"options": []map[string]any{
								{"id": "a", "text": "Внутри контейнера — так проще", "correct": false},
								{"id": "b", "text": "В именованном томе (volume)", "correct": true},
								{"id": "c", "text": "В образе, добавив их через COPY", "correct": false},
							},
							"explanation": "Контейнер эфемерен: удалили — данные пропали. Тома живут отдельно.",
						},
						{
							"id":   "q5",
							"text": "Почему тег latest — плохая практика в проде?",
							"options": []map[string]any{
								{"id": "a", "text": "Он замедляет загрузку образа", "correct": false},
								{"id": "b", "text": "Сборка перестаёт быть воспроизводимой: сегодня и завтра это разные образы", "correct": true},
								{"id": "c", "text": "latest недоступен в приватных реестрах", "correct": false},
							},
							"explanation": "Фиксируйте версию, чтобы одинаковый код всегда собирался в одинаковое окружение.",
						},
					},
				},
			},
		},
	}
}
