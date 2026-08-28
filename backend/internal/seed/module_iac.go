package seed

func moduleIaC() ModuleSeed {
	return ModuleSeed{
		Title:   "Инфраструктура как код",
		Summary: "Terraform и Ansible: воспроизводимая инфраструктура вместо ручных настроек",
		Lessons: []LessonSeed{
			{
				Title:       "Зачем описывать инфраструктуру кодом",
				Kind:        "text",
				Summary:     "Декларативный и императивный подходы, состояние, идемпотентность",
				DurationMin: 17,
				Content: map[string]any{
					"body": "## Проблема ручной настройки\n\n" +
						"Сервер, настроенный руками, невозможно повторить. Через полгода никто не помнит, " +
						"какие пакеты ставили и что правили в конфигах. Второй такой же сервер " +
						"получится похожим, но не таким же — и однажды это выстрелит.\n\n" +
						"Инфраструктура как код (IaC) решает это просто: всё описано в файлах, файлы лежат в Git, " +
						"изменения проходят ревью, а применение автоматическое и повторяемое.\n\n" +
						"## Два инструмента, две задачи\n\n" +
						"| Инструмент | Что делает | Подход |\n" +
						"|---|---|---|\n" +
						"| **Terraform** | создаёт ресурсы: серверы, сети, базы, кластеры | декларативный |\n" +
						"| **Ansible** | настраивает то, что уже создано: пакеты, конфиги, службы | процедурный, но идемпотентный |\n\n" +
						"Их часто используют вместе: Terraform поднимает машины, Ansible доводит их до нужного состояния.\n\n" +
						"## Terraform: состояние\n\n" +
						"Terraform хранит состояние — какие ресурсы он создал и с какими параметрами. " +
						"Из сравнения состояния с описанием рождается план изменений.\n\n" +
						"```bash\n" +
						"terraform init      # скачать провайдеры\n" +
						"terraform validate  # проверить синтаксис\n" +
						"terraform plan      # показать, что изменится\n" +
						"terraform apply     # применить\n" +
						"```\n\n" +
						"> `plan` перед `apply` — обязательная привычка. План показывает, что будет создано, " +
						"изменено и **удалено**. Строчка `1 to destroy` на проде должна останавливать руку.\n\n" +
						"Состояние нельзя хранить локально в командной работе: его кладут в общий бэкенд " +
						"(S3, Terraform Cloud) с блокировками, иначе двое инженеров затрут работу друг друга.\n\n" +
						"## Ansible: идемпотентность\n\n" +
						"Плейбук описывает желаемое состояние машины. Повторный запуск ничего не ломает: " +
						"если nginx уже установлен, задача просто отметится как `ok`, а не `changed`.\n\n" +
						"```yaml\n" +
						"- name: Configure web servers\n" +
						"  hosts: web\n" +
						"  become: true\n" +
						"  tasks:\n" +
						"    - name: Install nginx\n" +
						"      apt:\n" +
						"        name: nginx\n" +
						"        state: present\n" +
						"\n" +
						"    - name: Start nginx service\n" +
						"      service:\n" +
						"        name: nginx\n" +
						"        state: started\n" +
						"        enabled: true\n" +
						"```\n\n" +
						"Запуск с `--check` показывает, что изменилось бы, ничего не применяя — аналог `terraform plan`.\n\n" +
						"## Правила, которые экономят нервы\n\n" +
						"- Никаких ручных правок ресурсов, созданных через IaC: следующий `apply` их вернёт.\n" +
						"- Версии провайдеров и ролей фиксируются, иначе сборка перестанет быть повторяемой.\n" +
						"- Секреты не хранятся в `.tf` и плейбуках — только переменные окружения или Vault.\n" +
						"- Один модуль — одна ответственность: сеть, база, приложение описываются отдельно." +
						"\n\n## Terraform и OpenTofu\n\n" +
						"В 2023 году HashiCorp сменила лицензию Terraform на BSL, и сообщество создало " +
						"открытый форк — **OpenTofu** под фондом Linux Foundation. Язык и структура " +
						"конфигураций совместимы: `tofu plan` работает там же, где `terraform plan`. " +
						"Многие компании перешли на OpenTofu из-за лицензии, поэтому в вакансиях " +
						"встречаются оба названия.\n\n" +
						"## Что ещё стало нормой\n\n" +
						"- Проверка конфигураций в пайплайне: `tflint`, `checkov` и `trivy config` " +
						"ловят открытые наружу порты и незашифрованные диски до применения.\n" +
						"- Модули берут из реестра с закреплением версии, а не копируют между проектами.\n" +
						"- План сохраняют в файл и применяют именно его: `terraform plan -out=tfplan` " +
						"и затем `terraform apply tfplan` — иначе между просмотром и применением " +
						"мир может измениться.\n" +
						"- Доступ к облаку пайплайн получает через OIDC, а не через статические ключи.\n\n" +
						"> Отдельные окружения описывают отдельными состояниями. Общий стейт для dev и prod — " +
						"верный способ однажды снести прод при экспериментах в dev.",
					"resources": []map[string]any{
						{"title": "Terraform: документация", "url": "https://developer.hashicorp.com/terraform/docs", "note": "Язык, провайдеры, состояние и модули"},
						{"title": "OpenTofu", "url": "https://opentofu.org/docs/", "note": "Открытый форк Terraform под лицензией MPL — часто выбирают вместо него"},
						{"title": "Ansible: документация", "url": "https://docs.ansible.com/ansible/latest/", "note": "Модули, роли, инвентарь и переменные"},
						{"title": "Как хранить состояние в S3", "url": "https://developer.hashicorp.com/terraform/language/backend/s3", "note": "Общий бэкенд с блокировкой — обязателен для команды"},
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
