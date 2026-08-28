# Деплой платформы на сервер okvionsales.ru

Платформа ставится **рядом** с sales-app на том же сервере (111.88.216.61),
на отдельный поддомен `learn.okvionsales.ru`. sales-app не трогаем.

- sales-app: домен `okvionsales.ru`, контейнеры на портах 3000 и 8081.
- платформа: поддомен `learn.okvionsales.ru`, контейнеры на портах 8090 и 8091.

Порты и базы у каждого проекта свои — конфликтов нет.
HTTPS на сервере уже выпускается через host-nginx + certbot (как у sales-app),
поэтому платформа встраивается в ту же схему.

---

## 1. DNS: поддомен

В панели домена `okvionsales.ru` добавить A-запись на IP сервера:

| Тип | Имя | Значение |
|-----|-----|----------|
| A | `learn` | `111.88.216.61` |

Дождаться обновления: `ping learn.okvionsales.ru` должен отдавать этот IP.
Без этого certbot не выпустит сертификат.

## 2. Залить код на сервер

```bash
# с локальной машины, из папки platforma
rsync -av --exclude node_modules --exclude .git --exclude 'backend/.env' \
  ./ root@111.88.216.61:/opt/platforma/
```

## 3. Настроить секреты

```bash
ssh root@111.88.216.61
cd /opt/platforma

# пароль базы для docker-compose
echo "POSTGRES_PASSWORD=$(openssl rand -hex 16)" > .env

# боевые настройки бэкенда
cp backend/.env.production.example backend/.env.production
nano backend/.env.production
```

В `backend/.env.production` заполнить:

- `JWT_SECRET` — `openssl rand -hex 32` (свой, не такой, как у других проектов);
- `PUBLIC_BASE_URL=https://learn.okvionsales.ru`
- `CORS_ORIGINS=https://learn.okvionsales.ru`
- `EMAILJS_*` — те же значения, что в локальном `backend/.env`.

## 4. Поднять контейнеры

```bash
docker compose -f docker-compose.prod.yml up -d --build
```

Проверка (изнутри сервера):

```bash
curl http://127.0.0.1:8090/health   # {"status":"ok"}
curl -I http://127.0.0.1:8091        # 200
```

Миграции применяются сами при старте бэкенда.

## 5. Подключить поддомен к host-nginx

```bash
sudo cp deploy/nginx-platforma-https.conf /etc/nginx/sites-available/learn.okvionsales.ru
sudo ln -s /etc/nginx/sites-available/learn.okvionsales.ru /etc/nginx/sites-enabled/

# сертификат для поддомена (webroot тот же, что у sales-app)
sudo certbot certonly --webroot -w /var/www/certbot -d learn.okvionsales.ru

sudo nginx -t && sudo systemctl reload nginx
```

Если `nginx -t` ругается на `http2 on;` — nginx старый: заменить эту строку
на `listen 443 ssl http2;` в двух местах конфига.

## 6. Создать администратора

```bash
docker compose -f docker-compose.prod.yml exec backend ./createadmin \
  -email admin@okvionsales.ru -name "Ахмед" -password "СИЛЬНЫЙ-ПАРОЛЬ"
```

## 7. Загрузить курс

```bash
docker compose -f docker-compose.prod.yml exec backend ./seedcourse
```

## Готово

Открыть `https://learn.okvionsales.ru`, войти под администратором,
создать студентов в разделе «Студенты».

---

## Обновление платформы потом

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
