# CSINDER

Приложение в стиле Tinder для поиска тиммейтов по CS2 на Faceit. Свайпай игроков, отфильтрованных по ELO, и получай уведомления в реальном времени, когда кто-то лайкнул тебя в ответ.

## Стек

- **Backend:** Go 1.24, net/http
- **База данных:** PostgreSQL 17
- **Frontend:** Vanilla HTML/CSS/JS, Go templates
- **Real-time:** WebSocket (github.com/coder/websocket)
- **Авторизация:** FACEIT OAuth2 + PKCE
- **Деплой:** Docker, Docker Compose

## Возможности

- Вход через FACEIT OAuth2
- Карточки игроков с ELO, K/D, винрейтом, % хедшотов, длиннейшей серией побед
- Свайп-интерфейс с поддержкой drag и touch жестов
- Фильтрация кандидатов по ELO (±300)
- WebSocket-уведомления о лайках в реальном времени
- Фоновая синхронизация статистики игроков с FACEIT API

## Структура проекта

```
cmd/csinder/          Точка входа
internal/
  auth/               Контекст авторизации
  config/             Конфигурация из env
  faceit/             FACEIT API клиент и OAuth
  handler/            HTTP и WebSocket хендлеры
  model/              Модели данных
  server/             Роутер, мидлвари
  service/            Бизнес-логика (свайпы, синхронизация)
  storage/            Запросы к PostgreSQL
  ws/                 WebSocket хаб и клиенты
web/
  static/css/         Стили
  static/js/          Логика свайпов, WS клиент
  templates/          Go HTML шаблоны
migrations/           SQL схема
```

## Быстрый старт

### Docker Compose

```bash
cp .env.example .env
# заполни FACEIT_CLIENT_ID, FACEIT_API_KEY, COOKIE_SECRET
docker compose up
```

Приложение запускается на http://localhost:8080, PostgreSQL на порту 5432. Миграции применяются автоматически.

### Локально

Требуется Go 1.24+ и запущенный PostgreSQL.

```bash
export DATABASE_URL="postgres://user:pass@localhost:5432/csinder?sslmode=disable"
export FACEIT_CLIENT_ID="..."
export FACEIT_API_KEY="..."
export COOKIE_SECRET="..."

go run ./cmd/csinder
```

## Переменные окружения

| Переменная           | Обязательна | По умолчанию                           | Описание                        |
| -------------------- | ----------- | -------------------------------------- | ------------------------------- |
| `DATABASE_URL`       | да          | —                                      | Строка подключения к PostgreSQL |
| `FACEIT_CLIENT_ID`   | да          | —                                      | OAuth2 client ID от FACEIT      |
| `FACEIT_API_KEY`     | да          | —                                      | Серверный API ключ FACEIT       |
| `FACEIT_REDIRECT_URI`| нет         | `http://localhost:8080/auth/callback`  | URL для OAuth callback          |
| `COOKIE_SECRET`      | нет         | `dev-secret`                           | Ключ шифрования сессионных кук  |
| `BASE_URL`           | нет         | `http://localhost:8080`                | Публичный URL приложения        |
| `PORT`               | нет         | `8080`                                 | Порт HTTP сервера               |

## API

### Авторизация
- `GET /auth/faceit` — начать OAuth логин
- `GET /auth/callback` — OAuth callback
- `POST /auth/logout` — выход

### Страницы (требуется авторизация)
- `GET /swipe` — интерфейс свайпов
- `GET /profile` — твоя статистика
- `GET /matches` — кто тебя лайкнул

### REST (требуется авторизация)
- `GET /api/candidates` — следующая пачка игроков
- `POST /api/swipe` — записать свайп `{target_id, liked}`
- `GET /api/notifications` — твои уведомления
- `POST /api/notifications/seen` — отметить прочитанными `{ids: [...]}`

### WebSocket
- `GET /ws` — уведомления о лайках в реальном времени
