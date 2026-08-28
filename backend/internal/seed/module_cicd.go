package seed

func moduleCICD() ModuleSeed {
	return ModuleSeed{
		Title:   "CI/CD: автоматическая доставка",
		Summary: "Как код сам проверяется, собирается и попадает на сервер",
		Lessons: []LessonSeed{
			{
				Title:       "Конвейер: что происходит после коммита",
				Kind:        "text",
				Summary:     "Что такое CI и CD и из каких шагов состоит доставка",
				DurationMin: 12,
				Content: map[string]any{
					"body": "## Зачем это нужно\n\n" +
						"Пока команда маленькая, код на сервер можно положить руками. " +
						"Дальше начинаются проблемы: кто-то забыл прогнать тесты, кто-то выкатил не ту ветку, " +
						"кто-то не может повторить сборку коллеги.\n\n" +
						"Конвейер убирает человека из рутины. Машина делает одно и то же одинаково, каждый раз.\n\n" +
						"## Два слова, которые вечно путают\n\n" +
						"**CI** — непрерывная интеграция. После каждого коммита автоматика проверяет код: " +
						"прогоняет тесты, ищет ошибки стиля, собирает программу.\n\n" +
						"**CD** — непрерывная доставка. Проверенный код автоматически едет на сервер.\n\n" +
						"Проще так: **CI отвечает на вопрос «код рабочий?», CD — «код у пользователей?».**\n\n" +
						"## Из чего состоит конвейер\n\n" +
						"```\n" +
						"коммит → линтер → тесты → сборка образа → выкат на stage → проверка → выкат на прод\n" +
						"```\n\n" +
						"Каждый шаг запускается автоматически. Если шаг упал, дальше конвейер не идёт.\n\n" +
						"Порядок не случайный. Действует правило: **дешёвые проверки идут первыми.** " +
						"Линтер отрабатывает за секунды, сборка образа — за минуты. " +
						"Нет смысла тратить минуты на код, который не прошёл проверку за секунду.\n\n" +
						"## Окружения\n\n" +
						"| Окружение | Для чего | Кто пользуется |\n" +
						"|---|---|---|\n" +
						"| dev | ветки разработчиков, всё ломается | команда |\n" +
						"| stage | копия прода для проверки | тестировщики, заказчик |\n" +
						"| prod | боевой сервер | настоящие пользователи |\n\n" +
						"На прод попадает только то, что прошло stage.\n\n" +
						"## Как выглядит конвейер в файле\n\n" +
						"В GitHub это файл `.github/workflows/ci.yml`:\n\n" +
						"```yaml\n" +
						"name: ci\n" +
						"\n" +
						"on:\n" +
						"  push:\n" +
						"    branches: [main]     # запускать при push в main\n" +
						"\n" +
						"jobs:\n" +
						"  build:\n" +
						"    runs-on: ubuntu-latest\n" +
						"    steps:\n" +
						"      - uses: actions/checkout@v4    # скачать код\n" +
						"      - run: go test ./...           # прогнать тесты\n" +
						"```\n\n" +
						"Читается почти как обычный текст: «при push в main запусти сборку на Ubuntu, " +
						"скачай код, прогони тесты».\n\n" +
						"## Пароли в конвейере\n\n" +
						"Конвейеру часто нужен пароль: положить образ в хранилище, зайти на сервер.\n\n" +
						"**Пароли никогда не пишут прямо в файле.** Их кладут в секреты — отдельное " +
						"защищённое хранилище CI:\n\n" +
						"```yaml\n" +
						"run: docker login -u ci -p ${{ secrets.REGISTRY_PASSWORD }}\n" +
						"```\n\n" +
						"В логах на месте пароля будут звёздочки.\n\n" +
						"## Способы выката\n\n" +
						"- **Постепенный (rolling).** Копии программы обновляются по очереди, сервис не останавливается.\n" +
						"- **Сине-зелёный (blue/green).** Рядом поднимают новую версию и разом переключают на неё трафик.\n" +
						"- **Канареечный (canary).** Новую версию сначала видят 5% пользователей. Если всё хорошо — остальные.\n\n" +
						"У любого способа обязателен план отката. Нет отката — нет выката.\n\n" +
						"## Частые ошибки новичка\n\n" +
						"- **Сборка образа раньше тестов.** Тратите минуты впустую.\n" +
						"- **Пароль строкой в файле конвейера.** Он попадёт в репозиторий и в логи.\n" +
						"- **Выкат сразу в прод, минуя stage.** Ошибку увидят пользователи.\n\n" +
						"## Запомнить\n\n" +
						"1. CI проверяет код, CD доставляет его пользователям.\n" +
						"2. Дешёвые проверки — первыми, прод — последним.\n" +
						"3. Пароли только в секретах, никогда в файле.",
					"resources": []map[string]any{
						{
							"title": "GitHub Actions — документация",
							"url":   "https://docs.github.com/en/actions",
							"note":  "самый распространённый CI: примеры и справочник",
						},
						{
							"title": "GitLab CI/CD — документация",
							"url":   "https://docs.gitlab.com/ci/",
							"note":  "второй по популярности: принципы те же, синтаксис другой",
						},
					},
				},
			},
			{
				Title:       "GitOps: репозиторий как источник правды",
				Kind:        "text",
				Summary:     "Современный способ выката: описали в Git — кластер сам применил",
				DurationMin: 11,
				Content: map[string]any{
					"body": "## Зачем это нужно\n\n" +
						"Обычный CD работает так: конвейер получает доступ к серверу и что-то там меняет. " +
						"Значит, у конвейера есть ключи от прода, а что реально происходит на сервере — " +
						"вопрос доверия к логам.\n\n" +
						"GitOps переворачивает схему: **сервер сам забирает изменения из репозитория.**\n\n" +
						"## Как это работает\n\n" +
						"```\n" +
						"вы правите файл в Git  →  агент внутри кластера видит изменение\n" +
						"                       →  приводит кластер к описанному состоянию\n" +
						"```\n\n" +
						"В кластере живёт программа-агент (обычно Argo CD или Flux). " +
						"Она постоянно сравнивает: что описано в репозитории и что реально запущено. " +
						"Есть разница — устраняет.\n\n" +
						"## Что это даёт\n\n" +
						"- **История изменений.** Кто и когда поменял настройку — видно в Git.\n" +
						"- **Откат — обычный revert.** Вернули коммит, агент вернул старое состояние.\n" +
						"- **Меньше доступов.** У конвейера больше нет ключей от прода.\n" +
						"- **Никаких «поправил руками и забыл».** Ручное изменение агент вернёт обратно.\n\n" +
						"Последний пункт называется дрейфом конфигурации. Это когда на сервере " +
						"постепенно накапливаются ручные правки, о которых никто не помнит. " +
						"GitOps лечит это автоматически.\n\n" +
						"## Как выглядит на практике\n\n" +
						"Обычно заводят два репозитория:\n\n" +
						"1. **Код приложения.** Здесь конвейер собирает образ и публикует его.\n" +
						"2. **Описание окружения.** Здесь лежат манифесты: какая версия образа где запущена.\n\n" +
						"Выкат новой версии превращается в изменение одной строки во втором репозитории:\n\n" +
						"```yaml\n" +
						"image: registry.example.com/api:1.4.3   # было 1.4.2\n" +
						"```\n\n" +
						"Дальше работает агент.\n\n" +
						"## Четыре принципа GitOps\n\n" +
						"1. Состояние системы описано декларативно — в файлах, а не в командах.\n" +
						"2. Описание хранится в Git с историей версий.\n" +
						"3. Изменения применяются автоматически.\n" +
						"4. Агент постоянно проверяет и исправляет расхождения.\n\n" +
						"## Частые ошибки новичка\n\n" +
						"- **Правят кластер руками при GitOps.** Агент вернёт как было — и это правильно.\n" +
						"- **Держат секреты в том же репозитории.** Для них есть отдельные механизмы.\n" +
						"- **Смешивают код и манифесты в одном репозитории.** Тогда каждая правка кода трогает прод.\n\n" +
						"## Запомнить\n\n" +
						"1. Источник правды — репозиторий, а не сервер.\n" +
						"2. Выкат = изменение файла в Git, откат = revert.\n" +
						"3. Ручные правки на сервере агент отменит.",
					"resources": []map[string]any{
						{
							"title": "OpenGitOps — принципы",
							"url":   "https://opengitops.dev/",
							"note":  "те самые четыре принципа, коротко и официально",
						},
						{
							"title": "Argo CD — документация",
							"url":   "https://argo-cd.readthedocs.io/",
							"note":  "самый популярный GitOps-агент для Kubernetes",
						},
					},
				},
			},
			{
				Title:       "Практика: GitHub Actions",
				Kind:        "code",
				Summary:     "Соберите рабочий workflow для сборки и тестов",
				DurationMin: 25,
				Content: map[string]any{
					"resources": []map[string]any{
						{
							"title": "Синтаксис workflow GitHub Actions",
							"url":   "https://docs.github.com/en/actions/writing-workflows/workflow-syntax-for-github-actions",
							"note":  "полный справочник по ключам on, jobs, steps, matrix",
						},
						{
							"title": "Доступ к облаку без секретов: OIDC",
							"url":   "https://docs.github.com/en/actions/security-for-github-actions/security-hardening-your-deployments/about-security-hardening-with-openid-connect",
							"note":  "как выдать пайплайну короткоживущий токен вместо вечного ключа",
						},
					},
					"language": "yaml",
					"task": "Допишите workflow так, чтобы он:\n\n" +
						"1. запускался при push в ветку `main`;\n" +
						"2. выполнял шаг с `actions/checkout`;\n" +
						"3. запускал тесты командой `go test ./...`;\n" +
						"4. собирал Docker-образ на шаге сборки;\n" +
						"5. брал пароль реестра из `secrets`, а не из открытого текста.",
					"starter": "name: ci\n" +
						"\n" +
						"on:\n" +
						"  push:\n" +
						"    branches: []\n" +
						"\n" +
						"jobs:\n" +
						"  build:\n" +
						"    runs-on: ubuntu-latest\n" +
						"    steps:\n" +
						"      - name: Setup Go\n" +
						"        uses: actions/setup-go@v5\n" +
						"        with:\n" +
						"          go-version: '1.25'\n",
					"hint": "Секреты подставляются как ${{ secrets.ИМЯ }}.",
					"solution": "name: ci\n" +
						"\n" +
						"on:\n" +
						"  push:\n" +
						"    branches: [main]\n" +
						"\n" +
						"jobs:\n" +
						"  build:\n" +
						"    runs-on: ubuntu-latest\n" +
						"    steps:\n" +
						"      - uses: actions/checkout@v4\n" +
						"\n" +
						"      - name: Setup Go\n" +
						"        uses: actions/setup-go@v5\n" +
						"        with:\n" +
						"          go-version: '1.25'\n" +
						"\n" +
						"      - name: Test\n" +
						"        run: go test ./...\n" +
						"\n" +
						"      - name: Build image\n" +
						"        run: docker build -t registry.example.com/api:${{ github.sha }} .\n" +
						"\n" +
						"      - name: Login to registry\n" +
						"        run: echo \"${{ secrets.REGISTRY_PASSWORD }}\" | docker login registry.example.com -u ci --password-stdin\n",
					"checks": []map[string]any{
						{"type": "regex", "value": "branches:\\s*(\\[\\s*(main|\"main\"|'main')\\s*\\]|\\n\\s*-\\s*(main|\"main\"|'main'))", "message": "Workflow запускается на push в main"},
						{"type": "contains", "value": "actions/checkout", "message": "Код выгружается шагом checkout"},
						{"type": "contains", "value": "go test ./...", "message": "Тесты запускаются"},
						{"type": "regex", "value": "docker\\s+build", "message": "Образ собирается"},
						{"type": "regex", "value": "secrets\\.[A-Z_]+", "message": "Пароль берётся из секретов"},
						{"type": "notContains", "value": "password: ", "message": "Пароль не записан открытым текстом"},
					},
				},
			},
			{
				Title:       "Тренажёр: разбор конвейера",
				Kind:        "terminal",
				Summary:     "Проверьте, что собирается и выкатывается на самом деле",
				DurationMin: 18,
				Content: map[string]any{
					"resources": []map[string]any{
						{
							"title": "GitLab CI/CD — документация",
							"url":   "https://docs.gitlab.com/ci/",
							"note":  "второй по распространённости CI: те же принципы, другой синтаксис",
						},
						{
							"title": "Кэширование зависимостей в Actions",
							"url":   "https://docs.github.com/en/actions/writing-workflows/choosing-what-your-workflow-does/caching-dependencies-to-speed-up-workflows",
							"note":  "первое, что ускоряет медленный пайплайн",
						},
					},
					"intro": "Пайплайн упал на выкате. Проверьте, что происходит на сервере сборки.",
					"shell": "student@devops",
					"tasks": []map[string]any{
						{
							"id":     "c1",
							"prompt": "Посмотрите файл конвейера ~/projects/api/.github/workflows/ci.yml",
							"expected": []string{
								"cat ~/projects/api/.github/workflows/ci.yml",
								"cat /home/student/projects/api/.github/workflows/ci.yml",
							},
							"hint":    "cat и путь к файлу",
							"success": "Видно этапы: checkout и тесты.",
						},
						{
							"id":     "c2",
							"prompt": "Соберите образ api:1.4.2 из каталога ~/projects/api",
							"expected": []string{
								"docker build -t api:1.4.2 ~/projects/api",
								"docker build -t api:1.4.2 .",
							},
							"hint":    "docker build -t имя:тег путь",
							"success": "Образ собран с явным тегом версии.",
						},
						{
							"id":       "c3",
							"prompt":   "Посмотрите список локальных образов",
							"expected": []string{"docker images", "docker image ls"},
							"hint":     "docker images",
							"success":  "Образ на месте — сборка отработала.",
						},
						{
							"id":       "c4",
							"prompt":   "Проверьте историю коммитов в компактном виде, чтобы понять, что выкатывается",
							"expected": []string{"git log --oneline"},
							"hint":     "git log --oneline",
							"success":  "Понятно, какой коммит поедет в прод.",
						},
						{
							"id":     "c5",
							"prompt": "Проверьте статус выката деплоймента api в кластере",
							"expected": []string{
								"kubectl rollout status deploy/api",
								"kubectl rollout status deployment/api",
								"kubectl rollout status deployment api",
							},
							"hint":    "kubectl rollout status deploy/имя",
							"success": "Выкат завис на одной реплике — дальше смотрим поды и логи.",
						},
					},
				},
			},
			{
				Title:       "Проверка: CI/CD",
				Kind:        "quiz",
				Summary:     "Этапы конвейера, окружения и выкаты",
				DurationMin: 8,
				Content: map[string]any{
					"resources": []map[string]any{
						{
							"title": "SLSA — уровни зрелости цепочки поставки",
							"url":   "https://slsa.dev/",
							"note":  "что такое provenance и зачем подписывать сборку",
						},
						{
							"title": "Argo CD — документация",
							"url":   "https://argo-cd.readthedocs.io/",
							"note":  "если дальше идёте в GitOps: установка, приложения, синхронизация",
						},
					},
					"passScore": 70,
					"questions": []map[string]any{
						{
							"id":   "q1",
							"text": "Чем Continuous Delivery отличается от Continuous Deployment?",
							"options": []map[string]any{
								{"id": "a", "text": "При Delivery выкат в прод запускает человек, при Deployment — автоматика", "correct": true},
								{"id": "b", "text": "Это одно и то же", "correct": false},
								{"id": "c", "text": "Delivery работает только с контейнерами", "correct": false},
							},
							"explanation": "В обоих случаях сборка готова к выкату, разница — в наличии ручного шага.",
						},
						{
							"id":   "q2",
							"text": "Почему линтер ставят раньше сборки образа?",
							"options": []map[string]any{
								{"id": "a", "text": "Дешёвые проверки должны отсекать брак до дорогих шагов", "correct": true},
								{"id": "b", "text": "Линтер не умеет работать после сборки", "correct": false},
								{"id": "c", "text": "Так требует Docker", "correct": false},
							},
							"explanation": "Быстрая обратная связь: не тратим минуты сборки на код, который не проходит формальные проверки.",
						},
						{
							"id":   "q3",
							"text": "Что такое canary-выкат?",
							"options": []map[string]any{
								{"id": "a", "text": "Новая версия получает малую долю трафика, затем долю увеличивают", "correct": true},
								{"id": "b", "text": "Полное переключение трафика между двумя окружениями", "correct": false},
								{"id": "c", "text": "Откат к предыдущей версии", "correct": false},
							},
							"explanation": "Canary даёт увидеть проблему на небольшой части пользователей.",
						},
						{
							"id":       "q4",
							"text":     "Где допустимо хранить пароль от реестра образов?",
							"multiple": true,
							"options": []map[string]any{
								{"id": "a", "text": "В секретах CI-системы", "correct": true},
								{"id": "b", "text": "В Vault или другом хранилище секретов", "correct": true},
								{"id": "c", "text": "В репозитории в файле .env", "correct": false},
								{"id": "d", "text": "В открытом виде в workflow", "correct": false},
							},
							"explanation": "Всё, что попало в git, считается скомпрометированным.",
						},
					},
				},
			},
		},
	}
}
