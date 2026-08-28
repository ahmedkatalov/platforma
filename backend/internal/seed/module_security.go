package seed

func moduleSecurity() ModuleSeed {
	return ModuleSeed{
		Title:   "Безопасность и секреты",
		Summary: "Доступы, секреты, образы и базовая гигиена продакшена",
		Lessons: []LessonSeed{
			{
				Title:       "Секреты, доступы и цепочка поставки",
				Kind:        "text",
				Summary:     "Где хранить пароли, как выдавать права и почему важны версии образов",
				DurationMin: 18,
				Content: map[string]any{
					"body": "## Секреты\n\n" +
						"Правило простое: **всё, что попало в Git, считается скомпрометированным**. " +
						"Даже если коммит удалить, он останется в истории и у всех, кто делал clone.\n\n" +
						"Где хранить:\n\n" +
						"- секреты CI (GitHub Secrets, GitLab CI Variables) — для пайплайнов;\n" +
						"- Secret в Kubernetes — для приложений в кластере;\n" +
						"- Vault или облачный менеджер секретов — когда нужны ротация и аудит.\n\n" +
						"В приложение секрет попадает переменной окружения или смонтированным файлом, " +
						"но не строкой в коде.\n\n" +
						"> Утёкший ключ не «удаляют» — его **отзывают и выпускают заново**. " +
						"Удаление коммита ничего не решает.\n\n" +
						"## Права доступа\n\n" +
						"Принцип наименьших привилегий: у сервиса ровно те права, что нужны для работы. " +
						"Учётка для деплоя не должна уметь удалять базу, а разработчику не нужен root на проде.\n\n" +
						"```bash\n" +
						"chmod 600 ~/.ssh/id_ed25519   # приватный ключ читает только владелец\n" +
						"chmod 644 config.yml          # конфиг читают все, пишет владелец\n" +
						"```\n\n" +
						"Права `777` — почти всегда ошибка: так файл может изменить любой процесс на машине.\n\n" +
						"## Контейнеры\n\n" +
						"- Не запускайте процесс от root: в Dockerfile создайте пользователя и укажите `USER`.\n" +
						"- Фиксируйте версии образов: `alpine:3.20`, а не `alpine:latest` — иначе сборка невоспроизводима.\n" +
						"- Не кладите секреты в слои образа: они остаются в истории слоёв даже после удаления файла.\n" +
						"- Сканируйте образы на уязвимости в пайплайне (trivy, grype).\n\n" +
						"## Сеть и доступ снаружи\n\n" +
						"Наружу торчит только то, что должно: обычно 80 и 443. База данных, панель мониторинга " +
						"и админские эндпоинты — во внутренней сети или за VPN.\n\n" +
						"```bash\n" +
						"ss -tulpn            # что реально слушает наружу\n" +
						"openssl s_client -connect app.example.com:443   # проверить TLS\n" +
						"```\n\n" +
						"## Гигиена, которая закрывает большинство дыр\n\n" +
						"| Практика | Что предотвращает |\n" +
						"|---|---|\n" +
						"| Секреты вне репозитория | утечку доступов |\n" +
						"| Обновление зависимостей | эксплуатацию известных уязвимостей |\n" +
						"| Наименьшие привилегии | превращение мелкого инцидента в крупный |\n" +
						"| Ротация ключей | долгую жизнь утёкшего ключа |\n" +
						"| Аудит и логи доступа | незаметное присутствие злоумышленника |\n\n" +
						"Безопасность — не отдельный этап в конце, а требование к каждому шагу конвейера.",
					"resources": []map[string]any{
						{"title": "OWASP Top 10", "url": "https://owasp.org/www-project-top-ten/", "note": "Самые частые классы уязвимостей веб-приложений"},
						{"title": "Trivy: сканер уязвимостей", "url": "https://trivy.dev/latest/docs/", "note": "Проверка образов, файлов IaC и зависимостей в пайплайне"},
						{"title": "Sigstore и cosign", "url": "https://docs.sigstore.dev/", "note": "Подпись артефактов без хранения приватных ключей"},
						{"title": "SLSA: уровни защиты сборки", "url": "https://slsa.dev/", "note": "Требования к цепочке поставки от сборки до выката"},
						{"title": "HashiCorp Vault", "url": "https://developer.hashicorp.com/vault/docs", "note": "Хранение секретов, динамические доступы и ротация"},
					},
				},
			},
			{
				Title:       "Цепочка поставки: SBOM, подписи и доступ без секретов",
				Kind:        "text",
				Summary:     "Как убедиться, что в проде работает именно то, что собрал ваш пайплайн",
				DurationMin: 18,
				Content: map[string]any{
					"body": "## Что такое цепочка поставки\n\n" +
						"Между коммитом и запуском в проде стоит длинная цепочка: зависимости, базовый образ, " +
						"сборщик, реестр, кластер. Атака на любое звено даёт злоумышленнику прод, " +
						"не трогая ваш исходный код. Отсюда три вопроса, на которые нужен ответ:\n\n" +
						"1. **Что внутри?** — состав артефакта.\n" +
						"2. **Кто собрал?** — подпись и происхождение.\n" +
						"3. **Из чего собрал?** — какой коммит и какой пайплайн.\n\n" +
						"## SBOM: состав артефакта\n\n" +
						"SBOM — список всех компонентов образа с версиями (форматы SPDX и CycloneDX). " +
						"Он нужен, чтобы за минуты ответить на вопрос «есть ли у нас уязвимая библиотека X».\n\n" +
						"```bash\n" +
						"syft registry.example.com/api:1.4.2 -o spdx-json > sbom.json\n" +
						"grype sbom:sbom.json           # проверить состав на известные уязвимости\n" +
						"trivy image registry.example.com/api:1.4.2\n" +
						"```\n\n" +
						"SBOM генерируется в пайплайне и хранится рядом с образом — тогда он всегда соответствует " +
						"тому, что реально собрано.\n\n" +
						"## Подпись и происхождение\n\n" +
						"Подпись отвечает на вопрос «этот образ действительно собрал наш пайплайн». " +
						"Sigstore и cosign делают это без хранения приватных ключей: подпись выпускается " +
						"на короткоживущий сертификат, привязанный к личности пайплайна.\n\n" +
						"```bash\n" +
						"cosign sign --yes registry.example.com/api:1.4.2\n" +
						"cosign verify registry.example.com/api:1.4.2 \\\n" +
						"  --certificate-identity-regexp 'https://github.com/team/api/.*' \\\n" +
						"  --certificate-oidc-issuer https://token.actions.githubusercontent.com\n" +
						"```\n\n" +
						"Кластер можно настроить так, чтобы он запускал только подписанные образы " +
						"(политики Kyverno или Gatekeeper). Тогда подсунуть свой образ в реестр недостаточно.\n\n" +
						"## SLSA: уровни зрелости\n\n" +
						"SLSA описывает, насколько сборке можно доверять: от «просто собрали где-то» " +
						"до «сборка воспроизводима, изолирована и подтверждена происхождением (provenance)». " +
						"Практический минимум для команды — собирать только в CI, подписывать артефакты " +
						"и хранить provenance вместе с образом.\n\n" +
						"## Доступ без долгоживущих секретов\n\n" +
						"Токен реестра или облака в секретах CI — актив, который живёт годами и утекает вместе " +
						"с любым логом. Современный способ — федерация через OIDC: пайплайн предъявляет облаку " +
						"короткоживущий токен своей личности и получает временные права.\n\n" +
						"```yaml\n" +
						"permissions:\n" +
						"  id-token: write     # разрешить пайплайну получить OIDC-токен\n" +
						"  contents: read\n" +
						"\n" +
						"steps:\n" +
						"  - uses: aws-actions/configure-aws-credentials@v4\n" +
						"    with:\n" +
						"      role-to-assume: arn:aws:iam::123456789012:role/ci-deploy\n" +
						"      aws-region: eu-central-1\n" +
						"```\n\n" +
						"Ключей в секретах нет вообще: доступ выдаётся на время одного запуска и только той ветке, " +
						"которая указана в доверенной политике роли.\n\n" +
						"## Зависимости\n\n" +
						"- Локфайлы (`go.sum`, `package-lock.json`) обязательны: без них сборка тянет разные версии.\n" +
						"- Обновления зависимостей автоматизируют (Renovate, Dependabot) — иначе их не делают вовсе.\n" +
						"- Публичные действия CI закрепляют по хешу коммита, а не по тегу: тег можно переписать.\n\n" +
						"```yaml\n" +
						"- uses: actions/checkout@11bd71901bbe5b1630ceea73d27597364c9af683  # v4.2.2\n" +
						"```\n\n" +
						"> Минимальный набор на 2026 год: сканирование образа и зависимостей в пайплайне, " +
						"SBOM рядом с артефактом, подпись через cosign, доступ по OIDC вместо статических ключей.",
					"resources": []map[string]any{
						{"title": "SLSA: уровни защиты сборки", "url": "https://slsa.dev/spec/v1.0/levels", "note": "Что конкретно требуется на каждом уровне"},
						{"title": "Sigstore: подпись без ключей", "url": "https://docs.sigstore.dev/cosign/signing/overview/", "note": "cosign, прозрачный журнал и проверка личности"},
						{"title": "CycloneDX", "url": "https://cyclonedx.org/", "note": "Формат SBOM, поддерживается большинством сканеров"},
						{"title": "Syft и Grype", "url": "https://github.com/anchore/syft", "note": "Генерация SBOM и проверка состава на уязвимости"},
						{"title": "OIDC в GitHub Actions", "url": "https://docs.github.com/actions/deployment/security-hardening-your-deployments/about-security-hardening-with-openid-connect", "note": "Как убрать статические облачные ключи из CI"},
						{"title": "NIST SSDF", "url": "https://csrc.nist.gov/Projects/ssdf", "note": "Свод практик безопасной разработки, на него ссылаются регуляторы"},
					},
				},
			},
			{
				Title:       "Тренажёр: проверка доступов",
				Kind:        "terminal",
				Summary:     "Наведите порядок в правах и проверьте, что открыто наружу",
				DurationMin: 18,
				Content: map[string]any{
					"resources": []map[string]any{
						{
							"title": "OpenSSH — руководства",
							"url":   "https://www.openssh.com/manual.html",
							"note":  "ключи, агент, конфигурация клиента и сервера",
						},
						{
							"title": "Trivy — сканер уязвимостей",
							"url":   "https://trivy.dev/latest/docs/",
							"note":  "проверка образов, файловой системы и IaC-конфигураций",
						},
					},
					"intro": "Перед выкатом проверьте базовую гигиену сервера.",
					"shell": "student@devops",
					"tasks": []map[string]any{
						{
							"id":       "s1",
							"prompt":   "Посмотрите, какие порты открыты и какие процессы их слушают",
							"expected": []string{"ss -tulpn", "ss -lntp", "netstat -tulpn"},
							"hint":     "ss с флагами -tulpn",
							"success":  "Наружу смотрят 80 и 8080 — с этим списком уже можно работать.",
						},
						{
							"id":       "s2",
							"prompt":   "Сгенерируйте ssh-ключ типа ed25519",
							"expected": []string{"ssh-keygen -t ed25519", "ssh-keygen -t ed25519 -C student@devops"},
							"hint":     "ssh-keygen -t тип",
							"success":  "Ключ создан. Приватную часть никому не передают.",
						},
						{
							"id":       "s3",
							"prompt":   "Ограничьте права приватного ключа ~/.ssh/id_ed25519 до 600",
							"expected": []string{"chmod 600 ~/.ssh/id_ed25519", "chmod 600 /home/student/.ssh/id_ed25519"},
							"hint":     "chmod 600 и путь к файлу",
							"success":  "Теперь ключ читает только владелец — иначе ssh откажется его использовать.",
						},
						{
							"id":       "s4",
							"prompt":   "Сгенерируйте случайный секрет: 32 байта в hex",
							"expected": []string{"openssl rand -hex 32"},
							"hint":     "openssl rand -hex длина",
							"success":  "Такой секрет годится для подписи токенов.",
						},
						{
							"id":     "s5",
							"prompt": "Проверьте, нет ли слова password в файле ~/projects/api/Dockerfile",
							"expected": []string{
								"grep password ~/projects/api/Dockerfile",
								"grep -i password ~/projects/api/Dockerfile",
								"grep password /home/student/projects/api/Dockerfile",
							},
							"hint":    "grep по файлу",
							"success": "Совпадений нет — секретов в образе не будет.",
						},
						{
							"id":     "s6",
							"prompt": "Убедитесь, что в Dockerfile указан непривилегированный пользователь: найдите строку USER",
							"expected": []string{
								"grep USER ~/projects/api/Dockerfile",
								"grep -n USER ~/projects/api/Dockerfile",
								"grep USER /home/student/projects/api/Dockerfile",
							},
							"hint":    "grep USER и путь к Dockerfile",
							"success": "USER app — процесс не работает от root.",
						},
					},
				},
			},
			{
				Title:       "Практика: безопасный Dockerfile",
				Kind:        "code",
				Summary:     "Приведите сборку образа в порядок",
				DurationMin: 22,
				Content: map[string]any{
					"resources": []map[string]any{
						{
							"title": "Безопасность Docker",
							"url":   "https://docs.docker.com/engine/security/",
							"note":  "изоляция, возможности ядра, rootless-режим",
						},
						{
							"title": "Distroless — образы без лишнего",
							"url":   "https://github.com/GoogleContainerTools/distroless",
							"note":  "минимальная поверхность атаки: ни шелла, ни пакетного менеджера",
						},
					},
					"language": "dockerfile",
					"task": "Исправьте Dockerfile так, чтобы:\n\n" +
						"1. базовый образ был с зафиксированной версией, без `latest`;\n" +
						"2. появился непривилегированный пользователь и инструкция `USER`;\n" +
						"3. пароль не задавался через `ENV`, а приходил снаружи;\n" +
						"4. использовалась многоэтапная сборка (`FROM ... AS build` и второй `FROM`);\n" +
						"5. был объявлен `EXPOSE` с портом приложения.",
					"starter": "FROM golang:latest\n" +
						"WORKDIR /src\n" +
						"COPY . .\n" +
						"ENV DB_PASSWORD=supersecret\n" +
						"RUN go build -o /app ./cmd/api\n" +
						"CMD [\"/app\"]\n",
					"hint": "Пользователь в alpine создаётся так: RUN adduser -D -u 10001 app",
					"solution": "FROM golang:1.25-alpine AS build\n" +
						"WORKDIR /src\n" +
						"COPY . .\n" +
						"RUN go build -o /app ./cmd/api\n" +
						"\n" +
						"FROM alpine:3.20\n" +
						"RUN adduser -D -u 10001 app\n" +
						"USER app\n" +
						"COPY --from=build /app /app\n" +
						"EXPOSE 8080\n" +
						"ENTRYPOINT [\"/app\"]\n",
					"checks": []map[string]any{
						{"type": "notContains", "value": "latest", "message": "Версии образов зафиксированы"},
						{"type": "regex", "value": "FROM\\s+\\S+\\s+AS\\s+\\w+", "message": "Используется многоэтапная сборка"},
						{"type": "regex", "value": "(?m)^FROM[\\s\\S]*^FROM", "message": "Есть второй этап с чистым образом"},
						{"type": "regex", "value": "(?m)^USER\\s+\\w+", "message": "Процесс работает не от root"},
						{"type": "notContains", "value": "ENV DB_PASSWORD", "message": "Пароль не зашит в образ"},
						{"type": "regex", "value": "(?m)^EXPOSE\\s+\\d+", "message": "Объявлен порт приложения"},
					},
				},
			},
			{
				Title:       "Проверка: безопасность",
				Kind:        "quiz",
				Summary:     "Секреты, права и образы",
				DurationMin: 10,
				Content: map[string]any{
					"resources": []map[string]any{
						{
							"title": "OWASP Top 10",
							"url":   "https://owasp.org/www-project-top-ten/",
							"note":  "базовый список рисков, который спрашивают на собеседованиях",
						},
						{
							"title": "Sigstore / Cosign — подпись артефактов",
							"url":   "https://docs.sigstore.dev/",
							"note":  "подпись образов без управления ключами вручную",
						},
						{
							"title": "HashiCorp Vault — документация",
							"url":   "https://developer.hashicorp.com/vault/docs",
							"note":  "хранение и ротация секретов, динамические учётные данные",
						},
					},
					"passScore": 70,
					"questions": []map[string]any{
						{
							"id":   "q1",
							"text": "Токен случайно попал в коммит и был запушен. Что делать?",
							"options": []map[string]any{
								{"id": "a", "text": "Отозвать токен и выпустить новый", "correct": true},
								{"id": "b", "text": "Удалить коммит — этого достаточно", "correct": false},
								{"id": "c", "text": "Сделать репозиторий приватным", "correct": false},
							},
							"explanation": "Секрет уже мог быть скопирован. Единственный надёжный шаг — отзыв.",
						},
						{
							"id":   "q2",
							"text": "Почему процесс в контейнере не запускают от root?",
							"options": []map[string]any{
								{"id": "a", "text": "Чтобы взлом приложения не дал сразу максимальных прав", "correct": true},
								{"id": "b", "text": "Контейнер не запустится от root", "correct": false},
								{"id": "c", "text": "Так образ становится меньше", "correct": false},
							},
							"explanation": "Наименьшие привилегии ограничивают ущерб от компрометации.",
						},
						{
							"id":       "q3",
							"text":     "Где допустимо хранить пароль базы для приложения в кластере?",
							"multiple": true,
							"options": []map[string]any{
								{"id": "a", "text": "В Secret Kubernetes", "correct": true},
								{"id": "b", "text": "В Vault с ротацией", "correct": true},
								{"id": "c", "text": "В ENV внутри Dockerfile", "correct": false},
								{"id": "d", "text": "В репозитории в файле config.yml", "correct": false},
							},
							"explanation": "Слои образа и репозиторий читают многие — секретам там не место.",
						},
						{
							"id":   "q4",
							"text": "Что означают права 600 на приватном ssh-ключе?",
							"options": []map[string]any{
								{"id": "a", "text": "Читать и писать может только владелец", "correct": true},
								{"id": "b", "text": "Читать может любой пользователь", "correct": false},
								{"id": "c", "text": "Файл становится исполняемым", "correct": false},
							},
							"explanation": "При более широких правах ssh просто откажется использовать ключ.",
						},
						{
							"id":   "q5",
							"text": "Почему alpine:latest — плохой выбор для продакшена?",
							"options": []map[string]any{
								{"id": "a", "text": "Сборка перестаёт быть воспроизводимой: сегодня и завтра это разные образы", "correct": true},
								{"id": "b", "text": "latest всегда содержит уязвимости", "correct": false},
								{"id": "c", "text": "latest медленнее скачивается", "correct": false},
							},
							"explanation": "Фиксированная версия даёт одинаковый результат сборки в любой момент.",
						},
					},
				},
			},
		},
	}
}
