package seed

func moduleGit() ModuleSeed {
	return ModuleSeed{
		Title:   "Git и командная работа",
		Summary: "Сохранение версий, ветки и совместная работа — с нуля",
		Lessons: []LessonSeed{
			{
				Title:       "Git: зачем он нужен",
				Kind:        "text",
				Summary:     "Три состояния файла и путь от правки до коммита",
				DurationMin: 12,
				Content: map[string]any{
					"body": "## 🎯 Чему вы научитесь\n\n" +
						"- Провести изменения через git add и git commit и отправить их командой git push\n" +
						"- Создать ветку под задачу и отправить её на сервер для pull request\n" +
						"- Собрать .gitignore, чтобы секреты и лишние файлы не попадали в репозиторий\n" +
						"- Разрешить конфликт слияния и отменить правки через restore и revert\n" +
						"- Объяснить три состояния файла и разницу между revert и reset --hard\n" +
						"\n" +
						"## 📋 Что нужно знать заранее\n\n" +
						"- Работать в терминале Linux: запускать команды и переходить по папкам\n" +
						"- Понимать, что такое файлы, папки и путь к файлу\n" +
						"- Представлять, чем в целом занимается DevOps-инженер\n" +
						"\n" +
						"## Зачем это нужно\n\n" +
						"Знакомая ситуация: `отчёт.docx`, `отчёт_финал.docx`, `отчёт_финал2_точно.docx`. " +
						"Через неделю непонятно, где последняя версия.\n\n" +
						"Git решает это. Он **хранит все версии проекта** и помнит, кто, когда и что изменил.\n\n" +
						"Для инженера DevOps Git — рабочий инструмент каждый день. В нём лежит не только код, " +
						"но и настройки серверов, описания контейнеров, конвейеры сборки.\n\n" +
						"## Простыми словами\n\n" +
						"**Репозиторий** — папка проекта под присмотром Git.\n\n" +
						"**Коммит** — сохранённая версия. Как точка сохранения в игре: всегда можно вернуться.\n\n" +
						"**Ветка** — отдельная линия работы. Вы правите в своей ветке и не мешаете другим.\n\n" +
						"## Три состояния файла\n\n" +
						"Это главное, что нужно понять про Git:\n\n" +
						"```\n" +
						"1. Рабочая папка  →  2. Индекс  →  3. Репозиторий\n" +
						"   (вы правите)      (git add)     (git commit)\n" +
						"```\n\n" +
						"1. **Рабочая папка** — файлы, которые вы редактируете прямо сейчас.\n" +
						"2. **Индекс** — то, что вы отобрали для следующего сохранения.\n" +
						"3. **Репозиторий** — история сохранённых версий.\n\n" +
						"Зачем индекс? Чтобы сохранить не всё подряд, а только нужное. " +
						"Например, поправили два файла, а сохранить хотите пока один.\n\n" +
						"## Главная команда\n\n" +
						"```bash\n" +
						"git status\n" +
						"```\n\n" +
						"Она показывает, что изменилось и в каком состоянии находится. " +
						"**Набирайте её постоянно** — это бесплатно и спасает от ошибок.\n\n" +
						"## Обычный порядок работы\n\n" +
						"```bash\n" +
						"git status                  # что изменилось\n" +
						"git diff                    # посмотреть сами изменения\n" +
						"git add .                   # отобрать всё в индекс\n" +
						"git commit -m \"понятный текст\"   # сохранить версию\n" +
						"git push                    # отправить на сервер\n" +
						"```\n\n" +
						"А вот как весь цикл выглядит вживую — с реальным выводом Git:\n" +
						"\n" +
						"```bash\n" +
						"$ git status\n" +
						"On branch main\n" +
						"Changes not staged for commit:\n" +
						"  modified:   app.py\n" +
						"\n" +
						"Untracked files:\n" +
						"  config.yaml\n" +
						"\n" +
						"$ git add .\n" +
						"\n" +
						"$ git commit -m \"feat: добавить проверку доступности базы\"\n" +
						"[main a1b2c3d] feat: добавить проверку доступности базы\n" +
						" 2 files changed, 15 insertions(+)\n" +
						"\n" +
						"$ git status\n" +
						"On branch main\n" +
						"nothing to commit, working tree clean\n" +
						"```\n" +
						"\n" +
						"Сравните два `git status`: до коммита файлы «висят», после — «working tree clean».\n" +
						"\n" +
						"## Как писать сообщение коммита\n\n" +
						"Плохо: `фикс`, `правки`, `ещё раз`.\n\n" +
						"Хорошо: коротко о том, что изменится в поведении.\n\n" +
						"```\n" +
						"feat: добавить проверку доступности базы\n" +
						"fix: не терять платежи при таймауте\n" +
						"docs: описать порядок выката\n" +
						"```\n\n" +
						"Через полгода такое сообщение **спасёт вас же**.\n\n" +
						"## Частые ошибки новичка\n\n" +
						"- **Коммитят всё подряд командой `git add .` не глядя.** Сначала `git status`.\n" +
						"- **Сообщения вроде «правки».** Никто потом не поймёт, что там было.\n" +
						"- **Боятся сделать лишний коммит.** Не бойтесь: много маленьких коммитов лучше одного огромного.\n\n" +
						"## Как выглядит реальная сессия\n" +
						"\n" +
						"Обычный цикл, но с типичной заминкой — push отклонён:\n" +
						"\n" +
						"```bash\n" +
						"$ git status\n" +
						"On branch main\n" +
						"Changes not staged for commit:\n" +
						"  modified:   app.py\n" +
						"\n" +
						"$ git add app.py\n" +
						"$ git commit -m \"fix: обработка пустого заказа\"\n" +
						"[main a1b2c3d] fix: обработка пустого заказа\n" +
						" 1 file changed, 4 insertions(+)\n" +
						"\n" +
						"$ git push\n" +
						" ! [rejected]        main -> main (fetch first)\n" +
						"error: failed to push some refs\n" +
						"hint: Updates were rejected because the remote contains work that\n" +
						"hint: you do not have locally.\n" +
						"\n" +
						"$ git pull --rebase\n" +
						"$ git push\n" +
						"   9f8e7d6..a1b2c3d  main -> main\n" +
						"```\n" +
						"\n" +
						"`rejected ... fetch first` значит: на сервере есть чужие коммиты, которых нет у вас.\n" +
						"`git pull --rebase` подтягивает их и кладёт ваш коммит сверху, затем push проходит.\n" +
						"\n" +
						"## ❓ Проверьте себя\n\n" +
						"`git push` вернул `rejected ... fetch first`. Что это значит и что сделать?\n\n" +
						"*Ответ: на сервере есть чужие коммиты. Сначала `git pull --rebase`, потом снова `git push`.*\n\n" +
						"## Запомнить\n\n" +
						"1. Правки → `git add` → `git commit` → `git push`.\n" +
						"2. `git status` — самая полезная команда, набирайте её чаще.\n" +
						"3. Сообщение коммита пишут для будущего себя.",
					"resources": []map[string]any{
						{
							"title": "Pro Git — книга на русском",
							"url":   "https://git-scm.com/book/ru/v2",
							"note":  "первые три главы закрывают почти всю ежедневную работу",
						},
						{
							"title": "Learn Git Branching — тренажёр в браузере",
							"url":   "https://learngitbranching.js.org/?locale=ru_RU",
							"note":  "наглядно показывает, что происходит с ветками",
						},
						{
							"title": "Atlassian — учебник по Git",
							"url":   "https://www.atlassian.com/git/tutorials",
							"note":  "пошаговые статьи с картинками: от первого коммита до веток",
						},
						{
							"title": "git-scm — описание git commit",
							"url":   "https://git-scm.com/docs/git-commit",
							"note":  "официальный справочник по коммитам и их флагам",
						},
					},
				},
			},
			{
				Title:       "Квиз: основы Git",
				Kind:        "quiz",
				Summary:     "Три состояния файла и первые команды",
				DurationMin: 6,
				Content: map[string]any{
					"passScore": 70,
					"questions": []map[string]any{
						{
							"id":   "b1",
							"text": "В каком порядке изменения проходят путь до истории?",
							"options": []map[string]any{
								{"id": "a", "text": "Рабочая папка → индекс → репозиторий", "correct": true},
								{"id": "b", "text": "Индекс → рабочая папка → репозиторий", "correct": false},
								{"id": "c", "text": "Репозиторий → индекс → рабочая папка", "correct": false},
							},
							"explanation": "git add переносит в индекс, git commit — в репозиторий.",
						},
						{
							"id":   "b2",
							"text": "Зачем нужен индекс?",
							"options": []map[string]any{
								{"id": "a", "text": "Чтобы отобрать в коммит только часть изменений", "correct": true},
								{"id": "b", "text": "Чтобы ускорить работу Git", "correct": false},
								{"id": "c", "text": "Чтобы хранить резервную копию файлов", "correct": false},
							},
							"explanation": "Поправили пять файлов — закоммитить можно только два.",
						},
						{
							"id":   "b3",
							"text": "Какую команду стоит набирать чаще всего?",
							"options": []map[string]any{
								{"id": "a", "text": "git status — она показывает, что происходит, и ничего не меняет", "correct": true},
								{"id": "b", "text": "git push", "correct": false},
								{"id": "c", "text": "git reset --hard", "correct": false},
							},
							"explanation": "status безопасен: только смотрит.",
						},
						{
							"id":   "b4",
							"text": "Какое сообщение коммита лучше?",
							"options": []map[string]any{
								{"id": "a", "text": "fix: не терять платежи при таймауте", "correct": true},
								{"id": "b", "text": "правки", "correct": false},
								{"id": "c", "text": "ещё раз", "correct": false},
							},
							"explanation": "Сообщение пишут для того, кто прочитает его через полгода.",
						},
						{
							"id":     "b5",
							"review": true,
							"text":   "Повторение: какой командой найти строки со словом ERROR в файле лога?",
							"options": []map[string]any{
								{"id": "a", "text": "grep ERROR app.log", "correct": true},
								{"id": "b", "text": "cat ERROR app.log", "correct": false},
								{"id": "c", "text": "ls ERROR app.log", "correct": false},
							},
							"explanation": "grep ищет строки по слову внутри файла.",
						},
					},
					"resources": []map[string]any{
						{
							"title": "Pro Git — глава про основы",
							"url":   "https://git-scm.com/book/ru/v2",
							"note":  "первые главы на русском",
						},
					},
				},
			},
			{
				Title:       "Ветки и работа в команде",
				Kind:        "text",
				Summary:     "Зачем нужны ветки, что такое pull request и как откатить ошибку",
				DurationMin: 12,
				Content: map[string]any{
					"body": "## Зачем это нужно\n\n" +
						"Над проектом работает пять человек. Если все правят одни и те же файлы напрямую, " +
						"начинается хаос.\n\n" +
						"Ветки решают это: у каждого своя линия работы, а объединяют их осознанно.\n\n" +
						"## Что такое ветка на самом деле\n\n" +
						"Ветка — это не копия файлов. Это просто **указатель на коммит**.\n\n" +
						"Представьте закладку в книге: она отмечает страницу, но не копирует саму книгу.\n" +
						"\n" +
						"Ветка — такая же закладка на коммит. Перейти в другую ветку — просто переставить закладку.\n" +
						"\n" +
						"Поэтому создать ветку — операция мгновенная. Не экономьте на ветках: " +
						"новая задача — новая ветка.\n\n" +
						"```bash\n" +
						"git switch -c feature/health-check   # создать ветку и перейти в неё\n" +
						"git switch main                      # вернуться в основную\n" +
						"git branch                           # список веток\n" +
						"```\n\n" +
						"Главная ветка обычно называется `main`. В неё попадает **только проверенный код**.\n\n" +
						"## Как работает команда\n\n" +
						"1. Обновляете `main`: `git pull`.\n" +
						"2. Создаёте ветку под задачу.\n" +
						"3. Делаете коммиты.\n" +
						"4. Отправляете ветку: `git push -u origin имя-ветки`.\n" +
						"5. Открываете **pull request** — просьбу влить вашу работу в `main`.\n" +
						"6. Коллеги смотрят код, автоматика прогоняет тесты.\n" +
						"7. После одобрения ветку вливают в `main`.\n\n" +
						"Из `main` обычно и происходит выкат на серверы. Поэтому туда попадает только " +
						"проверенное.\n\n" +
						"## Если что-то пошло не так\n\n" +
						"| Ситуация | Что делать |\n" +
						"|---|---|\n" +
						"| Отменить правки в файле | `git restore файл` |\n" +
						"| Убрать файл из индекса | `git restore --staged файл` |\n" +
						"| Отменить уже отправленный коммит | `git revert хеш` |\n" +
						"| Отложить правки и вернуться позже | `git stash`, потом `git stash pop` |\n\n" +
						"Важное правило: **в общей ветке используйте `revert`, а не `reset --hard`.**\n\n" +
						"`revert` создаёт новый коммит, который отменяет старый. История сохраняется, " +
						"у всех всё продолжает работать.\n\n" +
						"`reset --hard` **стирает историю**. Если ветка общая, у коллег начнутся конфликты.\n\n" +
						"## Конфликт слияния — это не страшно\n\n" +
						"Конфликт возникает, когда двое поправили одну строку. Git не знает, чей вариант верный, " +
						"и помечает место так:\n\n" +
						"```\n" +
						"<<<<<<< HEAD\n" +
						"ваш вариант\n" +
						"=======\n" +
						"вариант коллеги\n" +
						">>>>>>> main\n" +
						"```\n\n" +
						"Вы открываете файл, оставляете правильный вариант, удаляете строки с `<<<`, `===`, `>>>`, " +
						"затем `git add` и `git commit`. Всё.\n\n" +
						"## Частые ошибки новичка\n\n" +
						"- **Работают прямо в `main`.** Заведите ветку — это одна команда.\n" +
						"- **Делают `git push --force` в общую ветку.** Так можно стереть работу коллег.\n" +
						"- **Копят изменения неделю в одной ветке.** Чем крупнее ветка, тем больнее конфликты.\n\n" +
						"## Разбор конфликта на практике\n" +
						"\n" +
						"Слили ветку, а Git не смог объединить одну строку:\n" +
						"\n" +
						"```bash\n" +
						"$ git merge feature/login\n" +
						"Auto-merging app.py\n" +
						"CONFLICT (content): Merge conflict in app.py\n" +
						"Automatic merge failed; fix conflicts and then commit the result.\n" +
						"```\n" +
						"\n" +
						"Открываем app.py — Git пометил спорное место:\n" +
						"\n" +
						"```\n" +
						"<<<<<<< HEAD\n" +
						"timeout = 30\n" +
						"=======\n" +
						"timeout = 60\n" +
						">>>>>>> feature/login\n" +
						"```\n" +
						"\n" +
						"Оставляем нужную строку, удаляем маркеры, затем:\n" +
						"\n" +
						"```bash\n" +
						"$ git add app.py\n" +
						"$ git commit\n" +
						"[main 7c2e5f9] Merge feature/login\n" +
						"```\n" +
						"\n" +
						"Конфликт — это не поломка, а вопрос Git: чей вариант оставить.\n" +
						"\n" +
						"## ❓ Проверьте себя\n\n" +
						"Сделали `git merge`, и в файле появились строки `<<<<<<<`, `=======`, `>>>>>>>`. Что теперь делать?\n\n" +
						"*Ответ: оставить правильный вариант, стереть строки-маркеры, затем `git add файл` и `git commit`. Так Git понимает, чей код вы выбрали, и завершает слияние.*\n\n" +
						"## Запомнить\n\n" +
						"1. Новая задача — новая ветка.\n" +
						"2. В `main` попадает только проверенный код, через pull request.\n" +
						"3. Отменять в общей ветке — только через `revert`.",
					"resources": []map[string]any{
						{
							"title": "GitHub Flow — простой рабочий процесс",
							"url":   "https://docs.github.com/en/get-started/using-github/github-flow",
							"note":  "тот самый порядок: ветка, pull request, проверки, слияние",
						},
						{
							"title": "Oh Shit, Git!?! — выход из типичных ситуаций",
							"url":   "https://ohshitgit.com/ru",
							"note":  "короткие рецепты на случай «я всё сломал»",
						},
						{
							"title": "Atlassian — конфликты слияния",
							"url":   "https://www.atlassian.com/git/tutorials/using-branches/merge-conflicts",
							"note":  "разбор, откуда берётся конфликт и как его закрыть",
						},
						{
							"title": "GitHub — разрешение конфликта в командной строке",
							"url":   "https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/addressing-merge-conflicts/resolving-a-merge-conflict-using-the-command-line",
							"note":  "официальная инструкция по шагам",
						},
					},
				},
			},
			{
				Title:       "Тренажёр: от правки до push",
				Kind:        "terminal",
				Summary:     "Пройдите весь путь изменения своими руками",
				DurationMin: 18,
				Content: map[string]any{
					"intro": "Вы в репозитории проекта, в файлах есть изменения. Проведите их через все шаги — от осмотра до отправки на сервер.",
					"shell": "student@devops",
					"tasks": []map[string]any{
						{
							"id":       "g1",
							"prompt":   "Посмотрите, что изменилось в проекте",
							"expected": []string{"git status"},
							"hint":     "git и слово status",
							"success":  "Видно изменённый файл и новый, который Git ещё не отслеживает.",
						},
						{
							"id":       "g2",
							"prompt":   "Посмотрите историю коммитов в коротком виде",
							"expected": []string{"git log --oneline"},
							"hint":     "git log с флагом --oneline",
							"success":  "Одна строка на коммит — так удобнее искать нужный.",
						},
						{
							"id":       "g3",
							"prompt":   "Создайте ветку feature/metrics и перейдите в неё",
							"expected": []string{"git switch -c feature/metrics", "git checkout -b feature/metrics"},
							"hint":     "git switch -c и название ветки",
							"success":  "Теперь работа идёт в отдельной ветке и не мешает main.",
						},
						{
							"id":       "g4",
							"prompt":   "Посмотрите сами изменения в файлах",
							"expected": []string{"git diff"},
							"hint":     "Одно короткое слово после git",
							"success":  "Перед коммитом всегда полезно перечитать diff.",
						},
						{
							"id":       "g5",
							"prompt":   "Добавьте все изменения в индекс",
							"expected": []string{"git add .", "git add -A", "git add --all"},
							"hint":     "git add и точка",
							"success":  "Файлы отобраны для коммита.",
						},
						{
							"id":       "g6",
							"prompt":   "Сохраните версию с сообщением: feat: add metrics endpoint",
							"expected": []string{`git commit -m "feat: add metrics endpoint"`, "git commit -m 'feat: add metrics endpoint'"},
							"hint":     "git commit -m и текст в кавычках",
							"success":  "Коммит создан. Теперь он в истории ветки.",
						},
						{
							"id":       "g7",
							"prompt":   "Отправьте ветку на сервер и свяжите её с удалённой",
							"expected": []string{"git push -u origin feature/metrics", "git push --set-upstream origin feature/metrics"},
							"hint":     "git push -u origin и имя ветки",
							"success":  "Ветка на сервере — можно открывать pull request.",
						},
					},
					"resources": []map[string]any{
						{
							"title": "Справочник команд Git",
							"url":   "https://git-scm.com/docs",
							"note":  "описание каждой команды с примерами",
						},
					},
				},
			},
			{
				Title:       "Практика: .gitignore для проекта",
				Kind:        "code",
				Summary:     "Закройте от Git всё, что не должно попасть в репозиторий",
				DurationMin: 15,
				Content: map[string]any{
					"language": "gitignore",
					"task": "Соберите .gitignore для проекта. В него должны попасть:\n\n" +
						"1. файл с секретами `.env`;\n" +
						"2. приватные ключи по шаблону `*.key`;\n" +
						"3. папка зависимостей `node_modules/`;\n" +
						"4. папка сборки `dist/`;\n" +
						"5. системный файл `.DS_Store`.\n\n" +
						"Файл `.env.example` наоборот должен попадать в репозиторий — добавьте для него " +
						"исключение строкой `!.env.example`.",
					"starter": "# Что не должно попасть в репозиторий\n" +
						"\n" +
						"# секреты\n" +
						"\n" +
						"# зависимости и сборка\n" +
						"\n" +
						"# системные файлы\n",
					"hint": "Восклицательный знак в начале строки отменяет игнорирование: !.env.example",
					"checks": []map[string]any{
						{"type": "regex", "value": "(?m)^\\.env$", "message": "Файл .env игнорируется"},
						{"type": "regex", "value": "(?m)^\\*\\.key$", "message": "Приватные ключи игнорируются"},
						{"type": "regex", "value": "(?m)^node_modules/?$", "message": "Папка зависимостей игнорируется"},
						{"type": "regex", "value": "(?m)^dist/?$", "message": "Папка сборки игнорируется"},
						{"type": "regex", "value": "(?m)^\\.DS_Store$", "message": "Системный файл игнорируется"},
						{"type": "regex", "value": "(?m)^!\\.env\\.example$", "message": "Шаблон .env.example не игнорируется"},
					},
					"solution": "# Что не должно попасть в репозиторий\n" +
						"\n" +
						"# секреты\n" +
						".env\n" +
						"*.key\n" +
						"!.env.example\n" +
						"\n" +
						"# зависимости и сборка\n" +
						"node_modules/\n" +
						"dist/\n" +
						"\n" +
						"# системные файлы\n" +
						".DS_Store\n",
					"resources": []map[string]any{
						{
							"title": "Документация по .gitignore",
							"url":   "https://git-scm.com/docs/gitignore",
							"note":  "правила шаблонов: звёздочки, слэши и исключения",
						},
						{
							"title": "Готовые .gitignore для языков и фреймворков",
							"url":   "https://github.com/github/gitignore",
							"note":  "не пишите с нуля — возьмите готовый и дополните",
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
								{"id": "a", "text": "Переносит изменения в индекс — отбирает их для коммита", "correct": true},
								{"id": "b", "text": "Отправляет изменения на сервер", "correct": false},
								{"id": "c", "text": "Создаёт новую ветку", "correct": false},
							},
							"explanation": "На сервер отправляет git push, а ветку создаёт git switch -c.",
						},
						{
							"id":   "q2",
							"text": "Какая команда покажет, что изменилось в проекте?",
							"options": []map[string]any{
								{"id": "a", "text": "git status", "correct": true},
								{"id": "b", "text": "git push", "correct": false},
								{"id": "c", "text": "git clone", "correct": false},
							},
							"explanation": "git status безопасен и ничего не меняет — набирайте его чаще.",
						},
						{
							"id":   "q3",
							"text": "Зачем создавать отдельную ветку под задачу?",
							"options": []map[string]any{
								{"id": "a", "text": "Чтобы работа не мешала остальным до того, как её проверят", "correct": true},
								{"id": "b", "text": "Чтобы код работал быстрее", "correct": false},
								{"id": "c", "text": "Чтобы файлы занимали меньше места", "correct": false},
							},
							"explanation": "В main попадает только проверенный код — через pull request.",
						},
						{
							"id":   "q4",
							"text": "Нужно отменить коммит, который уже отправлен в общую ветку. Что выбрать?",
							"options": []map[string]any{
								{"id": "a", "text": "git revert — создаст коммит-отмену и сохранит историю", "correct": true},
								{"id": "b", "text": "git reset --hard и push --force", "correct": false},
								{"id": "c", "text": "Удалить репозиторий и склонировать заново", "correct": false},
							},
							"explanation": "reset --hard переписывает историю и ломает работу коллегам.",
						},
						{
							"id":       "q5",
							"text":     "Что из этого хранят в Git инженеры DevOps?",
							"multiple": true,
							"options": []map[string]any{
								{"id": "a", "text": "Описания контейнеров (Dockerfile)", "correct": true},
								{"id": "b", "text": "Файлы конвейеров сборки", "correct": true},
								{"id": "c", "text": "Описания серверов и инфраструктуры", "correct": true},
								{"id": "d", "text": "Пароли от боевых серверов", "correct": false},
							},
							"explanation": "Всё, кроме секретов. Пароли хранят в отдельных хранилищах.",
						},
						{
							"id":   "q6",
							"text": "Появился конфликт слияния. Что это значит?",
							"options": []map[string]any{
								{"id": "a", "text": "Двое изменили одну строку, и нужно выбрать правильный вариант вручную", "correct": true},
								{"id": "b", "text": "Репозиторий повреждён", "correct": false},
								{"id": "c", "text": "Нужно удалить ветку и начать заново", "correct": false},
							},
							"explanation": "Открываете файл, оставляете нужный вариант, убираете маркеры, делаете add и commit.",
						},
						{
							"id":     "q7",
							"review": true,
							"text":   "Повторение: где хранятся настройки программ в Linux?",
							"options": []map[string]any{
								{"id": "a", "text": "/etc", "correct": true},
								{"id": "b", "text": "/var/log", "correct": false},
								{"id": "c", "text": "/home", "correct": false},
							},
							"explanation": "Настройки в /etc, логи в /var/log.",
						},
						{
							"id":   "q8z",
							"type": "order",
							"text": "Расставьте шаги по порядку, чтобы отправить изменение в общий репозиторий:",
							"items": []map[string]any{
								{"id": "s1", "text": "Изменить файлы"},
								{"id": "s2", "text": "Добавить их в индекс: git add"},
								{"id": "s3", "text": "Сохранить версию: git commit"},
								{"id": "s4", "text": "Отправить на сервер: git push"},
							},
							"explanation": "add готовит изменения, commit сохраняет их локально, push отправляет в общий репозиторий.",
						},
						{
							"id":          "q9z",
							"type":        "blank",
							"text":        "Впишите команду, которая создаёт сохранённую версию из добавленных изменений: git …",
							"accept":      []string{"commit", "git commit"},
							"hint":        "С этого начинается история проекта.",
							"explanation": "git commit фиксирует состояние проекта — точку, к которой всегда можно вернуться.",
						},
						{
							"id":   "m1",
							"type": "match",
							"text": "Сопоставьте команду Git с её действием:",
							"pairs": []map[string]any{
								{"id": "p1", "left": "git add", "right": "добавить изменения в индекс"},
								{"id": "p2", "left": "git commit", "right": "сохранить версию"},
								{"id": "p3", "left": "git push", "right": "отправить на сервер"},
								{"id": "p4", "left": "git revert", "right": "отменить коммит новым коммитом"},
							},
							"explanation": "Путь изменения: add → commit → push; revert — безопасная отмена в общей ветке.",
						},
					},
					"resources": []map[string]any{
						{
							"title": "Conventional Commits",
							"url":   "https://www.conventionalcommits.org/ru/v1.0.0/",
							"note":  "соглашение о сообщениях коммитов, принятое во многих командах",
						},
					},
				},
			},
		},
	}
}
