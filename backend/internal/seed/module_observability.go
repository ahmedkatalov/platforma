package seed

func moduleObservability() ModuleSeed {
	return ModuleSeed{
		Title:   "Мониторинг и логи",
		Summary: "Метрики, логи, алерты и SLO: как замечать проблемы раньше пользователей",
		Lessons: []LessonSeed{
			{
				Title:       "Метрики, логи и трассировки",
				Kind:        "text",
				Summary:     "Три источника данных о системе и как ими пользоваться",
				DurationMin: 18,
				Content: map[string]any{
					"body": "## Три взгляда на систему\n\n" +
						"- **Метрики** — числа во времени: запросов в секунду, доля ошибок, задержка, память. " +
						"Дешёвые, агрегируются, годятся для графиков и алертов.\n" +
						"- **Логи** — записи о событиях. Дороже метрик, зато содержат подробности: какой запрос, " +
						"какой пользователь, какая ошибка.\n" +
						"- **Трассировки** — путь одного запроса через все сервисы. Показывают, где именно время теряется.\n\n" +
						"Метрика говорит «что-то сломалось», лог — «что именно», трассировка — «где».\n\n" +
						"## Четыре золотых сигнала\n\n" +
						"1. **Задержка** — сколько занимает запрос (смотрите перцентили, не среднее).\n" +
						"2. **Трафик** — сколько запросов приходит.\n" +
						"3. **Ошибки** — доля неуспешных ответов.\n" +
						"4. **Насыщение** — насколько заполнены ресурсы: CPU, память, диск, пул соединений.\n\n" +
						"> Среднее время ответа скрывает проблемы. Если 95% запросов укладываются в 100 мс, " +
						"а 5% занимают 8 секунд, среднее будет приличным, а часть пользователей — недовольной. " +
						"Смотрите p95 и p99.\n\n" +
						"## Prometheus\n\n" +
						"Prometheus сам ходит по приложениям и забирает метрики с эндпоинта `/metrics`:\n\n" +
						"```\n" +
						"http_requests_total{code=\"200\",path=\"/api/orders\"} 18422\n" +
						"http_requests_total{code=\"500\",path=\"/api/orders\"} 37\n" +
						"```\n\n" +
						"Запросы пишутся на PromQL:\n\n" +
						"```promql\n" +
						"rate(http_requests_total{code=~\"5..\"}[5m])          # ошибок в секунду\n" +
						"histogram_quantile(0.95, http_request_duration_seconds_bucket)  # p95\n" +
						"```\n\n" +
						"## Алерты, которые не бесят\n\n" +
						"Плохой алерт срабатывает на каждый всплеск и приучает его игнорировать. Хороший — " +
						"означает, что страдают пользователи, и требует действия.\n\n" +
						"```yaml\n" +
						"- alert: HighErrorRate\n" +
						"  expr: rate(http_requests_total{code=~\"5..\"}[5m]) > 0.05\n" +
						"  for: 10m\n" +
						"  labels:\n" +
						"    severity: critical\n" +
						"  annotations:\n" +
						"    summary: Больше 5% ошибок за пять минут\n" +
						"```\n\n" +
						"Поле `for` обязательно: без него алерт сработает на секундный скачок.\n\n" +
						"## SLI, SLO и бюджет ошибок\n\n" +
						"- **SLI** — измеримый показатель: доля успешных запросов.\n" +
						"- **SLO** — цель: 99.9% успешных ответов за месяц.\n" +
						"- **Бюджет ошибок** — оставшаяся часть: при 99.9% это примерно 43 минуты в месяц.\n\n" +
						"Бюджет ошибок — инструмент договорённости: пока он есть, команда катит новые функции; " +
						"как только исчерпан — занимается надёжностью.\n\n" +
						"## Логи\n\n" +
						"Структурированные логи (JSON) удобнее искать и агрегировать. В каждой записи полезно иметь " +
						"уровень, время, идентификатор запроса и сервис. Пароли и токены в логи не пишут никогда.",
				},
			},
			{
				Title:       "Тренажёр: разбор инцидента",
				Kind:        "terminal",
				Summary:     "Найдите причину ошибок по метрикам и логам",
				DurationMin: 22,
				Content: map[string]any{
					"intro": "Алерт сообщил о росте ошибок. Разберитесь, что происходит с сервисом.",
					"shell": "student@devops",
					"tasks": []map[string]any{
						{
							"id":       "o1",
							"prompt":   "Заберите метрики приложения по адресу http://app:8080/metrics",
							"expected": []string{"curl http://app:8080/metrics", "curl -s http://app:8080/metrics"},
							"hint":     "curl и адрес эндпоинта",
							"success":  "Видно http_requests_total с кодом 500 — ошибки реальны.",
						},
						{
							"id":       "o2",
							"prompt":   "Найдите строки со словом ERROR в /var/log/app.log",
							"expected": []string{"grep ERROR /var/log/app.log", "grep -i error /var/log/app.log"},
							"hint":     "grep по файлу лога",
							"success":  "Ошибки связаны с платёжным сервисом: таймауты.",
						},
						{
							"id":     "o3",
							"prompt": "Посчитайте, сколько строк с ERROR в логе приложения",
							"expected": []string{
								"grep ERROR /var/log/app.log | wc -l",
								"cat /var/log/app.log | grep ERROR | wc -l",
								"grep -c ERROR /var/log/app.log",
							},
							"hint":    "Соедините grep и wc -l через конвейер",
							"success": "Теперь известен масштаб: столько ошибок за период лога.",
						},
						{
							"id":       "o4",
							"prompt":   "Посмотрите последние 20 строк системного журнала службы nginx",
							"expected": []string{"journalctl -u nginx -n 20", "journalctl -n 20 -u nginx"},
							"hint":     "journalctl -u имя-службы -n 20",
							"success":  "В журнале службы видно, что nginx работает штатно.",
						},
						{
							"id":       "o5",
							"prompt":   "Проверьте свободную память в человекочитаемом виде",
							"expected": []string{"free -h"},
							"hint":     "free и флаг -h",
							"success":  "Памяти достаточно — насыщение не причина.",
						},
						{
							"id":       "o6",
							"prompt":   "Проверьте, сколько свободного места на дисках",
							"expected": []string{"df -h"},
							"hint":     "df -h",
							"success":  "Диск не заполнен. Причина — внешний платёжный сервис, а не наш сервер.",
						},
						{
							"id":       "o7",
							"prompt":   "Посмотрите правила алертов в /etc/prometheus/alerts.yml",
							"expected": []string{"cat /etc/prometheus/alerts.yml"},
							"hint":     "cat и полный путь",
							"success":  "Алерт HighErrorRate сработал корректно: порог 5% и выдержка 10 минут.",
						},
					},
				},
			},
			{
				Title:       "Практика: правило алерта",
				Kind:        "code",
				Summary:     "Опишите алерт на рост задержки",
				DurationMin: 20,
				Content: map[string]any{
					"language": "yaml",
					"task": "Допишите правило так, чтобы:\n\n" +
						"1. алерт назывался `HighLatency`;\n" +
						"2. срабатывал по `histogram_quantile` от `http_request_duration_seconds_bucket`;\n" +
						"3. имел выдержку `for` не меньше пяти минут;\n" +
						"4. имел метку `severity`;\n" +
						"5. содержал понятное человеку описание в `annotations` с полем `summary`.",
					"starter": "groups:\n" +
						"  - name: api\n" +
						"    rules:\n" +
						"      - alert: HighErrorRate\n" +
						"        expr: rate(http_requests_total{code=~\"5..\"}[5m]) > 0.05\n" +
						"        for: 10m\n" +
						"        labels:\n" +
						"          severity: critical\n" +
						"        annotations:\n" +
						"          summary: Больше 5% ошибок за пять минут\n",
					"hint": "p95 считается так: histogram_quantile(0.95, rate(..._bucket[5m]))",
					"solution": "groups:\n" +
						"  - name: api\n" +
						"    rules:\n" +
						"      - alert: HighErrorRate\n" +
						"        expr: rate(http_requests_total{code=~\"5..\"}[5m]) > 0.05\n" +
						"        for: 10m\n" +
						"        labels:\n" +
						"          severity: critical\n" +
						"        annotations:\n" +
						"          summary: Больше 5% ошибок за пять минут\n" +
						"\n" +
						"      - alert: HighLatency\n" +
						"        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1\n" +
						"        for: 10m\n" +
						"        labels:\n" +
						"          severity: warning\n" +
						"        annotations:\n" +
						"          summary: p95 задержки превышает секунду\n",
					"checks": []map[string]any{
						{"type": "contains", "value": "alert: HighLatency", "message": "Добавлен алерт HighLatency"},
						{"type": "contains", "value": "histogram_quantile", "message": "Задержка считается через перцентиль"},
						{"type": "contains", "value": "http_request_duration_seconds_bucket", "message": "Используется гистограмма длительности"},
						{"type": "regex", "value": "for:\\s*([5-9]|[1-9][0-9])m", "message": "Есть выдержка не меньше пяти минут"},
						{"type": "regex", "value": "severity:\\s*\\w+", "message": "Указана важность алерта"},
						{"type": "regex", "value": "summary:\\s*\\S+", "message": "Есть описание для дежурного"},
					},
				},
			},
			{
				Title:       "Проверка: мониторинг",
				Kind:        "quiz",
				Summary:     "Сигналы, алерты и SLO",
				DurationMin: 12,
				Content: map[string]any{
					"passScore": 70,
					"questions": []map[string]any{
						{
							"id":   "q1",
							"text": "Почему для задержки смотрят p95 и p99, а не среднее?",
							"options": []map[string]any{
								{"id": "a", "text": "Среднее скрывает медленные запросы, от которых страдает часть пользователей", "correct": true},
								{"id": "b", "text": "Среднее сложнее вычислять", "correct": false},
								{"id": "c", "text": "Перцентили занимают меньше места", "correct": false},
							},
							"explanation": "Хвост распределения и есть та часть пользователей, которая уходит.",
						},
						{
							"id":       "q2",
							"text":     "Что относится к четырём золотым сигналам?",
							"multiple": true,
							"options": []map[string]any{
								{"id": "a", "text": "Задержка", "correct": true},
								{"id": "b", "text": "Трафик", "correct": true},
								{"id": "c", "text": "Ошибки", "correct": true},
								{"id": "d", "text": "Количество строк кода", "correct": false},
							},
							"explanation": "Четвёртый сигнал — насыщение ресурсов.",
						},
						{
							"id":   "q3",
							"text": "Зачем в правиле алерта поле for?",
							"options": []map[string]any{
								{"id": "a", "text": "Чтобы алерт срабатывал только при устойчивой проблеме, а не на секундном всплеске", "correct": true},
								{"id": "b", "text": "Чтобы ограничить время жизни алерта", "correct": false},
								{"id": "c", "text": "Чтобы задать интервал сбора метрик", "correct": false},
							},
							"explanation": "Без выдержки дежурного разбудит любой случайный скачок.",
						},
						{
							"id":   "q4",
							"text": "Что такое бюджет ошибок при SLO 99.9%?",
							"options": []map[string]any{
								{"id": "a", "text": "Допустимая доля неуспеха — около 43 минут недоступности в месяц", "correct": true},
								{"id": "b", "text": "Сумма денег на исправление инцидентов", "correct": false},
								{"id": "c", "text": "Количество разрешённых релизов в месяц", "correct": false},
							},
							"explanation": "Бюджет помогает решать, катить ли новые функции или заняться надёжностью.",
						},
						{
							"id":   "q5",
							"text": "Метрика показала рост ошибок. Что даст ответ на вопрос «что именно сломалось»?",
							"options": []map[string]any{
								{"id": "a", "text": "Логи с подробностями запроса и ошибки", "correct": true},
								{"id": "b", "text": "График загрузки CPU", "correct": false},
								{"id": "c", "text": "Список установленных пакетов", "correct": false},
							},
							"explanation": "Метрика сообщает о проблеме, лог объясняет её причину.",
						},
						{
							"id":   "q6",
							"text": "Что нельзя писать в логи?",
							"options": []map[string]any{
								{"id": "a", "text": "Пароли, токены и номера карт", "correct": true},
								{"id": "b", "text": "Идентификатор запроса", "correct": false},
								{"id": "c", "text": "Название сервиса и уровень записи", "correct": false},
							},
							"explanation": "Логи хранятся долго и доступны многим — секретам там не место.",
						},
					},
				},
			},
		},
	}
}
