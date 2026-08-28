package seed

func moduleCICD() ModuleSeed {
	return ModuleSeed{
		Title:   "CI/CD",
		Summary: "Конвейер сборки и доставки: от коммита до прода",
		Lessons: []LessonSeed{
			{
				Title:       "Как устроен конвейер",
				Kind:        "text",
				Summary:     "CI, CD, окружения и стратегии выката",
				DurationMin: 16,
				Content: map[string]any{
					"body": "## CI и CD — разные вещи\n\n" +
						"**CI (Continuous Integration)** — каждый коммит автоматически собирается и " +
						"проверяется тестами. Задача: узнать о поломке через минуты, а не на релизе.\n\n" +
						"**CD** расшифровывают двумя способами:\n\n" +
						"- *Continuous Delivery* — сборка всегда готова к выкату, но кнопку нажимает человек.\n" +
						"- *Continuous Deployment* — прошедшая проверки сборка едет в прод сама.\n\n" +
						"## Типичные этапы\n\n" +
						"```\n" +
						"commit → lint → tests → build image → push registry → deploy stage → smoke → deploy prod\n" +
						"```\n\n" +
						"Правило: чем дешевле проверка, тем раньше она стоит. Линтер отрабатывает за секунды — " +
						"нет смысла собирать образ, если код не проходит форматирование.\n\n" +
						"## Окружения\n\n" +
						"| Окружение | Зачем |\n" +
						"|---|---|\n" +
						"| dev | ветки разработчиков, всё ломается — и это нормально |\n" +
						"| stage | копия прода для приёмки |\n" +
						"| prod | боевое, сюда едет только проверенное |\n\n" +
						"## Стратегии выката\n\n" +
						"- **Rolling** — поды обновляются по очереди, старые остаются рабочими до готовности новых.\n" +
						"- **Blue/Green** — рядом с текущей версией поднимают новую и переключают трафик целиком.\n" +
						"- **Canary** — на новую версию направляют 5–10% трафика, смотрят метрики, потом остальное.\n\n" +
						"## Секреты\n\n" +
						"Пароли и токены не хранят в репозитории. Их держат в секретах CI " +
						"(GitHub Secrets, GitLab CI Variables) или в отдельном хранилище вроде Vault, " +
						"а в конвейер отдают переменными окружения.\n\n" +
						"> Хороший конвейер обязательно умеет откатываться. Если нет плана отката — " +
						"это не доставка, а лотерея." +
						"\n\n## Практики, которые ждут от инженера сегодня\n\n" +
						"- Доступ к облаку и реестру — через OIDC, статические ключи в секретах считаются устаревшими.\n" +
						"- Образы и артефакты подписывают, рядом кладут SBOM: без этого выкат в зрелых компаниях не пройдёт.\n" +
						"- Пайплайн кэширует зависимости и слои образа — сборка в десять минут убивает скорость команды.\n" +
						"- Внешние действия закрепляют по хешу коммита: тег можно переписать, хеш — нет.\n" +
						"- Выкат в прод отделяют окружением с ручным подтверждением или переносят в GitOps-репозиторий.\n" +
						"- Тесты гоняют параллельно матрицей версий, а нестабильные тесты чинят, а не перезапускают.",
					"resources": []map[string]any{
						{"title": "GitHub Actions: документация", "url": "https://docs.github.com/actions", "note": "Синтаксис workflow, матрицы, кэш и окружения"},
						{"title": "GitLab CI/CD", "url": "https://docs.gitlab.com/ee/ci/", "note": "Альтернативный синтаксис, часто встречается в компаниях"},
						{"title": "OIDC вместо долгих секретов", "url": "https://docs.github.com/actions/deployment/security-hardening-your-deployments/about-security-hardening-with-openid-connect", "note": "Как выдавать облачный доступ пайплайну без хранения ключей"},
						{"title": "Continuous Delivery: практики", "url": "https://continuousdelivery.com/", "note": "Первоисточник подхода от Джеза Хамбла"},
					},
				},
			},
			{
				Title:       "GitOps: репозиторий как источник правды",
				Kind:        "text",
				Summary:     "Pull-модель доставки, ArgoCD и Flux, отличия от push-выката",
				DurationMin: 16,
				Content: map[string]any{
					"body": "## Push и pull\n\n" +
						"Обычный пайплайн **push-модель**: CI сам ходит в кластер и применяет манифесты. " +
						"Для этого пайплайну нужны права на прод, а фактическое состояние кластера " +
						"никто не сверяет с репозиторием.\n\n" +
						"**GitOps** переворачивает схему: в кластере живёт агент (ArgoCD или Flux), " +
						"который сам следит за репозиторием и приводит кластер к описанному состоянию.\n\n" +
						"```\n" +
						"push:  CI ──kubectl apply──► кластер\n" +
						"pull:  CI ──коммит манифеста──► Git ◄──следит── агент в кластере\n" +
						"```\n\n" +
						"## Что это даёт\n\n" +
						"- У CI больше нет доступа к кластеру — только право коммита. Утечка токена CI перестаёт быть катастрофой.\n" +
						"- Состояние кластера всегда сверяется с Git: ручные правки видны как **drift** и откатываются.\n" +
						"- Откат — это `git revert`, а не отдельная процедура.\n" +
						"- История выкатов совпадает с историей репозитория: видно, кто и что выкатил.\n\n" +
						"## Два репозитория\n\n" +
						"Обычно разделяют код приложения и описание развёртывания:\n\n" +
						"| Репозиторий | Что внутри | Кто меняет |\n" +
						"|---|---|---|\n" +
						"| `api` | исходники, тесты, Dockerfile | разработчики |\n" +
						"| `api-deploy` | манифесты, values, версии образов | CI и дежурные |\n\n" +
						"CI собирает образ, проставляет тег и коммитит новую версию в репозиторий развёртывания — " +
						"дальше агент выкатывает сам.\n\n" +
						"## Пример описания приложения в ArgoCD\n\n" +
						"```yaml\n" +
						"apiVersion: argoproj.io/v1alpha1\n" +
						"kind: Application\n" +
						"metadata:\n" +
						"  name: api\n" +
						"  namespace: argocd\n" +
						"spec:\n" +
						"  project: default\n" +
						"  source:\n" +
						"    repoURL: https://github.com/team/api-deploy.git\n" +
						"    targetRevision: main\n" +
						"    path: overlays/production\n" +
						"  destination:\n" +
						"    server: https://kubernetes.default.svc\n" +
						"    namespace: production\n" +
						"  syncPolicy:\n" +
						"    automated:\n" +
						"      prune: true      # удалять то, чего больше нет в Git\n" +
						"      selfHeal: true   # возвращать ручные правки к состоянию из Git\n" +
						"```\n\n" +
						"> `selfHeal: true` означает, что правка через `kubectl edit` проживёт до следующей синхронизации. " +
						"Это и есть цель: единственный способ изменить прод — коммит.\n\n" +
						"## Где GitOps не нужен\n\n" +
						"Для одного сервера с docker compose агент только усложнит жизнь. " +
						"Подход окупается, когда окружений несколько, людей больше пяти, а выкаты идут ежедневно.\n\n" +
						"## Прогрессивная доставка\n\n" +
						"Поверх GitOps ставят контроллеры вроде Argo Rollouts или Flagger: они катят новую версию " +
						"постепенно, смотрят на метрики ошибок и задержки и сами откатывают выкат, " +
						"если показатели ухудшились. Так canary-выкат становится автоматическим, без дежурного у графиков.",
					"resources": []map[string]any{
						{"title": "Принципы OpenGitOps", "url": "https://opengitops.dev/", "note": "Четыре принципа подхода в одной странице"},
						{"title": "Argo CD: документация", "url": "https://argo-cd.readthedocs.io/en/stable/", "note": "Application, проекты, синхронизация и права"},
						{"title": "Flux", "url": "https://fluxcd.io/flux/", "note": "Альтернатива Argo CD, ближе к Kubernetes-нативному стилю"},
						{"title": "Argo Rollouts", "url": "https://argo-rollouts.readthedocs.io/en/stable/", "note": "Canary и blue/green с автоматическим откатом по метрикам"},
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
