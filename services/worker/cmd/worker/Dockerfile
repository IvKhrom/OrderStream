FROM golang:1.23.4-alpine AS build
WORKDIR /app

ARG GOPROXY=https://proxy.golang.org,direct
ENV GOPROXY=$GOPROXY

COPY go.work /app/go.work
COPY shared/go.mod shared/go.mod
COPY services/worker/go.mod services/worker/go.mod

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    cd services/worker && go mod download

COPY . /app

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    cd services/worker/cmd/worker && CGO_ENABLED=0 GOOS=linux go build -o /worker-service .

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=build /worker-service /usr/local/bin/worker-service
ENTRYPOINT ["/usr/local/bin/worker-service"]


