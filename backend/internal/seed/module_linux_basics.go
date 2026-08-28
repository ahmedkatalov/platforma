package seed

func moduleLinux() ModuleSeed {
	return ModuleSeed{
		Title:   "Linux и командная строка",
		Summary: "Файловая система, права, процессы и базовые навыки работы в терминале",
		Lessons: []LessonSeed{
			{
				Title:       "Файловая система и навигация",
				Kind:        "text",
				Summary:     "Как устроены каталоги Linux и чем ходить по ним",
				DurationMin: 15,
				Content: map[string]any{
					"body": "## Дерево каталогов\n\n" +
						"В Linux нет дисков `C:` и `D:` — всё растёт из одного корня `/`.\n\n" +
						"| Каталог | Что внутри |\n" +
						"|---|---|\n" +
						"| `/etc` | конфигурационные файлы системы и сервисов |\n" +
						"| `/var` | изменяемые данные: логи, кэши, базы |\n" +
						"| `/home` | домашние каталоги пользователей |\n" +
						"| `/usr/bin` | исполняемые файлы программ |\n" +
						"| `/opt` | стороннее ПО, установленное отдельно |\n" +
						"| `/tmp` | временные файлы, очищаются при перезагрузке |\n\n" +
						"## Базовые команды\n\n" +
						"```bash\n" +
						"pwd              # где я сейчас\n" +
						"ls -la           # список файлов, включая скрытые, с правами\n" +
						"cd /var/log      # перейти в каталог\n" +
						"cd ..            # на уровень вверх\n" +
						"cat app.log      # показать файл целиком\n" +
						"less app.log     # листать большой файл\n" +
						"tail -f app.log  # следить за новыми строками\n" +
						"```\n\n" +
						"## Поиск\n\n" +
						"```bash\n" +
						"grep -i error app.log        # строки с error, без учёта регистра\n" +
						"grep -rn \"timeout\" /etc      # рекурсивно, с номерами строк\n" +
						"find /var/log -name \"*.log\"  # файлы по маске имени\n" +
						"```\n\n" +
						"## Права доступа\n\n" +
						"Строка `-rw-r--r--` читается по три символа: владелец, группа, остальные.\n\n" +
						"```bash\n" +
						"chmod +x deploy.sh        # сделать скрипт исполняемым\n" +
						"chmod 640 secrets.env     # владельцу чтение+запись, группе чтение\n" +
						"chown app:app /srv/app    # сменить владельца\n" +
						"```\n\n" +
						"Числа — это биты: 4 (чтение) + 2 (запись) + 1 (выполнение).\n\n" +
						"## Процессы\n\n" +
						"```bash\n" +
						"ps aux | grep nginx   # найти процесс\n" +
						"kill -9 1234          # снять процесс по PID\n" +
						"systemctl status nginx\n" +
						"journalctl -u nginx -f\n" +
						"```\n\n" +
						"> Совет: конвейер `|` передаёт вывод одной команды на вход другой. " +
						"Именно из таких цепочек состоит ежедневная работа в терминале.",
				},
			},
			{
				Title:       "Тренажёр: первые команды",
				Kind:        "terminal",
				Summary:     "Выполните задания в эмуляторе терминала",
				DurationMin: 20,
				Content: map[string]any{
					"intro": "Перед вами учебный сервер. Выполняйте задания по очереди — " +
						"команды проверяются на сервере платформы.",
					"shell": "student@devops",
					"tasks": []map[string]any{
						{
							"id":       "t1",
							"prompt":   "Узнайте, в каком каталоге вы находитесь",
							"expected": []string{"pwd"},
							"hint":     "Print working directory",
							"success":  "Верно, pwd показывает текущий каталог.",
						},
						{
							"id":       "t2",
							"prompt":   "Покажите все файлы текущего каталога, включая скрытые, с правами доступа",
							"expected": []string{"ls -la", "ls -al", "ls -l -a", "ls -a -l"},
							"hint":     "Нужны два флага: -l и -a",
							"success":  "Отлично, -l даёт подробности, -a показывает скрытые файлы.",
						},
						{
							"id":       "t3",
							"prompt":   "Перейдите в каталог /var/log",
							"expected": []string{"cd /var/log"},
							"hint":     "cd и полный путь",
							"success":  "Вы в каталоге логов.",
						},
						{
							"id":       "t4",
							"prompt":   "Выведите последние строки файла app.log и следите за новыми",
							"expected": []string{"tail -f app.log", "tail -f ./app.log"},
							"hint":     "Флаг -f означает follow",
							"success":  "Так смотрят живые логи сервиса.",
						},
						{
							"id":       "t5",
							"prompt":   "Найдите в app.log строки со словом error без учёта регистра",
							"expected": []string{"grep -i error app.log", "grep -i \"error\" app.log"},
							"hint":     "grep с флагом -i",
							"success":  "grep -i — самый частый способ найти проблему в логе.",
						},
						{
							"id":       "t6",
							"prompt":   "Сделайте скрипт deploy.sh исполняемым",
							"expected": []string{"chmod +x deploy.sh", "chmod 755 deploy.sh", "chmod u+x deploy.sh"},
							"hint":     "chmod и бит выполнения",
							"success":  "Теперь скрипт можно запустить как ./deploy.sh",
						},
						{
							"id":       "t7",
							"prompt":   "Найдите запущенный процесс nginx среди всех процессов",
							"pattern":  "^ps aux \\| grep -?i? ?nginx$",
							"expected": []string{"ps aux | grep nginx"},
							"hint":     "ps aux и конвейер в grep",
							"success":  "Классическая связка для поиска процесса.",
						},
					},
				},
			},
			{
				Title:       "Проверка: Linux",
				Kind:        "quiz",
				Summary:     "Каталоги, права и процессы",
				DurationMin: 8,
				Content: map[string]any{
					"passScore": 70,
					"questions": []map[string]any{
						{
							"id":   "q1",
							"text": "Где по соглашению лежат конфигурационные файлы сервисов?",
							"options": []map[string]any{
								{"id": "a", "text": "/etc", "correct": true},
								{"id": "b", "text": "/var", "correct": false},
								{"id": "c", "text": "/tmp", "correct": false},
								{"id": "d", "text": "/proc", "correct": false},
							},
							"explanation": "/etc — конфигурация, /var — изменяемые данные и логи.",
						},
						{
							"id":   "q2",
							"text": "Что делает команда chmod 640 secrets.env?",
							"options": []map[string]any{
								{"id": "a", "text": "Владельцу чтение и запись, группе чтение, остальным ничего", "correct": true},
								{"id": "b", "text": "Всем полный доступ", "correct": false},
								{"id": "c", "text": "Делает файл исполняемым", "correct": false},
							},
							"explanation": "6 = 4+2 (чтение+запись), 4 = чтение, 0 = нет прав.",
						},
						{
							"id":   "q3",
							"text": "Какой командой удобно следить за логом в реальном времени?",
							"options": []map[string]any{
								{"id": "a", "text": "cat app.log", "correct": false},
								{"id": "b", "text": "tail -f app.log", "correct": true},
								{"id": "c", "text": "head app.log", "correct": false},
							},
							"explanation": "tail -f продолжает печатать новые строки по мере их появления.",
						},
						{
							"id":       "q4",
							"text":     "Что делает конвейер в команде ps aux | grep nginx?",
							"multiple": true,
							"options": []map[string]any{
								{"id": "a", "text": "Передаёт вывод ps на вход grep", "correct": true},
								{"id": "b", "text": "Фильтрует строки со словом nginx", "correct": true},
								{"id": "c", "text": "Перезапускает nginx", "correct": false},
							},
							"explanation": "Символ | соединяет команды: вывод левой становится входом правой.",
						},
					},
				},
			},
		},
	}
}
