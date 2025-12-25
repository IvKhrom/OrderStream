# OrderStream

Учебный проект в формате **mono-repo** с двумя Go-микросервисами и инфраструктурой (Kafka + Postgres + Redis).

### Архитектура и workflow

- **`api_service`** (HTTP):
  - принимает запросы,
  - сохраняет **сырые** данные заказа в своей Postgres (`postgres_api`),
  - публикует событие в Kafka **`orders.events`**,
  - читает ACK из Kafka **`orders.ack`** и (опционально) ждёт подтверждение,
  - для “полного результата” читает Redis по ключу `order_result:<order_id>` (результат кладёт worker).
- **`worker`**:
  - читает **`orders.events`**,
  - выполняет обработку/валидацию/upsert по логике из `@internal` (create/update/soft-delete + запрет апдейта удалённого),
  - пишет данные в свою Postgres (`postgres_worker`),
  - кладёт **полный результат** в Redis (`order_result:<order_id>`),
  - публикует ACK в Kafka **`orders.ack`**.

### Структура репозитория

- **Go workspace**: `go.work` (подключает два модуля)
- **`services/api_service/`**: модуль API сервиса
  - `services/api_service/api/` — `.proto` контракт
  - `services/api_service/cmd/api_service/` — точка входа (`main.go`)
  - `services/api_service/config/` — env-config
  - `services/api_service/internal/` — bootstrap/api/services/storage/pb
- **`services/worker/`**: модуль воркера
  - `services/worker/cmd/worker/` — точка входа (`main.go`)
  - `services/worker/config/` — env-config
  - `services/worker/internal/` — bootstrap/consumer/services/storage/models
- **`migrations/`** — SQL миграции (используются обоими Postgres в docker-compose)
- **`docker-compose.yml`** — инфраструктура + оба сервиса

### Быстрый старт (Docker Compose)

Поднять инфраструктуру и сервисы:

```powershell
docker compose up -d --build
```

Что откроется:
- **API**: `http://localhost:8080`
- **Swagger UI**: `http://localhost:8080/docs`
- **Kafka UI**: `http://localhost:8081`
- **Adminer**: `http://localhost:8082`
- **Redis UI (RedisInsight)**: `http://localhost:5540`

### Redis UI (RedisInsight):

Открой `http://localhost:5540` → **Add Redis database** и укажи:
- **Host**: `redis`
- **Port**: `6379`
- **Username/Password**: пусто (у нас Redis без пароля)

После подключения ключи результата обработки лежат в Redis по шаблону:
- `order_result:<order_id>`

Сбросить Postgres volume’ы и применить миграции заново:

```powershell
docker compose down -v
docker compose up -d --build
```

### Переменные окружения

Значения и дефолты см. в:
- `services/api_service/config/config.go`
- `services/worker/config/config.go`

Основные:
- **API**
  - `API_POSTGRES_DSN` (default: `postgres://postgres:upvel123@localhost:5433/orderstream_api?sslmode=disable`)
  - `API_PORT` (default: `8080`)
  - `KAFKA_BROKERS` (default: `localhost:29092`)
  - `REDIS_ADDR` (default: `localhost:6379`)
  - `ORDERS_EVENTS_TOPIC` (default: `orders.events`)
  - `ORDERS_ACK_TOPIC` (default: `orders.ack`)
- **Worker**
  - `WORKER_POSTGRES_DSN` (default: `postgres://postgres:upvel123@localhost:5434/orderstream_worker?sslmode=disable`)
  - `WORKER_GROUP` (default: `order-workers`)
  - `KAFKA_BROKERS`, `REDIS_ADDR`, `ORDERS_EVENTS_TOPIC`, `ORDERS_ACK_TOPIC` — аналогично

### HTTP API

- **Health**: `GET /health` → `{ "status": "ok" }`

- **Upsert**: `POST /orders`
  - create: `order_id` пустой или `"0"`
  - soft-delete: `"status": "deleted"`
  - payload можно передавать:
    - как `payload_json` (сырое json),
    - или как `payload` (объект) — сервер сам сериализует.

- **Get by id**: `GET /orders/{order_id}`
- **Get by external id**: `GET /orders/by-external/{external_id}?user_id=<uuid>`
  - `external_id` берётся из `payload.id` (в Postgres это запрос вида `payload->>'id'`).
  - если worker уже обработал заказ, API подмешивает “полный результат” из Redis.

### Kafka topics

- **`orders.events`**: API → Worker
- **`orders.ack`**: Worker → API

### Проверка данных в Postgres (Adminer)

Adminer: `http://localhost:8082`

- **Для API DB**
  - Server: `postgres_api` (внутри docker-compose сети)
  - Database: `orderstream_api`
- **Для Worker DB**
  - Server: `postgres_worker`
  - Database: `orderstream_worker`
- Username/Password: `postgres` / `upvel123`


### Тесты, моки, покрытие

#### Worker: suite-тесты с mockery-style моками

- Моки генерируются `mockery` по `.mockery.yaml` и лежат рядом с интерфейсами:
  - `services/worker/internal/services/worker/mocks/`

Установить mockery:

```powershell
go install github.com/vektra/mockery/v2@latest
```

Сгенерировать моки:

```powershell
mockery
```

Запуск тестов (из корня репозитория, где лежит `go.work`):

```powershell
go test ./... -count=1
```

Покрытие для worker (пример):

```powershell
go test ./internal/services/worker -count=1 -coverprofile=cover
go tool cover -func cover > cover.func.txt
Select-String -Path cover.func.txt -Pattern "internal/services/worker/handle_order_event.go" -SimpleMatch
```