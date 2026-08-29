package seed

func moduleLinux() ModuleSeed {
	return ModuleSeed{
		Title:   "Linux и командная строка",
		Summary: "База профессии: файлы, права, процессы и логи — с нуля и на практике",
		Lessons: []LessonSeed{
			{
				Title:       "Терминал: первое знакомство",
				Kind:        "text",
				Summary:     "Что такое командная строка и как в ней не потеряться",
				DurationMin: 10,
				Content: map[string]any{
					"body": "## Зачем это нужно\n\n" +
						"Серверы работают без мышки и окон. Всё управление — текстовыми командами.\n\n" +
						"Это кажется неудобным ровно до того момента, пока не нужно сделать одно и то же " +
						"на сотне серверов. Команду можно записать в файл и повторить. Клики мышкой — нет.\n\n" +
						"## Как читать приглашение\n\n" +
						"Строка, куда вы вводите команды, выглядит так:\n\n" +
						"```\n" +
						"student@devops:~$\n" +
						"```\n\n" +
						"- `student` — кто вы;\n" +
						"- `devops` — имя сервера;\n" +
						"- `~` — где вы находитесь (`~` означает домашнюю папку);\n" +
						"- `$` — конец приглашения, дальше пишете команду.\n\n" +
						"## Файлы и папки\n\n" +
						"В Linux всё начинается от корня — `/`. Внутри папки, внутри них другие папки.\n\n" +
						"Путь бывает двух видов:\n\n" +
						"- **Абсолютный** — от корня: `/var/log/app.log`. Работает откуда угодно.\n" +
						"- **Относительный** — от текущей папки: `log/app.log`.\n\n" +
						"Два полезных сокращения: `.` — текущая папка, `..` — папка выше.\n\n" +
						"## Куда что кладут\n\n" +
						"| Папка | Что внутри |\n" +
						"|---|---|\n" +
						"| `/home/student` | ваши личные файлы |\n" +
						"| `/etc` | настройки программ |\n" +
						"| `/var/log` | логи — записи о том, что происходило |\n" +
						"| `/tmp` | временные файлы, чистятся сами |\n" +
						"| `/usr/bin` | сами программы |\n\n" +
						"Запомнить просто: **настройки в `/etc`, логи в `/var/log`.** Это две папки, " +
						"куда вы будете заходить чаще всего.\n\n" +
						"## Первые пять команд\n\n" +
						"```bash\n" +
						"pwd              # где я сейчас\n" +
						"ls               # что здесь лежит\n" +
						"cd /var/log      # перейти в папку\n" +
						"cat app.log      # показать содержимое файла\n" +
						"cd ..            # подняться на уровень выше\n" +
						"```\n\n" +
						"Этих пяти команд хватит, чтобы осмотреться на незнакомом сервере.\n\n" +
						"А теперь та же навигация как живая сессия. Заметьте, как меняется приглашение:\n" +
						"\n" +
						"```\n" +
						"student@devops:~$ pwd\n" +
						"/home/student\n" +
						"student@devops:~$ ls\n" +
						"notes.txt  projects  photo.jpg\n" +
						"student@devops:~$ cd projects\n" +
						"student@devops:~/projects$ pwd\n" +
						"/home/student/projects\n" +
						"```\n" +
						"\n" +
						"После `cd` приглашение стало `~/projects` — оно всегда показывает, где вы сейчас.\n" +
						"\n" +
						"Реальная сессия почти всегда начинается с пары ошибок. Вот как они выглядят и читаются:\n" +
						"\n" +
						"```\n" +
						"student@devops:~$ cd /var/logs\n" +
						"-bash: cd: /var/logs: No such file or directory\n" +
						"student@devops:~$ cd /var/log\n" +
						"student@devops:/var/log$ cat app.log\n" +
						"cat: app.log: Permission denied\n" +
						"student@devops:/var/log$ ls -l app.log\n" +
						"-rw-r----- 1 root adm 48213 Aug 29 09:14 app.log\n" +
						"```\n" +
						"\n" +
						"Две ошибки подряд. Первая: `/var/logs` не существует — опечатка, папка называется `/var/log`.\n" +
						"\n" +
						"Вторая: `Permission denied`. Файл принадлежит `root`, группе `adm`, а мы не в ней.\n" +
						"\n" +
						"`ls -l` это подтверждает: у «остальных» нет даже `r`. Прочитать лог сможет только владелец или root.\n" +
						"\n" +
						"## Частые ошибки новичка\n\n" +
						"- **Забыли, где находитесь.** Наберите `pwd` — это ничего не сломает.\n" +
						"- **Команда не найдена.** Проверьте раскладку и опечатки: Linux различает большие и маленькие буквы.\n" +
						"- **Боитесь что-то сломать.** `pwd`, `ls`, `cat` только читают. Опасны команды, которые удаляют.\n\n" +
						"## Запомнить\n\n" +
						"1. `pwd` показывает, где вы. `ls` — что вокруг.\n" +
						"2. Настройки живут в `/etc`, логи — в `/var/log`.\n" +
						"3. Linux различает регистр букв: `App.log` и `app.log` — разные файлы.",
					"resources": []map[string]any{
						{
							"title": "The Linux Command Line — книга для начинающих",
							"url":   "https://linuxcommand.org/tlcl.php",
							"note":  "бесплатно, с нуля, самый мягкий вход в терминал",
						},
						{
							"title": "ExplainShell — разбор команды по частям",
							"url":   "https://explainshell.com/",
							"note":  "вставьте непонятную команду и увидите, что делает каждый флаг",
						},
						{
							"title": "GNU Coreutils — руководство",
							"url":   "https://www.gnu.org/software/coreutils/manual/coreutils.html",
							"note":  "официальные описания pwd, ls, cat и других базовых команд",
						},
						{
							"title": "Filesystem Hierarchy Standard (FHS 3.0)",
							"url":   "https://refspecs.linuxfoundation.org/FHS_3.0/fhs/index.html",
							"note":  "стандарт, объясняющий назначение /etc, /var, /tmp, /usr",
						},
						{
							"title": "ls — официальное описание",
							"url":   "https://man7.org/linux/man-pages/man1/ls.1.html",
							"note":  "все флаги ls, включая -l и -a",
						},
					},
				},
			},
			{
				Title:       "Тренажёр: осматриваемся на сервере",
				Kind:        "terminal",
				Summary:     "Пять базовых команд на практике",
				DurationMin: 15,
				Content: map[string]any{
					"intro": "Вы подключились к учебному серверу. Осмотритесь: где вы, что лежит вокруг, что внутри файлов.",
					"shell": "student@devops",
					"tasks": []map[string]any{
						{
							"id":       "l1",
							"prompt":   "Узнайте, в какой папке вы находитесь",
							"expected": []string{"pwd"},
							"hint":     "Три буквы: print working directory",
							"success":  "Вы в /home/student — это ваша домашняя папка.",
						},
						{
							"id":       "l2",
							"prompt":   "Посмотрите, какие файлы лежат в текущей папке",
							"expected": []string{"ls"},
							"hint":     "Две буквы: list",
							"success":  "Видно файлы и папки. Пока без подробностей.",
						},
						{
							"id":       "l3",
							"prompt":   "Покажите те же файлы подробно и вместе со скрытыми",
							"expected": []string{"ls -la", "ls -al", "ls -a -l"},
							"hint":     "К ls добавьте флаги -l (подробно) и -a (все)",
							"success":  "Теперь видны права, владелец, размер и дата. Файлы с точки в начале — скрытые.",
						},
						{
							"id":       "l4",
							"prompt":   "Прочитайте файл notes.txt",
							"expected": []string{"cat notes.txt", "cat ./notes.txt"},
							"hint":     "cat и имя файла",
							"success":  "cat выводит файл целиком — годится для коротких файлов.",
						},
						{
							"id":       "l5",
							"prompt":   "Перейдите в папку /var/log",
							"expected": []string{"cd /var/log"},
							"hint":     "cd и полный путь от корня",
							"success":  "Вы в папке с логами. Приглашение изменилось — оно всегда показывает, где вы.",
						},
						{
							"id":       "l6",
							"prompt":   "Покажите последние 10 строк файла app.log",
							"expected": []string{"tail app.log", "tail -n 10 app.log"},
							"hint":     "tail показывает конец файла",
							"success":  "Свежие записи всегда в конце — поэтому логи читают с хвоста.",
						},
						{
							"id":       "l7",
							"prompt":   "Вернитесь в домашнюю папку",
							"expected": []string{"cd", "cd ~", "cd /home/student"},
							"hint":     "cd без аргументов возвращает домой",
							"success":  "Готово. Вы освоили базовую навигацию.",
						},
					},
					"resources": []map[string]any{
						{
							"title": "Шпаргалка по командам Linux",
							"url":   "https://man7.org/linux/man-pages/dir_section_1.html",
							"note":  "официальные описания всех команд из этого урока",
						},
					},
				},
			},
			{
				Title:       "Квиз: навигация и файлы",
				Kind:        "quiz",
				Summary:     "Проверим команды из первой темы",
				DurationMin: 6,
				Content: map[string]any{
					"passScore": 70,
					"questions": []map[string]any{
						{
							"id":   "n1",
							"text": "Какая команда показывает, в какой папке вы находитесь?",
							"options": []map[string]any{
								{"id": "a", "text": "pwd", "correct": true},
								{"id": "b", "text": "ls", "correct": false},
								{"id": "c", "text": "cat", "correct": false},
							},
							"explanation": "ls показывает содержимое папки, cat — содержимое файла.",
						},
						{
							"id":   "n2",
							"text": "Что делает команда cd .. ?",
							"options": []map[string]any{
								{"id": "a", "text": "Поднимает на одну папку вверх", "correct": true},
								{"id": "b", "text": "Возвращает в домашнюю папку", "correct": false},
								{"id": "c", "text": "Показывает скрытые файлы", "correct": false},
							},
							"explanation": "Домой возвращает cd без аргументов или cd ~.",
						},
						{
							"id":   "n3",
							"text": "Где по традиции лежат логи?",
							"options": []map[string]any{
								{"id": "a", "text": "/var/log", "correct": true},
								{"id": "b", "text": "/etc", "correct": false},
								{"id": "c", "text": "/tmp", "correct": false},
							},
							"explanation": "В /etc лежат настройки, в /tmp — временные файлы.",
						},
						{
							"id":   "n4",
							"text": "Файл лога занимает 200 тысяч строк. Чем его открыть, чтобы посмотреть свежие записи?",
							"options": []map[string]any{
								{"id": "a", "text": "tail — покажет конец файла", "correct": true},
								{"id": "b", "text": "cat — выведет весь файл", "correct": false},
								{"id": "c", "text": "head — покажет начало", "correct": false},
							},
							"explanation": "Свежие записи всегда в конце, поэтому логи читают с хвоста.",
						},
						{
							"id":       "n5",
							"text":     "Что показывают флаги в команде ls -la?",
							"multiple": true,
							"options": []map[string]any{
								{"id": "a", "text": "-l выводит подробности: права, владельца, размер", "correct": true},
								{"id": "b", "text": "-a показывает скрытые файлы", "correct": true},
								{"id": "c", "text": "-l сортирует файлы по размеру", "correct": false},
							},
							"explanation": "Скрытые файлы начинаются с точки: .env, .gitignore.",
						},
						{
							"id":   "n6",
							"text": "app.log и App.log — это один файл или разные?",
							"options": []map[string]any{
								{"id": "a", "text": "Разные: Linux различает большие и маленькие буквы", "correct": true},
								{"id": "b", "text": "Один и тот же файл", "correct": false},
								{"id": "c", "text": "Зависит от настроек терминала", "correct": false},
							},
							"explanation": "Частая причина ошибки «файл не найден» — регистр букв.",
						},
					},
					"resources": []map[string]any{
						{
							"title": "ExplainShell — разбор команды по частям",
							"url":   "https://explainshell.com/",
							"note":  "вставьте команду и увидите, что делает каждый флаг",
						},
					},
				},
			},
			{
				Title:       "Права доступа: кто что может",
				Kind:        "text",
				Summary:     "Как читать rwx и почему 777 — плохая идея",
				DurationMin: 12,
				Content: map[string]any{
					"body": "## Зачем это нужно\n\n" +
						"На сервере работает много программ и людей. Права решают, кто какой файл может " +
						"прочитать, изменить или запустить.\n\n" +
						"Половина проблем новичка на сервере — это «permission denied». Разберёмся, что это значит.\n\n" +
						"## Как читать строку прав\n\n" +
						"Команда `ls -l` показывает такое:\n\n" +
						"```\n" +
						"-rw-r--r--  1 student student  471 Aug 27 09:00 app.log\n" +
						"```\n\n" +
						"Первый символ — тип: `-` файл, `d` папка.\n\n" +
						"Дальше три группы по три символа:\n\n" +
						"```\n" +
						"rw-   r--   r--\n" +
						" |     |     |\n" +
						" |     |     └── остальные: только читать\n" +
						" |     └──────── группа: только читать\n" +
						" └────────────── владелец: читать и писать\n" +
						"```\n\n" +
						"Буквы означают:\n\n" +
						"- `r` — читать (read);\n" +
						"- `w` — изменять (write);\n" +
						"- `x` — запускать (execute);\n" +
						"- `-` — права нет.\n\n" +
						"## Цифры вместо букв\n\n" +
						"Те же права записывают числами. Каждая буква — своё число:\n\n" +
						"| Право | Число |\n" +
						"|---|---|\n" +
						"| `r` читать | 4 |\n" +
						"| `w` писать | 2 |\n" +
						"| `x` запускать | 1 |\n\n" +
						"Числа складывают. Получается одна цифра на группу:\n\n" +
						"- `7` = 4+2+1 = читать, писать, запускать;\n" +
						"- `6` = 4+2 = читать и писать;\n" +
						"- `4` = только читать.\n\n" +
						"Три цифры — три группы: владелец, группа, остальные.\n\n" +
						"```bash\n" +
						"chmod 644 app.log     # владелец пишет, остальные читают\n" +
						"chmod 600 secret.key  # только владелец, больше никто\n" +
						"chmod 755 deploy.sh   # все могут запускать, менять — владелец\n" +
						"```\n\n" +
						"Посмотрите, как chmod меняет права. Слева от файла — та самая строка rwx:\n" +
						"\n" +
						"```\n" +
						"student@devops:~$ ls -l secret.key\n" +
						"-rw-r--r-- 1 student student 1675 Aug 27 09:00 secret.key\n" +
						"student@devops:~$ chmod 600 secret.key\n" +
						"student@devops:~$ ls -l secret.key\n" +
						"-rw------- 1 student student 1675 Aug 27 09:00 secret.key\n" +
						"```\n" +
						"\n" +
						"До chmod остальные могли читать файл, после `600` — только владелец.\n" +
						"\n" +
						"## Почему 777 — плохо\n\n" +
						"`chmod 777` даёт всем полные права. Файл сможет изменить любая программа на сервере, " +
						"включая ту, которую взломали.\n\n" +
						"Совет из практики: если хочется поставить 777, чтобы «наконец заработало», " +
						"почти всегда правильный ответ — сменить владельца файла через `chown`.\n\n" +
						"## Три режима, которые покрывают почти всё\n\n" +
						"| Режим | Для чего |\n" +
						"|---|---|\n" +
						"| `600` | приватные ключи, файлы с паролями |\n" +
						"| `644` | обычные файлы и конфиги |\n" +
						"| `755` | скрипты и папки |\n\n" +
						"Классика, на которой спотыкаются все: ssh отказывается брать ключ со слишком широкими правами.\n" +
						"\n" +
						"```\n" +
						"student@devops:~$ ssh -i id_rsa deploy@10.0.0.5\n" +
						"@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@\n" +
						"@    WARNING: UNPROTECTED PRIVATE KEY FILE!        @\n" +
						"@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@\n" +
						"Permissions 0644 for 'id_rsa' are too open.\n" +
						"This private key will be ignored.\n" +
						"student@devops:~$ chmod 600 id_rsa\n" +
						"student@devops:~$ ssh -i id_rsa deploy@10.0.0.5\n" +
						"deploy@web-1:~$\n" +
						"```\n" +
						"\n" +
						"Ключ читался всеми (`0644`), и ssh его проигнорировал — это защита, а не сбой.\n" +
						"\n" +
						"После `chmod 600` файл виден только владельцу, и вход проходит. Отсюда правило `600` для ключей.\n" +
						"\n" +
						"## Частые ошибки новичка\n\n" +
						"- **Ставят 777 при любой ошибке доступа.** Это не решение, а дыра в безопасности.\n" +
						"- **Забывают про x у папок.** Без `x` в папку нельзя войти, даже если есть `r`.\n" +
						"- **Правят права вместо владельца.** Иногда нужен `chown`, а не `chmod`.\n\n" +
						"Ещё одна частая ловушка — папка без права `x`. Внутрь не войти, хотя `r` на месте.\n" +
						"\n" +
						"```\n" +
						"student@devops:~$ ls -ld backups\n" +
						"drw-r--r-- 2 student student 4096 Aug 29 10:00 backups\n" +
						"student@devops:~$ cd backups\n" +
						"-bash: cd: backups: Permission denied\n" +
						"student@devops:~$ chmod 755 backups\n" +
						"student@devops:~$ cd backups\n" +
						"student@devops:~/backups$\n" +
						"```\n" +
						"\n" +
						"У папки права `rw-` без `x`. Первый символ `d` — это каталог, а не файл.\n" +
						"\n" +
						"`chmod 755` добавляет `x`, и `cd` срабатывает. Для папок `x` означает «можно войти внутрь».\n" +
						"\n" +
						"## Запомнить\n\n" +
						"1. Читать 4, писать 2, запускать 1 — складываем и получаем цифру.\n" +
						"2. Три цифры: владелец, группа, остальные.\n" +
						"3. `600` для секретов, `644` для файлов, `755` для скриптов. 777 — никогда.",
					"resources": []map[string]any{
						{
							"title": "chmod — официальное описание",
							"url":   "https://man7.org/linux/man-pages/man1/chmod.1.html",
							"note":  "числовые и буквенные режимы, если понадобятся тонкости",
						},
						{
							"title": "chown — смена владельца файла",
							"url":   "https://man7.org/linux/man-pages/man1/chown.1.html",
							"note":  "когда вместо chmod правильнее поменять владельца",
						},
						{
							"title": "GNU Coreutils — права на файлы",
							"url":   "https://www.gnu.org/software/coreutils/manual/html_node/File-permissions.html",
							"note":  "числовая и буквенная запись режимов рядом",
						},
						{
							"title": "Arch Wiki — права и атрибуты файлов",
							"url":   "https://wiki.archlinux.org/title/File_permissions_and_attributes",
							"note":  "подробный разбор rwx, прав на папки и особых битов",
						},
					},
				},
			},
			{
				Title:       "Квиз: права доступа",
				Kind:        "quiz",
				Summary:     "Числа, буквы и типичные режимы",
				DurationMin: 7,
				Content: map[string]any{
					"passScore": 70,
					"questions": []map[string]any{
						{
							"id":   "r1",
							"text": "Что означает буква w в строке прав?",
							"options": []map[string]any{
								{"id": "a", "text": "Право изменять файл", "correct": true},
								{"id": "b", "text": "Право читать файл", "correct": false},
								{"id": "c", "text": "Право запускать файл", "correct": false},
							},
							"explanation": "r — читать, w — писать, x — запускать.",
						},
						{
							"id":   "r2",
							"text": "Сколько будет 4 + 2 в правах и что это значит?",
							"options": []map[string]any{
								{"id": "a", "text": "6 — читать и писать", "correct": true},
								{"id": "b", "text": "6 — читать и запускать", "correct": false},
								{"id": "c", "text": "42 — полный доступ", "correct": false},
							},
							"explanation": "Читать 4, писать 2, запускать 1. Числа складываются.",
						},
						{
							"id":   "r3",
							"text": "Какой режим поставить приватному ключу ssh?",
							"options": []map[string]any{
								{"id": "a", "text": "600 — только владелец читает и пишет", "correct": true},
								{"id": "b", "text": "644 — остальные тоже читают", "correct": false},
								{"id": "c", "text": "777 — чтобы точно работало", "correct": false},
							},
							"explanation": "При более широких правах ssh откажется использовать ключ.",
						},
						{
							"id":   "r4",
							"text": "Скрипт не запускается: «permission denied». Что скорее всего не так?",
							"options": []map[string]any{
								{"id": "a", "text": "У файла нет права на запуск — нужен chmod +x или 755", "correct": true},
								{"id": "b", "text": "Файл слишком большой", "correct": false},
								{"id": "c", "text": "Неправильное имя файла", "correct": false},
							},
							"explanation": "Право x у скриптов обязательно.",
						},
						{
							"id":   "r5",
							"text": "Почему chmod 777 — плохое решение?",
							"options": []map[string]any{
								{"id": "a", "text": "Файл сможет изменить любая программа на сервере, включая взломанную", "correct": true},
								{"id": "b", "text": "Файл станет недоступен владельцу", "correct": false},
								{"id": "c", "text": "Такой режим не существует", "correct": false},
							},
							"explanation": "Чаще всего нужно поменять владельца через chown, а не раздавать права всем.",
						},
						{
							"id":     "r6",
							"review": true,
							"text":   "Повторение: какой командой посмотреть права на файлы в текущей папке?",
							"options": []map[string]any{
								{"id": "a", "text": "ls -l", "correct": true},
								{"id": "b", "text": "pwd", "correct": false},
								{"id": "c", "text": "cat", "correct": false},
							},
							"explanation": "Флаг -l выводит подробности, включая права.",
						},
					},
					"resources": []map[string]any{
						{
							"title": "chmod — официальное описание",
							"url":   "https://man7.org/linux/man-pages/man1/chmod.1.html",
							"note":  "числовые и буквенные режимы",
						},
					},
				},
			},
			{
				Title:       "Процессы, службы и логи",
				Kind:        "text",
				Summary:     "Кто занимает ресурсы и где искать причину сбоя",
				DurationMin: 13,
				Content: map[string]any{
					"body": "## Зачем это нужно\n\n" +
						"Сервис упал или тормозит. Нужно быстро понять причину. " +
						"Для этого хватает четырёх вещей: процессы, службы, логи, ресурсы.\n\n" +
						"## Процессы\n\n" +
						"Процесс — это запущенная программа. У каждой есть номер (PID) и владелец.\n\n" +
						"```bash\n" +
						"ps aux                # все запущенные программы\n" +
						"ps aux | grep nginx   # найти конкретную\n" +
						"kill 1043             # попросить программу завершиться\n" +
						"kill -9 1043          # убить принудительно\n" +
						"```\n\n" +
						"Вертикальная черта `|` называется конвейером. Она передаёт вывод одной команды " +
						"другой: «покажи все процессы» → «оставь только строки со словом nginx».\n\n" +
						"> Сначала всегда обычный `kill`. Он даёт программе закрыть файлы и доделать дела. " +
						"`kill -9` обрывает мгновенно — это крайняя мера.\n\n" +
						"## Службы\n\n" +
						"Программы, которые работают постоянно, запускает systemd. Он же перезапускает их " +
						"после падения.\n\n" +
						"```bash\n" +
						"systemctl status nginx    # работает или нет\n" +
						"systemctl restart nginx   # перезапустить\n" +
						"systemctl enable nginx    # запускать при загрузке сервера\n" +
						"```\n\n" +
						"Вот как выглядит ответ `systemctl status`:\n" +
						"\n" +
						"```\n" +
						"student@devops:~$ systemctl status nginx\n" +
						"● nginx.service - A high performance web server\n" +
						"     Loaded: loaded (/lib/systemd/system/nginx.service; enabled)\n" +
						"     Active: active (running) since Wed 2026-08-27 09:00:14; 2h ago\n" +
						"   Main PID: 1043 (nginx: master process)\n" +
						"```\n" +
						"\n" +
						"Главная строка — `Active`. `active (running)` — работает, `failed` — упала, `inactive (dead)` — остановлена.\n" +
						"\n" +
						"Важная деталь: `start` запускает сейчас, `enable` — включает автозапуск. " +
						"Это разные вещи. Забыли `enable` — после перезагрузки сервера служба не поднимется.\n\n" +
						"## Логи\n\n" +
						"Логи — записи о том, что происходило. Два основных места:\n\n" +
						"```bash\n" +
						"tail -f /var/log/app.log      # файл лога, следим за новыми строками\n" +
						"journalctl -u nginx -n 50     # журнал службы, последние 50 строк\n" +
						"```\n\n" +
						"Полезный приём — искать по слову:\n\n" +
						"```bash\n" +
						"grep ERROR /var/log/app.log            # найти строки с ошибками\n" +
						"grep ERROR /var/log/app.log | wc -l    # посчитать, сколько их\n" +
						"```\n\n" +
						"## Ресурсы\n\n" +
						"```bash\n" +
						"df -h      # сколько места на дисках\n" +
						"free -h    # сколько свободной памяти\n" +
						"uptime     # сколько работает сервер и какая нагрузка\n" +
						"```\n\n" +
						"Флаг `-h` означает «human readable» — покажет гигабайты вместо длинных чисел.\n\n" +
						"Вот как выглядит разбор упавшей службы: от статуса к журналу и к настоящей причине.\n" +
						"\n" +
						"```\n" +
						"student@devops:~$ systemctl status app\n" +
						"● app.service - Payment API\n" +
						"     Loaded: loaded (/etc/systemd/system/app.service; enabled)\n" +
						"     Active: failed (Result: exit-code) since Sat 2026-08-29 11:02:53; 20s ago\n" +
						"   Main PID: 2841 (code=exited, status=1/FAILURE)\n" +
						"student@devops:~$ journalctl -u app -n 3 --no-pager\n" +
						"Aug 29 11:02:53 devops app[2841]: Error: ENOENT: no such file or directory, open '/etc/app/config.yml'\n" +
						"Aug 29 11:02:53 devops systemd[1]: app.service: Main process exited, code=exited, status=1/FAILURE\n" +
						"Aug 29 11:02:53 devops systemd[1]: app.service: Failed with result 'exit-code'.\n" +
						"student@devops:~$ ls -l /etc/app/\n" +
						"total 0\n" +
						"```\n" +
						"\n" +
						"`Active: failed` и `status=1/FAILURE` — служба стартовала и тут же упала.\n" +
						"\n" +
						"`journalctl` показывает причину: `ENOENT` — программа не нашла файл `/etc/app/config.yml`.\n" +
						"\n" +
						"`ls` это подтверждает: папка `/etc/app` пустая. Чинить нужно конфиг, а не службу.\n" +
						"\n" +
						"## Порядок разбора сбоя\n\n" +
						"Запомните эту последовательность — она закрывает большинство случаев:\n\n" +
						"1. Служба вообще запущена? → `systemctl status`\n" +
						"2. Что в её журнале? → `journalctl -u`\n" +
						"3. Есть место на диске? → `df -h`\n" +
						"4. Хватает памяти? → `free -h`\n\n" +
						"Забитый диск — очень частая причина странных поломок. Программа не может записать " +
						"лог или временный файл и падает без понятной ошибки.\n\n" +
						"## Частые ошибки новичка\n\n" +
						"- **Сразу `kill -9`.** Программа не успевает закрыть файлы, данные теряются.\n" +
						"- **Читают начало лога вместо конца.** Свежие записи внизу: `tail`, а не `head`.\n" +
						"- **Проверяют всё, кроме диска.** `df -h` — первое, что стоит посмотреть.\n\n" +
						"И самая недооценённая причина поломок — забитый диск. Проверяется в две команды.\n" +
						"\n" +
						"```\n" +
						"student@devops:~$ df -h\n" +
						"Filesystem      Size  Used Avail Use% Mounted on\n" +
						"/dev/root        20G   20G  0.2G  99% /\n" +
						"tmpfs           2.0G  1.1M  2.0G   1% /run\n" +
						"student@devops:~$ du -sh /var/log/* | sort -h | tail -3\n" +
						"120M  /var/log/journal\n" +
						"1.8G  /var/log/nginx\n" +
						"8.9G  /var/log/app.log\n" +
						"```\n" +
						"\n" +
						"`Use% 99%` и `Avail 0.2G` — корневой раздел почти полон, писать уже некуда.\n" +
						"\n" +
						"`du` с сортировкой находит виновника: `app.log` разросся до `8.9G`. Его и чистим первым.\n" +
						"\n" +
						"## Запомнить\n\n" +
						"1. `ps aux` — что запущено, `systemctl status` — как чувствует себя служба.\n" +
						"2. Логи читают с конца: `tail`, `journalctl -n`.\n" +
						"3. При непонятной поломке первым делом проверьте `df -h`.",
					"resources": []map[string]any{
						{
							"title": "systemd — справочник по службам",
							"url":   "https://www.freedesktop.org/software/systemd/man/latest/",
							"note":  "systemctl, journalctl и формат файлов служб",
						},
						{
							"title": "Метод USE — быстрый поиск узкого места",
							"url":   "https://www.brendangregg.com/usemethod.html",
							"note":  "когда захочется разбираться в производительности глубже",
						},
						{
							"title": "journalctl — чтение системного журнала",
							"url":   "https://www.freedesktop.org/software/systemd/man/latest/journalctl.html",
							"note":  "фильтры -u, -n, -f и выборка по времени",
						},
						{
							"title": "ps — список процессов",
							"url":   "https://man7.org/linux/man-pages/man1/ps.1.html",
							"note":  "форматы вывода и что означают колонки ps aux",
						},
						{
							"title": "signal — сигналы процессам",
							"url":   "https://man7.org/linux/man-pages/man7/signal.7.html",
							"note":  "разница между SIGTERM (kill) и SIGKILL (kill -9)",
						},
					},
				},
			},
			{
				Title:       "Тренажёр: права, поиск и диагностика",
				Kind:        "terminal",
				Summary:     "Разбираемся с правами и ищем причину ошибки в логах",
				DurationMin: 18,
				Content: map[string]any{
					"intro": "Теперь задачи ближе к реальным: посмотреть права, найти ошибки в логе, проверить ресурсы сервера.",
					"shell": "student@devops",
					"tasks": []map[string]any{
						{
							"id":       "p1",
							"prompt":   "Посмотрите подробный список файлов, чтобы увидеть права",
							"expected": []string{"ls -l", "ls -la", "ls -al"},
							"hint":     "ls с флагом -l",
							"success":  "Первый столбец — те самые rwx.",
						},
						{
							"id":       "p2",
							"prompt":   "Сделайте файл deploy.sh запускаемым для всех (режим 755)",
							"expected": []string{"chmod 755 deploy.sh", "chmod 755 ./deploy.sh"},
							"hint":     "chmod 755 и имя файла",
							"success":  "Теперь скрипт можно запустить.",
						},
						{
							"id":       "p3",
							"prompt":   "Закройте доступ к notes.txt для всех, кроме владельца (режим 600)",
							"expected": []string{"chmod 600 notes.txt", "chmod 600 ./notes.txt"},
							"hint":     "chmod 600 и имя файла",
							"success":  "Так защищают файлы с чувствительными данными.",
						},
						{
							"id":       "p4",
							"prompt":   "Найдите строки со словом ERROR в файле app.log",
							"expected": []string{"grep ERROR app.log", "grep -i error app.log", "grep ERROR ./app.log"},
							"hint":     "grep слово файл",
							"success":  "Видно ошибки платёжного сервиса.",
						},
						{
							"id":     "p5",
							"prompt": "Посчитайте, сколько строк с ERROR в этом файле",
							"expected": []string{
								"grep ERROR app.log | wc -l",
								"grep -c ERROR app.log",
								"cat app.log | grep ERROR | wc -l",
							},
							"hint":    "Соедините grep и wc -l вертикальной чертой",
							"success": "Конвейер работает: одна команда передала вывод другой.",
						},
						{
							"id":       "p6",
							"prompt":   "Проверьте, сколько свободного места на дисках",
							"expected": []string{"df -h"},
							"hint":     "df с флагом -h",
							"success":  "Место есть — эту причину сбоя можно исключить.",
						},
						{
							"id":       "p7",
							"prompt":   "Проверьте состояние службы nginx",
							"expected": []string{"systemctl status nginx"},
							"hint":     "systemctl status и имя службы",
							"success":  "Служба активна. Диагностика закончена.",
						},
					},
					"resources": []map[string]any{
						{
							"title": "grep — поиск по файлам",
							"url":   "https://man7.org/linux/man-pages/man1/grep.1.html",
							"note":  "флаги -i, -c, -n и регулярные выражения",
						},
					},
				},
			},
			{
				Title:       "Практика: скрипт проверки сервера",
				Kind:        "code",
				Summary:     "Соберите короткий bash-скрипт для быстрой диагностики",
				DurationMin: 20,
				Content: map[string]any{
					"language": "bash",
					"task": "Допишите скрипт, который выводит состояние сервера. В нём должно быть:\n\n" +
						"1. первая строка `#!/bin/bash`;\n" +
						"2. проверка свободного места командой `df -h`;\n" +
						"3. проверка памяти командой `free -h`;\n" +
						"4. подсчёт строк с `ERROR` в `/var/log/app.log` через `grep` и `wc -l`;\n" +
						"5. проверка службы командой `systemctl status nginx`.",
					"starter": "# Скрипт быстрой проверки сервера\n" +
						"\n" +
						"echo \"=== Диск ===\"\n" +
						"\n" +
						"echo \"=== Память ===\"\n" +
						"\n" +
						"echo \"=== Ошибки в логе ===\"\n" +
						"\n" +
						"echo \"=== Служба nginx ===\"\n",
					"hint": "Первая строка скрипта называется shebang и выглядит так: #!/bin/bash",
					"checks": []map[string]any{
						{"type": "regex", "value": "^#!/bin/bash", "message": "Скрипт начинается с #!/bin/bash"},
						{"type": "regex", "value": "df\\s+-h", "message": "Проверяется свободное место"},
						{"type": "regex", "value": "free\\s+-h", "message": "Проверяется память"},
						{"type": "regex", "value": "grep\\s+(-\\w+\\s+)*ERROR", "message": "Ищутся ошибки в логе"},
						{"type": "regex", "value": "wc\\s+-l|grep\\s+-c", "message": "Ошибки подсчитываются"},
						{"type": "contains", "value": "systemctl status nginx", "message": "Проверяется служба nginx"},
					},
					"solution": "#!/bin/bash\n" +
						"# Скрипт быстрой проверки сервера\n" +
						"\n" +
						"echo \"=== Диск ===\"\n" +
						"df -h\n" +
						"\n" +
						"echo \"=== Память ===\"\n" +
						"free -h\n" +
						"\n" +
						"echo \"=== Ошибки в логе ===\"\n" +
						"grep ERROR /var/log/app.log | wc -l\n" +
						"\n" +
						"echo \"=== Служба nginx ===\"\n" +
						"systemctl status nginx\n",
					"resources": []map[string]any{
						{
							"title": "Bash Reference Manual",
							"url":   "https://www.gnu.org/software/bash/manual/bash.html",
							"note":  "переменные, условия и циклы, когда скрипт вырастет",
						},
						{
							"title": "ShellCheck — проверка скриптов на ошибки",
							"url":   "https://www.shellcheck.net/",
							"note":  "вставьте скрипт и увидите проблемы, которые легко пропустить",
						},
					},
				},
			},
			{
				Title:       "Проверка: Linux",
				Kind:        "quiz",
				Summary:     "Навигация, права, процессы и логи",
				DurationMin: 10,
				Content: map[string]any{
					"passScore": 70,
					"questions": []map[string]any{
						{
							"id":   "q1",
							"text": "Какая команда покажет, в какой папке вы находитесь?",
							"options": []map[string]any{
								{"id": "a", "text": "pwd", "correct": true},
								{"id": "b", "text": "ls", "correct": false},
								{"id": "c", "text": "cd", "correct": false},
							},
							"explanation": "pwd — print working directory, показывает текущий путь.",
						},
						{
							"id":   "q2",
							"text": "Где по традиции лежат настройки программ?",
							"options": []map[string]any{
								{"id": "a", "text": "/etc", "correct": true},
								{"id": "b", "text": "/tmp", "correct": false},
								{"id": "c", "text": "/usr/bin", "correct": false},
							},
							"explanation": "Настройки в /etc, логи в /var/log — это два самых частых адреса.",
						},
						{
							"id":   "q3",
							"text": "Что означает режим 644?",
							"options": []map[string]any{
								{"id": "a", "text": "Владелец читает и пишет, остальные только читают", "correct": true},
								{"id": "b", "text": "Все могут всё", "correct": false},
								{"id": "c", "text": "Файл можно только запускать", "correct": false},
							},
							"explanation": "6 = 4+2 (читать и писать), 4 = только читать.",
						},
						{
							"id":   "q4",
							"text": "Почему chmod 777 — плохая идея?",
							"options": []map[string]any{
								{"id": "a", "text": "Файл сможет изменить любая программа на сервере", "correct": true},
								{"id": "b", "text": "Файл станет недоступен владельцу", "correct": false},
								{"id": "c", "text": "Так файл занимает больше места", "correct": false},
							},
							"explanation": "Обычно проблема не в правах, а во владельце файла — помогает chown.",
						},
						{
							"id":   "q5",
							"text": "Какой командой удобно следить за новыми строками лога?",
							"options": []map[string]any{
								{"id": "a", "text": "tail -f app.log", "correct": true},
								{"id": "b", "text": "cat app.log", "correct": false},
								{"id": "c", "text": "head app.log", "correct": false},
							},
							"explanation": "tail показывает конец файла, флаг -f оставляет его открытым.",
						},
						{
							"id":   "q6",
							"text": "Что делает вертикальная черта в команде ps aux | grep nginx?",
							"options": []map[string]any{
								{"id": "a", "text": "Передаёт вывод первой команды на вход второй", "correct": true},
								{"id": "b", "text": "Запускает обе команды одновременно", "correct": false},
								{"id": "c", "text": "Сохраняет результат в файл", "correct": false},
							},
							"explanation": "Это конвейер: список процессов фильтруется по слову nginx.",
						},
						{
							"id":   "q7",
							"text": "Сервис ведёт себя странно, ошибка непонятная. Что проверить в первую очередь?",
							"options": []map[string]any{
								{"id": "a", "text": "Свободное место на диске командой df -h", "correct": true},
								{"id": "b", "text": "Версию ядра через uname", "correct": false},
								{"id": "c", "text": "Список пользователей", "correct": false},
							},
							"explanation": "Заполненный диск ломает запись логов и временных файлов — очень частая причина.",
						},
						{
							"id":     "q8",
							"review": true,
							"text":   "Повторение: что означает режим 600?",
							"options": []map[string]any{
								{"id": "a", "text": "Читать и писать может только владелец", "correct": true},
								{"id": "b", "text": "Все могут читать", "correct": false},
								{"id": "c", "text": "Файл можно запускать", "correct": false},
							},
							"explanation": "Так защищают ключи и файлы с паролями.",
						},
					},
					"resources": []map[string]any{
						{
							"title": "Linux Journey — интерактивный курс",
							"url":   "https://linuxjourney.com/",
							"note":  "если хочется закрепить тему отдельными упражнениями",
						},
					},
				},
			},
		},
	}
}
