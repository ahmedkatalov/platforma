# Деплой платформы на сервер

Платформа живёт рядом с cmf на одном сервере и не пересекается с ним:
свои контейнеры, своя база, свои порты (8090–8091; у cmf — 8080–8082).

## 1. Залить код

```bash
# на сервере
git clone <адрес-репозитория> ~/platforma
cd ~/platforma
```

Или скопировать папку целиком: `rsync -av --exclude node_modules --exclude .git ./ user@server:~/platforma/`

## 2. Настроить окружение

```bash
cd ~/platforma

# пароль базы для docker-compose
echo "POSTGRES_PASSWORD=$(openssl rand -hex 16)" > .env

# настройки бэкенда
cp backend/.env.production.example backend/.env.production
nano backend/.env.production
```

В `backend/.env.production` обязательно заполнить:

- `JWT_SECRET` — сгенерировать: `openssl rand -hex 32` (НЕ тот же, что у cmf);
- `PUBLIC_BASE_URL` и `CORS_ORIGINS` — домен платформы;
- ключи `EMAILJS_*` — те же значения, что в локальном `backend/.env`.

## 3. Запустить

```bash
docker compose -f docker-compose.prod.yml up -d --build
```

Проверка:

```bash
curl http://127.0.0.1:8090/health     # {"status":"ok"}
curl -I http://127.0.0.1:8091         # 200
```

Миграции применяются сами при старте бэкенда.

## 4. Создать администратора

```bash
docker compose -f docker-compose.prod.yml exec backend ./createadmin \
  -email admin@example.com -name "Ахмед" -password "СИЛЬНЫЙ-ПАРОЛЬ"
```

## 5. Загрузить курс

```bash
docker compose -f docker-compose.prod.yml exec backend ./seedcourse
```

## 6. Домен на хостовом nginx

Вариант А — поддомен (рекомендую): `learn.ваш-домен.ru`.

```nginx
server {
    listen 80;
    server_name learn.example.ru;

    location / {
        proxy_pass http://127.0.0.1:8091;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Файл кладут в `/etc/nginx/sites-available/`, включают симлинком в `sites-enabled/`, затем:

```bash
sudo nginx -t && sudo nginx -s reload
sudo certbot --nginx -d learn.example.ru   # HTTPS, certbot сам продлевает
```

`/api` отдельно проксировать не нужно: nginx внутри контейнера фронтенда
сам передаёт его бэкенду.

Не забудьте A-запись поддомена на IP сервера у регистратора домена.

## Обновление

```bash
cd ~/platforma
git pull
docker compose -f docker-compose.prod.yml up -d --build
```

База и её данные при этом не трогаются (том `platforma_pgdata`).

## Резервная копия базы

```bash
docker compose -f docker-compose.prod.yml exec db \
  pg_dump -U platforma platforma > backup-$(date +%F).sql
```

Восстановление: `cat backup.sql | docker compose -f docker-compose.prod.yml exec -T db psql -U platforma platforma`

## Если что-то не так

```bash
docker compose -f docker-compose.prod.yml ps            # статусы контейнеров
docker compose -f docker-compose.prod.yml logs backend  # логи API
docker compose -f docker-compose.prod.yml logs db       # логи базы
```
