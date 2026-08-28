package seed

func moduleExam() ModuleSeed {
	return ModuleSeed{
		Title:   "Итоговая аттестация",
		Summary: "Проверяем всё вместе: теория, диагностика и практика",
		Lessons: []LessonSeed{
			{
				Title:       "Что вас ждёт в аттестации",
				Kind:        "text",
				Summary:     "Как устроен финал курса и что повторить перед ним",
				DurationMin: 7,
				Content: map[string]any{
					"body": "## Зачем это нужно\n\n" +
						"Аттестация не для оценки, а для проверки себя. " +
						"Она показывает, какие темы уже в руках, а к каким стоит вернуться.\n\n" +
						"## Три части\n\n" +
						"**1. Экзамен.** Двенадцать вопросов по всем модулям. Порог выше обычного — 80%. " +
						"Пересдавать можно сколько угодно раз.\n\n" +
						"**2. Диагностика.** Тренажёр: пришла жалоба на ошибки, нужно найти причину. " +
						"Работаете как на настоящем дежурстве.\n\n" +
						"**3. Практика.** Собираете манифест для боевого запуска со всеми требованиями.\n\n" +
						"## Что повторить\n\n" +
						"Пробегитесь по блокам «Запомнить» в конце каждого урока. Если по какой-то теме " +
						"осталась неуверенность — вернитесь к её тренажёру, а не к тексту. " +
						"Руки запоминают лучше, чем глаза.\n\n" +
						"Короткий список того, что должно быть в голове:\n\n" +
						"| Тема | Главное |\n" +
						"|---|---|\n" +
						"| Linux | права 600/644/755, чтение логов, `df -h` при странных сбоях |\n" +
						"| Git | ветка → коммит → push, откат через `revert` |\n" +
						"| Сеть | 500 против 502, `nginx -t` перед перезагрузкой |\n" +
						"| Docker | образ и контейнер, тома, версии вместо `latest` |\n" +
						"| CI/CD | порядок этапов, секреты не в репозитории |\n" +
						"| IaC | `plan` перед `apply`, ручные правки затираются |\n" +
						"| Kubernetes | под, деплоймент, сервис, пробы, `logs` при падении |\n" +
						"| Мониторинг | метрика против лога, процентили вместо среднего |\n" +
						"| Безопасность | секрет в Git — отзывать, не запускать от root |\n\n" +
						"## Сертификат\n\n" +
						"Когда все уроки курса пройдены, платформа сама выдаёт сертификат: " +
						"уникальный номер и страница проверки, ссылку на которую можно отправить работодателю.\n\n" +
						"## Что делать после курса\n\n" +
						"1. Соберите свой маленький проект: приложение, база, Dockerfile, конвейер сборки.\n" +
						"2. Выложите его на GitHub — это лучшее портфолио для junior.\n" +
						"3. Поднимите Kubernetes локально (minikube или kind) и выкатите туда своё приложение.\n" +
						"4. Дальше можно готовиться к сертификации CKA — она признаётся работодателями.\n\n" +
						"Не пытайтесь выучить всё сразу. Возьмите один инструмент и доведите его " +
						"до уверенного владения. Для первой работы этого достаточно.\n\n" +
						"## Запомнить\n\n" +
						"1. Порог экзамена — 80%, пересдавать можно сколько угодно раз.\n" +
						"2. Перед экзаменом перечитайте блоки «Запомнить» в уроках.\n" +
						"3. После курса соберите свой проект — он весит больше сертификата.",
					"resources": []map[string]any{
						{
							"title": "roadmap.sh — DevOps Roadmap",
							"url":   "https://roadmap.sh/devops",
							"note":  "карта тем: что учить дальше и в каком порядке",
						},
						{
							"title": "Программа сертификации CKA",
							"url":   "https://github.com/cncf/curriculum",
							"note":  "официальный список тем экзамена по Kubernetes",
						},
						{
							"title": "Google SRE Book",
							"url":   "https://sre.google/books/",
							"note":  "две книги целиком бесплатно — следующий уровень",
						},
					},
				},
			},
			{
				Title:       "Экзамен по курсу",
				Kind:        "quiz",
				Summary:     "Двенадцать вопросов по всем модулям, порог 80%",
				DurationMin: 25,
				Content: map[string]any{
					"resources": []map[string]any{
						{
							"title": "Google SRE Book и SRE Workbook",
							"url":   "https://sre.google/books/",
							"note":  "две книги целиком онлайн — следующий уровень после курса",
						},
					},
					"passScore": 80,
					"intro":     "Отвечайте вдумчиво: порог прохождения выше обычного — 80%.",
					"questions": []map[string]any{
						{
							"id":   "e1",
							"text": "Команда сокращает время от коммита до продакшена. Какая практика влияет на это сильнее всего?",
							"options": []map[string]any{
								{"id": "a", "text": "Автоматический конвейер сборки, тестов и выката", "correct": true},
								{"id": "b", "text": "Еженедельные совещания по релизам", "correct": false},
								{"id": "c", "text": "Более подробная документация", "correct": false},
							},
							"explanation": "Ручные шаги — главный источник задержек и ошибок при доставке.",
						},
						{
							"id":   "e2",
							"text": "Файл имеет права 640 и владельца root. Пользователь app в группе root открывает его. Что он может?",
							"options": []map[string]any{
								{"id": "a", "text": "Только читать", "correct": true},
								{"id": "b", "text": "Читать и изменять", "correct": false},
								{"id": "c", "text": "Ничего", "correct": false},
							},
							"explanation": "6 — владелец: чтение и запись, 4 — группа: только чтение, 0 — остальные.",
						},
						{
							"id":   "e3",
							"text": "Вы отправили коммит в общую ветку и хотите отменить изменение. Что выбрать?",
							"options": []map[string]any{
								{"id": "a", "text": "git revert — он создаст коммит-отмену и не перепишет историю", "correct": true},
								{"id": "b", "text": "git reset --hard и принудительный push", "correct": false},
								{"id": "c", "text": "Удалить ветку и создать заново", "correct": false},
							},
							"explanation": "Переписывание общей истории ломает работу всей команде.",
						},
						{
							"id":   "e4",
							"text": "Пользователи получают 502, в логах nginx — connection refused к upstream. Где искать причину?",
							"options": []map[string]any{
								{"id": "a", "text": "Приложение не слушает свой порт: упало или не стартовало", "correct": true},
								{"id": "b", "text": "Истёк TLS-сертификат", "correct": false},
								{"id": "c", "text": "Неверная запись DNS", "correct": false},
							},
							"explanation": "Прокси дошёл до сервера, но соединение с приложением отклонено.",
						},
						{
							"id":       "e5",
							"text":     "Что сохраняется при перезапуске контейнера, а что теряется?",
							"multiple": true,
							"options": []map[string]any{
								{"id": "a", "text": "Данные в томе сохраняются", "correct": true},
								{"id": "b", "text": "Изменения в слое записи контейнера теряются при пересоздании", "correct": true},
								{"id": "c", "text": "Всё содержимое контейнера сохраняется навсегда", "correct": false},
							},
							"explanation": "Именно поэтому базы данных всегда работают с томами.",
						},
						{
							"id":   "e6",
							"text": "В каком порядке разумно расположить этапы конвейера?",
							"options": []map[string]any{
								{"id": "a", "text": "Линтер → тесты → сборка образа → выкат на stage → выкат на прод", "correct": true},
								{"id": "b", "text": "Сборка образа → выкат на прод → тесты", "correct": false},
								{"id": "c", "text": "Тесты → выкат на прод → линтер", "correct": false},
							},
							"explanation": "Дешёвые проверки идут первыми, прод — последним и только после stage.",
						},
						{
							"id":   "e7",
							"text": "Terraform plan показывает «2 to change, 1 to destroy» на проде. Ваши действия?",
							"options": []map[string]any{
								{"id": "a", "text": "Разобраться, какой ресурс удаляется и почему, до применения", "correct": true},
								{"id": "b", "text": "Применить и посмотреть, что получится", "correct": false},
								{"id": "c", "text": "Удалить файл состояния и начать заново", "correct": false},
							},
							"explanation": "Удаление ресурса на проде — потенциальный простой.",
						},
						{
							"id":   "e8",
							"text": "Под в Kubernetes перезапускается каждые несколько минут. С чего начать разбор?",
							"options": []map[string]any{
								{"id": "a", "text": "kubectl logs и kubectl describe pod — посмотреть ошибки и события", "correct": true},
								{"id": "b", "text": "Увеличить число реплик", "correct": false},
								{"id": "c", "text": "Пересоздать кластер", "correct": false},
							},
							"explanation": "Причина перезапусков почти всегда видна в логах или в событиях пода.",
						},
						{
							"id":   "e9",
							"text": "Зачем нужна readinessProbe при обновлении деплоймента?",
							"options": []map[string]any{
								{"id": "a", "text": "Чтобы трафик не шёл в ещё не готовый под", "correct": true},
								{"id": "b", "text": "Чтобы перезапускать зависшие контейнеры", "correct": false},
								{"id": "c", "text": "Чтобы ограничить потребление памяти", "correct": false},
							},
							"explanation": "Без readiness обновление даёт всплеск ошибок у пользователей.",
						},
						{
							"id":   "e10",
							"text": "Средняя задержка 120 мс, p99 — 9 секунд. Что это значит?",
							"options": []map[string]any{
								{"id": "a", "text": "Часть пользователей регулярно ждёт очень долго — нужно разбираться с хвостом", "correct": true},
								{"id": "b", "text": "Всё в порядке, среднее в норме", "correct": false},
								{"id": "c", "text": "Метрика собрана неверно", "correct": false},
							},
							"explanation": "Хвост распределения — это реальные пользователи, а не статистический шум.",
						},
						{
							"id":       "e11",
							"text":     "Какие практики относятся к безопасной работе с секретами?",
							"multiple": true,
							"options": []map[string]any{
								{"id": "a", "text": "Хранить секреты в менеджере секретов или переменных CI", "correct": true},
								{"id": "b", "text": "Отзывать и перевыпускать утёкшие ключи", "correct": true},
								{"id": "c", "text": "Регулярно ротировать доступы", "correct": true},
								{"id": "d", "text": "Хранить .env с паролями в репозитории для удобства команды", "correct": false},
							},
							"explanation": "Удобство не оправдывает хранение секретов в истории Git.",
						},
						{
							"id":   "e12",
							"text": "Что означает исчерпанный бюджет ошибок при SLO 99.9%?",
							"options": []map[string]any{
								{"id": "a", "text": "Пора остановить новые функции и заняться надёжностью", "correct": true},
								{"id": "b", "text": "Нужно понизить SLO, чтобы бюджет появился", "correct": false},
								{"id": "c", "text": "Можно катить релизы быстрее обычного", "correct": false},
							},
							"explanation": "Бюджет ошибок — это договорённость между скоростью и стабильностью.",
						},
					},
				},
			},
			{
				Title:       "Итоговая диагностика",
				Kind:        "terminal",
				Summary:     "Пройдите путь от жалобы пользователя до причины сбоя",
				DurationMin: 25,
				Content: map[string]any{
					"resources": []map[string]any{
						{
							"title": "Диагностика приложений в Kubernetes",
							"url":   "https://kubernetes.io/docs/tasks/debug/debug-application/",
							"note":  "официальный порядок разбора проблем с подами и сервисами",
						},
					},
					"intro": "Поступила жалоба: часть запросов к API завершается ошибкой. Проведите разбор целиком — от внешней проверки до логов и кластера.",
					"shell": "student@devops",
					"tasks": []map[string]any{
						{
							"id":       "x1",
							"prompt":   "Проверьте, отвечает ли приложение: запросите заголовки http://app:8080/health",
							"expected": []string{"curl -I http://app:8080/health", "curl --head http://app:8080/health"},
							"hint":     "curl -I и адрес",
							"success":  "Само приложение живо — идём дальше по цепочке.",
						},
						{
							"id":       "x2",
							"prompt":   "Посмотрите, какие порты слушает сервер",
							"expected": []string{"ss -tulpn", "ss -lntp", "netstat -tulpn"},
							"hint":     "ss -tulpn",
							"success":  "Порты на месте, процессы запущены.",
						},
						{
							"id":       "x3",
							"prompt":   "Найдите ошибки в логе приложения /var/log/app.log",
							"expected": []string{"grep ERROR /var/log/app.log", "grep -i error /var/log/app.log"},
							"hint":     "grep ERROR по файлу",
							"success":  "Видно таймауты к платёжному сервису.",
						},
						{
							"id":       "x4",
							"prompt":   "Посмотрите список запущенных контейнеров",
							"expected": []string{"docker ps"},
							"hint":     "docker ps",
							"success":  "Контейнеры работают, значит дело не в них.",
						},
						{
							"id":       "x5",
							"prompt":   "Посмотрите поды в кластере",
							"expected": []string{"kubectl get pods", "kubectl get pod", "kubectl get po"},
							"hint":     "kubectl get pods",
							"success":  "Один под не в состоянии Running — вот и виновник.",
						},
						{
							"id":       "x6",
							"prompt":   "Изучите события проблемного пода api-7d9f",
							"expected": []string{"kubectl describe pod api-7d9f", "kubectl describe po api-7d9f"},
							"hint":     "kubectl describe pod имя",
							"success":  "В событиях видна причина падения.",
						},
						{
							"id":       "x7",
							"prompt":   "Откатите деплоймент api к предыдущей версии",
							"expected": []string{"kubectl rollout undo deploy/api", "kubectl rollout undo deployment/api"},
							"hint":     "kubectl rollout undo deploy/имя",
							"success":  "Сервис вернулся к рабочей версии — инцидент закрыт, дальше разбор причин.",
						},
					},
				},
			},
			{
				Title:       "Итоговая практика: манифест развёртывания",
				Kind:        "code",
				Summary:     "Соберите Deployment, готовый к продакшену",
				DurationMin: 30,
				Content: map[string]any{
					"resources": []map[string]any{
						{
							"title": "Kubernetes production best practices",
							"url":   "https://learnk8s.io/production-best-practices",
							"note":  "чек-лист: пробы, ресурсы, безопасность, отказоустойчивость",
						},
						{
							"title": "Конфигурация подов: securityContext",
							"url":   "https://kubernetes.io/docs/tasks/configure-pod-container/security-context/",
							"note":  "runAsNonRoot, readOnlyRootFilesystem и другие ограничения",
						},
					},
					"language": "yaml",
					"task": "Доведите манифест до продакшен-качества. В нём должно быть:\n\n" +
						"1. не меньше двух реплик;\n" +
						"2. образ с зафиксированной версией, без `latest`;\n" +
						"3. `readinessProbe` и `livenessProbe`;\n" +
						"4. `resources` с `requests` и `limits`;\n" +
						"5. пароль из `secretKeyRef`, а не открытым текстом;\n" +
						"6. запуск не от root: `runAsNonRoot: true`.",
					"starter": "apiVersion: apps/v1\n" +
						"kind: Deployment\n" +
						"metadata:\n" +
						"  name: api\n" +
						"spec:\n" +
						"  replicas: 1\n" +
						"  selector:\n" +
						"    matchLabels:\n" +
						"      app: api\n" +
						"  template:\n" +
						"    metadata:\n" +
						"      labels:\n" +
						"        app: api\n" +
						"    spec:\n" +
						"      containers:\n" +
						"        - name: api\n" +
						"          image: registry.example.com/api:latest\n" +
						"          ports:\n" +
						"            - containerPort: 8080\n" +
						"          env:\n" +
						"            - name: DB_PASSWORD\n" +
						"              value: supersecret\n",
					"hint": "Секрет подключается так:\nvalueFrom:\n  secretKeyRef:\n    name: api-secrets\n    key: db-password",
					"solution": "apiVersion: apps/v1\n" +
						"kind: Deployment\n" +
						"metadata:\n" +
						"  name: api\n" +
						"spec:\n" +
						"  replicas: 3\n" +
						"  selector:\n" +
						"    matchLabels:\n" +
						"      app: api\n" +
						"  template:\n" +
						"    metadata:\n" +
						"      labels:\n" +
						"        app: api\n" +
						"    spec:\n" +
						"      securityContext:\n" +
						"        runAsNonRoot: true\n" +
						"      containers:\n" +
						"        - name: api\n" +
						"          image: registry.example.com/api:1.4.2\n" +
						"          ports:\n" +
						"            - containerPort: 8080\n" +
						"          readinessProbe:\n" +
						"            httpGet:\n" +
						"              path: /health\n" +
						"              port: 8080\n" +
						"            initialDelaySeconds: 5\n" +
						"          livenessProbe:\n" +
						"            httpGet:\n" +
						"              path: /health\n" +
						"              port: 8080\n" +
						"            periodSeconds: 10\n" +
						"          resources:\n" +
						"            requests:\n" +
						"              cpu: 100m\n" +
						"              memory: 128Mi\n" +
						"            limits:\n" +
						"              cpu: 500m\n" +
						"              memory: 256Mi\n" +
						"          env:\n" +
						"            - name: DB_PASSWORD\n" +
						"              valueFrom:\n" +
						"                secretKeyRef:\n" +
						"                  name: api-secrets\n" +
						"                  key: db-password\n",
					"checks": []map[string]any{
						{"type": "regex", "value": "replicas:\\s*[2-9]", "message": "Реплик не меньше двух"},
						{"type": "notContains", "value": "latest", "message": "Версия образа зафиксирована"},
						{"type": "contains", "value": "readinessProbe", "message": "Есть проверка готовности"},
						{"type": "contains", "value": "livenessProbe", "message": "Есть проверка живости"},
						{"type": "regex", "value": "(?s)resources:.*requests:", "message": "Заданы requests"},
						{"type": "regex", "value": "(?s)resources:.*limits:", "message": "Заданы limits"},
						{"type": "contains", "value": "secretKeyRef", "message": "Пароль берётся из секрета"},
						{"type": "notContains", "value": "value: supersecret", "message": "Пароля нет в открытом виде"},
						{"type": "regex", "value": "runAsNonRoot:\\s*true", "message": "Контейнер работает не от root"},
					},
				},
			},
		},
	}
}
