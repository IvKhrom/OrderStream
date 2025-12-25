FROM golang:1.23.4-alpine AS build
WORKDIR /app

ARG GOPROXY=https://proxy.golang.org,direct
ENV GOPROXY=$GOPROXY

# go.work + go.mod'ы нужны для сборки в mono-repo
COPY go.work /app/go.work
COPY shared/go.mod shared/go.mod
COPY services/api_service/go.mod services/api_service/go.mod

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    cd services/api_service && go mod download

COPY . /app

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    cd services/api_service/cmd/api_service && CGO_ENABLED=0 GOOS=linux go build -o /api-service .

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=build /api-service /usr/local/bin/api-service
ENTRYPOINT ["/usr/local/bin/api-service"]


