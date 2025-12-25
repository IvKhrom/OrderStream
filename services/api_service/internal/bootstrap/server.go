package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	ordersackconsumer "github.com/ivkhr/orderstream/services/api_service/internal/consumer/orders_ack_consumer"
	ordershttp "github.com/ivkhr/orderstream/services/api_service/internal/api/orders_http"
	"github.com/ivkhr/orderstream/services/api_service/internal/api/swagger"
)

type Router interface {
	Route(pattern string, fn func(r chi.Router))
}

func RunHTTP(httpPort string, ordersHandler *ordershttp.Handler, ackConsumer *ordersackconsumer.Consumer) error {
	if ackConsumer != nil {
		go ackConsumer.Consume(context.Background())
	}

	r := chi.NewRouter()
	sw := swagger.NewHTTP(swagger.NewEmbeddedProvider())
	sw.Register(r)
	ordersHandler.Routes(r)

	addr := fmt.Sprintf(":%s", httpPort)
	slog.Info("HTTP server listening on " + addr)
	return http.ListenAndServe(addr, r)
}


