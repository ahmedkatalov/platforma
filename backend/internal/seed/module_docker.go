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
						"## Как это выглядит\n" +
						"\n" +
						"Запустим nginx и посмотрим на результат. Флаг `--name web` задаёт имя — по нему потом читаем логи и останавливаем контейнер.\n" +
						"\n" +
						"```bash\n" +
						"$ docker run -d -p 8080:80 --name web nginx:1.27\n" +
						"a3f1c2b9e7d4f60b2a1c\n" +
						"\n" +
						"$ docker ps\n" +
						"CONTAINER ID   IMAGE        STATUS         PORTS                  NAMES\n" +
						"a3f1c2b9e7d4   nginx:1.27   Up 6 seconds   0.0.0.0:8080->80/tcp   web\n" +
						"```\n" +
						"\n" +
						"Колонка PORTS подтверждает проброс: снаружи 8080, внутри контейнера 80.\n" +
						"\n" +
						"## Когда контейнер сразу падает\n" +
						"\n" +
						"Бывает, контейнер запустился и тут же остановился. В `docker ps` его уже нет, но `docker ps -a` и логи всё покажут.\n" +
						"\n" +
						"```bash\n" +
						"$ docker run -d --name api myapp:1.0\n" +
						"9f3a1b2c4d5e6a7b8c9d\n" +
						"\n" +
						"$ docker ps\n" +
						"CONTAINER ID   IMAGE   STATUS   PORTS   NAMES\n" +
						"\n" +
						"$ docker ps -a\n" +
						"CONTAINER ID   IMAGE       STATUS                     NAMES\n" +
						"9f3a1b2c4d5e   myapp:1.0   Exited (1) 4 seconds ago   api\n" +
						"\n" +
						"$ docker logs api\n" +
						"panic: open /config/app.yaml: no such file or directory\n" +
						"```\n" +
						"\n" +
						"Порядок чтения простой. `docker ps` пуст, значит смотрим `docker ps -a`.\n" +
						"\n" +
						"Статус `Exited (1)` — код выхода 1, программа упала. Причину даёт `docker logs`.\n" +
						"\n" +
						"Здесь не хватило файла конфига внутри контейнера: его забыли положить в образ или пробросить томом.\n" +
						"\n" +
						"\n" +
						"## Важно: данные внутри контейнера не вечны\n\n" +
						"Удалили контейнер — исчезло всё, что он записал внутрь себя.\n\n" +
						"Поэтому базы данных и загруженные файлы хранят в **томах** — отдельном хранилище, " +
						"которое живёт независимо от контейнера:\n\n" +
						"```bash\n" +
						"docker run -v pgdata:/var/lib/postgresql/data postgres:16\n" +
						"```\n\n" +
						"## Когда порт уже занят\n" +
						"\n" +
						"Частая ошибка при запуске — порт хоста уже держит другой контейнер или процесс.\n" +
						"\n" +
						"```bash\n" +
						"$ docker run -d -p 8080:80 --name web nginx:1.27\n" +
						"docker: Error response from daemon: driver failed programming external\n" +
						"connectivity on endpoint web: Bind for 0.0.0.0:8080 failed: port is already allocated.\n" +
						"\n" +
						"$ docker ps\n" +
						"CONTAINER ID   IMAGE        STATUS         PORTS                  NAMES\n" +
						"c4d1a9f2b8e7   nginx:1.27   Up 3 minutes   0.0.0.0:8080->80/tcp   web-old\n" +
						"\n" +
						"$ docker run -d -p 8081:80 --name web nginx:1.27\n" +
						"b7e2c1a4d9f0e3a5c6d7\n" +
						"```\n" +
						"\n" +
						"Текст `port is already allocated` говорит прямо: порт 8080 занят.\n" +
						"\n" +
						"`docker ps` находит виновника — старый контейнер `web-old` уже слушает 8080.\n" +
						"\n" +
						"Решение — освободить порт через `docker stop web-old` либо взять свободный, как здесь 8081.\n" +
						"\n" +
						"\n" +
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
						{
							"title": "Docker overview: образ, контейнер, реестр",
							"url":   "https://docs.docker.com/get-started/docker-overview/",
							"note":  "официальное объяснение образа, контейнера и registry",
						},
						{
							"title": "Справочник docker run",
							"url":   "https://docs.docker.com/reference/cli/docker/container/run/",
							"note":  "все флаги запуска: -d, -p, --name, -v",
						},
						{
							"title": "Docker volumes",
							"url":   "https://docs.docker.com/engine/storage/volumes/",
							"note":  "почему данные держат в томах, а не внутри контейнера",
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
						"Так выглядит повторная сборка. Строки CACHED — это слои, которые Docker не пересобирал:\n" +
						"\n" +
						"```bash\n" +
						"$ docker build -t myapp:1.0 .\n" +
						"[+] Building 6.1s (9/9) FINISHED\n" +
						" => [1/4] FROM golang:1.25-alpine              0.0s\n" +
						" => CACHED [2/4] COPY go.mod go.sum ./         0.0s\n" +
						" => CACHED [3/4] RUN go mod download           0.0s\n" +
						" => [4/4] COPY . .                             0.5s\n" +
						" => exporting to image                         0.4s\n" +
						" => => naming to docker.io/library/myapp:1.0   0.0s\n" +
						"```\n" +
						"\n" +
						"Зависимости не менялись — шаги 2 и 3 взяты из кэша, сборка заняла секунды.\n" +
						"\n" +
						"Разница между этими вариантами — минуты на каждой сборке.\n\n" +
						"## Когда сборка падает\n" +
						"\n" +
						"Сборка обрывается на конкретном шаге. Docker показывает номер шага, инструкцию и вывод команды.\n" +
						"\n" +
						"```bash\n" +
						"$ docker build -t myapp:1.0 .\n" +
						"[+] Building 4.2s (8/9)\n" +
						" => [1/5] FROM golang:1.25-alpine                  0.0s\n" +
						" => [2/5] WORKDIR /src                             0.1s\n" +
						" => [3/5] COPY go.mod go.sum ./                    0.1s\n" +
						" => [4/5] RUN go mod download                      2.0s\n" +
						" => ERROR [5/5] RUN go build -o /app ./cmd/api     1.9s\n" +
						"------\n" +
						" > [5/5] RUN go build -o /app ./cmd/api:\n" +
						"0.9 cmd/api/main.go:12:2: undefined: hndler\n" +
						"------\n" +
						"Dockerfile:6\n" +
						"--------------------\n" +
						"   5 |     COPY . .\n" +
						"   6 | >>> RUN go build -o /app ./cmd/api\n" +
						"   7 |\n" +
						"--------------------\n" +
						"ERROR: failed to solve: process \"/bin/sh -c go build\" did not complete successfully: exit code: 1\n" +
						"```\n" +
						"\n" +
						"Читаем снизу вверх. `exit code: 1` — команда внутри `RUN` завершилась с ошибкой.\n" +
						"\n" +
						"Стрелки `>>>` указывают на проблемную инструкцию: строка 6 Dockerfile.\n" +
						"\n" +
						"Сама причина выше: компилятор нашёл опечатку `undefined: hndler` в файле и строке.\n" +
						"\n" +
						"\n" +
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
						"## Когда сервис в Compose не поднимается\n" +
						"\n" +
						"Запустили `docker compose up`, но один сервис крутится в перезапуске. Смотрим статусы и логи.\n" +
						"\n" +
						"```bash\n" +
						"$ docker compose up -d\n" +
						"[+] Running 2/2\n" +
						" Container proj-db-1   Started\n" +
						" Container proj-app-1  Started\n" +
						"\n" +
						"$ docker compose ps\n" +
						"NAME         IMAGE         STATUS                       PORTS\n" +
						"proj-app-1   proj-app      Restarting (1) 5 seconds ago\n" +
						"proj-db-1    postgres:16   Up 9 seconds                 5432/tcp\n" +
						"\n" +
						"$ docker compose logs app --tail 2\n" +
						"proj-app-1  | dial tcp 172.18.0.2:5432: connect: connection refused\n" +
						"proj-app-1  | exit status 1\n" +
						"```\n" +
						"\n" +
						"Статус `Restarting` — контейнер падает, и Compose поднимает его снова.\n" +
						"\n" +
						"Лог объясняет причину: `connection refused` на порт 5432, база ещё не готова к соединениям.\n" +
						"\n" +
						"Важный нюанс: `depends_on` ждёт запуска контейнера базы, но не её готовности.\n" +
						"\n" +
						"Лечится проверкой `healthcheck` у базы или повторами подключения в самом приложении.\n" +
						"\n" +
						"\n" +
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
						{
							"title": "Dockerfile reference",
							"url":   "https://docs.docker.com/reference/dockerfile/",
							"note":  "все инструкции: FROM, COPY, RUN, CMD, ENTRYPOINT",
						},
						{
							"title": "Docker build cache",
							"url":   "https://docs.docker.com/build/cache/",
							"note":  "как работает кэш слоёв и что его инвалидирует",
						},
						{
							"title": "Multi-stage builds",
							"url":   "https://docs.docker.com/build/building/multi-stage/",
							"note":  "официально про многоэтапную сборку и уменьшение образа",
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
			{
				Title:       "Отладка контейнеров: заглянуть внутрь",
				Kind:        "text",
				Summary:     "logs, exec, inspect и stats — как понять, почему контейнер упал или тормозит",
				DurationMin: 14,
				Content: map[string]any{
					"body": "## Зачем это нужно\n\n" +
						"Контейнер упал или тормозит. **Догадки не помогут — нужно заглянуть внутрь.** " +
						"Для этого хватает четырёх команд.\n\n" +
						"## 1. Статус — docker ps\n\n" +
						"Первый вопрос: контейнер вообще работает?\n\n" +
						"```\n" +
						"student@devops:~$ docker ps -a\n" +
						"CONTAINER ID   IMAGE        STATUS                      NAMES\n" +
						"3f1c9a4b7e21   nginx:1.27   Up 12 minutes               web\n" +
						"9b2d8c5f0a13   api:1.5.0    Exited (1) 30 seconds ago   api\n" +
						"```\n\n" +
						"Флаг `-a` показывает и остановленные. **`Exited (1)`** значит: контейнер `api` " +
						"упал с кодом 1. Дальше выясняем почему.\n\n" +
						"## 2. Логи — docker logs\n\n" +
						"```\n" +
						"student@devops:~$ docker logs api\n" +
						"2026-08-30T09:15:02Z ERROR failed to connect to db: connection refused\n" +
						"```\n\n" +
						"Логи почти всегда называют причину. Здесь приложение не достучалось до базы.\n\n" +
						"## 3. Заглянуть внутрь — docker exec\n\n" +
						"Иногда нужно проверить что-то прямо в контейнере — файл, переменную, сеть:\n\n" +
						"```\n" +
						"student@devops:~$ docker exec -it web sh\n" +
						"/ # (вы внутри контейнера, для выхода наберите exit)\n" +
						"```\n\n" +
						"`-it` даёт интерактивную оболочку. Это как «зайти по ssh», но в контейнер.\n\n" +
						"## 4. Подробности — docker inspect\n\n" +
						"`inspect` выдаёт всё о контейнере в JSON. Два самых полезных поля:\n\n" +
						"```\n" +
						"student@devops:~$ docker inspect web\n" +
						"        \"State\": { \"Status\": \"running\", \"ExitCode\": 0, \"OOMKilled\": false },\n" +
						"        \"NetworkSettings\": { \"IPAddress\": \"172.17.0.2\", \"Ports\": { \"80/tcp\": \"8080\" } }\n" +
						"```\n\n" +
						"> **`OOMKilled: true`** — важнейший признак: контейнеру не хватило памяти, и его убило ядро. " +
						"Тогда причина не в коде, а в лимите памяти.\n\n" +
						"## 5. Ресурсы — docker stats\n\n" +
						"```\n" +
						"student@devops:~$ docker stats\n" +
						"CONTAINER   CPU %   MEM USAGE / LIMIT   MEM %\n" +
						"web         0.03%   12.4MiB / 512MiB    2.42%\n" +
						"```\n\n" +
						"Показывает нагрузку в реальном времени — сразу видно, кто ест память или процессор.\n\n" +
						"## Порядок отладки\n\n" +
						"1. Работает? → `docker ps -a` (смотрим STATUS).\n" +
						"2. Что в логах? → `docker logs`.\n" +
						"3. Хватает памяти? → `docker inspect` (OOMKilled) и `docker stats`.\n" +
						"4. Нужно проверить руками? → `docker exec -it ... sh`.\n\n" +
						"## Запомнить\n\n" +
						"1. `docker ps -a` показывает упавшие контейнеры и код выхода.\n" +
						"2. Причина почти всегда в `docker logs`.\n" +
						"3. `OOMKilled: true` в inspect — контейнеру не хватило памяти.",
					"resources": []map[string]any{
						{
							"title": "Docker — просмотр логов контейнера",
							"url":   "https://docs.docker.com/reference/cli/docker/container/logs/",
							"note":  "флаги -f, --tail, --since для чтения логов",
						},
						{
							"title": "docker exec — команды в работающем контейнере",
							"url":   "https://docs.docker.com/reference/cli/docker/container/exec/",
							"note":  "как зайти внутрь и что там можно",
						},
						{
							"title": "docker inspect — метаданные контейнера",
							"url":   "https://docs.docker.com/reference/cli/docker/inspect/",
							"note":  "State, OOMKilled, сеть и тома одним ответом",
						},
					},
				},
			},
			{
				Title:       "Тренажёр: чиним упавший контейнер",
				Kind:        "terminal",
				Summary:     "Пройдите отладку по шагам: статус, логи, внутренности",
				DurationMin: 16,
				Content: map[string]any{
					"intro": "Контейнер api упал. Проведите разбор: статус, логи, ресурсы, загляните внутрь работающего web.",
					"shell": "student@devops",
					"tasks": []map[string]any{
						{"id": "d1", "prompt": "Покажите все контейнеры, включая остановленные", "expected": []string{"docker ps -a", "docker ps --all"}, "hint": "docker ps с флагом -a", "success": "Видно, что api в статусе Exited (1) — упал."},
						{"id": "d2", "prompt": "Прочитайте логи контейнера api", "expected": []string{"docker logs api"}, "hint": "docker logs и имя", "success": "В логах причина: не достучался до базы."},
						{"id": "d3", "prompt": "Посмотрите подробности контейнера web (inspect)", "expected": []string{"docker inspect web"}, "hint": "docker inspect и имя", "success": "В State виден ExitCode и OOMKilled, в конце — сеть и порты."},
						{"id": "d4", "prompt": "Проверьте потребление ресурсов контейнерами", "expected": []string{"docker stats"}, "hint": "одна команда docker", "success": "Память и CPU в норме — дело не в ресурсах."},
						{"id": "d5", "prompt": "Зайдите внутрь работающего контейнера web через sh", "expected": []string{"docker exec -it web sh", "docker exec web sh", "docker exec -it web /bin/sh"}, "hint": "docker exec -it имя sh", "success": "Вы внутри контейнера — можно проверить файлы и сеть."},
					},
					"resources": []map[string]any{
						{"title": "Docker — отладка контейнеров", "url": "https://docs.docker.com/engine/containers/run/", "note": "жизненный цикл и диагностика контейнера"},
					},
				},
			},
			{
				Title:       "Квиз: отладка контейнеров",
				Kind:        "quiz",
				Summary:     "Статус, логи, inspect и нехватка памяти",
				DurationMin: 8,
				Content: map[string]any{
					"passScore": 70,
					"questions": []map[string]any{
						{"id": "q1", "text": "Контейнер в статусе Exited (1). Что это значит?", "options": []map[string]any{{"id": "a", "text": "Он упал с кодом ошибки 1", "correct": true}, {"id": "b", "text": "Он успешно завершился", "correct": false}, {"id": "c", "text": "Он приостановлен", "correct": false}}, "explanation": "Код 0 — успех, любой другой — ошибка. Дальше смотрим docker logs."},
						{"id": "q2", "text": "Где быстрее всего найти причину падения контейнера?", "options": []map[string]any{{"id": "a", "text": "В docker logs", "correct": true}, {"id": "b", "text": "В docker images", "correct": false}, {"id": "c", "text": "В docker pull", "correct": false}}, "explanation": "Логи почти всегда называют причину прямым текстом."},
						{"id": "q3", "text": "В docker inspect видно OOMKilled: true. О чём это?", "options": []map[string]any{{"id": "a", "text": "Контейнеру не хватило памяти, его убило ядро", "correct": true}, {"id": "b", "text": "Контейнер удалён вручную", "correct": false}, {"id": "c", "text": "Ошибка в образе", "correct": false}}, "explanation": "Причина не в коде, а в лимите памяти — его надо поднять или найти утечку."},
						{"id": "q4", "text": "Зачем docker exec -it web sh?", "options": []map[string]any{{"id": "a", "text": "Зайти внутрь работающего контейнера и проверить всё руками", "correct": true}, {"id": "b", "text": "Пересобрать образ", "correct": false}, {"id": "c", "text": "Остановить контейнер", "correct": false}}, "explanation": "Это как ssh, но в контейнер: файлы, переменные, сеть."},
						{"id": "q5", "review": true, "text": "Повторение: почему данные базы хранят в томе, а не внутри контейнера?", "options": []map[string]any{{"id": "a", "text": "Без тома данные пропадут при пересоздании контейнера", "correct": true}, {"id": "b", "text": "С томом контейнер быстрее", "correct": false}, {"id": "c", "text": "Тома обязательны для любого контейнера", "correct": false}}, "explanation": "Слой контейнера недолговечен — базы всегда работают с томами."},
					},
					"resources": []map[string]any{
						{"title": "Docker — управление ресурсами контейнера", "url": "https://docs.docker.com/engine/containers/resource_constraints/", "note": "лимиты памяти и CPU, чтобы избежать OOMKilled"},
					},
				},
			},
		},
	}
}
