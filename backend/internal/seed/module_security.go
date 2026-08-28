package seed

func moduleSecurity() ModuleSeed {
	return ModuleSeed{
		Title:   "Безопасность и секреты",
		Summary: "Доступы, секреты, образы и базовая гигиена продакшена",
		Lessons: []LessonSeed{
			{
				Title:       "Секреты, доступы и цепочка поставки",
				Kind:        "text",
				Summary:     "Где хранить пароли, как выдавать права и почему важны версии образов",
				DurationMin: 18,
				Content: map[string]any{
					"body": "## Секреты\n\n" +
						"Правило простое: **всё, что попало в Git, считается скомпрометированным**. " +
						"Даже если коммит удалить, он останется в истории и у всех, кто делал clone.\n\n" +
						"Где хранить:\n\n" +
						"- секреты CI (GitHub Secrets, GitLab CI Variables) — для пайплайнов;\n" +
						"- Secret в Kubernetes — для приложений в кластере;\n" +
						"- Vault или облачный менеджер секретов — когда нужны ротация и аудит.\n\n" +
						"В приложение секрет попадает переменной окружения или смонтированным файлом, " +
						"но не строкой в коде.\n\n" +
						"> Утёкший ключ не «удаляют» — его **отзывают и выпускают заново**. " +
						"Удаление коммита ничего не решает.\n\n" +
						"## Права доступа\n\n" +
						"Принцип наименьших привилегий: у сервиса ровно те права, что нужны для работы. " +
						"Учётка для деплоя не должна уметь удалять базу, а разработчику не нужен root на проде.\n\n" +
						"```bash\n" +
						"chmod 600 ~/.ssh/id_ed25519   # приватный ключ читает только владелец\n" +
						"chmod 644 config.yml          # конфиг читают все, пишет владелец\n" +
						"```\n\n" +
						"Права `777` — почти всегда ошибка: так файл может изменить любой процесс на машине.\n\n" +
						"## Контейнеры\n\n" +
						"- Не запускайте процесс от root: в Dockerfile создайте пользователя и укажите `USER`.\n" +
						"- Фиксируйте версии образов: `alpine:3.20`, а не `alpine:latest` — иначе сборка невоспроизводима.\n" +
						"- Не кладите секреты в слои образа: они остаются в истории слоёв даже после удаления файла.\n" +
						"- Сканируйте образы на уязвимости в пайплайне (trivy, grype).\n\n" +
						"## Сеть и доступ снаружи\n\n" +
						"Наружу торчит только то, что должно: обычно 80 и 443. База данных, панель мониторинга " +
						"и админские эндпоинты — во внутренней сети или за VPN.\n\n" +
						"```bash\n" +
						"ss -tulpn            # что реально слушает наружу\n" +
						"openssl s_client -connect app.example.com:443   # проверить TLS\n" +
						"```\n\n" +
						"## Гигиена, которая закрывает большинство дыр\n\n" +
						"| Практика | Что предотвращает |\n" +
						"|---|---|\n" +
						"| Секреты вне репозитория | утечку доступов |\n" +
						"| Обновление зависимостей | эксплуатацию известных уязвимостей |\n" +
						"| Наименьшие привилегии | превращение мелкого инцидента в крупный |\n" +
						"| Ротация ключей | долгую жизнь утёкшего ключа |\n" +
						"| Аудит и логи доступа | незаметное присутствие злоумышленника |\n\n" +
						"Безопасность — не отдельный этап в конце, а требование к каждому шагу конвейера.",
				},
			},
			{
				Title:       "Тренажёр: проверка доступов",
				Kind:        "terminal",
				Summary:     "Наведите порядок в правах и проверьте, что открыто наружу",
				DurationMin: 18,
				Content: map[string]any{
					"intro": "Перед выкатом проверьте базовую гигиену сервера.",
					"shell": "student@devops",
					"tasks": []map[string]any{
						{
							"id":       "s1",
							"prompt":   "Посмотрите, какие порты открыты и какие процессы их слушают",
							"expected": []string{"ss -tulpn", "ss -lntp", "netstat -tulpn"},
							"hint":     "ss с флагами -tulpn",
							"success":  "Наружу смотрят 80 и 8080 — с этим списком уже можно работать.",
						},
						{
							"id":       "s2",
							"prompt":   "Сгенерируйте ssh-ключ типа ed25519",
							"expected": []string{"ssh-keygen -t ed25519", "ssh-keygen -t ed25519 -C student@devops"},
							"hint":     "ssh-keygen -t тип",
							"success":  "Ключ создан. Приватную часть никому не передают.",
						},
						{
							"id":       "s3",
							"prompt":   "Ограничьте права приватного ключа ~/.ssh/id_ed25519 до 600",
							"expected": []string{"chmod 600 ~/.ssh/id_ed25519", "chmod 600 /home/student/.ssh/id_ed25519"},
							"hint":     "chmod 600 и путь к файлу",
							"success":  "Теперь ключ читает только владелец — иначе ssh откажется его использовать.",
						},
						{
							"id":       "s4",
							"prompt":   "Сгенерируйте случайный секрет: 32 байта в hex",
							"expected": []string{"openssl rand -hex 32"},
							"hint":     "openssl rand -hex длина",
							"success":  "Такой секрет годится для подписи токенов.",
						},
						{
							"id":     "s5",
							"prompt": "Проверьте, нет ли слова password в файле ~/projects/api/Dockerfile",
							"expected": []string{
								"grep password ~/projects/api/Dockerfile",
								"grep -i password ~/projects/api/Dockerfile",
								"grep password /home/student/projects/api/Dockerfile",
							},
							"hint":    "grep по файлу",
							"success": "Совпадений нет — секретов в образе не будет.",
						},
						{
							"id":     "s6",
							"prompt": "Убедитесь, что в Dockerfile указан непривилегированный пользователь: найдите строку USER",
							"expected": []string{
								"grep USER ~/projects/api/Dockerfile",
								"grep -n USER ~/projects/api/Dockerfile",
								"grep USER /home/student/projects/api/Dockerfile",
							},
							"hint":    "grep USER и путь к Dockerfile",
							"success": "USER app — процесс не работает от root.",
						},
					},
				},
			},
			{
				Title:       "Практика: безопасный Dockerfile",
				Kind:        "code",
				Summary:     "Приведите сборку образа в порядок",
				DurationMin: 22,
				Content: map[string]any{
					"language": "dockerfile",
					"task": "Исправьте Dockerfile так, чтобы:\n\n" +
						"1. базовый образ был с зафиксированной версией, без `latest`;\n" +
						"2. появился непривилегированный пользователь и инструкция `USER`;\n" +
						"3. пароль не задавался через `ENV`, а приходил снаружи;\n" +
						"4. использовалась многоэтапная сборка (`FROM ... AS build` и второй `FROM`);\n" +
						"5. был объявлен `EXPOSE` с портом приложения.",
					"starter": "FROM golang:latest\n" +
						"WORKDIR /src\n" +
						"COPY . .\n" +
						"ENV DB_PASSWORD=supersecret\n" +
						"RUN go build -o /app ./cmd/api\n" +
						"CMD [\"/app\"]\n",
					"hint": "Пользователь в alpine создаётся так: RUN adduser -D -u 10001 app",
					"solution": "FROM golang:1.25-alpine AS build\n" +
						"WORKDIR /src\n" +
						"COPY . .\n" +
						"RUN go build -o /app ./cmd/api\n" +
						"\n" +
						"FROM alpine:3.20\n" +
						"RUN adduser -D -u 10001 app\n" +
						"USER app\n" +
						"COPY --from=build /app /app\n" +
						"EXPOSE 8080\n" +
						"ENTRYPOINT [\"/app\"]\n",
					"checks": []map[string]any{
						{"type": "notContains", "value": "latest", "message": "Версии образов зафиксированы"},
						{"type": "regex", "value": "FROM\\s+\\S+\\s+AS\\s+\\w+", "message": "Используется многоэтапная сборка"},
						{"type": "regex", "value": "(?m)^FROM[\\s\\S]*^FROM", "message": "Есть второй этап с чистым образом"},
						{"type": "regex", "value": "(?m)^USER\\s+\\w+", "message": "Процесс работает не от root"},
						{"type": "notContains", "value": "ENV DB_PASSWORD", "message": "Пароль не зашит в образ"},
						{"type": "regex", "value": "(?m)^EXPOSE\\s+\\d+", "message": "Объявлен порт приложения"},
					},
				},
			},
			{
				Title:       "Проверка: безопасность",
				Kind:        "quiz",
				Summary:     "Секреты, права и образы",
				DurationMin: 10,
				Content: map[string]any{
					"passScore": 70,
					"questions": []map[string]any{
						{
							"id":   "q1",
							"text": "Токен случайно попал в коммит и был запушен. Что делать?",
							"options": []map[string]any{
								{"id": "a", "text": "Отозвать токен и выпустить новый", "correct": true},
								{"id": "b", "text": "Удалить коммит — этого достаточно", "correct": false},
								{"id": "c", "text": "Сделать репозиторий приватным", "correct": false},
							},
							"explanation": "Секрет уже мог быть скопирован. Единственный надёжный шаг — отзыв.",
						},
						{
							"id":   "q2",
							"text": "Почему процесс в контейнере не запускают от root?",
							"options": []map[string]any{
								{"id": "a", "text": "Чтобы взлом приложения не дал сразу максимальных прав", "correct": true},
								{"id": "b", "text": "Контейнер не запустится от root", "correct": false},
								{"id": "c", "text": "Так образ становится меньше", "correct": false},
							},
							"explanation": "Наименьшие привилегии ограничивают ущерб от компрометации.",
						},
						{
							"id":       "q3",
							"text":     "Где допустимо хранить пароль базы для приложения в кластере?",
							"multiple": true,
							"options": []map[string]any{
								{"id": "a", "text": "В Secret Kubernetes", "correct": true},
								{"id": "b", "text": "В Vault с ротацией", "correct": true},
								{"id": "c", "text": "В ENV внутри Dockerfile", "correct": false},
								{"id": "d", "text": "В репозитории в файле config.yml", "correct": false},
							},
							"explanation": "Слои образа и репозиторий читают многие — секретам там не место.",
						},
						{
							"id":   "q4",
							"text": "Что означают права 600 на приватном ssh-ключе?",
							"options": []map[string]any{
								{"id": "a", "text": "Читать и писать может только владелец", "correct": true},
								{"id": "b", "text": "Читать может любой пользователь", "correct": false},
								{"id": "c", "text": "Файл становится исполняемым", "correct": false},
							},
							"explanation": "При более широких правах ssh просто откажется использовать ключ.",
						},
						{
							"id":   "q5",
							"text": "Почему alpine:latest — плохой выбор для продакшена?",
							"options": []map[string]any{
								{"id": "a", "text": "Сборка перестаёт быть воспроизводимой: сегодня и завтра это разные образы", "correct": true},
								{"id": "b", "text": "latest всегда содержит уязвимости", "correct": false},
								{"id": "c", "text": "latest медленнее скачивается", "correct": false},
							},
							"explanation": "Фиксированная версия даёт одинаковый результат сборки в любой момент.",
						},
					},
				},
			},
		},
	}
}
