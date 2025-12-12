FROM golang:1.23.4-alpine AS build
WORKDIR /app
COPY go.mod go.sum /app/
RUN go env -w GOPROXY=https://proxy.golang.org
COPY . /app
RUN cd cmd/api && CGO_ENABLED=0 GOOS=linux go build -o /orderstream-api .

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=build /orderstream-api /usr/local/bin/orderstream-api
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/orderstream-api"]
