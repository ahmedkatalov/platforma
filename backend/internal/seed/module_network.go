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
						"Вот что увидите — слева адрес и порт, справа программа за этой «дверью»:\n" +
						"\n" +
						"```bash\n" +
						"$ ss -tulpn\n" +
						"State   Local Address:Port   Process\n" +
						"LISTEN  0.0.0.0:80           nginx\n" +
						"LISTEN  0.0.0.0:443          nginx\n" +
						"LISTEN  127.0.0.1:8080       app\n" +
						"LISTEN  127.0.0.1:5432       postgres\n" +
						"```\n" +
						"\n" +
						"`0.0.0.0` — дверь открыта наружу, `127.0.0.1` — только внутри сервера.\n" +
						"\n" +
						"Так и должно быть: 80 и 443 смотрят в интернет, а 8080 и база — нет.\n" +
						"\n" +
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
						"## Разбор боевого 502\n" +
						"\n" +
						"Сайт отдаёт 502. Пройдём путь запроса сверху вниз и найдём причину.\n" +
						"\n" +
						"```bash\n" +
						"$ curl -I https://app.example.com/\n" +
						"HTTP/1.1 502 Bad Gateway\n" +
						"Server: nginx/1.24.0\n" +
						"\n" +
						"# Прокси жив, но ответа от приложения нет. Слушает ли оно порт 8080?\n" +
						"$ ss -tulpn | grep 8080\n" +
						"$\n" +
						"\n" +
						"# Пусто — на 8080 никто не слушает. Смотрим службу приложения.\n" +
						"$ systemctl status app\n" +
						"● app.service - Orders API\n" +
						"     Loaded: loaded (/etc/systemd/system/app.service; enabled)\n" +
						"     Active: failed (Result: exit-code) since Sat 2026-08-29 10:12:03 UTC\n" +
						"    Process: 4821 ExecStart=/usr/bin/app (code=exited, status=1/FAILURE)\n" +
						"\n" +
						"# Служба упала. Читаем последние строки журнала.\n" +
						"$ journalctl -u app -n 5 --no-pager\n" +
						"app[4821]: panic: dial tcp 10.0.0.5:5432: connect: connection refused\n" +
						"app[4821]: goroutine 1 [running]:\n" +
						"app[4821]: main.main()\n" +
						"```\n" +
						"\n" +
						"Приложение не достучалось до базы на 5432 и упало. 502 был симптомом, а не причиной.\n" +
						"\n" +
						"Вывод: не лезьте сразу в конфиг nginx — идите по цепочке до реального источника.\n" +
						"\n" +
						"\n" +
						"## Частые ошибки новичка\n\n" +
						"- **Видят 502 и лезут в настройки nginx.** Чаще всего дело в упавшей программе.\n" +
						"- **Перезагружают nginx без проверки.** Одна опечатка — и сайт лежит.\n" +
						"- **Открывают наружу лишние порты.** База данных не должна смотреть в интернет.\n\n" +
						"## Как выглядит провал nginx -t\n" +
						"\n" +
						"Проверка ловит опечатку до перезагрузки — иначе сайт ляжет целиком.\n" +
						"\n" +
						"```bash\n" +
						"$ nginx -t\n" +
						"nginx: [emerg] unknown directive \"proxy_passs\" in /etc/nginx/conf.d/app.conf:14\n" +
						"nginx: configuration file /etc/nginx/nginx.conf test failed\n" +
						"```\n" +
						"\n" +
						"nginx назвал файл и строку 14 — там лишняя `s` в `proxy_pass`. Правим и проверяем снова.\n" +
						"\n" +
						"```bash\n" +
						"$ nginx -t\n" +
						"nginx: the configuration file /etc/nginx/nginx.conf syntax is ok\n" +
						"nginx: configuration file /etc/nginx/nginx.conf test is successful\n" +
						"\n" +
						"$ nginx -s reload\n" +
						"```\n" +
						"\n" +
						"Перезагрузка безопасна только после строки «test is successful».\n" +
						"\n" +
						"\n" +
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
						{
							"title": "RFC 9110 — семантика HTTP и коды ответов",
							"url":   "https://www.rfc-editor.org/rfc/rfc9110.html",
							"note":  "первоисточник: что официально означает каждый код (200, 404, 500, 502)",
						},
						{
							"title": "man ss(8) — сокеты и открытые порты",
							"url":   "https://man7.org/linux/man-pages/man8/ss.8.html",
							"note":  "точное значение флагов -tulpn и как читать столбцы вывода",
						},
						{
							"title": "Cloudflare Learning — What is DNS?",
							"url":   "https://www.cloudflare.com/learning/dns/what-is-dns/",
							"note":  "наглядно про разрешение имени в IP и цепочку DNS-запросов",
						},
					},
				},
			},
			{
				Title:       "Квиз: как работает запрос",
				Kind:        "quiz",
				Summary:     "DNS, порты и коды ответов",
				DurationMin: 6,
				Content: map[string]any{
					"passScore": 70,
					"questions": []map[string]any{
						{
							"id":   "w1",
							"text": "Что делает DNS?",
							"options": []map[string]any{
								{"id": "a", "text": "Превращает имя сайта в IP-адрес", "correct": true},
								{"id": "b", "text": "Шифрует соединение", "correct": false},
								{"id": "c", "text": "Распределяет нагрузку между серверами", "correct": false},
							},
							"explanation": "Это как телефонный справочник: по имени находим номер.",
						},
						{
							"id":   "w2",
							"text": "На каком порту обычно работает сайт по https?",
							"options": []map[string]any{
								{"id": "a", "text": "443", "correct": true},
								{"id": "b", "text": "80", "correct": false},
								{"id": "c", "text": "22", "correct": false},
							},
							"explanation": "80 — http, 22 — ssh, 443 — https.",
						},
						{
							"id":   "w3",
							"text": "Сервер вернул 500. О чём это говорит?",
							"options": []map[string]any{
								{"id": "a", "text": "Программа ответила ошибкой — сломалось внутри неё", "correct": true},
								{"id": "b", "text": "Такой страницы не существует", "correct": false},
								{"id": "c", "text": "Пользователь не авторизован", "correct": false},
							},
							"explanation": "404 — нет страницы, 401 — нет авторизации.",
						},
						{
							"id":   "w4",
							"text": "А если пришёл 502?",
							"options": []map[string]any{
								{"id": "a", "text": "Прокси не смог достучаться до программы — скорее всего она не работает", "correct": true},
								{"id": "b", "text": "Программа вернула ошибку в ответе", "correct": false},
								{"id": "c", "text": "Истёк сертификат", "correct": false},
							},
							"explanation": "500 — программа ответила ошибкой, 502 — не ответила вовсе.",
						},
						{
							"id":       "w5",
							"text":     "Зачем перед приложением ставят nginx?",
							"multiple": true,
							"options": []map[string]any{
								{"id": "a", "text": "Он берёт на себя шифрование", "correct": true},
								{"id": "b", "text": "Раздаёт статику быстрее приложения", "correct": true},
								{"id": "c", "text": "Распределяет запросы между копиями приложения", "correct": true},
								{"id": "d", "text": "Ускоряет работу базы данных", "correct": false},
							},
							"explanation": "Прокси снимает с приложения инфраструктурные задачи.",
						},
						{
							"id":   "w6",
							"text": "Что нужно сделать перед перезагрузкой nginx?",
							"options": []map[string]any{
								{"id": "a", "text": "Проверить настройки командой nginx -t", "correct": true},
								{"id": "b", "text": "Перезапустить сервер целиком", "correct": false},
								{"id": "c", "text": "Очистить логи", "correct": false},
							},
							"explanation": "Одна опечатка в файле настроек — и сайт ляжет целиком.",
						},
					},
					"resources": []map[string]any{
						{
							"title": "MDN — коды ответа HTTP",
							"url":   "https://developer.mozilla.org/ru/docs/Web/HTTP/Status",
							"note":  "полный список на русском",
						},
					},
				},
			},
			{
				Title:       "HTTPS и сертификаты",
				Kind:        "text",
				Summary:     "Что такое TLS, откуда берутся сертификаты и почему сайт вдруг «небезопасен»",
				DurationMin: 11,
				Content: map[string]any{
					"body": "## Зачем это нужно\n\n" +
						"Без HTTPS данные идут по сети открытым текстом. Логин, пароль, номер карты " +
						"может прочитать любой на пути запроса.\n\n" +
						"HTTPS — это тот же HTTP, но внутри защищённого канала. " +
						"Браузеры давно помечают сайты без него как небезопасные.\n\n" +
						"## Что такое сертификат\n\n" +
						"Сертификат — документ, который подтверждает: этот сервер действительно " +
						"обслуживает домен `app.example.com`.\n\n" +
						"Выдаёт его **удостоверяющий центр** — организация, которой доверяют браузеры. " +
						"Самый известный бесплатный центр — Let's Encrypt.\n\n" +
						"В сертификате три важные вещи:\n\n" +
						"- для какого домена он выдан;\n" +
						"- кто его выдал;\n" +
						"- до какого числа он действует.\n\n" +
						"## Как выглядит установка\n\n" +
						"У вас два файла: сертификат и приватный ключ. Ключ — секрет, его никому не отдают.\n\n" +
						"```nginx\n" +
						"server {\n" +
						"    listen 443 ssl;\n" +
						"    server_name app.example.com;\n" +
						"\n" +
						"    ssl_certificate     /etc/ssl/app.crt;   # сам сертификат\n" +
						"    ssl_certificate_key /etc/ssl/app.key;   # приватный ключ\n" +
						"}\n" +
						"```\n\n" +
						"Права на ключ ставят строгие: `chmod 600`.\n\n" +
						"## Три ошибки, которые видит пользователь\n\n" +
						"| Что видит пользователь | Причина |\n" +
						"|---|---|\n" +
						"| «Сертификат истёк» | забыли продлить |\n" +
						"| «Не совпадает имя» | сертификат для другого домена |\n" +
						"| «Не удалось проверить» | сервер отдал не всю цепочку сертификатов |\n\n" +
						"Первая — самая частая. Сертификаты Let's Encrypt живут 90 дней, " +
						"поэтому продление автоматизируют. Программа certbot делает это сама.\n\n" +
						"## Как проверить\n\n" +
						"```bash\n" +
						"curl -I https://app.example.com                    # отвечает ли сайт\n" +
						"openssl s_client -connect app.example.com:443      # подробности соединения\n" +
						"```\n\n" +
						"В ответе `openssl` смотрят на строки `notBefore` и `notAfter` — срок действия.\n\n" +
						"Полный вывод `openssl` большой. Чтобы увидеть только сроки, добавьте `x509 -dates`:\n" +
						"\n" +
						"```bash\n" +
						"$ echo | openssl s_client -connect app.example.com:443 2>/dev/null | openssl x509 -noout -dates\n" +
						"notBefore=Jun  1 00:00:00 2026 GMT\n" +
						"notAfter=Aug 30 23:59:59 2026 GMT\n" +
						"```\n" +
						"\n" +
						"`notAfter` — дата, после которой браузер начнёт ругаться. Здесь сертификат жив до 30 августа.\n" +
						"\n" +
						"## Где заканчивается шифрование\n\n" +
						"Обычно HTTPS обрабатывает nginx, а дальше внутри своей сети он ходит к приложению " +
						"по обычному HTTP. Это называется «терминировать TLS».\n\n" +
						"Так делают, чтобы приложение не занималось шифрованием. Главное, чтобы участок " +
						"от nginx до приложения был во внутренней сети, а не в интернете.\n\n" +
						"## TLS-рукопожатие вживую\n" +
						"\n" +
						"Флаг `-v` у curl показывает обмен сертификатами. Так виден момент обрыва.\n" +
						"\n" +
						"```bash\n" +
						"$ curl -v https://expired.example.com/ 2>&1 | head -n 10\n" +
						"*   Trying 93.184.216.34:443...\n" +
						"* Connected to expired.example.com (93.184.216.34) port 443\n" +
						"* TLSv1.3 (OUT), TLS handshake, Client hello (1):\n" +
						"* TLSv1.3 (IN), TLS handshake, Server hello (2):\n" +
						"* TLSv1.3 (IN), TLS handshake, Certificate (11):\n" +
						"* SSL certificate problem: certificate has expired\n" +
						"* Closing connection\n" +
						"curl: (60) SSL certificate problem: certificate has expired\n" +
						"```\n" +
						"\n" +
						"Соединение оборвалось на шаге Certificate. Код curl `(60)` — это ошибка проверки сертификата.\n" +
						"\n" +
						"Проверяем срок напрямую, не доверяя браузеру.\n" +
						"\n" +
						"```bash\n" +
						"$ echo | openssl s_client -connect expired.example.com:443 2>/dev/null \\\n" +
						"    | openssl x509 -noout -enddate\n" +
						"notAfter=Jul 15 12:00:00 2026 GMT\n" +
						"```\n" +
						"\n" +
						"Дата в прошлом — сертификат просрочен. Лечится продлением: `certbot renew`.\n" +
						"\n" +
						"\n" +
						"## Частые ошибки новичка\n\n" +
						"- **Забывают про продление.** Поставьте автопродление сразу, не «потом».\n" +
						"- **Кладут приватный ключ в репозиторий.** Это утечка, ключ придётся перевыпускать.\n" +
						"- **Оставляют HTTP без перенаправления.** Добавьте `return 301 https://...`.\n\n" +
						"## Запомнить\n\n" +
						"1. Сертификат подтверждает домен и имеет срок действия.\n" +
						"2. Продление автоматизируют — иначе сайт однажды упадёт сам.\n" +
						"3. Приватный ключ хранят с правами 600 и никогда не коммитят.",
					"resources": []map[string]any{
						{
							"title": "Let's Encrypt — как это работает",
							"url":   "https://letsencrypt.org/ru/how-it-works/",
							"note":  "на русском: бесплатные сертификаты и автопродление",
						},
						{
							"title": "Генератор конфигурации TLS от Mozilla",
							"url":   "https://ssl-config.mozilla.org/",
							"note":  "готовые настройки для nginx — не подбирайте шифры вручную",
						},
						{
							"title": "certbot — установка и автопродление",
							"url":   "https://eff-certbot.readthedocs.io/en/stable/using.html",
							"note":  "как настроить certbot renew, чтобы сертификат не протух молча",
						},
						{
							"title": "RFC 8446 — протокол TLS 1.3",
							"url":   "https://www.rfc-editor.org/rfc/rfc8446.html",
							"note":  "первоисточник о том, как устроено TLS-рукопожатие и обмен сертификатами",
						},
						{
							"title": "Qualys SSL Labs — тест сервера",
							"url":   "https://www.ssllabs.com/ssltest/",
							"note":  "онлайн проверяет цепочку сертификатов и настройки TLS, ставит оценку",
						},
					},
				},
			},
			{
				Title:       "Квиз: HTTPS и сертификаты",
				Kind:        "quiz",
				Summary:     "Сертификаты, сроки и типичные ошибки",
				DurationMin: 7,
				Content: map[string]any{
					"passScore": 70,
					"questions": []map[string]any{
						{
							"id":   "q1",
							"text": "Что подтверждает сертификат сайта?",
							"options": []map[string]any{
								{"id": "a", "text": "Что сервер действительно обслуживает этот домен", "correct": true},
								{"id": "b", "text": "Что сайт не содержит вирусов", "correct": false},
								{"id": "c", "text": "Что владелец сайта заплатил за хостинг", "correct": false},
							},
							"explanation": "Сертификат связывает домен и сервер, а доверие к нему даёт удостоверяющий центр.",
						},
						{
							"id":   "q2",
							"text": "Браузер пишет «сертификат истёк». Что произошло?",
							"options": []map[string]any{
								{"id": "a", "text": "Забыли продлить сертификат — у него закончился срок действия", "correct": true},
								{"id": "b", "text": "Сервер перегружен", "correct": false},
								{"id": "c", "text": "Домен купил другой человек", "correct": false},
							},
							"explanation": "Поэтому продление всегда автоматизируют: у Let's Encrypt срок всего 90 дней.",
						},
						{
							"id":   "q3",
							"text": "Где хранят приватный ключ от сертификата?",
							"options": []map[string]any{
								{"id": "a", "text": "На сервере, с правами 600, вне репозитория", "correct": true},
								{"id": "b", "text": "В Git рядом с конфигурацией nginx", "correct": false},
								{"id": "c", "text": "В переменной окружения браузера", "correct": false},
							},
							"explanation": "Ключ в репозитории = скомпрометированный ключ, его придётся перевыпускать.",
						},
						{
							"id":   "q4",
							"text": "Что значит «терминировать TLS на nginx»?",
							"options": []map[string]any{
								{"id": "a", "text": "nginx расшифровывает HTTPS, а к приложению идёт по внутренней сети", "correct": true},
								{"id": "b", "text": "nginx запрещает HTTPS-соединения", "correct": false},
								{"id": "c", "text": "nginx выдаёт сертификаты сам", "correct": false},
							},
							"explanation": "Приложение освобождается от шифрования, но участок до него должен быть внутренним.",
						},
						{
							"id":     "q5",
							"review": true,
							"text":   "Повторение: сайт отдаёт 502 Bad Gateway. Где причина?",
							"options": []map[string]any{
								{"id": "a", "text": "Приложение не отвечает: упало или не слушает свой порт", "correct": true},
								{"id": "b", "text": "Истёк сертификат", "correct": false},
								{"id": "c", "text": "Пользователь отправил некорректный запрос", "correct": false},
							},
							"explanation": "502 — прокси не смог получить ответ от приложения. 500 — приложение ответило ошибкой.",
						},
					},
					"resources": []map[string]any{
						{
							"title": "MDN — что такое HTTPS",
							"url":   "https://developer.mozilla.org/ru/docs/Glossary/HTTPS",
							"note":  "короткое объяснение на русском",
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
						{
							"id":     "q6",
							"review": true,
							"text":   "Повторение: как отменить коммит, который уже отправлен в общую ветку?",
							"options": []map[string]any{
								{"id": "a", "text": "git revert — он создаст коммит-отмену", "correct": true},
								{"id": "b", "text": "git reset --hard и push --force", "correct": false},
								{"id": "c", "text": "Удалить ветку", "correct": false},
							},
							"explanation": "В общей ветке историю не переписывают.",
						},
					},
				},
			},
		},
	}
}
