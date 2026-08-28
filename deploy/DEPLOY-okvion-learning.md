# Деплой платформы на okvion-learning.ru

Платформа ставится на сервер 111.88.216.61 (там же живёт sales-app),
но на своём домене `okvion-learning.ru` и в своих контейнерах.
sales-app не трогаем.

- sales-app: домен `okvionsales.ru`, порты 3000 и 8081.
- платформа: домен `okvion-learning.ru`, порты 8090 и 8091.

Порты и базы у каждого свои — конфликтов нет.
HTTPS на сервере выпускается через host-nginx + certbot (как у sales-app).

---

## 1. DNS

У регистратора домена `okvion-learning.ru` создать A-записи на IP сервера:

| Тип | Имя | Значение |
|-----|-----|----------|
| A | `@` | `111.88.216.61` |
| A | `www` | `111.88.216.61` |

Дождаться обновления: `ping okvion-learning.ru` должен отдавать этот IP.
Без этого certbot не выпустит сертификат.

## 2. Залить код на сервер

```bash
# с локальной машины, из папки platforma
rsync -av --exclude node_modules --exclude .git --exclude 'backend/.env' \
  ./ root@111.88.216.61:/opt/platforma/
```

`backend/.env.production` уже подготовлен локально (с доменом и ключами почты)
и уедет вместе с остальными файлами — отдельно настраивать не нужно.

Проверить, что секреты на месте: `docker-compose.prod.yml` читает пароль базы
из файла `.env` рядом с собой — он тоже создаётся ниже.

## 3. Секрет базы

```bash
ssh root@111.88.216.61
cd /opt/platforma
echo "POSTGRES_PASSWORD=$(openssl rand -hex 16)" > .env
```

(Если `backend/.env.production` не залился — скопировать из примера и заполнить
`JWT_SECRET`, `PUBLIC_BASE_URL=https://okvion-learning.ru`,
`CORS_ORIGINS=https://okvion-learning.ru` и ключи `EMAILJS_*`.)

## 4. Поднять контейнеры

```bash
docker compose -f docker-compose.prod.yml up -d --build
```

Проверка изнутри сервера:

```bash
curl http://127.0.0.1:8090/health   # {"status":"ok"}
curl -I http://127.0.0.1:8091        # 200
```

Миграции применяются сами при старте бэкенда.

## 5. Домен в host-nginx

```bash
sudo cp deploy/nginx-okvion-learning-https.conf \
  /etc/nginx/sites-available/okvion-learning.ru
sudo ln -s /etc/nginx/sites-available/okvion-learning.ru /etc/nginx/sites-enabled/

sudo certbot certonly --webroot -w /var/www/certbot \
  -d okvion-learning.ru -d www.okvion-learning.ru

sudo nginx -t && sudo systemctl reload nginx
```

Если `nginx -t` ругается на `http2 on;` — nginx старый: заменить строку
на `listen 443 ssl http2;` во всех трёх server-блоках.

## 6. Администратор

```bash
docker compose -f docker-compose.prod.yml exec backend ./createadmin \
  -email admin@okvion-learning.ru -name "Ахмед" -password "СИЛЬНЫЙ-ПАРОЛЬ"
```

## 7. Курс

```bash
docker compose -f docker-compose.prod.yml exec backend ./seedcourse
```

## Готово

`https://okvion-learning.ru` → войти админом → создать студентов.

---

## Обновление

```bash
cd /opt/platforma
git pull            # либо повторный rsync
docker compose -f docker-compose.prod.yml up -d --build
```

## Бэкап базы

```bash
docker compose -f docker-compose.prod.yml exec db \
  pg_dump -U platforma platforma > ~/platforma-$(date +%F).sql
```
