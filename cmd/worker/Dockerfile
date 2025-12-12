FROM golang:1.23.4-alpine AS build
WORKDIR /app
COPY go.mod go.sum /app/
RUN go env -w GOPROXY=https://proxy.golang.org
COPY . /app
RUN cd cmd/worker && CGO_ENABLED=0 GOOS=linux go build -o /orderstream-worker .

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=build /orderstream-worker /usr/local/bin/orderstream-worker
ENTRYPOINT ["/usr/local/bin/orderstream-worker"]
