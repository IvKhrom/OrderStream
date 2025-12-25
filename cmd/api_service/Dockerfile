FROM golang:1.23.4-alpine AS build
WORKDIR /app

ARG GOPROXY=https://proxy.golang.org,direct
ENV GOPROXY=$GOPROXY

COPY go.mod go.sum /app/

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


