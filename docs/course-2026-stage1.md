# Этап 1 — гигиена: отчёт об изменениях (сентябрь 2026)

Файл курса: `devops-engineer.course.2026-stage1.json` (platforma-course v1, все id сохранены и совпадают с `devops-engineer.course (5).json`). Резервная копия seed до этапа — `backup_stage1/` (18 файлов). Правки внесены в seed-контент платформы (источник правды), применены в БД через `seedcourse -force` без потери прогресса студентов; экспорт регенерирован `exportcourse -ids` (см. «Инструменты»).

## Резюме

- Версии инструментов переведены на сентябрь 2026 по таблице промпта: **146 замен**; ни одного значения из колонки «Было» в файле не осталось.
- Только `docker compose` и `compose.yaml`; legacy-упоминание — одной фразой в M9.L15.
- Ссылки: **31 URL** заменён на первоисточник, **12** честно переподписаны «Википедия: …»; неверно подписанных Википедий — 0.
- 4.2: капстоун 18.10 (`truncate -s 0` вместо `rm`, `: > файл` как второй вариант, работа с Deployment через `deploy/app`, архитектура стенда nginx→k3s); Dockerfile 9.5 (`package-lock.json`, `npm ci --omit=dev`, `USER node`) согласованно с текстом 9.7; GitHub Actions 10.7/10.9/10.10 (формат SHA-пина в тексте, в решениях — теги); 2.8 `ubuntu-24.04`; карта 18.12.
- 4.4: Helm 4 (15.18), Ingress NGINX → Gateway API (15.7, 15.8), Prometheus 3.x + Zabbix/VictoriaMetrics (16.10), зеркало HashiCorp и OpenTofu (13.1), Docker Hub из РФ (9.13), DORA (1.6), инциденты 2025 (17.4).
- Гигиена схемы: удалены 2 лишних ключа `task` в квизах M6.L6 и M6.L11 (не рендерились).
- Итог: уроков 248 → 248, минут 2894 → 2894, по типам {'text': 100, 'quiz': 61, 'terminal': 38, 'code': 49} — новых уроков нет (по правилам этапа).

## Изменённые уроки

| Урок | Название | Тип | Что изменено |
|---|---|---|---|
| M18.L10 | Капстоун: находим и чиним намеренные поломки | terminal | t12: rm заменён на truncate -s 0 (альтернатива `: > файл`), в hints/success/intro/debug объяснено, почему rm открытого файла не освобождает место и что системное решение — logrotate. t2–t3 переведены на Deployment (describe deploy app, logs deploy/app), t1 дополнительно принимает `kubectl get pods -l app=app`; в intro добавлен раздел «Как устроен стенд» (nginx — обратный прокси на хосте перед k3s, приложение — Deployment в k3s) и обновлён список команд; в debug добавлены разборы ошибок NotFound для пода и «удалил лог, а df 100%». (поля: debug, intro, tasks) |
| M9.L5 | Пишем Dockerfile: от простого к реальному приложению | code | COPY package.json package-lock.json ./, RUN npm ci --omit=dev, USER node перед CMD; task дополнен объяснением npm ci и «не root — норма с первого Dockerfile». Checks: +package-lock.json, +npm ci, +--omit=dev, +USER node, +порядок USER после npm ci; starter проходит только FROM и WORKDIR, solution — все 11. (поля: checks, hint, solution, starter, task) |
| M9.L7 | Инструкции Dockerfile и оптимизация образа | text | В multi-stage-примере на Node: COPY package.json package-lock.json ./ и RUN npm ci (без --omit=dev — с пояснением, что сборщик в devDependencies); две фразы про non-root (USER node уже есть в образе; «норма с первого Dockerfile»). Версии node:24, golang:1.27, alpine:3.22 уже соответствуют фактам, не трогал. (поля: body) |
| M9.L15 | Первый Compose: зачем, основные секции и реальный стек | code | В абзац про Docker Compose добавлены два коротких предложения: стандарт — плагин `docker compose` и compose.yaml, имя docker-compose.yml тоже работает, а команда docker-compose (V1, Python) вне поддержки с 2023 и упоминается лишь как legacy. Остальные поля (starter, hint, checks, solution) не тронуты — solution по-прежнему проходит все checks, starter — нет. (поля: task) |
| M9.L13 | Реестры образов: теги, push, pull и приватные реестры | terminal | Пункт 4.4 размещён здесь (урок про реестры и Docker Hub, в заголовке есть «приватные реестры»): после первого абзаца добавлен блок «Если Docker Hub недоступен из РФ» — 6 предложений про блокировки 2024/2026, registry-mirrors с фрагментом daemon.json в блоке кода, pull-through cache на registry:2, Yandex CR / Selectel / GitLab Registry, Harbor и Nexus. Задачи t1–t5, hints, expected, debug, challenge не менялись. M9.L14 не трогал. (поля: intro) |
| M10.L7 | Первый workflow: YAML, структура и прогон тестов | code | В task добавлен абзац про закрепление actions по полному SHA (формат uses: owner/action@<40-символьный-sha> # v5.0.0, инциденты 2025 с подменой тегов, «как бывает неправильно»); node-version в пункте 2 исправлен с 20 на 24 (согласовано со starter/solution, с датой EOL Node 20); одно длинное предложение «Словарь…» разбито на два. В hint добавлена короткая фраза про SHA. Checks и solution без изменений — solution проходит все checks. (поля: hint, task) |
| M10.L9 | Триггеры и CI на pull request | code | В task перед формулировкой задания добавлен абзац про закрепление actions по полному SHA коммита (формат с плейсхолдером, инциденты 2025 с подменой тегов, «как бывает неправильно»). В hint — короткая фраза про SHA с указанием оставить тег @v5. Checks/solution/starter не менялись — solution проходит все checks, starter нет. (поля: hint, task) |
| M10.L10 | Секреты, артефакты и кэш | code | В task добавлен абзац «Закрепление версий actions» в стиле существующего блока про needs (формат с SHA, инциденты 2025, связь с секретами, «как бывает неправильно»); длинное предложение-задание (358 знаков) переоформлено в список с дефисами, предложение про ключ кэша разбито на два. В hint — короткая фраза про SHA. Checks/solution/starter не менялись — solution проходит все checks, starter нет. (поля: hint, task) |
| M15.L18 | Helm: чарты, values и релизы | text | Добавлен абзац про Helm 4 (GA 12.11.2025: --atomic → --rollback-on-failure, --force → --force-replace, --wait требует RBAC watch, релизы Helm 3 совместимы) в раздел «Как это устроено внутри»; «Helm (версии 3)» → «версии 3 и 4». Пример bitnami/postgresql заменён на установку OCI-чарта с плейсхолдерами и явным --version, добавлена оговорка про bitnamilegacy / Bitnami Secure Images; тег nginx 1.27 → 1.30; две новые строки в «Частые ошибки» и одна в «Запомнить». Объём 9964 знака. (поля: body) |
| M15.L7 | Namespace и Service: изоляция и стабильный адрес | text | В раздел «Как это работает» после схемы типов Service добавлен подраздел «Взгляд вперёд: Ingress и Gateway API»: Ingress остаётся в API, контроллер Ingress NGINX выведен из поддержки в марте 2026 (заявление Steering/SRC 29.01.2026), стандарт — Gateway API (GatewayClass, Gateway, HTTPRoute), для миграции ingress2gateway 1.0 (03.2026). Объём 7 767 знаков. (поля: body) |
| M15.L8 | Service и Ingress: маршрутизация HTTP-трафика | code | В task добавлен абзац «Ingress в 2026 году» (Ingress в API остаётся, контроллер Ingress NGINX выведен из поддержки в марте 2026, стандарт — Gateway API, миграция через ingress2gateway 1.0). Первое предложение абзаца про TLS сокращено до ≤155 знаков без изменения смысла; starter, solution, hint, checks не тронуты, solution проходит все checks, starter — нет. (поля: task) |
| M16.L10 | Grafana, Loki, OpenTelemetry и алерты | text | В раздел «Как всё встаёт рядом: поток данных» добавлены два абзаца: Prometheus 3.x (2.x без поддержки с 03.12.2024, образ v3.14.0, новый UI, UTF-8 в именах метрик, native histograms, антипример со старым образом) и соседи — Zabbix и VictoriaMetrics; в шпаргалку добавлен один пункт про них. Попутно исправлена опечатка «каждый알ерт» → «каждый алерт» в «Частых ошибках». Объём 9243 знака. (поля: body) |
| M13.L1 | Зачем IaC и как устроен Terraform | text | В раздел «Что будет, если…» добавлен четвёртый сценарий про 403 из России: зеркало Yandex Cloud с блоком provider_installation/network_mirror в ~/.terraformrc (hcl) и абзац про OpenTofu (MPL после BSL 2023, совместимые команды, блокировка РФ с 28.08.2024). Добавлено по одному пункту в «Частые ошибки новичка» и «Запомнить». Объём 9838 знаков. (поля: body) |
| M1.L6 | Автоматизация и CI/CD | text | В раздел «Как это работает» после flow-схемы добавлен блок о четырёх метриках DORA (частота деплоев, lead time, change failure rate, time to restore/MTTR) с пояснением «скорость + надёжность смотрят вместе» и примером «как бывает неправильно»; в «Запомнить» добавлен пункт про DORA. Нормативных цифр нет. Объём: 7481 → 8396 знаков. (поля: body) |
| M2.L8 | Описываем нужный сервер в конфиге | code | ubuntu-24.04 согласован в task/anatomy/solution. Единственная несогласованность: урок и «сломанный пример» учат, что после двоеточия обязателен пробел, а regex `[ ]*` принимал `name:web-1`; ужесточил до `[ ]+` и уточнил сообщения. Проверено re.search: solution проходит все 6 проверок, starter — только notContains, вариант без пробелов теперь не проходит. (поля: checks) |
| M17.L4 | Безопасность образов и цепочки поставки: контейнеры, сканирование, SBOM, CI/CD | text | В раздел «Секреты и защита пайплайна» добавлен блок с кейсами 2025 года (tj-actions/changed-files CVE-2025-30066, Nx s1ngularity через pull_request_target, Shai-Hulud и Shai-Hulud 2.0) и пять выводов: пин actions по SHA, минимальные permissions, OIDC, Dependabot/Renovate, secret scanning. В «Запомнить» добавлен один пункт. Объём 9 089 знаков. (поля: body) |
| M18.L12 | Карта участка и тренды 2026: платформы, self-service и ИИ в эксплуатации | text | В схеме «Карта всего участка» строка «docker compose» заменена на «Docker Compose» (название инструмента); длина строки та же, выравнивание колонок сохранено. Остальной текст без изменений. (поля: body) |
| M6.L6, M6.L11 | Квизы главы Bash | quiz | удалён лишний ключ `task` (гигиена схемы) |
| M9.L15 | Первый Compose | code | legacy-фраза про `docker-compose.yml` и `docker-compose` (V1) |
| M18.L10 | Капстоун (терминал) | terminal | контраст «logs по имени пода» перефразирован без литерала `kubectl logs app` |
| все главы | — | — | замены версий (таблица 4.1) и ссылок (4.3) |

## 4.1 Версии — таблица замен (вхождений в seed)

| Было | Стало | Вхождений |
|---|---|---|
| `node:20-slim` | `node:24-slim` | 7 |
| `node:20-alpine` | `node:24-alpine` | 5 |
| `node:20` | `node:24` | 3 |
| `node-version: "20"` | `node-version: "24"` | 2 |
| `prom/prometheus:v2.53.0` | `prom/prometheus:v3.14.0` | 5 |
| `grafana/grafana:11.1.0` | `grafana/grafana:13.2.0` | 3 |
| `nginx:1.27` | `nginx:1.30` | 24 |
| `nginx:1.28` | `nginx:1.30` | 4 |
| `postgres:16` | `postgres:17` | 17 |
| `golang:1.22` | `golang:1.27` | 2 |
| `alpine:3.20` | `alpine:3.22` | 3 |
| `python:3.12-slim` | `python:3.13-slim` | 4 |
| `python:3.12` | `python:3.13` | 1 |
| `busybox:1.36` | `busybox:1.37` | 2 |
| `redis:7` | `redis:8` | 2 |
| `Terraform v1.7.0` | `Terraform v1.16.0` | 1 |
| `required_version = ">= 1.6"` | `required_version = ">= 1.10"` | 2 |
| `Client Version: v1.30.0` | `Client Version: v1.37.0` | 1 |
| `ubuntu-22.04` | `ubuntu-24.04` | 5 |
| `actions/checkout@v4` | `actions/checkout@v5` | 15 |
| `actions/setup-node@v4` | `actions/setup-node@v5` | 3 |
| `docker-compose.yml` | `compose.yaml` | 20 |
| `docker-compose (команда)` | `docker compose` | 15 |
| **Итого** | | **146** |

Не найдено в курсе (замены не потребовались): `hashicorp/aws ~> 5.0`, `Helm 3` как строка. Оставлены `actions/upload-artifact@v4`, `actions/cache@v4` (мажор без сети не подтверждён — по правилу таблицы).

## 4.3 Ссылки — было → стало

### Заменён URL на первоисточник

| Урок | Подпись | Было | Стало | Источник |
|---|---|---|---|---|
| 001.L7 | YAML — официальный сайт | https://ru.wikipedia.org/wiki/YAML | https://yaml.org/spec/1.2.2/ | из промпта |
| 001.L8 | Terraform: введение в инфраструктуру как код | https://ru.wikipedia.org/wiki/Terraform | https://developer.hashicorp.com/terraform/intro | канонический (ручная проверка) |
| 003.L1 | Ubuntu Server documentation | https://ru.wikipedia.org/wiki/Filesystem_Hierarchy_Standard | https://man7.org/linux/man-pages/man7/hier.7.html | из промпта |
| 003.L2 | hier(7) — описание иерархии файловой системы Linux | https://ru.wikipedia.org/wiki/Filesystem_Hierarchy_Standard | https://man7.org/linux/man-pages/man7/hier.7.html | из промпта |
| 003.L2 | Filesystem Hierarchy Standard (Ubuntu) | https://ru.wikipedia.org/wiki/Filesystem_Hierarchy_Standard | https://man7.org/linux/man-pages/man7/hier.7.html | из промпта |
| 003.L4 | hier(7) — иерархия файловой системы | https://ru.wikipedia.org/wiki/Filesystem_Hierarchy_Standard | https://man7.org/linux/man-pages/man7/hier.7.html | из промпта |
| 003.L11 | chmod — изменение прав (man) | https://ru.wikipedia.org/wiki/Chmod | https://man7.org/linux/man-pages/man1/chmod.1.html | из промпта |
| 003.L12 | signal — man7.org | https://ru.wikipedia.org/wiki/Сигнал_(Unix) | https://man7.org/linux/man-pages/man7/signal.7.html | канонический (ручная проверка) |
| 003.L14 | chmod — man7.org | https://ru.wikipedia.org/wiki/Chmod | https://man7.org/linux/man-pages/man1/chmod.1.html | из промпта |
| 003.L17 | Права доступа Linux (chmod, man7.org) | https://ru.wikipedia.org/wiki/Chmod | https://man7.org/linux/man-pages/man1/chmod.1.html | из промпта |
| 003.L17 | Сигналы процессам (man7.org) | https://ru.wikipedia.org/wiki/Сигнал_(Unix) | https://man7.org/linux/man-pages/man7/signal.7.html | канонический (ручная проверка) |
| 004.L12 | MDN: Transport Layer Security (TLS) | https://ru.wikipedia.org/wiki/TLS | https://developer.mozilla.org/ru/docs/Web/Security/Transport_Layer_Security | из промпта |
| 006.L1 | chmod (man7.org) | https://ru.wikipedia.org/wiki/Chmod | https://man7.org/linux/man-pages/man1/chmod.1.html | из промпта |
| 006.L7 | chmod (man7.org) | https://ru.wikipedia.org/wiki/Chmod | https://man7.org/linux/man-pages/man1/chmod.1.html | из промпта |
| 006.L9 | man crontab.5 (формат файла) | https://ru.wikipedia.org/wiki/Cron | https://man7.org/linux/man-pages/man5/crontab.5.html | канонический (ручная проверка) |
| 006.L10 | man crontab.5 | https://ru.wikipedia.org/wiki/Cron | https://man7.org/linux/man-pages/man5/crontab.5.html | канонический (ручная проверка) |
| 006.L11 | man crontab.5 | https://ru.wikipedia.org/wiki/Cron | https://man7.org/linux/man-pages/man5/crontab.5.html | канонический (ручная проверка) |
| 007.L9 | GraphQL — официальный сайт | https://ru.wikipedia.org/wiki/GraphQL | https://graphql.org/learn/ | канонический (ручная проверка) |
| 007.L9 | gRPC — Introduction | https://ru.wikipedia.org/wiki/GRPC | https://grpc.io/docs/what-is-grpc/introduction/ | канонический (ручная проверка) |
| 007.L10 | Introducing JSON | https://ru.wikipedia.org/wiki/JSON | https://www.json.org/json-ru.html | канонический (ручная проверка) |
| 008.L2 | PostgreSQL: индексы | https://ru.wikipedia.org/wiki/Индекс_(базы_данных) | https://www.postgresql.org/docs/current/indexes.html | из промпта |
| 008.L3 | PostgreSQL: Indexes | https://ru.wikipedia.org/wiki/Индекс_(базы_данных) | https://www.postgresql.org/docs/current/indexes.html | из промпта |
| 008.L8 | PostgreSQL: высокая доступность и репликация | https://ru.wikipedia.org/wiki/Репликация_(вычислительная_техника) | https://www.postgresql.org/docs/current/high-availability.html | канонический (ручная проверка) |
| 008.L10 | PostgreSQL: Indexes | https://ru.wikipedia.org/wiki/Индекс_(базы_данных) | https://www.postgresql.org/docs/current/indexes.html | из промпта |
| 009.L1 | Docker overview | https://ru.wikipedia.org/wiki/Docker | https://docs.docker.com/get-started/docker-overview/ | из промпта |
| 011.L1 | MDN: Proxy servers and tunneling | https://ru.wikipedia.org/wiki/Прокси-сервер | https://developer.mozilla.org/en-US/docs/Web/HTTP/Proxy_servers_and_tunneling | канонический (ручная проверка) |
| 012.L10 | Отказоустойчивость и репликация (PostgreSQL) | https://ru.wikipedia.org/wiki/Репликация_(вычислительная_техника) | https://www.postgresql.org/docs/current/high-availability.html | канонический (ручная проверка) |
| 013.L1 | Terraform: введение и назначение | https://ru.wikipedia.org/wiki/Terraform | https://developer.hashicorp.com/terraform/intro | канонический (ручная проверка) |
| 013.L6 | Terraform: введение | https://ru.wikipedia.org/wiki/Terraform | https://developer.hashicorp.com/terraform/intro | канонический (ручная проверка) |
| 013.L13 | Terraform: обзор языка и рабочего процесса | https://ru.wikipedia.org/wiki/Terraform | https://developer.hashicorp.com/terraform/intro | канонический (ручная проверка) |
| 014.L10 | Ansible: введение (документация) | https://ru.wikipedia.org/wiki/Ansible | https://docs.ansible.com/ansible/latest/getting_started/introduction.html | канонический (ручная проверка) |

### Честно переподписано «Википедия: …» (URL сохранён; кандидаты на первоисточник — этап 2)

| Урок | Было | Стало | URL |
|---|---|---|---|
| 002.L1 | Аппаратное обеспечение | Википедия: Аппаратное обеспечение | https://ru.wikipedia.org/wiki/Аппаратное_обеспечение |
| 002.L3 | SSD: твердотельный накопитель | Википедия: SSD: твердотельный накопитель | https://ru.wikipedia.org/wiki/Твердотельный_накопитель |
| 002.L4 | Операционная система | Википедия: Операционная система | https://ru.wikipedia.org/wiki/Операционная_система |
| 002.L6 | Клиент — сервер | Википедия: Клиент — сервер | https://ru.wikipedia.org/wiki/Клиент_—_сервер |
| 002.L7 | Облачные вычисления | Википедия: Облачные вычисления | https://ru.wikipedia.org/wiki/Облачные_вычисления |
| 002.L8 | Инфраструктура как код | Википедия: Инфраструктура как код | https://ru.wikipedia.org/wiki/Инфраструктура_как_код |
| 003.L1 | The GNU Operating System | Википедия: The GNU Operating System | https://ru.wikipedia.org/wiki/GNU |
| 012.L1 | Облачные вычисления: модели IaaS/PaaS/SaaS | Википедия: Облачные вычисления: модели IaaS/PaaS/SaaS | https://ru.wikipedia.org/wiki/Облачные_вычисления |
| 012.L9 | DNS: как имя домена превращается в IP | Википедия: DNS: как имя домена превращается в IP | https://ru.wikipedia.org/wiki/DNS |
| 012.L12 | Многоуровневая архитектура веб-приложения | Википедия: Многоуровневая архитектура веб-приложения | https://ru.wikipedia.org/wiki/Многоуровневая_архитектура |
| 013.L1 | Сценарии использования IaC | Википедия: Сценарии использования IaC | https://ru.wikipedia.org/wiki/Инфраструктура_как_код |
| 018.L8 | CI/CD: непрерывная доставка | Википедия: CI/CD: непрерывная доставка | https://ru.wikipedia.org/wiki/CI/CD |

## Результаты `check_course.py` (экспорт, полный вывод)

```
Проверка: /private/tmp/claude-501/-Users-aaaaaakk12123gmail-com-Desktop-platforma/653c8f1a-a367-44df-940f-1e613bda2e89/scratchpad/stage1/devops-engineer.course.2026-stage1.json (export) — модулей 18, уроков 248, по типам {'text': 100, 'quiz': 61, 'terminal': 38, 'code': 49}

[1.схема] нарушений: 0

[2.id] нарушений: 0

[3.код] нарушений: 6
   - M1.L7: starter проходит ВСЕ checks
   - M3.L9: starter проходит ВСЕ checks
   - M6.L1: starter проходит ВСЕ checks
   - M6.L5: starter проходит ВСЕ checks
   - M6.L10: starter проходит ВСЕ checks
   - M13.L3: starter проходит ВСЕ checks

[4.терминал] нарушений: 0

[5.квиз] нарушений: 0

[6.текст] нарушений: 0

[7.устаревшее] нарушений: 0

[8.ссылки] нарушений: 1
   - M15.L23: «Kubernetes Roadmap» → roadmap.sh (подпись не соответствует домену)

[8.ссылки] всего 537, уникальных URL 354, URL с разными подписями: 81

[9.итог] durationMin по модулям:
     120 мин · 11 ур. · DevOps и мир доставки ПО
     103 мин · 10 ур. · Компьютер, ОС и серверы
     198 мин · 17 ур. · Linux с нуля
     210 мин · 18 ур. · Сети с нуля
     144 мин · 13 ур. · Git и контроль версий
     132 мин · 11 ур. · Bash и автоматизация
     130 мин · 12 ур. · Как устроены приложения
     124 мин · 10 ур. · Базы данных для DevOps
     218 мин · 18 ур. · Docker и контейнеры
     177 мин · 16 ур. · Тестирование и CI/CD
     134 мин · 11 ур. · Веб-серверы, обратный прокси, деплой
     139 мин · 12 ур. · Облако: чужая инфраструктура, которую вы арендуете
     151 мин · 13 ур. · Terraform: инфраструктура как код
     127 мин · 11 ур. · Ansible: автоматическая настройка серверов
     293 мин · 24 ур. · Kubernetes: оркестрация контейнеров
     133 мин · 12 ур. · Наблюдаемость и мониторинг
     188 мин · 16 ур. · Безопасность (DevSecOps) и надёжность (SRE)
     173 мин · 13 ур. · Проекты по нарастанию и капстоун
   ИТОГО: 2894 мин, уроков 248: {'text': 100, 'quiz': 61, 'terminal': 38, 'code': 49}

ИТОГО НАРУШЕНИЙ: 7
```

### Разбор оставшихся 7 нарушений

- **[3.код] 6 уроков (M1.L7, M3.L9, M6.L1, M6.L5, M6.L10, M13.L3) — ранее существовавшее.** Проверки там по ключевым словам (`for `, `awk`, `#!/…bash`), а starter-каркас уже содержит эти слова в комментариях/структуре — формально проходит все checks. Не входят в перечень 4.2; исправление требует ужесточить regex (привязать к коду, не к комментариям) или утоньшить starter — рекомендую на этап 2, чтобы не менять педагогику в гигиене.
- **[8.ссылки] M15.L23 «Kubernetes Roadmap» → roadmap.sh — намеренно.** Легитимный ресурс-дорожная карта; подпись соответствует сайту, а правило «Kubernetes → kubernetes.io» даёт лишь формальное несоответствие.

## Требует ручной проверки

- **Канонические URL не из промпта (9 адресов):** crontab(5) и signal(7) на man7.org; PostgreSQL high-availability; developer.hashicorp.com/terraform/intro; docs.ansible.com getting_started; graphql.org/learn; grpc.io introduction; json.org/json-ru; MDN Proxy servers — адреса стабильные и общеизвестные, но по правилу 5 отмечаю: сверить, что резолвятся (запуск `check_course.py --net`).
- **12 ссылок переподписаны «Википедия: …»** — URL сохранён; на этапе 2 подобрать первоисточники.
- **`actions/upload-artifact@v4`, `actions/cache@v4`** оставлены (мажор без сети не подтверждён).
- **Эмулятор терминала капстоуна (18.10):** `truncate -s 0 …` и `: > …` учебным эмулятором не исполняются (напечатает «команда не найдена»), но засчитываются сервером по точному совпадению — как и другие неэмулируемые команды. Добавить `truncate` в эмулятор — этап 2.
- **Метка `app=app` и имя деплоя `app` (18.10)** — сверить с манифестами модуля 18, когда появятся на этапах 2/3.
- **Helm 4 (15.18):** описание `--force-replace` дано по семантике `--force` Helm 3; `helm show values` с `oci://` и совместимость плагина helm-diff с Helm 4 — сверить по релиз-нотам. Объём body ≈ 9 960 знаков — у верхней границы.
- **Gateway API (15.7/15.8):** расшифровка SRC (Security Response Committee) и даты марта 2026 — сверить с заявлением от 29.01.2026; ссылки на gateway-api.sigs.k8s.io и ingress2gateway не добавлял (ресурсы по правилам не трогал) — добавить на этапе 2.
- **VictoriaMetrics (16.10):** «те же PromQL-запросы подходят» — MetricsQL совместим не на 100 %.
- **DORA (1.6):** расшифровка «DevOps Research and Assessment» — из общих знаний, не из таблицы фактов.
- **Shai-Hulud (17.4):** механизм «сам распространял себя» выведен из слова «червь»; конкретика (кража npm-токенов, автопубликация) в фактах не задана.
- **M2.L8:** regex ужесточены `[ ]*` → `[ ]+` (пробел после двоеточия обязателен — смысл урока); прежние ответы без пробела теперь неверны.
- **Загрузка бинарника Terraform из РФ (13.1):** URL зеркала релизов в фактах нет — не указан.
- **Длинные предложения:** 18 фрагментов >160 знаков — ранее существовавшие пункты списков и ответы «Проверьте себя», не новый текст; платформенный `checkcourse` проходит.

## Инструменты

- `check_course.py` — 9 проверок раздела 6; принимает экспорт-JSON или папку seed; `--net` проверяет HTTP-статусы.
- `cmd/exportcourse` пакует **встроенный seed** (БД не читает), поэтому сам по себе id не знает. Добавлен флаг `-ids <прошлый экспорт>`: переносит id курса/глав/уроков по позиции с проверкой совпадения заголовков и числа уроков (иначе отказ). Так экспорт остаётся импортируемым «на месте» без потери прогресса.

## Проверки платформы

- `go run ./cmd/checkcourse` — ✓ Замечаний нет.
- `go run ./cmd/seedcourse -force` — обновлено на месте, прогресс сохранён (+0 ~248 −0).

## Zip

`stage1_2026.zip`: `devops-engineer.course.2026-stage1.json`, `backup_stage1/`, `export_before.json` (без id, для сравнения объёмов), `check_course.py`, `check_seed.txt`, `check_export.txt`, `changes_stage1.md`, `links_changes.json`.

**Этап 1 завершён. Жду «дальше» для этапа 2 (новые уроки).**