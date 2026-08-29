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
						"Вот как выглядит вывод `kubectl get pods` на практике:\n" +
						"\n" +
						"```bash\n" +
						"kubectl get pods\n" +
						"NAME            READY   STATUS             RESTARTS   AGE\n" +
						"api-7d9f-2k4l   1/1     Running            0          6m\n" +
						"api-7d9f-8x1m   1/1     Running            0          6m\n" +
						"api-7d9f-qp3z   0/1     CrashLoopBackOff   5          3m\n" +
						"```\n" +
						"\n" +
						"`READY 1/1` — под готов. Третий под: `0/1`, статус `CrashLoopBackOff` и растущий `RESTARTS` — значит, падает при старте.\n" +
						"\n" +
						"Частые состояния подов:\n\n" +
						"| Состояние | Что значит |\n" +
						"|---|---|\n" +
						"| `Running` | работает |\n" +
						"| `Pending` | ждёт места на серверах |\n" +
						"| `CrashLoopBackOff` | падает при старте, кластер пытается снова |\n" +
						"| `ImagePullBackOff` | не может скачать образ: опечатка в имени или нет доступа |\n\n" +
						"При `CrashLoopBackOff` смотрите `kubectl logs` — причина почти всегда там.\n\n" +
						"Ещё одно частое состояние — `Pending`. Под создан, но планировщик не нашёл узел с нужными ресурсами. Причина видна в событиях:\n" +
						"\n" +
						"```bash\n" +
						"kubectl get pods\n" +
						"NAME            READY   STATUS    RESTARTS   AGE\n" +
						"web-5c8d-h2k9   0/1     Pending   0          45s\n" +
						"\n" +
						"kubectl describe pod web-5c8d-h2k9\n" +
						"...\n" +
						"Events:\n" +
						"  Type     Reason            Age   From               Message\n" +
						"  ----     ------            ----  ----               -------\n" +
						"  Warning  FailedScheduling  40s   default-scheduler  0/3 nodes are available: 3 Insufficient cpu.\n" +
						"```\n" +
						"\n" +
						"Читаем последнюю строку: `Insufficient cpu` — ни на одном из трёх узлов нет свободного CPU под `requests` пода. Лечится снижением `requests` или добавлением узлов, а не правкой кода.\n" +
						"\n" +
						"## Частые ошибки новичка\n\n" +
						"- **Нет проб.** Обновление даёт всплеск ошибок у пользователей.\n" +
						"- **Нет `limits`.** Один сервис с утечкой памяти роняет соседей.\n" +
						"- **`latest` в образе.** Кластер не поймёт, что версия изменилась.\n" +
						"- **Правят прод через `kubectl edit`.** Изменение потеряется при следующем применении файла.\n\n" +
						"А так выглядит разбор `ImagePullBackOff`. Сам список подов только называет проблему — настоящая причина в конце `describe`, в разделе Events:\n" +
						"\n" +
						"```bash\n" +
						"kubectl get pods\n" +
						"NAME            READY   STATUS             RESTARTS   AGE\n" +
						"api-7d9f-qp3z   0/1     ImagePullBackOff   0          90s\n" +
						"\n" +
						"kubectl describe pod api-7d9f-qp3z\n" +
						"Name:    api-7d9f-qp3z\n" +
						"Status:  Pending\n" +
						"Containers:\n" +
						"  api:\n" +
						"    Image:   registry.example.com/api:1.4.2\n" +
						"    State:   Waiting\n" +
						"      Reason:  ImagePullBackOff\n" +
						"Events:\n" +
						"  Type     Reason     Age                From     Message\n" +
						"  ----     ------     ----               ----     -------\n" +
						"  Normal   Pulling    88s (x3 over 90s)  kubelet  Pulling image \"registry.example.com/api:1.4.2\"\n" +
						"  Warning  Failed     86s (x3 over 90s)  kubelet  Failed to pull image \"registry.example.com/api:1.4.2\": manifest unknown\n" +
						"  Warning  Failed     86s (x3 over 90s)  kubelet  Error: ErrImagePull\n" +
						"  Normal   BackOff    60s (x5 over 90s)  kubelet  Back-off pulling image \"registry.example.com/api:1.4.2\"\n" +
						"```\n" +
						"\n" +
						"Строка `manifest unknown` означает: такого тега нет в реестре. Частая причина — опечатка в версии образа. Если бы было `unauthorized`, дело в доступе к приватному реестру.\n" +
						"\n" +
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
						{
							"title": "Pod — концепция",
							"url":   "https://kubernetes.io/docs/concepts/workloads/pods/",
							"note":  "что такое под, почему он одноразовый и когда пересоздаётся",
						},
						{
							"title": "Deployment — управление репликами",
							"url":   "https://kubernetes.io/docs/concepts/workloads/controllers/deployment/",
							"note":  "желаемое состояние, количество реплик и обновления",
						},
						{
							"title": "Service — постоянный адрес для подов",
							"url":   "https://kubernetes.io/docs/concepts/services-networking/service/",
							"note":  "как Service находит поды по меткам и балансирует трафик",
						},
					},
				},
			},
			{
				Title:       "Квиз: объекты кластера",
				Kind:        "quiz",
				Summary:     "Поды, деплойменты, сервисы и пробы",
				DurationMin: 7,
				Content: map[string]any{
					"passScore": 70,
					"questions": []map[string]any{
						{
							"id":   "k1",
							"text": "Что такое под?",
							"options": []map[string]any{
								{"id": "a", "text": "Самая маленькая единица запуска, обычно с одним контейнером внутри", "correct": true},
								{"id": "b", "text": "Физический сервер кластера", "correct": false},
								{"id": "c", "text": "Другое название образа", "correct": false},
							},
							"explanation": "Сервер кластера называется узлом (node).",
						},
						{
							"id":   "k2",
							"text": "Зачем нужен Service, если у пода есть свой адрес?",
							"options": []map[string]any{
								{"id": "a", "text": "Поды пересоздаются и меняют адреса, а Service даёт постоянный", "correct": true},
								{"id": "b", "text": "Service ускоряет работу приложения", "correct": false},
								{"id": "c", "text": "Без Service под не запустится", "correct": false},
							},
							"explanation": "Service находит поды по меткам и распределяет между ними запросы.",
						},
						{
							"id":   "k3",
							"text": "Что описывает Deployment?",
							"options": []map[string]any{
								{"id": "a", "text": "Какое приложение, какой версии и сколько копий должно работать", "correct": true},
								{"id": "b", "text": "Правила доступа в кластер", "correct": false},
								{"id": "c", "text": "Настройки сети между узлами", "correct": false},
							},
							"explanation": "Кластер сам поддерживает описанное состояние.",
						},
						{
							"id":   "k4",
							"text": "Под в состоянии CrashLoopBackOff. Что это значит?",
							"options": []map[string]any{
								{"id": "a", "text": "Приложение падает при старте, кластер пробует запустить снова", "correct": true},
								{"id": "b", "text": "Не хватает места в кластере", "correct": false},
								{"id": "c", "text": "Образ не скачивается", "correct": false},
							},
							"explanation": "Проблемы со скачиванием образа — это ImagePullBackOff.",
						},
						{
							"id":   "k5",
							"text": "Зачем нужна readinessProbe?",
							"options": []map[string]any{
								{"id": "a", "text": "Чтобы запросы не шли в под, который ещё не готов их принимать", "correct": true},
								{"id": "b", "text": "Чтобы перезапускать зависшие контейнеры", "correct": false},
								{"id": "c", "text": "Чтобы ограничить память", "correct": false},
							},
							"explanation": "Перезапуском занимается livenessProbe.",
						},
						{
							"id":     "k6",
							"review": true,
							"text":   "Повторение: что произойдёт с данными приложения при пересоздании контейнера, если не использовать том?",
							"options": []map[string]any{
								{"id": "a", "text": "Они пропадут", "correct": true},
								{"id": "b", "text": "Сохранятся в образе", "correct": false},
								{"id": "c", "text": "Переедут на другой сервер", "correct": false},
							},
							"explanation": "Данные хранят в томах — это правило работает и в Kubernetes.",
						},
					},
					"resources": []map[string]any{
						{
							"title": "Основы Kubernetes — интерактивный курс",
							"url":   "https://kubernetes.io/ru/docs/tutorials/kubernetes-basics/",
							"note":  "официальный, есть на русском",
						},
					},
				},
			},
			{
				Title:       "Обновления без простоя: пробы, реплики, откат",
				Kind:        "text",
				Summary:     "Как Kubernetes катит новую версию и почему иногда ломает сайт",
				DurationMin: 12,
				Content: map[string]any{
					"body": "## Зачем это нужно\n\n" +
						"Выкат новой версии не должен быть заметен пользователю. " +
						"Kubernetes умеет это сам, но только если вы правильно всё описали.\n\n" +
						"Большинство «сайт лёг на две минуты при релизе» — это ненастроенные пробы.\n\n" +
						"## Как идёт обновление\n\n" +
						"Вы поменяли версию образа и применили файл. Дальше кластер:\n\n" +
						"1. поднимает под с новой версией;\n" +
						"2. ждёт, пока он станет готов;\n" +
						"3. переводит на него трафик;\n" +
						"4. гасит один старый под;\n" +
						"5. повторяет, пока не обновятся все.\n\n" +
						"Такой способ называется **rolling update**. Старая версия работает, " +
						"пока новая не готова.\n\n" +
						"Во время выката старые и новые поды какое-то время живут вместе. Это видно по возрасту (`AGE`):\n" +
						"\n" +
						"```bash\n" +
						"kubectl get pods\n" +
						"NAME            READY   STATUS    RESTARTS   AGE\n" +
						"api-6b4c-abcd   1/1     Running   0          20m\n" +
						"api-6b4c-efgh   1/1     Running   0          20m\n" +
						"api-9f2d-ijkl   0/1     Running   0          8s\n" +
						"```\n" +
						"\n" +
						"Старые поды (`api-6b4c`, 20m) держат трафик. Новый (`api-9f2d`, 8s) уже `Running`, но `0/1` — readinessProbe ещё не пустила к нему трафик.\n" +
						"\n" +
						"## Ключевой момент: что значит «готов»\n\n" +
						"Кластер не знает, когда ваше приложение действительно готово принимать запросы. " +
						"Он узнаёт это из **readinessProbe**.\n\n" +
						"```yaml\n" +
						"readinessProbe:\n" +
						"  httpGet:\n" +
						"    path: /health\n" +
						"    port: 8080\n" +
						"  initialDelaySeconds: 5    # подождать перед первой проверкой\n" +
						"  periodSeconds: 5          # как часто проверять\n" +
						"```\n\n" +
						"Без этой пробы кластер считает под готовым сразу после старта контейнера. " +
						"Приложение ещё подключается к базе, а на него уже льётся трафик. " +
						"Пользователи видят ошибки.\n\n" +
						"## Три пробы и их роли\n\n" +
						"| Проба | Вопрос | Что делает кластер |\n" +
						"|---|---|---|\n" +
						"| `startup` | приложение вообще запустилось? | ждёт дольше обычного при старте |\n" +
						"| `readiness` | готово принимать трафик? | убирает из балансировки, если нет |\n" +
						"| `liveness` | ещё живо? | перезапускает контейнер |\n\n" +
						"Осторожно с `liveness`: если сделать её слишком строгой, кластер начнёт " +
						"перезапускать здоровое, но медленное приложение. Тогда простой станет " +
						"бесконечным циклом перезапусков.\n\n" +
						"## Сколько подов держать\n\n" +
						"Одна реплика = простой при любом обновлении и падении. Минимум для боевого " +
						"сервиса — две, лучше три.\n\n" +
						"```yaml\n" +
						"spec:\n" +
						"  replicas: 3\n" +
						"  strategy:\n" +
						"    rollingUpdate:\n" +
						"      maxUnavailable: 0    # не гасить старые, пока нет новых\n" +
						"      maxSurge: 1          # поднимать по одному сверх нормы\n" +
						"```\n\n" +
						"С `maxUnavailable: 0` мощность сервиса не проседает во время выката.\n\n" +
						"## Когда всё пошло не так\n\n" +
						"```bash\n" +
						"kubectl rollout status deploy/api    # следить за ходом выката\n" +
						"kubectl rollout undo deploy/api      # вернуть предыдущую версию\n" +
						"kubectl rollout history deploy/api   # список версий\n" +
						"```\n\n" +
						"Откат — обычная операция, а не признак провала. Сначала возвращаем рабочую " +
						"версию, потом спокойно разбираемся в причинах.\n\n" +
						"## Частые ошибки новичка\n\n" +
						"- **Нет readinessProbe.** Каждый релиз даёт всплеск ошибок.\n" +
						"- **Одна реплика на проде.** Любое обновление — простой.\n" +
						"- **Слишком агрессивная livenessProbe.** Приложение перезапускается под нагрузкой.\n" +
						"- **Разбираются в причинах вместо отката.** Сначала вернуть работу пользователям.\n\n" +
						"Как выглядит зависший выкат целиком. Катим новую версию, но новый под падает, а старые продолжают держать трафик:\n" +
						"\n" +
						"```bash\n" +
						"kubectl set image deploy/api api=registry.example.com/api:1.5.0\n" +
						"deployment.apps/api image updated\n" +
						"\n" +
						"kubectl rollout status deploy/api\n" +
						"Waiting for deployment \"api\" rollout to finish: 1 out of 3 new replicas have been updated...\n" +
						"\n" +
						"kubectl get pods\n" +
						"NAME            READY   STATUS             RESTARTS   AGE\n" +
						"api-6b4c-abcd   1/1     Running            0          30m\n" +
						"api-6b4c-efgh   1/1     Running            0          30m\n" +
						"api-9f2d-ijkl   0/1     CrashLoopBackOff   3          80s\n" +
						"\n" +
						"kubectl logs api-9f2d-ijkl\n" +
						"2026/08/30 10:14:02 starting api v1.5.0\n" +
						"2026/08/30 10:14:02 fatal: DB_PASSWORD is not set\n" +
						"```\n" +
						"\n" +
						"Читаем вывод: старые поды (`AGE 30m`) остались `Running` — благодаря `maxUnavailable: 0` простоя нет. Новый упал: в новой версии забыли переменную. Сначала откат, разбор потом:\n" +
						"\n" +
						"```bash\n" +
						"kubectl rollout undo deploy/api\n" +
						"deployment.apps/api rolled back\n" +
						"```\n" +
						"\n" +
						"## Запомнить\n\n" +
						"1. readinessProbe решает, идёт ли трафик в под. Без неё релизы ломают сайт.\n" +
						"2. Реплик должно быть минимум две.\n" +
						"3. При проблеме сначала `rollout undo`, разбор — потом.",
					"resources": []map[string]any{
						{
							"title": "Пробы: liveness, readiness, startup",
							"url":   "https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/",
							"note":  "официальное руководство с примерами настройки",
						},
						{
							"title": "Стратегии обновления Deployment",
							"url":   "https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#strategy",
							"note":  "maxSurge, maxUnavailable и откат",
						},
						{
							"title": "Обновление приложения — учебник",
							"url":   "https://kubernetes.io/docs/tutorials/kubernetes-basics/update/update-intro/",
							"note":  "rolling update по шагам в интерактивном тренажёре",
						},
						{
							"title": "Нарушения работы подов (Disruptions)",
							"url":   "https://kubernetes.io/docs/concepts/workloads/pods/disruptions/",
							"note":  "PodDisruptionBudget: сколько подов можно ронять при обслуживании узлов",
						},
					},
				},
			},
			{
				Title:       "Квиз: обновления и надёжность",
				Kind:        "quiz",
				Summary:     "Пробы, реплики и откат выката",
				DurationMin: 8,
				Content: map[string]any{
					"passScore": 70,
					"questions": []map[string]any{
						{
							"id":   "q1",
							"text": "Зачем нужна readinessProbe?",
							"options": []map[string]any{
								{"id": "a", "text": "Чтобы трафик не шёл в под, который ещё не готов работать", "correct": true},
								{"id": "b", "text": "Чтобы перезапускать зависшие контейнеры", "correct": false},
								{"id": "c", "text": "Чтобы ограничить память контейнера", "correct": false},
							},
							"explanation": "Перезапуском занимается livenessProbe, а readiness управляет трафиком.",
						},
						{
							"id":   "q2",
							"text": "Сколько реплик нужно боевому сервису как минимум?",
							"options": []map[string]any{
								{"id": "a", "text": "Две, лучше три", "correct": true},
								{"id": "b", "text": "Одна — Kubernetes сам перезапустит её при падении", "correct": false},
								{"id": "c", "text": "Не меньше десяти", "correct": false},
							},
							"explanation": "С одной репликой любое обновление или падение означает простой.",
						},
						{
							"id":   "q3",
							"text": "Выкат сломал прод. Что делать первым делом?",
							"options": []map[string]any{
								{"id": "a", "text": "kubectl rollout undo — вернуть рабочую версию", "correct": true},
								{"id": "b", "text": "Искать причину в коде, пока пользователи ждут", "correct": false},
								{"id": "c", "text": "Увеличить количество реплик", "correct": false},
							},
							"explanation": "Сначала возвращаем сервис пользователям, разбор причин — следующим шагом.",
						},
						{
							"id":   "q4",
							"text": "Что делает настройка maxUnavailable: 0 при обновлении?",
							"options": []map[string]any{
								{"id": "a", "text": "Не даёт гасить старые поды, пока не готовы новые", "correct": true},
								{"id": "b", "text": "Запрещает обновление совсем", "correct": false},
								{"id": "c", "text": "Обновляет все поды одновременно", "correct": false},
							},
							"explanation": "Так мощность сервиса не проседает во время выката.",
						},
						{
							"id":       "q5",
							"review":   true,
							"text":     "Повторение: под в состоянии CrashLoopBackOff. С чего начать?",
							"multiple": true,
							"options": []map[string]any{
								{"id": "a", "text": "Посмотреть kubectl logs", "correct": true},
								{"id": "b", "text": "Посмотреть события через kubectl describe pod", "correct": true},
								{"id": "c", "text": "Сразу пересоздать кластер", "correct": false},
							},
							"explanation": "Приложение падает при старте — причина почти всегда в логах или событиях.",
						},
					},
					"resources": []map[string]any{
						{
							"title": "Production best practices",
							"url":   "https://learnk8s.io/production-best-practices",
							"note":  "чек-лист перед выкатом в прод",
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
						{
							"title": "Отладка запущенного пода",
							"url":   "https://kubernetes.io/docs/tasks/debug/debug-application/debug-running-pod/",
							"note":  "kubectl exec, ephemeral-контейнеры и чтение логов внутри пода",
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
						{
							"id":     "q6",
							"review": true,
							"text":   "Повторение: как в GitOps выкатывают новую версию?",
							"options": []map[string]any{
								{"id": "a", "text": "Меняют версию в репозитории, агент применяет изменение сам", "correct": true},
								{"id": "b", "text": "Заходят на сервер и выполняют команду вручную", "correct": false},
								{"id": "c", "text": "Пересоздают кластер", "correct": false},
							},
							"explanation": "Выкат = коммит, откат = revert.",
						},
					},
				},
			},
		},
	}
}
