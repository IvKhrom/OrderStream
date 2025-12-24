FROM golang:1.23.4-alpine AS build
WORKDIR /app

# Позволяем переопределять GOPROXY при сборке (актуально если proxy.golang.org недоступен/блокируется).
ARG GOPROXY=https://proxy.golang.org,direct
ENV GOPROXY=$GOPROXY

COPY go.mod go.sum /app/

# Отдельно скачиваем модули (проще отлаживать, плюс лучше кешируется).
# Используем кеши BuildKit для ускорения повторных сборок.
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

COPY . /app
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    cd cmd/api_service && CGO_ENABLED=0 GOOS=linux go build -o /orderstream-api-service .

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=build /orderstream-api-service /usr/local/bin/orderstream-api-service
ENTRYPOINT ["/usr/local/bin/orderstream-api-service"]


