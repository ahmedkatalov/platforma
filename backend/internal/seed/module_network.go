package seed

func moduleNetwork() ModuleSeed {
	return ModuleSeed{
		Title:   "Сети и веб-сервер",
		Summary: "Как запрос от пользователя доходит до программы и что делать, когда не доходит",
		Lessons: []LessonSeed{
			{
				Title:       "Путь запроса: от адреса до ответа",
				Kind:        "text",
				Summary:     "DNS, порты, коды ответов — простыми словами",
				DurationMin: 12,
				Content: map[string]any{
					"body": "## Зачем это нужно\n\n" +
						"Пользователь пишет «сайт не работает». Чтобы найти причину за минуты, а не за часы, " +
						"нужно понимать, через что проходит его запрос.\n\n" +
						"## Что происходит, когда открывают сайт\n\n" +
						"1. **DNS.** Имя `example.com` превращается в адрес вида `93.184.216.34`. " +
						"Это как телефонный справочник: по имени находим номер.\n" +
						"2. **Соединение.** Браузер стучится по этому адресу в нужный порт.\n" +
						"3. **Шифрование.** Для `https` сервер показывает сертификат, дальше разговор идёт в шифре.\n" +
						"4. **Запрос и ответ.** Браузер просит страницу, сервер отвечает.\n\n" +
						"## Что такое порт\n\n" +
						"Сервер один, а программ на нём много. Порт — номер двери, за которой ждёт нужная программа.\n\n" +
						"| Порт | Кто обычно за ним |\n" +
						"|---|---|\n" +
						"| 80 | сайт по http |\n" +
						"| 443 | сайт по https |\n" +
						"| 22 | подключение по ssh |\n" +
						"| 5432 | база PostgreSQL |\n\n" +
						"Посмотреть, какие двери открыты:\n\n" +
						"```bash\n" +
						"ss -tulpn\n" +
						"```\n\n" +
						"## Коды ответа: что говорит сервер\n\n" +
						"| Код | Смысл | Кто виноват |\n" +
						"|---|---|---|\n" +
						"| 200 | всё хорошо | — |\n" +
						"| 301, 302 | перенаправление на другой адрес | настройка |\n" +
						"| 401, 403 | не представился / нет прав | доступы |\n" +
						"| 404 | такой страницы нет | адрес или маршрут |\n" +
						"| 500 | программа сломалась внутри | приложение |\n" +
						"| 502, 504 | сервер не дозвонился до программы | программа не отвечает |\n\n" +
						"Запомните главное различие: **500 — программа ответила ошибкой, 502 — программа вообще не ответила.**\n\n" +
						"## Обратный прокси\n\n" +
						"Перед программой обычно стоит nginx. Он принимает запросы и передаёт их дальше.\n\n" +
						"Зачем так:\n\n" +
						"- обрабатывает шифрование, программе это делать не нужно;\n" +
						"- раздаёт картинки и стили быстрее;\n" +
						"- распределяет нагрузку между несколькими копиями программы;\n" +
						"- прячет внутреннее устройство от интернета.\n\n" +
						"Простая настройка выглядит так:\n\n" +
						"```nginx\n" +
						"server {\n" +
						"    listen 443 ssl;\n" +
						"    server_name app.example.com;\n" +
						"\n" +
						"    location / {\n" +
						"        proxy_pass http://app:8080;\n" +
						"        proxy_set_header X-Real-IP $remote_addr;\n" +
						"    }\n" +
						"}\n" +
						"```\n\n" +
						"Строка с `X-Real-IP` важна: без неё программа будет видеть адрес прокси вместо адреса " +
						"настоящего посетителя. В логах окажется один и тот же IP у всех.\n\n" +
						"## Правило перед перезагрузкой nginx\n\n" +
						"```bash\n" +
						"nginx -t          # проверить настройки\n" +
						"nginx -s reload   # применить\n" +
						"```\n\n" +
						"Сначала проверка, потом перезагрузка. Перезапуск со сломанным файлом настроек " +
						"положит сайт целиком.\n\n" +
						"## Инструменты для проверки\n\n" +
						"```bash\n" +
						"curl -I http://app:8080/health   # какой код отвечает программа\n" +
						"dig +short example.com           # какой адрес у домена\n" +
						"ss -tulpn                        # какие порты открыты\n" +
						"```\n\n" +
						"## Частые ошибки новичка\n\n" +
						"- **Видят 502 и лезут в настройки nginx.** Чаще всего дело в упавшей программе.\n" +
						"- **Перезагружают nginx без проверки.** Одна опечатка — и сайт лежит.\n" +
						"- **Открывают наружу лишние порты.** База данных не должна смотреть в интернет.\n\n" +
						"## Запомнить\n\n" +
						"1. 500 — программа ответила ошибкой, 502 — не ответила вовсе.\n" +
						"2. `nginx -t` перед каждой перезагрузкой.\n" +
						"3. Наружу открыты только 80 и 443, остальное — внутри.",
					"resources": []map[string]any{
						{
							"title": "MDN — коды ответа HTTP на русском",
							"url":   "https://developer.mozilla.org/ru/docs/Web/HTTP/Status",
							"note":  "полный список с объяснением, когда какой уместен",
						},
						{
							"title": "nginx — руководство для начинающих",
							"url":   "https://nginx.org/ru/docs/beginners_guide.html",
							"note":  "структура файла настроек на русском",
						},
					},
				},
			},
			{
				Title:       "Тренажёр: диагностика запроса",
				Kind:        "terminal",
				Summary:     "Найдите, почему сайт отвечает ошибкой",
				DurationMin: 20,
				Content: map[string]any{
					"resources": []map[string]any{
						{
							"title": "curl — руководство с примерами",
							"url":   "https://curl.se/docs/manual.html",
							"note":  "-I, -v, --resolve и другие флаги для диагностики",
						},
						{
							"title": "Отладка nginx: логи и трассировка",
							"url":   "https://nginx.org/en/docs/debugging_log.html",
							"note":  "как включить подробный лог, когда причина неочевидна",
						},
					},
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
					"resources": []map[string]any{
						{
							"title": "nginx — руководство для начинающих",
							"url":   "https://nginx.org/ru/docs/beginners_guide.html",
							"note":  "структура конфигурации на русском",
						},
						{
							"title": "Mozilla SSL Configuration Generator",
							"url":   "https://ssl-config.mozilla.org/",
							"note":  "готовая TLS-конфигурация под нужную версию nginx — не выдумывайте шифры сами",
						},
					},
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
					"resources": []map[string]any{
						{
							"title": "MDN — коды ответа HTTP",
							"url":   "https://developer.mozilla.org/ru/docs/Web/HTTP/Status",
							"note":  "полный список с объяснением, когда какой уместен",
						},
						{
							"title": "High Performance Browser Networking",
							"url":   "https://hpbn.co/",
							"note":  "книга целиком онлайн: TCP, TLS, HTTP/2 и HTTP/3 без упрощений",
						},
					},
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
