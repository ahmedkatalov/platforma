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
						"это не доставка, а лотерея.",
				},
			},
			{
				Title:       "Практика: GitHub Actions",
				Kind:        "code",
				Summary:     "Соберите рабочий workflow для сборки и тестов",
				DurationMin: 25,
				Content: map[string]any{
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
				Title:       "Проверка: CI/CD",
				Kind:        "quiz",
				Summary:     "Этапы конвейера, окружения и выкаты",
				DurationMin: 8,
				Content: map[string]any{
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
