# OrderStream

Учебный микросервис обработки заказов в event-driven архитектуре:

- **API** принимает HTTP-запросы, сохраняет заказ в PostgreSQL и публикует событие в Kafka (`orders.events`)
- **Worker** читает события из Kafka, обновляет данные в PostgreSQL и публикует ACK в Kafka (`orders.ack`)
- **API** читает ACK из Kafka и может дождаться подтверждения обработки

## Компоненты

- **API**: `cmd/api_service`
- **Worker**: `cmd/worker`
- **Контракт API**: `api/orders_api/orders.proto`
- **API реализация**: `internal/api/orders_service_api`
- **Сервисный слой**: `internal/services/ordersService`
- **PostgreSQL storage**: `internal/storage/pgstorage`
- **Kafka publisher**: `internal/storage/kafkastorage`
- **Миграции**: `migrations`

## Быстрый старт (Docker Compose)

Поднять инфраструктуру и сервисы:

```powershell
docker compose up -d --build
```

Что откроется:

- **API**: `http://localhost:8080`
- **Kafka UI**: `http://localhost:8081`
- **Postgres UI (Adminer)**: `http://localhost:8082`

Если нужно пересоздать Postgres “с нуля” (volume + миграции):

```powershell
docker compose down -v
docker compose up -d --build
```

## Переменные окружения

Загрузка переменных описана в `config/config.go`. Основные:

- `POSTGRES_DSN`
- `KAFKA_BROKERS`
- `API_PORT`
- `WORKER_GROUP`

## HTTP API

- `GET /health` → `{ "status": "ok" }`

- `POST /orders` — создать/обновить/soft-delete.
  - Создание: `order_id` пустой или `"0"`
  - Soft-delete: `"status": "deleted"`
  - `payload` можно передавать объектом — сервер сам сериализует в строку `payload_json`

- `GET /orders/{order_id}` — получить заказ по внутреннему UUID
- `GET /orders/by-external/{external_id}?user_id=<uuid>` — поиск по внешнему id (берётся из `payload->>'id'`)

## Топики Kafka

- `orders.events` — события заказов от API к Worker
- `orders.ack` — подтверждения (ACK) от Worker к API

## Проверка данных в Postgres (Adminer)

Adminer: `http://localhost:8082`

- System: `PostgreSQL`
- Server: `postgres`
- Username: `postgres`
- Password: `upvel123`
- Database: `orderstream`


## Тесты, моки, покрытие

Моки генерируются `mockery` по конфигу `.mockery.yaml` и лежат в `internal/services/ordersService/mocks`.

Установить mockery:

```powershell
go install github.com/vektra/mockery/v2@latest
```

Сгенерировать моки:

```powershell
make mock
```

Запуск тестов:

```powershell
go test ./... -count=1
```

Покрытие:

```powershell
go test ./... -count=1 -coverprofile=cover.out
go tool cover -func=cover.out
```

Проверка 100% покрытия по конкретным файлам сервисного слоя:

```powershell
# собрать профиль покрытия только для пакета сервисного слоя
go test ./internal/services/ordersService -count=1 -coverprofile=cover

# вывести процент по файлу upsert_order.go
go tool cover -func cover > cover.func.txt
Select-String -Path cover.func.txt -Pattern "internal/services/ordersService/upsert_order.go" -SimpleMatch

# вывести процент по файлу get_order_by_external_id.go
Select-String -Path cover.func.txt -Pattern "internal/services/ordersService/get_order_by_external_id.go" -SimpleMatch
```
