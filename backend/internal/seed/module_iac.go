package seed

func moduleIaC() ModuleSeed {
	return ModuleSeed{
		Title:   "Инфраструктура как код",
		Summary: "Серверы, описанные файлами: Terraform создаёт, Ansible настраивает",
		Lessons: []LessonSeed{
			{
				Title:       "Зачем описывать серверы кодом",
				Kind:        "text",
				Summary:     "Terraform и Ansible: что делают и чем отличаются",
				DurationMin: 13,
				Content: map[string]any{
					"body": "## Зачем это нужно\n\n" +
						"Сервер настроили руками. Через полгода он сломался, и нужен такой же.\n\n" +
						"Никто не помнит, какие пакеты ставили, что правили в настройках, какие права выдавали. " +
						"Новый сервер получится похожим, но не таким же. Однажды это выстрелит.\n\n" +
						"Решение: описать сервер файлом. Файл лежит в Git, его читают, проверяют и применяют " +
						"сколько угодно раз с одинаковым результатом.\n\n" +
						"## Два инструмента, две задачи\n\n" +
						"| | Terraform | Ansible |\n" +
						"|---|---|---|\n" +
						"| Что делает | **создаёт** серверы, сети, базы | **настраивает** то, что уже создано |\n" +
						"| Аналогия | построить дом | завезти мебель |\n" +
						"| Формат | файлы `.tf` | файлы `.yml` |\n\n" +
						"Часто используют вместе: Terraform поднимает машины, Ansible ставит на них программы.\n\n" +
						"## Terraform: описываем результат\n\n" +
						"Вы не пишете «создай сервер». Вы описываете, что должно быть:\n\n" +
						"```hcl\n" +
						"resource \"aws_instance\" \"app\" {\n" +
						"  instance_type = \"t3.micro\"\n" +
						"\n" +
						"  tags = {\n" +
						"    Name = \"app-server\"\n" +
						"  }\n" +
						"}\n" +
						"```\n\n" +
						"Дальше Terraform сам сравнивает описание с реальностью и делает недостающее.\n\n" +
						"## Четыре команды\n\n" +
						"```bash\n" +
						"terraform init       # подготовить папку, скачать нужное\n" +
						"terraform validate   # проверить, нет ли ошибок в файле\n" +
						"terraform plan       # показать, что изменится\n" +
						"terraform apply      # применить\n" +
						"```\n\n" +
						"**`plan` перед `apply` — обязательная привычка.** План показывает три числа: " +
						"сколько создать, изменить и удалить.\n\n" +
						"```\n" +
						"Plan: 1 to add, 0 to change, 0 to destroy.\n" +
						"```\n\n" +
						"Если на боевой инфраструктуре видите `to destroy` — остановитесь и разберитесь. " +
						"Это удаление настоящего сервера.\n\n" +
						"## Состояние\n\n" +
						"Terraform запоминает, что он создал. Эта память называется **состоянием**.\n\n" +
						"В команде состояние хранят не на ноутбуке, а в общем месте с блокировкой. Иначе двое " +
						"применят изменения одновременно и затрут работу друг друга.\n\n" +
						"> Есть открытый форк Terraform — **OpenTofu**. Команды те же, файлы те же. " +
						"Если увидите `tofu plan` — это он.\n\n" +
						"## Ansible: настраиваем сервер\n\n" +
						"Ansible читает список задач и выполняет их на сервере:\n\n" +
						"```yaml\n" +
						"- name: Настроить веб-серверы\n" +
						"  hosts: web\n" +
						"  tasks:\n" +
						"    - name: Установить nginx\n" +
						"      apt:\n" +
						"        name: nginx\n" +
						"        state: present\n" +
						"```\n\n" +
						"Ключевое слово — `state: present`, то есть «должен быть установлен».\n\n" +
						"Запустите этот файл десять раз подряд — ничего не сломается. Если nginx уже стоит, " +
						"Ansible просто пропустит задачу. Это свойство называется **идемпотентность**: " +
						"повторный запуск не меняет результат.\n\n" +
						"Проверить, что изменится, без применения:\n\n" +
						"```bash\n" +
						"ansible-playbook playbook.yml --check\n" +
						"```\n\n" +
						"## Частые ошибки новичка\n\n" +
						"- **`apply` без `plan`.** Однажды удалите боевую базу.\n" +
						"- **Правят сервер руками поверх Terraform.** Следующий `apply` вернёт как было.\n" +
						"- **Пароли прямо в `.tf` файле.** Файл лежит в Git — это утечка.\n" +
						"- **Состояние на своём ноутбуке.** Коллега применит изменения и всё сломает.\n\n" +
						"## Запомнить\n\n" +
						"1. Terraform создаёт, Ansible настраивает.\n" +
						"2. Всегда `plan` перед `apply`, особенно смотрите на `to destroy`.\n" +
						"3. Источник правды — файлы, а не ручные правки на сервере.",
					"resources": []map[string]any{
						{
							"title": "Terraform — руководство для начинающих",
							"url":   "https://developer.hashicorp.com/terraform/tutorials",
							"note":  "пошаговые уроки с бесплатной песочницей",
						},
						{
							"title": "OpenTofu — документация",
							"url":   "https://opentofu.org/docs/",
							"note":  "открытый форк Terraform, команды совместимы",
						},
						{
							"title": "Ansible — документация",
							"url":   "https://docs.ansible.com/ansible/latest/",
							"note":  "модули, роли и примеры плейбуков",
						},
					},
				},
			},
			{
				Title:       "Тренажёр: Terraform и Ansible",
				Kind:        "terminal",
				Summary:     "Проведите изменение инфраструктуры от плана до применения",
				DurationMin: 20,
				Content: map[string]any{
					"resources": []map[string]any{
						{
							"title": "Terraform CLI — справочник команд",
							"url":   "https://developer.hashicorp.com/terraform/cli",
							"note":  "init, plan, apply, state и работа с рабочими пространствами",
						},
						{
							"title": "OpenTofu — документация",
							"url":   "https://opentofu.org/docs/",
							"note":  "открытый форк Terraform под лицензией MPL, команды совместимы",
						},
					},
					"intro": "В каталоге ~/infra лежит описание сервера, в ~/ansible — плейбук. Проведите изменение по всем шагам.",
					"shell": "student@devops",
					"tasks": []map[string]any{
						{
							"id":       "i1",
							"prompt":   "Перейдите в каталог infra",
							"expected": []string{"cd infra", "cd ~/infra", "cd /home/student/infra"},
							"hint":     "cd и имя каталога",
							"success":  "Вы в каталоге с описанием инфраструктуры.",
						},
						{
							"id":       "i2",
							"prompt":   "Посмотрите содержимое main.tf",
							"expected": []string{"cat main.tf"},
							"hint":     "cat и имя файла",
							"success":  "Ресурс aws_instance описан декларативно — Terraform сам приведёт реальность к нему.",
						},
						{
							"id":       "i3",
							"prompt":   "Инициализируйте рабочий каталог Terraform",
							"expected": []string{"terraform init"},
							"hint":     "terraform и подкоманда инициализации",
							"success":  "Провайдеры скачаны, можно работать.",
						},
						{
							"id":       "i4",
							"prompt":   "Проверьте, что конфигурация корректна",
							"expected": []string{"terraform validate"},
							"hint":     "terraform validate",
							"success":  "Синтаксис в порядке.",
						},
						{
							"id":       "i5",
							"prompt":   "Посмотрите план изменений",
							"expected": []string{"terraform plan"},
							"hint":     "terraform plan",
							"success":  "План показывает: 1 to add. Ничего не удаляется — применять безопасно.",
						},
						{
							"id":     "i6",
							"prompt": "Прогоните плейбук ~/ansible/playbook.yml в режиме проверки, без применения",
							"expected": []string{
								"ansible-playbook ~/ansible/playbook.yml --check",
								"ansible-playbook --check ~/ansible/playbook.yml",
								"ansible-playbook playbook.yml --check",
							},
							"hint":    "ansible-playbook, путь к файлу и флаг --check",
							"success": "Режим проверки показывает изменения, не трогая сервер.",
						},
					},
				},
			},
			{
				Title:       "Практика: описание сервера",
				Kind:        "code",
				Summary:     "Опишите ресурс с переменными и выводом",
				DurationMin: 22,
				Content: map[string]any{
					"resources": []map[string]any{
						{
							"title": "Язык конфигурации Terraform",
							"url":   "https://developer.hashicorp.com/terraform/language",
							"note":  "переменные, выражения, модули, зависимости между ресурсами",
						},
						{
							"title": "Terraform Style Guide",
							"url":   "https://developer.hashicorp.com/terraform/language/style",
							"note":  "официальные соглашения об именовании и структуре файлов",
						},
					},
					"language": "hcl",
					"task": "Допишите конфигурацию Terraform так, чтобы:\n\n" +
						"1. тип машины задавался переменной `var.instance_type`, а не строкой;\n" +
						"2. у ресурса был тег `Environment` со значением `production`;\n" +
						"3. был блок `variable \"instance_type\"` со значением по умолчанию `t3.micro`;\n" +
						"4. был `output` с адресом машины;\n" +
						"5. в конфигурации не осталось захардкоженных паролей (`password =`).",
					"starter": "resource \"aws_instance\" \"app\" {\n" +
						"  ami           = var.ami_id\n" +
						"  instance_type = \"t3.micro\"\n" +
						"\n" +
						"  tags = {\n" +
						"    Name = \"app-server\"\n" +
						"  }\n" +
						"}\n" +
						"\n" +
						"variable \"ami_id\" {\n" +
						"  type = string\n" +
						"}\n",
					"hint": "output описывается блоком output \"имя\" { value = ... }",
					"solution": "resource \"aws_instance\" \"app\" {\n" +
						"  ami           = var.ami_id\n" +
						"  instance_type = var.instance_type\n" +
						"\n" +
						"  tags = {\n" +
						"    Name        = \"app-server\"\n" +
						"    Environment = \"production\"\n" +
						"  }\n" +
						"}\n" +
						"\n" +
						"variable \"ami_id\" {\n" +
						"  type = string\n" +
						"}\n" +
						"\n" +
						"variable \"instance_type\" {\n" +
						"  type    = string\n" +
						"  default = \"t3.micro\"\n" +
						"}\n" +
						"\n" +
						"output \"app_ip\" {\n" +
						"  value = aws_instance.app.private_ip\n" +
						"}\n",
					"checks": []map[string]any{
						{"type": "regex", "value": "instance_type\\s*=\\s*var\\.instance_type", "message": "Тип машины вынесен в переменную"},
						{"type": "regex", "value": "Environment\\s*=\\s*\"production\"", "message": "Задан тег Environment"},
						{"type": "regex", "value": "variable\\s+\"instance_type\"", "message": "Объявлена переменная instance_type"},
						{"type": "regex", "value": "default\\s*=\\s*\"t3\\.micro\"", "message": "У переменной есть значение по умолчанию"},
						{"type": "regex", "value": "output\\s+\"[a-z_]+\"", "message": "Есть output с адресом машины"},
						{"type": "notContains", "value": "password =", "message": "В коде нет паролей"},
					},
				},
			},
			{
				Title:       "Проверка: инфраструктура как код",
				Kind:        "quiz",
				Summary:     "Состояние, план, идемпотентность",
				DurationMin: 10,
				Content: map[string]any{
					"resources": []map[string]any{
						{
							"title": "Ansible — документация",
							"url":   "https://docs.ansible.com/ansible/latest/",
							"note":  "модули, роли, инвентарь и переменные",
						},
						{
							"title": "Хранение состояния Terraform в удалённом бэкенде",
							"url":   "https://developer.hashicorp.com/terraform/language/backend",
							"note":  "блокировки и совместная работа над инфраструктурой",
						},
					},
					"passScore": 70,
					"questions": []map[string]any{
						{
							"id":   "q1",
							"text": "Зачем Terraform хранит состояние?",
							"options": []map[string]any{
								{"id": "a", "text": "Чтобы знать, какие ресурсы он создал, и вычислять разницу с описанием", "correct": true},
								{"id": "b", "text": "Чтобы ускорить скачивание провайдеров", "correct": false},
								{"id": "c", "text": "Чтобы хранить логи применения", "correct": false},
							},
							"explanation": "Из сравнения состояния и конфигурации рождается план изменений.",
						},
						{
							"id":   "q2",
							"text": "В плане вы видите строку «1 to destroy» на продовой инфраструктуре. Что делать?",
							"options": []map[string]any{
								{"id": "a", "text": "Остановиться и разобраться, какой ресурс и почему удаляется", "correct": true},
								{"id": "b", "text": "Применить: Terraform знает лучше", "correct": false},
								{"id": "c", "text": "Удалить файл состояния", "correct": false},
							},
							"explanation": "Удаление ресурса на проде — это простой. План для того и нужен, чтобы это заметить.",
						},
						{
							"id":   "q3",
							"text": "Что означает идемпотентность плейбука Ansible?",
							"options": []map[string]any{
								{"id": "a", "text": "Повторный запуск не меняет систему, если она уже в нужном состоянии", "correct": true},
								{"id": "b", "text": "Плейбук можно запускать только один раз", "correct": false},
								{"id": "c", "text": "Плейбук работает без подключения к серверу", "correct": false},
							},
							"explanation": "Задача сравнивает текущее состояние с желаемым и меняет только то, что нужно.",
						},
						{
							"id":       "q4",
							"text":     "Почему состояние Terraform не хранят локально в командной работе?",
							"multiple": true,
							"options": []map[string]any{
								{"id": "a", "text": "Двое инженеров затрут изменения друг друга", "correct": true},
								{"id": "b", "text": "Нет блокировок при параллельном применении", "correct": true},
								{"id": "c", "text": "В состоянии могут быть чувствительные данные, которым не место на ноутбуке", "correct": true},
								{"id": "d", "text": "Локальный файл работает медленнее", "correct": false},
							},
							"explanation": "Общий бэкенд с блокировкой — обязательное условие командной работы.",
						},
						{
							"id":   "q5",
							"text": "Вы вручную поменяли настройку на сервере, созданном через Terraform. Что произойдёт дальше?",
							"options": []map[string]any{
								{"id": "a", "text": "Следующий apply вернёт ресурс к описанному состоянию", "correct": true},
								{"id": "b", "text": "Terraform подхватит изменение в код", "correct": false},
								{"id": "c", "text": "Ничего, изменения сохранятся навсегда", "correct": false},
							},
							"explanation": "Источник правды — код. Ручные правки перезатираются, поэтому менять нужно описание.",
						},
					},
				},
			},
		},
	}
}
