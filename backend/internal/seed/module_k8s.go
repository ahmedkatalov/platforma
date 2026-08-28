package seed

func moduleKubernetes() ModuleSeed {
	return ModuleSeed{
		Title:   "Kubernetes",
		Summary: "Оркестратор контейнеров: что это, из чего состоит и как чинить",
		Lessons: []LessonSeed{
			{
				Title:       "Kubernetes: зачем нужен оркестратор",
				Kind:        "text",
				Summary:     "Поды, деплойменты и сервисы — на понятных примерах",
				DurationMin: 14,
				Content: map[string]any{
					"body": "## Зачем это нужно\n\n" +
						"Docker запускает контейнеры на одном сервере. Пока сервер один — всё просто.\n\n" +
						"Но что если серверов десять? Кто решит, на каком запускать? Кто перезапустит " +
						"упавший контейнер ночью? Кто обновит версию без остановки сервиса?\n\n" +
						"Этим занимается Kubernetes. Его называют оркестратором: он дирижирует контейнерами " +
						"на многих серверах.\n\n" +
						"## Главная идея: вы описываете результат\n\n" +
						"Вы не говорите «запусти контейнер на сервере номер три».\n\n" +
						"Вы говорите: «нужно три копии приложения версии 1.4.2». Кластер сам решает, где их " +
						"разместить, и следит, чтобы их всегда было три. Упала одна — поднимет новую.\n\n" +
						"## Четыре слова, которые надо знать\n\n" +
						"**Pod (под)** — самая маленькая единица запуска. Обычно внутри один контейнер. " +
						"Поды приходят и уходят: их пересоздают при обновлении и падении.\n\n" +
						"**Deployment** — описание: какое приложение, какая версия, сколько копий. " +
						"Именно с ним вы работаете чаще всего.\n\n" +
						"**Service** — постоянный адрес для группы подов. Поды меняются, адрес остаётся.\n\n" +
						"**Namespace** — папка внутри кластера. Разделяет команды и окружения.\n\n" +
						"## Как это связано\n\n" +
						"```\n" +
						"Deployment (нужно 3 копии)\n" +
						"     ↓ создаёт\n" +
						"  Pod   Pod   Pod\n" +
						"     ↑ находит по метке\n" +
						"  Service (постоянный адрес)\n" +
						"```\n\n" +
						"## Как выглядит описание\n\n" +
						"```yaml\n" +
						"apiVersion: apps/v1\n" +
						"kind: Deployment\n" +
						"metadata:\n" +
						"  name: api\n" +
						"spec:\n" +
						"  replicas: 3                 # сколько копий\n" +
						"  template:\n" +
						"    spec:\n" +
						"      containers:\n" +
						"        - name: api\n" +
						"          image: myapp:1.4.2  # какой образ\n" +
						"          ports:\n" +
						"            - containerPort: 8080\n" +
						"```\n\n" +
						"Файл описывает желаемое состояние. Применяется командой `kubectl apply -f файл.yaml`.\n\n" +
						"## Две проверки, без которых обновление ломает сайт\n\n" +
						"**readinessProbe** — готов ли под принимать запросы. Пока не готов, трафик к нему не идёт.\n\n" +
						"**livenessProbe** — жив ли контейнер. Не отвечает — кластер перезапустит.\n\n" +
						"```yaml\n" +
						"readinessProbe:\n" +
						"  httpGet:\n" +
						"    path: /health\n" +
						"    port: 8080\n" +
						"```\n\n" +
						"Без readiness при обновлении трафик пойдёт в ещё не запустившееся приложение, " +
						"и пользователи увидят ошибки.\n\n" +
						"## Сколько ресурсов просить\n\n" +
						"```yaml\n" +
						"resources:\n" +
						"  requests:      # сколько нужно гарантированно\n" +
						"    cpu: 100m\n" +
						"    memory: 128Mi\n" +
						"  limits:        # больше этого не дадим\n" +
						"    cpu: 500m\n" +
						"    memory: 256Mi\n" +
						"```\n\n" +
						"`requests` помогает кластеру выбрать сервер. `limits` не даёт одному сервису " +
						"съесть память всего узла.\n\n" +
						"## Команды для разбора проблем\n\n" +
						"```bash\n" +
						"kubectl get pods                  # что запущено и в каком состоянии\n" +
						"kubectl describe pod api-7d9f     # подробности и события\n" +
						"kubectl logs api-7d9f             # логи приложения\n" +
						"kubectl rollout undo deploy/api   # откатить обновление\n" +
						"```\n\n" +
						"Частые состояния подов:\n\n" +
						"| Состояние | Что значит |\n" +
						"|---|---|\n" +
						"| `Running` | работает |\n" +
						"| `Pending` | ждёт места на серверах |\n" +
						"| `CrashLoopBackOff` | падает при старте, кластер пытается снова |\n" +
						"| `ImagePullBackOff` | не может скачать образ: опечатка в имени или нет доступа |\n\n" +
						"При `CrashLoopBackOff` смотрите `kubectl logs` — причина почти всегда там.\n\n" +
						"## Частые ошибки новичка\n\n" +
						"- **Нет проб.** Обновление даёт всплеск ошибок у пользователей.\n" +
						"- **Нет `limits`.** Один сервис с утечкой памяти роняет соседей.\n" +
						"- **`latest` в образе.** Кластер не поймёт, что версия изменилась.\n" +
						"- **Правят прод через `kubectl edit`.** Изменение потеряется при следующем применении файла.\n\n" +
						"## Запомнить\n\n" +
						"1. Вы описываете желаемое состояние, кластер его поддерживает.\n" +
						"2. Pod живёт недолго, Service даёт постоянный адрес.\n" +
						"3. При проблеме: `kubectl get pods` → `describe` → `logs`.",
					"resources": []map[string]any{
						{
							"title": "Kubernetes — основы для начинающих",
							"url":   "https://kubernetes.io/ru/docs/tutorials/kubernetes-basics/",
							"note":  "официальный интерактивный курс, есть на русском",
						},
						{
							"title": "kubectl — шпаргалка команд",
							"url":   "https://kubernetes.io/docs/reference/kubectl/quick-reference/",
							"note":  "самые нужные команды одним листом",
						},
						{
							"title": "Диагностика приложений в кластере",
							"url":   "https://kubernetes.io/docs/tasks/debug/debug-application/",
							"note":  "официальный порядок разбора: под не стартует, падает, недоступен",
						},
					},
				},
			},
			{
				Title:       "Тренажёр: kubectl",
				Kind:        "terminal",
				Summary:     "Диагностика приложения в кластере",
				DurationMin: 20,
				Content: map[string]any{
					"resources": []map[string]any{
						{
							"title": "kubectl — шпаргалка",
							"url":   "https://kubernetes.io/docs/reference/kubectl/quick-reference/",
							"note":  "самые частые команды одним листом",
						},
						{
							"title": "Диагностика приложений в кластере",
							"url":   "https://kubernetes.io/docs/tasks/debug/debug-application/",
							"note":  "официальный порядок разбора: под не стартует, падает, недоступен",
						},
					},
					"intro": "В кластере развёрнут деплоймент api. Разберитесь с его состоянием.",
					"shell": "student@devops",
					"tasks": []map[string]any{
						{
							"id":       "k1",
							"prompt":   "Выведите список подов",
							"expected": []string{"kubectl get pods", "kubectl get pod", "kubectl get po"},
							"hint":     "kubectl get и тип ресурса",
							"success":  "Так проверяют, что запущено в текущем namespace.",
						},
						{
							"id":       "k2",
							"prompt":   "Посмотрите подробности и события пода api-7d9f (describe)",
							"expected": []string{"kubectl describe pod api-7d9f", "kubectl describe po api-7d9f"},
							"hint":     "kubectl describe pod имя",
							"success":  "В разделе Events видно, почему под не стартует.",
						},
						{
							"id":       "k3",
							"prompt":   "Прочитайте логи пода api-7d9f в режиме слежения",
							"expected": []string{"kubectl logs -f api-7d9f", "kubectl logs --follow api-7d9f"},
							"hint":     "kubectl logs с флагом -f",
							"success":  "Логи — первое место, куда смотрят при CrashLoopBackOff.",
						},
						{
							"id":       "k4",
							"prompt":   "Проверьте статус выката деплоймента api",
							"expected": []string{"kubectl rollout status deploy/api", "kubectl rollout status deployment/api", "kubectl rollout status deployment api"},
							"hint":     "kubectl rollout status deploy/имя",
							"success":  "Команда ждёт завершения обновления и покажет, если оно застряло.",
						},
						{
							"id":       "k5",
							"prompt":   "Откатите деплоймент api к предыдущей версии",
							"expected": []string{"kubectl rollout undo deploy/api", "kubectl rollout undo deployment/api"},
							"hint":     "rollout undo",
							"success":  "Быстрый откат — обязательная часть плана выката.",
						},
						{
							"id":       "k6",
							"prompt":   "Отмасштабируйте деплоймент api до 3 реплик",
							"expected": []string{"kubectl scale deploy/api --replicas=3", "kubectl scale deployment/api --replicas=3", "kubectl scale deployment api --replicas=3"},
							"hint":     "kubectl scale и флаг --replicas",
							"success":  "Реплик стало три — нагрузка распределится между ними.",
						},
					},
				},
			},
			{
				Title:       "Практика: конфигурация и секреты",
				Kind:        "code",
				Summary:     "Вынесите настройки из образа в ConfigMap и Secret",
				DurationMin: 22,
				Content: map[string]any{
					"resources": []map[string]any{
						{
							"title": "ConfigMap и Secret",
							"url":   "https://kubernetes.io/docs/concepts/configuration/secret/",
							"note":  "как подключать настройки переменными и файлами",
						},
						{
							"title": "Управление ресурсами контейнеров",
							"url":   "https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/",
							"note":  "requests, limits и что происходит при их превышении",
						},
					},
					"language": "yaml",
					"task": "Допишите манифест так, чтобы:\n\n" +
						"1. переменная `LOG_LEVEL` приходила из ConfigMap через `configMapKeyRef`;\n" +
						"2. пароль базы приходил из Secret через `secretKeyRef`;\n" +
						"3. в манифесте не осталось пароля открытым текстом;\n" +
						"4. были заданы `resources` с `requests`;\n" +
						"5. образ был с зафиксированной версией, без `latest`.",
					"starter": "apiVersion: apps/v1\n" +
						"kind: Deployment\n" +
						"metadata:\n" +
						"  name: api\n" +
						"spec:\n" +
						"  replicas: 2\n" +
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
						"          env:\n" +
						"            - name: LOG_LEVEL\n" +
						"              value: debug\n" +
						"            - name: DB_PASSWORD\n" +
						"              value: supersecret\n",
					"hint": "Ссылка на ConfigMap:\nvalueFrom:\n  configMapKeyRef:\n    name: api-config\n    key: log-level",
					"solution": "apiVersion: apps/v1\n" +
						"kind: Deployment\n" +
						"metadata:\n" +
						"  name: api\n" +
						"spec:\n" +
						"  replicas: 2\n" +
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
						"          image: registry.example.com/api:1.4.2\n" +
						"          resources:\n" +
						"            requests:\n" +
						"              cpu: 100m\n" +
						"              memory: 128Mi\n" +
						"            limits:\n" +
						"              cpu: 500m\n" +
						"              memory: 256Mi\n" +
						"          env:\n" +
						"            - name: LOG_LEVEL\n" +
						"              valueFrom:\n" +
						"                configMapKeyRef:\n" +
						"                  name: api-config\n" +
						"                  key: log-level\n" +
						"            - name: DB_PASSWORD\n" +
						"              valueFrom:\n" +
						"                secretKeyRef:\n" +
						"                  name: api-secrets\n" +
						"                  key: db-password\n",
					"checks": []map[string]any{
						{"type": "contains", "value": "configMapKeyRef", "message": "LOG_LEVEL берётся из ConfigMap"},
						{"type": "contains", "value": "secretKeyRef", "message": "Пароль берётся из Secret"},
						{"type": "notContains", "value": "value: supersecret", "message": "Пароля нет в открытом виде"},
						{"type": "regex", "value": "(?s)resources:.*requests:", "message": "Заданы запрошенные ресурсы"},
						{"type": "notContains", "value": "latest", "message": "Версия образа зафиксирована"},
					},
				},
			},
			{
				Title:       "Проверка: Kubernetes",
				Kind:        "quiz",
				Summary:     "Объекты, пробы и диагностика",
				DurationMin: 10,
				Content: map[string]any{
					"resources": []map[string]any{
						{
							"title": "Пробы: liveness, readiness, startup",
							"url":   "https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/",
							"note":  "разбор с примерами и типичными ошибками настройки",
						},
						{
							"title": "Gateway API",
							"url":   "https://gateway-api.sigs.k8s.io/",
							"note":  "преемник Ingress: маршрутизация, разделение ролей, поддержка не только HTTP",
						},
						{
							"title": "Production best practices (learnk8s)",
							"url":   "https://learnk8s.io/production-best-practices",
							"note":  "чек-лист перед выкатом в прод",
						},
					},
					"passScore": 70,
					"questions": []map[string]any{
						{
							"id":   "q1",
							"text": "Что такое Pod?",
							"options": []map[string]any{
								{"id": "a", "text": "Минимальная единица запуска: один или несколько контейнеров с общей сетью", "correct": true},
								{"id": "b", "text": "Физический сервер кластера", "correct": false},
								{"id": "c", "text": "Другое название образа", "correct": false},
							},
							"explanation": "Узел кластера — это Node, а Pod — единица запуска на нём.",
						},
						{
							"id":   "q2",
							"text": "Зачем нужен Service, если у пода уже есть IP-адрес?",
							"options": []map[string]any{
								{"id": "a", "text": "Поды пересоздаются и меняют адреса, Service даёт стабильную точку входа", "correct": true},
								{"id": "b", "text": "Service ускоряет работу контейнера", "correct": false},
								{"id": "c", "text": "Без Service под не запустится", "correct": false},
							},
							"explanation": "Service балансирует трафик между подами по меткам и не меняет адрес.",
						},
						{
							"id":       "q3",
							"text":     "Чем readinessProbe отличается от livenessProbe?",
							"multiple": true,
							"options": []map[string]any{
								{"id": "a", "text": "readiness решает, слать ли поду трафик", "correct": true},
								{"id": "b", "text": "liveness решает, надо ли перезапустить контейнер", "correct": true},
								{"id": "c", "text": "Обе пробы делают одно и то же", "correct": false},
							},
							"explanation": "Readiness управляет трафиком, liveness — перезапуском.",
						},
						{
							"id":   "q4",
							"text": "Под в статусе CrashLoopBackOff. С чего начать?",
							"options": []map[string]any{
								{"id": "a", "text": "Посмотреть kubectl logs и события в describe", "correct": true},
								{"id": "b", "text": "Сразу увеличить количество реплик", "correct": false},
								{"id": "c", "text": "Перезагрузить кластер", "correct": false},
							},
							"explanation": "Приложение падает при старте — причина почти всегда в логах и событиях.",
						},
						{
							"id":   "q5",
							"text": "Зачем задавать requests и limits для контейнера?",
							"options": []map[string]any{
								{"id": "a", "text": "Чтобы планировщик знал, сколько ресурсов нужно, и один сервис не съел весь узел", "correct": true},
								{"id": "b", "text": "Это обязательное поле, без него манифест невалиден", "correct": false},
								{"id": "c", "text": "Чтобы ускорить скачивание образа", "correct": false},
							},
							"explanation": "requests участвуют в планировании, limits ограничивают потребление.",
						},
					},
				},
			},
		},
	}
}
