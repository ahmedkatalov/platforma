package seed

func moduleCICD() ModuleSeed {
	return ModuleSeed{
		Title:   "CI/CD: автоматическая доставка",
		Summary: "Как код сам проходит проверки и попадает на серверы",
		Lessons: []LessonSeed{
			{
				Title:       "Конвейер: что происходит после коммита",
				Kind:        "text",
				Summary:     "CI, CD, этапы сборки и окружения — простыми словами",
				DurationMin: 12,
				Content: map[string]any{
					"body": "## Зачем это нужно\n\n" +
						"Без автоматики выкат выглядит так: разработчик собирает проект у себя, копирует файлы " +
						"на сервер, перезапускает программу. Каждый раз руками.\n\n" +
						"Проблема **не в скорости, а в ошибках**. Забыл прогнать тесты, скопировал не ту версию, " +
						"перепутал сервер. Конвейер убирает эти ошибки: делает одно и то же одинаково.\n\n" +
						"## Два слова, которые путают\n\n" +
						"**CI** — автоматическая проверка. После каждого коммита запускаются тесты и сборка. " +
						"Задача: узнать о поломке через минуты, а не через неделю.\n\n" +
						"**CD** — автоматическая доставка проверенного кода на серверы.\n\n" +
						"У CD есть две версии:\n\n" +
						"- кнопку выката нажимает человек;\n" +
						"- выкат идёт сам, если все проверки прошли.\n\n" +
						"## Как выглядит конвейер\n\n" +
						"```\n" +
						"коммит → проверка стиля → тесты → сборка образа\n" +
						"       → выкат на stage → проверка → выкат на прод\n" +
						"```\n\n" +
						"Порядок не случайный. Правило простое: **дешёвые проверки идут первыми.**\n\n" +
						"Проверка стиля занимает секунды, сборка образа — минуты. Нет смысла тратить минуты " +
						"на код, в котором лишний пробел.\n\n" +
						"## Окружения\n\n" +
						"| Окружение | Для чего | Кто пострадает от поломки |\n" +
						"|---|---|---|\n" +
						"| dev | ветки разработчиков | никто |\n" +
						"| stage | копия прода для проверки | никто |\n" +
						"| prod | боевые серверы | все пользователи |\n\n" +
						"На прод попадает **только то, что уже поработало на stage**.\n\n" +
						"## Как выглядит файл конвейера\n\n" +
						"Пример для GitHub Actions:\n\n" +
						"```yaml\n" +
						"name: ci\n" +
						"\n" +
						"on:\n" +
						"  push:\n" +
						"    branches: [main]      # когда запускать\n" +
						"\n" +
						"jobs:\n" +
						"  build:\n" +
						"    runs-on: ubuntu-latest\n" +
						"    steps:\n" +
						"      - uses: actions/checkout@v4   # забрать код\n" +
						"      - run: go test ./...          # прогнать тесты\n" +
						"```\n\n" +
						"Читается сверху вниз: когда запускать, где выполнять, какие шаги делать.\n\n" +
						"Когда кто-то ломает тест, конвейер краснеет и не пускает код дальше:\n" +
						"\n" +
						"```\n" +
						"push 9f3c1a2 → ci\n" +
						"\n" +
						"✓ checkout\n" +
						"✓ setup-go\n" +
						"✗ go test ./...\n" +
						"  --- FAIL: TestSum (0.00s)\n" +
						"      want 5, got 4\n" +
						"  FAIL  api/internal/math   0.31s\n" +
						"\n" +
						"конвейер остановлен, на прод ничего не поехало\n" +
						"```\n" +
						"\n" +
						"Поломку видно **сразу после коммита**, а не через неделю на проде.\n" +
						"\n" +
						"## Секреты\n\n" +
						"Конвейеру нужны пароли: от реестра образов, от серверов. В файл их не пишут — " +
						"он лежит в репозитории и виден всем.\n\n" +
						"Пароли хранят в настройках CI и подставляют так:\n\n" +
						"```yaml\n" +
						"run: docker login -u ci -p ${{ secrets.REGISTRY_PASSWORD }}\n" +
						"```\n\n" +
						"Современный подход — вообще без паролей: CI получает временный токен на несколько минут " +
						"по протоколу OIDC. **Украсть такой токен бесполезно**, он истечёт.\n\n" +
						"Когда пайплайн краснеет, необязательно открывать браузер. Упавший шаг видно из терминала через `gh`:\n" +
						"\n" +
						"```bash\n" +
						"$ gh run list --limit 3\n" +
						"STATUS  TITLE      WORKFLOW  BRANCH  EVENT  ID          ELAPSED\n" +
						"X       api 1.4.3  ci        main    push   9823471021  1m12s\n" +
						"✓       fix logs   ci        main    push   9821203344  1m40s\n" +
						"✓       bump deps  ci        main    push   9820556610  1m38s\n" +
						"\n" +
						"$ gh run view 9823471021 --log-failed\n" +
						"build  Test  --- FAIL: TestCharge (0.00s)\n" +
						"build  Test      charge_test.go:41: want 1200, got 0\n" +
						"build  Test  FAIL  api/internal/billing  0.28s\n" +
						"build  Test  Error: Process completed with exit code 1\n" +
						"```\n" +
						"\n" +
						"`X` — упавший запуск. `--log-failed` печатает только красный шаг: тест `TestCharge` ждал 1200, а получил 0.\n" +
						"\n" +
						"## Как выкатывают, чтобы не уронить прод\n\n" +
						"- **По очереди.** Обновляем сервер за сервером, старые пока работают.\n" +
						"- **Две среды.** Поднимаем новую версию рядом и переключаем трафик целиком.\n" +
						"- **Понемногу (canary).** Пускаем на новую версию 5% пользователей, смотрим ошибки, потом остальных.\n\n" +
						"И главное правило: **план отката должен быть всегда.** Если откатиться нельзя, " +
						"это не доставка, а лотерея.\n\n" +
						"Сборка образа падает на конкретном слое. По выводу видно, какие слои взялись из кэша, а какой упал:\n" +
						"\n" +
						"```bash\n" +
						"$ docker build -t api:1.4.3 .\n" +
						"[+] Building 8.2s (10/12)\n" +
						" => [1/6] FROM golang:1.25                        0.0s\n" +
						" => CACHED [2/6] WORKDIR /src                      0.0s\n" +
						" => CACHED [3/6] COPY go.mod go.sum ./             0.0s\n" +
						" => CACHED [4/6] RUN go mod download               0.0s\n" +
						" => [5/6] COPY . .                                 0.3s\n" +
						" => ERROR [6/6] RUN go build -o /api ./cmd/api     6.1s\n" +
						"------\n" +
						" > [6/6] RUN go build -o /api ./cmd/api:\n" +
						"0.412 cmd/api/main.go:12:2: undefined: startServer\n" +
						"------\n" +
						"ERROR: failed to solve: process \"/bin/sh -c go build -o /api ./cmd/api\" did not complete successfully: exit code: 1\n" +
						"```\n" +
						"\n" +
						"`CACHED` — слой не пересобирали, взяли из кэша. Упал слой 6: код не компилируется, `undefined: startServer`. Правим код и собираем заново.\n" +
						"\n" +
						"## Частые ошибки новичка\n\n" +
						"- **Пароль прямо в файле конвейера.** Он попадёт в историю Git навсегда.\n" +
						"- **Выкат сразу на прод, минуя stage.** Проверять на пользователях — дорого.\n" +
						"- **Тесты в конце, после сборки.** Порядок наоборот: сначала быстрое, потом долгое.\n\n" +
						"## Запомнить\n\n" +
						"1. CI проверяет код, CD доставляет его на серверы.\n" +
						"2. Дешёвые проверки — в начало конвейера.\n" +
						"3. Секреты — в настройках CI, никогда в файлах репозитория.",
					"resources": []map[string]any{
						{
							"title": "GitHub Actions — документация",
							"url":   "https://docs.github.com/ru/actions",
							"note":  "есть на русском: как написать первый конвейер",
						},
						{
							"title": "Доступ без паролей: OIDC в Actions",
							"url":   "https://docs.github.com/en/actions/security-for-github-actions/security-hardening-your-deployments/about-security-hardening-with-openid-connect",
							"note":  "как выдать конвейеру короткоживущий токен вместо вечного ключа",
						},
						{
							"title": "Continuous Integration — статья Мартина Фаулера",
							"url":   "https://martinfowler.com/articles/continuousIntegration.html",
							"note":  "классическое определение CI и почему поломку ловят за минуты",
						},
						{
							"title": "Continuous Delivery — Мартин Фаулер",
							"url":   "https://martinfowler.com/bliki/ContinuousDelivery.html",
							"note":  "чем Delivery отличается от Deployment: наличие ручного шага выката",
						},
						{
							"title": "Стратегии выката приложений — Google Cloud",
							"url":   "https://cloud.google.com/architecture/application-deployment-and-testing-strategies",
							"note":  "rolling, blue-green и canary: когда что применять",
						},
					},
				},
			},
			{
				Title:       "Квиз: конвейер",
				Kind:        "quiz",
				Summary:     "Этапы, окружения и секреты",
				DurationMin: 6,
				Content: map[string]any{
					"passScore": 70,
					"questions": []map[string]any{
						{
							"id":   "p1",
							"text": "Чем CI отличается от CD?",
							"options": []map[string]any{
								{"id": "a", "text": "CI проверяет код, CD доставляет его на серверы", "correct": true},
								{"id": "b", "text": "CI для тестов, CD для документации", "correct": false},
								{"id": "c", "text": "Это два названия одного и того же", "correct": false},
							},
							"explanation": "Сначала проверка, потом доставка.",
						},
						{
							"id":   "p2",
							"text": "Почему проверку стиля ставят раньше сборки образа?",
							"options": []map[string]any{
								{"id": "a", "text": "Она быстрая: незачем тратить минуты сборки на код с опечаткой", "correct": true},
								{"id": "b", "text": "Линтер не умеет работать после сборки", "correct": false},
								{"id": "c", "text": "Так требует Docker", "correct": false},
							},
							"explanation": "Правило: дешёвые проверки идут первыми.",
						},
						{
							"id":   "p3",
							"text": "Зачем нужен stage, если есть прод?",
							"options": []map[string]any{
								{"id": "a", "text": "Чтобы проверить изменения на копии прода, где нет пользователей", "correct": true},
								{"id": "b", "text": "Чтобы хранить резервные копии", "correct": false},
								{"id": "c", "text": "Чтобы ускорить сборку", "correct": false},
							},
							"explanation": "На проде ошибку увидят все пользователи сразу.",
						},
						{
							"id":   "p4",
							"text": "Где хранить пароль от реестра образов?",
							"options": []map[string]any{
								{"id": "a", "text": "В секретах CI-системы", "correct": true},
								{"id": "b", "text": "Прямо в файле конвейера", "correct": false},
								{"id": "c", "text": "В README проекта", "correct": false},
							},
							"explanation": "Всё, что попало в репозиторий, считается скомпрометированным.",
						},
						{
							"id":   "p5",
							"text": "Что такое canary-выкат?",
							"options": []map[string]any{
								{"id": "a", "text": "Новую версию сначала показывают небольшой части пользователей", "correct": true},
								{"id": "b", "text": "Выкат сразу на все серверы одновременно", "correct": false},
								{"id": "c", "text": "Откат к прошлой версии", "correct": false},
							},
							"explanation": "Если что-то не так, проблему увидят 5% пользователей, а не все.",
						},
						{
							"id":     "p6",
							"review": true,
							"text":   "Повторение: что нужно сделать перед тем, как влить свою ветку в main?",
							"options": []map[string]any{
								{"id": "a", "text": "Открыть pull request и дождаться проверок и ревью", "correct": true},
								{"id": "b", "text": "Удалить ветку", "correct": false},
								{"id": "c", "text": "Сделать git reset --hard", "correct": false},
							},
							"explanation": "В main попадает только проверенный код.",
						},
					},
					"resources": []map[string]any{
						{
							"title": "GitHub Actions — документация на русском",
							"url":   "https://docs.github.com/ru/actions",
							"note":  "как написать свой первый конвейер",
						},
					},
				},
			},
			{
				Title:       "GitOps: репозиторий вместо ручных команд",
				Kind:        "text",
				Summary:     "Почему выкат делают через Git, а не через kubectl",
				DurationMin: 11,
				Content: map[string]any{
					"body": "## Зачем это нужно\n\n" +
						"Инженер зашёл на сервер и поправил настройку руками. Сервис заработал. " +
						"Через месяц сервер пересоздали — **правка пропала**, и никто не помнит, что там было.\n\n" +
						"GitOps убирает такие истории. Правило одно: **всё, что работает в проде, описано в Git.**\n\n" +
						"## Простыми словами\n\n" +
						"В репозитории лежит описание: какие сервисы запущены, какие версии, сколько копий.\n\n" +
						"В кластере работает программа-агент. Она постоянно сравнивает описание с реальностью " +
						"и приводит реальность к описанию.\n\n" +
						"Хотите выкатить новую версию — меняете строку в файле и делаете коммит. " +
						"Дальше агент всё сделает сам.\n\n" +
						"## Что это даёт\n\n" +
						"- **История.** Кто, когда и что менял — видно в Git.\n" +
						"- **Откат.** Вернуть прошлую версию — это `git revert`.\n" +
						"- **Ревью.** Изменение прода проходит проверку коллегой, как обычный код.\n" +
						"- **Порядок.** Ручные правки на серверах агент вернёт обратно.\n\n" +
						"Последний пункт кажется неудобным, пока не поймёшь: это защита. " +
						"**Никто не может втихую поменять прод.**\n\n" +
						"## Как это выглядит\n\n" +
						"Два репозитория:\n\n" +
						"1. **Код приложения** — тут разработка, тесты, сборка образа.\n" +
						"2. **Описание окружений** — тут манифесты: какая версия образа где запущена.\n\n" +
						"Конвейер собрал образ `api:1.4.3` и поменял версию во втором репозитории. " +
						"Агент увидел изменение и обновил приложение в кластере.\n\n" +
						"Популярные агенты: **Argo CD** и **Flux**.\n\n" +
						"Вот как выглядит сам «выкат». Меняется одна строка в манифесте:\n" +
						"\n" +
						"```diff\n" +
						" # deploy/api.yaml\n" +
						" spec:\n" +
						"   containers:\n" +
						"     - name: api\n" +
						"-      image: registry.example.com/api:1.4.2\n" +
						"+      image: registry.example.com/api:1.4.3\n" +
						"```\n" +
						"\n" +
						"Дальше — обычный коммит:\n" +
						"\n" +
						"```bash\n" +
						"$ git commit -am \"api 1.4.3\"\n" +
						"$ git push\n" +
						"```\n" +
						"\n" +
						"Больше ничего делать не нужно: **агент сам подтянет новую версию**.\n" +
						"\n" +
						"Агент постоянно сравнивает Git с кластером. Если прод поправили руками, приложение становится OutOfSync:\n" +
						"\n" +
						"```bash\n" +
						"$ argocd app get api\n" +
						"Name:           api\n" +
						"Project:        default\n" +
						"Sync Status:    OutOfSync from main (7d1f9a2)\n" +
						"Health Status:  Healthy\n" +
						"\n" +
						"GROUP  KIND        NAMESPACE  NAME  STATUS     HEALTH\n" +
						"apps   Deployment  prod       api   OutOfSync  Healthy\n" +
						"\n" +
						"$ argocd app diff api\n" +
						"===== /Deployment prod/api ======\n" +
						"27c27\n" +
						"<         image: registry.example.com/api:1.4.3\n" +
						"---\n" +
						">         image: registry.example.com/api:1.4.2\n" +
						"```\n" +
						"\n" +
						"`<` — как записано в Git (нужно 1.4.3), `>` — как сейчас в кластере (кто-то откатил на 1.4.2). Агент сам вернёт 1.4.3.\n" +
						"\n" +
						"## Разница на практике\n\n" +
						"| | Ручной выкат | GitOps |\n" +
						"|---|---|---|\n" +
						"| Как выкатывают | команда `kubectl apply` | коммит в репозиторий |\n" +
						"| Кто имеет доступ к проду | все инженеры | только агент |\n" +
						"| Как откатить | вспомнить прошлую версию | `git revert` |\n" +
						"| Видно ли, что менялось | иногда | всегда |\n\n" +
						"Выкат 1.4.3 оказался плохим. Откат в GitOps — это не паника, а один коммит:\n" +
						"\n" +
						"```bash\n" +
						"$ git log --oneline -3\n" +
						"7d1f9a2 (HEAD -> main) api 1.4.3\n" +
						"a1b2c3d api 1.4.2\n" +
						"5e8c0b1 tune limits\n" +
						"\n" +
						"$ git revert --no-edit 7d1f9a2\n" +
						"[main 3c4d5e6] Revert \"api 1.4.3\"\n" +
						" 1 file changed, 1 insertion(+), 1 deletion(-)\n" +
						"\n" +
						"$ git push\n" +
						"```\n" +
						"\n" +
						"`revert` создаёт новый коммит, который отменяет прошлый. Агент увидит его и вернёт прод на 1.4.2.\n" +
						"\n" +
						"## Частые ошибки новичка\n\n" +
						"- **Правят кластер руками при живом GitOps.** Агент вернёт как было — и это правильно.\n" +
						"- **Кладут секреты в репозиторий описаний.** Секреты подключают отдельно, через хранилище.\n" +
						"- **Один репозиторий на код и описание.** Работает, но сборка начинает запускать сама себя.\n\n" +
						"## Запомнить\n\n" +
						"1. Источник правды — Git, а не то, что кто-то набрал в терминале.\n" +
						"2. Выкат = коммит. Откат = revert.\n" +
						"3. Ручные правки на серверах агент откатит обратно.",
					"resources": []map[string]any{
						{
							"title": "OpenGitOps — принципы подхода",
							"url":   "https://opengitops.dev/",
							"note":  "четыре принципа, на которых работают Argo CD и Flux",
						},
						{
							"title": "Argo CD — документация",
							"url":   "https://argo-cd.readthedocs.io/",
							"note":  "самый распространённый агент: установка и первое приложение",
						},
						{
							"title": "Flux — GitOps-агент (CNCF)",
							"url":   "https://fluxcd.io/flux/",
							"note":  "второй по популярности агент после Argo CD, принципы GitOps",
						},
						{
							"title": "git revert — документация Git",
							"url":   "https://git-scm.com/docs/git-revert",
							"note":  "как откат оформляется отдельным коммитом, без переписывания истории",
						},
						{
							"title": "Argo CD — основные понятия",
							"url":   "https://argo-cd.readthedocs.io/en/stable/core_concepts/",
							"note":  "что такое sync и drift, как агент приводит кластер к описанию",
						},
					},
				},
			},
			{
				Title:       "Квиз: GitOps",
				Kind:        "quiz",
				Summary:     "Репозиторий как источник правды",
				DurationMin: 6,
				Content: map[string]any{
					"passScore": 70,
					"questions": []map[string]any{
						{
							"id":   "g1",
							"text": "Как выкатывают новую версию при GitOps?",
							"options": []map[string]any{
								{"id": "a", "text": "Меняют версию в файле и делают коммит, дальше агент применяет сам", "correct": true},
								{"id": "b", "text": "Заходят на сервер и выполняют команду вручную", "correct": false},
								{"id": "c", "text": "Пересобирают кластер с нуля", "correct": false},
							},
							"explanation": "Выкат = коммит. Ручные команды не нужны.",
						},
						{
							"id":   "g2",
							"text": "Как откатить неудачный выкат?",
							"options": []map[string]any{
								{"id": "a", "text": "Сделать git revert — агент вернёт прошлое состояние", "correct": true},
								{"id": "b", "text": "Восстановить кластер из резервной копии", "correct": false},
								{"id": "c", "text": "Никак, придётся выкатывать заново вручную", "correct": false},
							},
							"explanation": "История изменений прода лежит в Git.",
						},
						{
							"id":   "g3",
							"text": "Инженер поправил настройку в кластере руками. Что произойдёт?",
							"options": []map[string]any{
								{"id": "a", "text": "Агент вернёт как описано в репозитории", "correct": true},
								{"id": "b", "text": "Изменение попадёт в репозиторий автоматически", "correct": false},
								{"id": "c", "text": "Ничего, правка останется навсегда", "correct": false},
							},
							"explanation": "Это защита: никто не меняет прод втихую.",
						},
						{
							"id":       "g4",
							"text":     "Что даёт GitOps по сравнению с ручным выкатом?",
							"multiple": true,
							"options": []map[string]any{
								{"id": "a", "text": "Видно, кто и когда менял прод", "correct": true},
								{"id": "b", "text": "Изменения проходят ревью, как обычный код", "correct": true},
								{"id": "c", "text": "Откат становится обычным revert", "correct": true},
								{"id": "d", "text": "Приложение начинает работать быстрее", "correct": false},
							},
							"explanation": "Это про порядок и предсказуемость, а не про скорость работы приложения.",
						},
						{
							"id":     "g5",
							"review": true,
							"text":   "Повторение: почему в образе фиксируют версию вместо latest?",
							"options": []map[string]any{
								{"id": "a", "text": "Чтобы сборка была воспроизводимой и понятно, что именно работает", "correct": true},
								{"id": "b", "text": "Чтобы образ занимал меньше места", "correct": false},
								{"id": "c", "text": "Чтобы ускорить скачивание", "correct": false},
							},
							"explanation": "С latest сегодня и завтра разворачиваются разные образы.",
						},
					},
					"resources": []map[string]any{
						{
							"title": "OpenGitOps — принципы подхода",
							"url":   "https://opengitops.dev/",
							"note":  "четыре принципа простым языком",
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
						{
							"id":     "q5",
							"review": true,
							"text":   "Повторение: зачем в Dockerfile создают отдельного пользователя и пишут USER?",
							"options": []map[string]any{
								{"id": "a", "text": "Чтобы процесс не работал от root и взлом не дал сразу полных прав", "correct": true},
								{"id": "b", "text": "Чтобы образ занимал меньше места", "correct": false},
								{"id": "c", "text": "Без этого контейнер не запустится", "correct": false},
							},
							"explanation": "Наименьшие привилегии ограничивают ущерб.",
						},
					},
				},
			},
		},
	}
}
