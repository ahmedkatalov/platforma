package seed

func moduleDocker() ModuleSeed {
	return ModuleSeed{
		Title:   "Docker и контейнеры",
		Summary: "Как упаковать программу так, чтобы она одинаково работала везде",
		Lessons: []LessonSeed{
			{
				Title:       "Контейнеры: зачем они нужны",
				Kind:        "text",
				Summary:     "Проблема «на моей машине работает» и как её решает Docker",
				DurationMin: 12,
				Content: map[string]any{
					"body": "## Зачем это нужно\n\n" +
						"Классическая история. Программист написал программу, у него всё работает. " +
						"Отправил на сервер — не запускается.\n\n" +
						"Причина обычно одна: на сервере другая версия языка, не хватает библиотеки " +
						"или отличается настройка. Часы уходят на поиск различий.\n\n" +
						"Контейнер решает это раз и навсегда: программа едет на сервер **вместе со своим окружением**.\n\n" +
						"## Простыми словами\n\n" +
						"Представьте переезд. Можно перевозить вещи по отдельности и надеяться, что на новом месте " +
						"найдётся нужная розетка. А можно упаковать всё в коробку вместе с переходниками.\n\n" +
						"Контейнер — такая коробка. Внутри программа, библиотеки, настройки. " +
						"Она открывается одинаково на ноутбуке, тестовом сервере и в проде.\n\n" +
						"## Чем контейнер отличается от виртуальной машины\n\n" +
						"| | Виртуальная машина | Контейнер |\n" +
						"|---|---|---|\n" +
						"| Что внутри | целая операционная система | только программа и её зависимости |\n" +
						"| Размер | гигабайты | десятки мегабайт |\n" +
						"| Запуск | минуты | секунды |\n\n" +
						"Контейнер использует ядро того сервера, на котором работает. " +
						"Поэтому он лёгкий и быстрый.\n\n" +
						"## Образ и контейнер — разные вещи\n\n" +
						"**Образ** — шаблон, слепок. Как установочный файл программы.\n\n" +
						"**Контейнер** — запущенный экземпляр образа.\n\n" +
						"Из одного образа можно запустить сколько угодно контейнеров. " +
						"Как из одного установочного файла — установить программу на десять компьютеров.\n\n" +
						"## Первые команды\n\n" +
						"```bash\n" +
						"docker ps                        # какие контейнеры сейчас работают\n" +
						"docker images                    # какие образы скачаны\n" +
						"docker pull nginx:1.27           # скачать образ\n" +
						"docker run -d -p 8080:80 nginx   # запустить в фоне\n" +
						"docker logs web                  # посмотреть логи контейнера\n" +
						"docker stop web                  # остановить\n" +
						"```\n\n" +
						"Разберём флаги команды `run`:\n\n" +
						"- `-d` — запустить в фоне, не занимая терминал;\n" +
						"- `-p 8080:80` — проброс порта: обращение к порту 8080 сервера попадёт в порт 80 контейнера.\n\n" +
						"## Важно: данные внутри контейнера не вечны\n\n" +
						"Удалили контейнер — исчезло всё, что он записал внутрь себя.\n\n" +
						"Поэтому базы данных и загруженные файлы хранят в **томах** — отдельном хранилище, " +
						"которое живёт независимо от контейнера:\n\n" +
						"```bash\n" +
						"docker run -v pgdata:/var/lib/postgresql/data postgres:16\n" +
						"```\n\n" +
						"## Частые ошибки новичка\n\n" +
						"- **Пишут `latest` вместо версии.** Сегодня скачается одно, завтра другое. Указывайте версию: `nginx:1.27`.\n" +
						"- **Хранят данные внутри контейнера.** После пересоздания они пропадут.\n" +
						"- **Забывают проброс порта.** Контейнер работает, но снаружи недоступен.\n\n" +
						"## Запомнить\n\n" +
						"1. Образ — шаблон, контейнер — запущенный экземпляр.\n" +
						"2. Данные храните в томах, а не внутри контейнера.\n" +
						"3. Всегда указывайте версию образа вместо `latest`.",
					"resources": []map[string]any{
						{
							"title": "Docker — руководство для начинающих",
							"url":   "https://docs.docker.com/get-started/",
							"note":  "официальный старт: установка и первые контейнеры",
						},
						{
							"title": "Play with Docker — песочница в браузере",
							"url":   "https://labs.play-with-docker.com/",
							"note":  "можно потрогать настоящий Docker без установки",
						},
					},
				},
			},
			{
				Title:       "Квиз: контейнеры и образы",
				Kind:        "quiz",
				Summary:     "Чем образ отличается от контейнера и где живут данные",
				DurationMin: 6,
				Content: map[string]any{
					"passScore": 70,
					"questions": []map[string]any{
						{
							"id":   "c1",
							"text": "Какую проблему решает контейнер?",
							"options": []map[string]any{
								{"id": "a", "text": "«На моей машине работает, а на сервере нет»", "correct": true},
								{"id": "b", "text": "Медленную работу программы", "correct": false},
								{"id": "c", "text": "Нехватку места на диске", "correct": false},
							},
							"explanation": "Программа едет на сервер вместе со своим окружением.",
						},
						{
							"id":   "c2",
							"text": "Чем контейнер отличается от виртуальной машины?",
							"options": []map[string]any{
								{"id": "a", "text": "Контейнер не тащит с собой целую операционную систему, поэтому лёгкий и быстрый", "correct": true},
								{"id": "b", "text": "Контейнер работает только на Linux, а виртуальная машина — везде", "correct": false},
								{"id": "c", "text": "Разницы нет, это одно и то же", "correct": false},
							},
							"explanation": "Контейнер использует ядро того сервера, где запущен.",
						},
						{
							"id":   "c3",
							"text": "Что произойдёт с данными, если удалить контейнер?",
							"options": []map[string]any{
								{"id": "a", "text": "Всё, что он записал внутрь себя, пропадёт", "correct": true},
								{"id": "b", "text": "Данные сохранятся автоматически", "correct": false},
								{"id": "c", "text": "Данные переедут в образ", "correct": false},
							},
							"explanation": "Поэтому базы и загруженные файлы держат в томах.",
						},
						{
							"id":   "c4",
							"text": "Что делает флаг -p 8080:80 в команде docker run?",
							"options": []map[string]any{
								{"id": "a", "text": "Запросы к порту 8080 сервера попадут в порт 80 контейнера", "correct": true},
								{"id": "b", "text": "Ограничивает контейнер по памяти", "correct": false},
								{"id": "c", "text": "Запускает контейнер в фоне", "correct": false},
							},
							"explanation": "В фоне запускает флаг -d.",
						},
						{
							"id":   "c5",
							"text": "Почему не стоит писать nginx:latest?",
							"options": []map[string]any{
								{"id": "a", "text": "Сегодня скачается одна версия, завтра другая — поведение изменится незаметно", "correct": true},
								{"id": "b", "text": "latest всегда содержит ошибки", "correct": false},
								{"id": "c", "text": "Такой образ скачивается дольше", "correct": false},
							},
							"explanation": "Фиксированная версия делает сборку предсказуемой.",
						},
						{
							"id":     "c6",
							"review": true,
							"text":   "Повторение: какой командой посмотреть, какие порты слушает сервер?",
							"options": []map[string]any{
								{"id": "a", "text": "ss -tulpn", "correct": true},
								{"id": "b", "text": "df -h", "correct": false},
								{"id": "c", "text": "ls -la", "correct": false},
							},
							"explanation": "ss показывает открытые порты и процессы, которые их держат.",
						},
					},
					"resources": []map[string]any{
						{
							"title": "Docker — руководство для начинающих",
							"url":   "https://docs.docker.com/get-started/",
							"note":  "официальный старт с первыми контейнерами",
						},
					},
				},
			},
			{
				Title:       "Dockerfile: собираем свой образ",
				Kind:        "text",
				Summary:     "Как из кода получается образ и почему важны слои",
				DurationMin: 12,
				Content: map[string]any{
					"body": "## Зачем это нужно\n\n" +
						"Готовые образы вроде nginx кто-то собрал. Свою программу придётся упаковать самому.\n\n" +
						"Для этого пишут **Dockerfile** — пошаговый рецепт сборки образа.\n\n" +
						"## Как читать Dockerfile\n\n" +
						"```dockerfile\n" +
						"FROM golang:1.25-alpine     # с какого образа начинаем\n" +
						"WORKDIR /src                # рабочая папка внутри образа\n" +
						"COPY . .                    # копируем файлы проекта внутрь\n" +
						"RUN go build -o /app        # выполняем команду при сборке\n" +
						"EXPOSE 8080                 # подсказка: программа слушает этот порт\n" +
						"CMD [\"/app\"]                # что запустить при старте контейнера\n" +
						"```\n\n" +
						"Каждая строка — шаг. Собирается командой:\n\n" +
						"```bash\n" +
						"docker build -t myapp:1.0 .\n" +
						"```\n\n" +
						"`-t myapp:1.0` — имя и версия образа, точка в конце — папка со сборкой.\n\n" +
						"## Слои и кэш\n\n" +
						"Каждая инструкция создаёт слой. Docker запоминает слои и при повторной сборке " +
						"переиспользует неизменившиеся.\n\n" +
						"Отсюда важное правило: **сначала то, что меняется редко, потом то, что часто.**\n\n" +
						"Плохо:\n\n" +
						"```dockerfile\n" +
						"COPY . .\n" +
						"RUN go mod download    # пересобирается при любой правке кода\n" +
						"```\n\n" +
						"Хорошо:\n\n" +
						"```dockerfile\n" +
						"COPY go.mod go.sum ./\n" +
						"RUN go mod download    # выполнится заново только при смене зависимостей\n" +
						"COPY . .\n" +
						"```\n\n" +
						"Разница между этими вариантами — минуты на каждой сборке.\n\n" +
						"## Многоэтапная сборка\n\n" +
						"Для сборки нужен компилятор, а для запуска — нет. Зачем тащить его в готовый образ?\n\n" +
						"```dockerfile\n" +
						"FROM golang:1.25-alpine AS build\n" +
						"WORKDIR /src\n" +
						"COPY . .\n" +
						"RUN go build -o /app ./cmd/api\n" +
						"\n" +
						"FROM alpine:3.20\n" +
						"COPY --from=build /app /app\n" +
						"ENTRYPOINT [\"/app\"]\n" +
						"```\n\n" +
						"Первый этап собирает, второй берёт только результат. " +
						"Образ уменьшается с сотен мегабайт до десятков.\n\n" +
						"## Несколько сервисов сразу: Compose\n\n" +
						"Приложению обычно нужна база. Запускать два контейнера руками неудобно — " +
						"описывают их в одном файле `docker-compose.yml`:\n\n" +
						"```yaml\n" +
						"services:\n" +
						"  app:\n" +
						"    build: .\n" +
						"    ports:\n" +
						"      - \"8080:8080\"\n" +
						"    depends_on:\n" +
						"      - db\n" +
						"\n" +
						"  db:\n" +
						"    image: postgres:16\n" +
						"    environment:\n" +
						"      POSTGRES_PASSWORD: secret\n" +
						"    volumes:\n" +
						"      - pgdata:/var/lib/postgresql/data\n" +
						"\n" +
						"volumes:\n" +
						"  pgdata:\n" +
						"```\n\n" +
						"Запуск одной командой: `docker compose up -d`. Остановка: `docker compose down`.\n\n" +
						"Обратите внимание: `docker compose` пишется без дефиса. Старый вариант " +
						"`docker-compose` встречается в статьях, но это предыдущая версия.\n\n" +
						"Внутри Compose сервисы видят друг друга по имени: приложение подключается к базе " +
						"по адресу `db`, а не по IP.\n\n" +
						"## Частые ошибки новичка\n\n" +
						"- **`COPY . .` в начале файла.** Ломает кэш, каждая сборка идёт с нуля.\n" +
						"- **Секреты внутри Dockerfile.** Они остаются в слоях образа даже после удаления файла.\n" +
						"- **Запуск от root.** Создайте пользователя и добавьте `USER app`.\n\n" +
						"## Запомнить\n\n" +
						"1. Dockerfile — рецепт сборки, каждая строка создаёт слой.\n" +
						"2. Редко меняющееся ставьте выше — сборка будет быстрее.\n" +
						"3. Многоэтапная сборка убирает из образа всё лишнее.",
					"resources": []map[string]any{
						{
							"title": "Рекомендации по Dockerfile",
							"url":   "https://docs.docker.com/build/building/best-practices/",
							"note":  "официальные советы про слои, кэш и размер образа",
						},
						{
							"title": "Справочник docker-compose.yml",
							"url":   "https://docs.docker.com/reference/compose-file/",
							"note":  "все поля файла с примерами",
						},
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
						{
							"id":     "q7",
							"review": true,
							"text":   "Повторение: сервер вернул 502. Где искать причину?",
							"options": []map[string]any{
								{"id": "a", "text": "Приложение не отвечает: упало или не слушает свой порт", "correct": true},
								{"id": "b", "text": "В настройках DNS", "correct": false},
								{"id": "c", "text": "В правах на файлы статики", "correct": false},
							},
							"explanation": "502 означает, что прокси не дозвонился до приложения.",
						},
					},
				},
			},
		},
	}
}
