package seed

func moduleGit() ModuleSeed {
	return ModuleSeed{
		Title:   "Git и командная работа",
		Summary: "Ветки, коммиты, ревью и разбор конфликтов — основа любой доставки",
		Lessons: []LessonSeed{
			{
				Title:       "Как работает Git",
				Kind:        "text",
				Summary:     "Коммиты, ветки, удалённый репозиторий и рабочий процесс команды",
				DurationMin: 16,
				Content: map[string]any{
					"body": "## Зачем инженеру DevOps знать Git\n\n" +
						"В Git лежит не только код приложения, но и всё, чем занимается DevOps: " +
						"Dockerfile, манифесты Kubernetes, конфигурация Terraform, пайплайны CI. " +
						"Любое изменение инфраструктуры проходит через ветку и ревью — " +
						"это и есть подход «инфраструктура как код».\n\n" +
						"## Три состояния файла\n\n" +
						"```\n" +
						"рабочий каталог  →  индекс (staging)  →  репозиторий\n" +
						"     правки            git add            git commit\n" +
						"```\n\n" +
						"`git status` всегда показывает, где что находится. Читайте его перед каждым коммитом — " +
						"это самая частая причина «я закоммитил не то».\n\n" +
						"## Ветки\n\n" +
						"Ветка — это подвижный указатель на коммит, а не копия файлов. Поэтому ветки дешёвые: " +
						"создавать их можно свободно.\n\n" +
						"```bash\n" +
						"git switch -c feature/healthcheck   # создать ветку и перейти в неё\n" +
						"git switch main                     # вернуться\n" +
						"git branch                          # список веток\n" +
						"```\n\n" +
						"## Типичный процесс команды\n\n" +
						"1. Обновить `main`: `git pull`.\n" +
						"2. Создать ветку под задачу.\n" +
						"3. Небольшие осмысленные коммиты.\n" +
						"4. Отправить ветку: `git push -u origin feature/...`.\n" +
						"5. Открыть pull request, пройти ревью и проверки CI.\n" +
						"6. Влить в `main` — оттуда запускается доставка.\n\n" +
						"## Сообщения коммитов\n\n" +
						"Хорошее сообщение отвечает на вопрос «что изменится в поведении», а не «какие строки поправлены».\n\n" +
						"```\n" +
						"feat: добавить healthcheck с проверкой базы\n" +
						"fix: не терять платежи при таймауте\n" +
						"chore: обновить зависимости\n" +
						"```\n\n" +
						"## Если что-то пошло не так\n\n" +
						"| Ситуация | Команда |\n" +
						"|---|---|\n" +
						"| Убрать файл из индекса | `git reset HEAD файл` |\n" +
						"| Отменить правки в файле | `git restore файл` |\n" +
						"| Отменить уже отправленный коммит | `git revert хеш` |\n" +
						"| Отложить правки | `git stash`, вернуть — `git stash pop` |\n\n" +
						"> `git revert` создаёт новый коммит-отмену и безопасен для общей ветки. " +
						"`git reset --hard` переписывает историю — в общей ветке так делать нельзя.",
				},
			},
			{
				Title:       "Тренажёр: ветка и коммит",
				Kind:        "terminal",
				Summary:     "Пройдите путь от правки до отправки ветки",
				DurationMin: 18,
				Content: map[string]any{
					"intro": "Вы в репозитории проекта. Проведите изменение через все стадии — от ветки до push.",
					"shell": "student@devops",
					"tasks": []map[string]any{
						{
							"id":       "g1",
							"prompt":   "Посмотрите состояние рабочего каталога",
							"expected": []string{"git status"},
							"hint":     "Команда из двух слов",
							"success":  "Видно изменённый файл и неотслеживаемый — это разные состояния.",
						},
						{
							"id":       "g2",
							"prompt":   "Посмотрите историю коммитов в компактном виде",
							"expected": []string{"git log --oneline"},
							"hint":     "К git log добавьте флаг --oneline",
							"success":  "Одна строка на коммит — так удобнее искать нужный.",
						},
						{
							"id":       "g3",
							"prompt":   "Создайте ветку feature/metrics и сразу перейдите в неё",
							"expected": []string{"git switch -c feature/metrics", "git checkout -b feature/metrics"},
							"hint":     "git switch -c имя-ветки",
							"success":  "Ветка создана, работа не мешает main.",
						},
						{
							"id":       "g4",
							"prompt":   "Посмотрите, что именно изменилось в файлах",
							"expected": []string{"git diff"},
							"hint":     "Одна короткая команда",
							"success":  "Перед добавлением в индекс всегда полезно перечитать diff.",
						},
						{
							"id":       "g5",
							"prompt":   "Добавьте все изменения в индекс",
							"expected": []string{"git add .", "git add -A", "git add --all"},
							"hint":     "git add и точка",
							"success":  "Файлы в индексе и готовы к коммиту.",
						},
						{
							"id":       "g6",
							"prompt":   "Сделайте коммит с сообщением: feat: add metrics endpoint",
							"expected": []string{`git commit -m "feat: add metrics endpoint"`, "git commit -m 'feat: add metrics endpoint'"},
							"hint":     "git commit -m \"текст\"",
							"success":  "Коммит записан в историю ветки.",
						},
						{
							"id":       "g7",
							"prompt":   "Отправьте ветку feature/metrics в origin, связав её с удалённой",
							"expected": []string{"git push -u origin feature/metrics", "git push --set-upstream origin feature/metrics"},
							"hint":     "git push -u origin имя-ветки",
							"success":  "Ветка на сервере — можно открывать pull request.",
						},
					},
				},
			},
			{
				Title:       "Проверка: Git",
				Kind:        "quiz",
				Summary:     "Состояния файлов, ветки и отмена изменений",
				DurationMin: 10,
				Content: map[string]any{
					"passScore": 70,
					"questions": []map[string]any{
						{
							"id":   "q1",
							"text": "Что делает git add?",
							"options": []map[string]any{
								{"id": "a", "text": "Переносит изменения из рабочего каталога в индекс", "correct": true},
								{"id": "b", "text": "Отправляет изменения на сервер", "correct": false},
								{"id": "c", "text": "Создаёт новую ветку", "correct": false},
							},
							"explanation": "Индекс — промежуточная область: туда попадает то, что войдёт в следующий коммит.",
						},
						{
							"id":   "q2",
							"text": "Чем git revert отличается от git reset --hard?",
							"options": []map[string]any{
								{"id": "a", "text": "revert создаёт новый коммит-отмену, reset --hard переписывает историю", "correct": true},
								{"id": "b", "text": "Ничем, это синонимы", "correct": false},
								{"id": "c", "text": "revert работает только с ветками", "correct": false},
							},
							"explanation": "В общей ветке безопасен только revert: история остаётся неизменной для всех.",
						},
						{
							"id":   "q3",
							"text": "Почему ветки в Git считаются дешёвыми?",
							"options": []map[string]any{
								{"id": "a", "text": "Ветка — это указатель на коммит, а не копия файлов", "correct": true},
								{"id": "b", "text": "Git сжимает файлы при создании ветки", "correct": false},
								{"id": "c", "text": "Ветки хранятся только на сервере", "correct": false},
							},
							"explanation": "Создание ветки — это запись одного указателя, поэтому оно мгновенное.",
						},
						{
							"id":       "q4",
							"text":     "Что стоит хранить в Git у DevOps-инженера?",
							"multiple": true,
							"options": []map[string]any{
								{"id": "a", "text": "Dockerfile и манифесты Kubernetes", "correct": true},
								{"id": "b", "text": "Конфигурацию Terraform", "correct": true},
								{"id": "c", "text": "Файлы пайплайнов CI", "correct": true},
								{"id": "d", "text": "Пароли от прода в открытом виде", "correct": false},
							},
							"explanation": "Всё описание инфраструктуры — в репозитории, секреты — в хранилище секретов.",
						},
						{
							"id":   "q5",
							"text": "Вы случайно закоммитили лишний файл, но ещё не сделали push. Что проще всего?",
							"options": []map[string]any{
								{"id": "a", "text": "Поправить последний коммит: git reset HEAD~1 или git commit --amend", "correct": true},
								{"id": "b", "text": "Удалить репозиторий и клонировать заново", "correct": false},
								{"id": "c", "text": "Ничего сделать нельзя", "correct": false},
							},
							"explanation": "Пока коммит локальный, историю можно спокойно переписать.",
						},
					},
				},
			},
		},
	}
}
