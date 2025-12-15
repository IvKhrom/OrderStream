# OrderStream

Небольшой учебный микросервис для обработки заказов в event-driven архитектуре: HTTP API (`orders`) + Worker, использующие PostgreSQL (партиционирование), Kafka (events + ACK) и Redis (кэш).

**Кратко:** API записывает заказ в БД и публикует событие в Kafka; Worker обрабатывает событие, обновляет статус заказа и публикует ACK — API при необходимости ожидает ACK перед ответом клиенту.

**Компоненты**
- **API:** API-сервис ([cmd/api_service](cmd/api_service)) — принимает запрос, пишет заказ в БД, публикует событие и (опционально) ждёт ACK.
- **Worker:** фоновый обработчик ([cmd/worker](cmd/worker)) — читает `orders.events`, обрабатывает, обновляет Postgres и публикует ACK в `orders.ack`.
- **Postgres:** таблица `orders` с хеш-партиционированием на 4 партиции (см. `migrations`).
- **Kafka:** топики `orders.events` и `orders.ack`.
- **Redis:** опциональный кэш (в текущей версии не используется).

**Где смотреть код**
- Точка входа API: [cmd/api_service/main.go](cmd/api_service/main.go)
- Точка входа Worker: [cmd/worker/main.go](cmd/worker/main.go)
- API контракт: [api/orders_api/orders.proto](api/orders_api/orders.proto)
- Реализация API: [internal/api/orders_service_api](internal/api/orders_service_api)
- Конфигурация загрузки: [config/config.go](config/config.go)
- Сервисный слой: [internal/services/ordersService](internal/services/ordersService)
- Работа с данными (Postgres): [internal/storage/pgstorage](internal/storage/pgstorage)
- Миграции: [migrations](migrations)

**Быстрый старт (локально с Docker Compose)**

1. Поднимите инфраструктуру и сервисы:

```powershell
docker compose up -d --build
```

2. API будет доступен по `http://localhost:8080`.

3. Остановить/перезапустить только API (если нужно запускать локально через `go run`):

```powershell
docker compose stop api
cd cmd/api_service
go run .

# в другом терминале
cd cmd/worker
go run .
```

Или используйте Makefile-цели:

```powershell
make run-api
make run-worker
```

**Переменные окружения (важные)**
- `POSTGRES_DSN` — DSN для Postgres (по умолчанию без пароля указан в `internal/config/config.go`).
- `KAFKA_BROKERS` — адрес брокера Kafka (например, `localhost:9092`).
- `API_PORT` — порт API (по умолчанию `8080`).
- `WORKER_GROUP` — consumer group для Worker.
- `REDIS_ADDR` — адрес Redis.

Файл с загрузкой переменных: [config/config.go](config/config.go).

**Миграции**
- Файлы миграций находятся в папке `migrations/`.
- `0001_init.up.sql` — создание таблицы `orders` и партиций.
- `0002_remove_external_id.up.sql` — удаление устаревшего столбца `external_id` (и индекса).

Пример применения миграции вручную (в контейнере Postgres):

```powershell
docker cp migrations/0001_init.up.sql <postgres_container>:/tmp/0001_init.up.sql
docker exec -it <postgres_container> psql -U postgres -d orderstream -f /tmp/0001_init.up.sql
```

**HTTP API**

- `GET /health` — возвращает `ok`.

- `POST /orders` — создание / обновление / soft-delete.
  - Тело запроса (пример создания):

```json
{
  "order_id": "0",
  "user_id": "<uuid>",
  "payload": { "id": "external-123", "items": [...] }
}
```

  - Правила:
    - `order_id` пустой или `"0"` — создаётся новый внутренний `order_id` (UUID).
    - Для обновления указывайте внутренний `order_id` (UUID).
    - Soft-delete: укажите `"status": "deleted"` — запись помечается `deleted`.

  - Ответы:
    - `201 Created` при создании: `{ "order_id": "<uuid>", "status": "created" }`.
    - `200 OK` при успешном обновлении: `{ "order_id": "<uuid>", "status": "updated" }`.
    - `200 OK` при удалении: `{ "status": "deleted" }`.
    - `409 Conflict` при попытке обновить уже удалённый заказ.

- `GET /orders/{id}` — получить заказ по внутреннему `order_id` (UUID).

- `GET /orders/by-external/{external}?user_id=<uuid>` — поиск по внешнему marketplace id, которое берётся из `payload->>'id'`.

**Поток обработки (упрощённо)**
- Клиент -> POST /orders -> API записывает заказ и публикует событие в `orders.events`.
- Worker читает `orders.events`, обрабатывает заказ, обновляет Postgres и публикует ACK в `orders.ack`.
- API, при включённом ack-consumer, ждёт ACK и только после этого возвращает результат клиенту.

**Тесты и моки**
- Запуск тестов: `go test ./... -v`.
 

**Возможное развитие**
- Добавить CI (GitHub Actions): `go generate`, `go test -cover` и публикация coverage.
- Добавить интеграционные тесты с помощью Testcontainers / docker-compose для проверки взаимодействия с Kafka/Postgres/Redis.
- В production учесть безопасное хранение секретов (не хранить пароли в `docker-compose.yml`).
