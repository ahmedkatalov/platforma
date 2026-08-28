package seed

func moduleNetwork() ModuleSeed {
	return ModuleSeed{
		Title:   "Сети и веб-сервер",
		Summary: "HTTP, DNS, TLS и nginx — как запрос доходит до приложения",
		Lessons: []LessonSeed{
			{
				Title:       "Путь запроса: DNS, TCP, HTTP",
				Kind:        "text",
				Summary:     "Что происходит между браузером и вашим приложением",
				DurationMin: 18,
				Content: map[string]any{
					"body": "## Что происходит, когда пользователь открывает сайт\n\n" +
						"1. **DNS.** Имя `app.example.com` превращается в IP-адрес. Отвечает DNS-сервер, " +
						"ответ кэшируется на время TTL.\n" +
						"2. **TCP.** Клиент устанавливает соединение с этим адресом на порт 80 или 443.\n" +
						"3. **TLS.** Для HTTPS происходит рукопожатие: сервер показывает сертификат, стороны " +
						"договариваются о шифровании.\n" +
						"4. **HTTP.** Клиент отправляет запрос, сервер отвечает статусом и телом.\n" +
						"5. **Обратный прокси.** Обычно первым отвечает nginx, а уже он передаёт запрос приложению.\n\n" +
						"## Коды ответа, которые нужно различать\n\n" +
						"| Код | Значение | Кто виноват |\n" +
						"|---|---|---|\n" +
						"| 200 | всё хорошо | — |\n" +
						"| 301 / 302 | перенаправление | конфигурация |\n" +
						"| 400 | некорректный запрос | клиент |\n" +
						"| 401 / 403 | нет аутентификации / нет прав | клиент или права |\n" +
						"| 404 | ресурс не найден | маршрут или клиент |\n" +
						"| 500 | ошибка приложения | приложение |\n" +
						"| 502 / 504 | прокси не достучался до приложения | приложение лежит или тайм-аут |\n\n" +
						"> 502 в логах nginx почти всегда означает: приложение не отвечает на своём порту. " +
						"Проверяйте, слушает ли процесс порт и жив ли контейнер.\n\n" +
						"## Обратный прокси\n\n" +
						"nginx стоит перед приложением и решает несколько задач: терминирует TLS, " +
						"раздаёт статику, балансирует нагрузку, ограничивает частоту запросов.\n\n" +
						"```nginx\n" +
						"server {\n" +
						"    listen 443 ssl;\n" +
						"    server_name app.example.com;\n" +
						"\n" +
						"    ssl_certificate     /etc/ssl/app.crt;\n" +
						"    ssl_certificate_key /etc/ssl/app.key;\n" +
						"\n" +
						"    location / {\n" +
						"        proxy_pass http://app:8080;\n" +
						"        proxy_set_header Host $host;\n" +
						"        proxy_set_header X-Real-IP $remote_addr;\n" +
						"    }\n" +
						"}\n" +
						"```\n\n" +
						"Заголовок `X-Real-IP` важен: без него приложение увидит адрес прокси вместо адреса клиента.\n\n" +
						"## Порты и процессы\n\n" +
						"```bash\n" +
						"ss -tulpn          # кто какие порты слушает\n" +
						"curl -I http://app # только заголовки ответа\n" +
						"dig +short app.example.com\n" +
						"nginx -t           # проверить конфигурацию перед перезагрузкой\n" +
						"```\n\n" +
						"Правило: **сначала `nginx -t`, потом `nginx -s reload`**. Перезагрузка со сломанным " +
						"конфигом уронит сайт.",
				},
			},
			{
				Title:       "Тренажёр: диагностика запроса",
				Kind:        "terminal",
				Summary:     "Найдите, почему сайт отвечает ошибкой",
				DurationMin: 20,
				Content: map[string]any{
					"intro": "Пользователи жалуются на ошибки. Пройдите путь запроса и найдите причину.",
					"shell": "student@devops",
					"tasks": []map[string]any{
						{
							"id":       "n1",
							"prompt":   "Узнайте IP-адрес домена example.com коротким выводом",
							"expected": []string{"dig +short example.com", "dig example.com +short"},
							"hint":     "dig и флаг +short",
							"success":  "DNS отвечает — значит имя разрешается корректно.",
						},
						{
							"id":       "n2",
							"prompt":   "Посмотрите, какие порты слушает сервер",
							"expected": []string{"ss -tulpn", "ss -tuln", "ss -lntp", "netstat -tulpn"},
							"hint":     "ss с флагами -tulpn",
							"success":  "Видно, что nginx слушает 80, приложение — 8080.",
						},
						{
							"id":       "n3",
							"prompt":   "Запросите только заголовки ответа приложения по адресу http://app:8080/health",
							"expected": []string{"curl -I http://app:8080/health", "curl --head http://app:8080/health"},
							"hint":     "curl с флагом -I",
							"success":  "Приложение отвечает 200 — значит проблема выше по цепочке.",
						},
						{
							"id":       "n4",
							"prompt":   "Найдите строки с кодом 500 в логе nginx /var/log/nginx/access.log",
							"expected": []string{"grep 500 /var/log/nginx/access.log", "grep \" 500 \" /var/log/nginx/access.log"},
							"hint":     "grep по файлу лога",
							"success":  "Ошибки приходят на POST /api/orders — сузили область поиска.",
						},
						{
							"id":       "n5",
							"prompt":   "Посчитайте количество строк в логе ошибок nginx",
							"expected": []string{"wc -l /var/log/nginx/error.log", "cat /var/log/nginx/error.log | wc -l"},
							"hint":     "wc -l и путь к файлу",
							"success":  "Немного строк — можно прочитать их целиком.",
						},
						{
							"id":       "n6",
							"prompt":   "Проверьте синтаксис конфигурации nginx",
							"expected": []string{"nginx -t"},
							"hint":     "Одна команда с флагом -t",
							"success":  "Конфигурация корректна — перезагружать безопасно.",
						},
					},
				},
			},
			{
				Title:       "Практика: конфигурация nginx",
				Kind:        "code",
				Summary:     "Соберите обратный прокси с TLS и проверкой здоровья",
				DurationMin: 22,
				Content: map[string]any{
					"language": "nginx",
					"task": "Допишите конфигурацию так, чтобы:\n\n" +
						"1. сервер слушал порт `443` с `ssl`;\n" +
						"2. запросы уходили в приложение через `proxy_pass http://app:8080`;\n" +
						"3. приложение получало настоящий адрес клиента через `X-Real-IP`;\n" +
						"4. запросы на `/health` не писались в лог доступа (`access_log off`);\n" +
						"5. был отдельный `server`-блок, перенаправляющий HTTP на HTTPS (`return 301`).",
					"starter": "server {\n" +
						"    listen 80;\n" +
						"    server_name app.example.com;\n" +
						"}\n" +
						"\n" +
						"server {\n" +
						"    server_name app.example.com;\n" +
						"\n" +
						"    ssl_certificate     /etc/ssl/app.crt;\n" +
						"    ssl_certificate_key /etc/ssl/app.key;\n" +
						"\n" +
						"    location / {\n" +
						"    }\n" +
						"\n" +
						"    location /health {\n" +
						"        proxy_pass http://app:8080/health;\n" +
						"    }\n" +
						"}\n",
					"hint": "Перенаправление пишется как return 301 https://$host$request_uri;",
					"solution": "server {\n" +
						"    listen 80;\n" +
						"    server_name app.example.com;\n" +
						"    return 301 https://$host$request_uri;\n" +
						"}\n" +
						"\n" +
						"server {\n" +
						"    listen 443 ssl;\n" +
						"    server_name app.example.com;\n" +
						"\n" +
						"    ssl_certificate     /etc/ssl/app.crt;\n" +
						"    ssl_certificate_key /etc/ssl/app.key;\n" +
						"\n" +
						"    location / {\n" +
						"        proxy_pass http://app:8080;\n" +
						"        proxy_set_header Host $host;\n" +
						"        proxy_set_header X-Real-IP $remote_addr;\n" +
						"    }\n" +
						"\n" +
						"    location /health {\n" +
						"        access_log off;\n" +
						"        proxy_pass http://app:8080/health;\n" +
						"    }\n" +
						"}\n",
					"checks": []map[string]any{
						{"type": "regex", "value": "listen\\s+443\\s+ssl", "message": "Сервер слушает 443 с ssl"},
						{"type": "contains", "value": "proxy_pass http://app:8080", "message": "Запросы уходят в приложение"},
						{"type": "regex", "value": "proxy_set_header\\s+X-Real-IP\\s+\\$remote_addr", "message": "Приложение получает адрес клиента"},
						{"type": "regex", "value": "access_log\\s+off", "message": "Проверки здоровья не засоряют лог"},
						{"type": "regex", "value": "return\\s+301\\s+https://", "message": "HTTP перенаправляется на HTTPS"},
					},
				},
			},
			{
				Title:       "Проверка: сети и nginx",
				Kind:        "quiz",
				Summary:     "Коды ответа, прокси и диагностика",
				DurationMin: 10,
				Content: map[string]any{
					"passScore": 70,
					"questions": []map[string]any{
						{
							"id":   "q1",
							"text": "nginx отдаёт 502 Bad Gateway. Что проверять первым делом?",
							"options": []map[string]any{
								{"id": "a", "text": "Жив ли процесс приложения и слушает ли он свой порт", "correct": true},
								{"id": "b", "text": "Срок действия домена", "correct": false},
								{"id": "c", "text": "Права на файлы статики", "correct": false},
							},
							"explanation": "502 означает: прокси не смог получить ответ от upstream.",
						},
						{
							"id":   "q2",
							"text": "Зачем передавать заголовок X-Real-IP или X-Forwarded-For?",
							"options": []map[string]any{
								{"id": "a", "text": "Иначе приложение увидит адрес прокси вместо адреса клиента", "correct": true},
								{"id": "b", "text": "Без него не работает HTTPS", "correct": false},
								{"id": "c", "text": "Это ускоряет проксирование", "correct": false},
							},
							"explanation": "Логи и ограничения по IP без этого заголовка бесполезны.",
						},
						{
							"id":   "q3",
							"text": "Что делает команда nginx -t?",
							"options": []map[string]any{
								{"id": "a", "text": "Проверяет синтаксис конфигурации, ничего не применяя", "correct": true},
								{"id": "b", "text": "Перезапускает nginx", "correct": false},
								{"id": "c", "text": "Показывает список активных соединений", "correct": false},
							},
							"explanation": "Проверка перед перезагрузкой спасает от падения сайта из-за опечатки.",
						},
						{
							"id":       "q4",
							"text":     "Какие задачи обычно решает обратный прокси?",
							"multiple": true,
							"options": []map[string]any{
								{"id": "a", "text": "Терминирует TLS", "correct": true},
								{"id": "b", "text": "Раздаёт статику", "correct": true},
								{"id": "c", "text": "Балансирует нагрузку между экземплярами", "correct": true},
								{"id": "d", "text": "Компилирует исходный код приложения", "correct": false},
							},
							"explanation": "Прокси снимает с приложения инфраструктурные задачи.",
						},
						{
							"id":   "q5",
							"text": "Какой командой посмотреть, какой процесс занимает порт 8080?",
							"options": []map[string]any{
								{"id": "a", "text": "ss -tulpn", "correct": true},
								{"id": "b", "text": "df -h", "correct": false},
								{"id": "c", "text": "uname -a", "correct": false},
							},
							"explanation": "ss показывает сокеты вместе с процессами, которые их держат.",
						},
					},
				},
			},
		},
	}
}
