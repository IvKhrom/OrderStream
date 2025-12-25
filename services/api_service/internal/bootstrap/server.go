package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	api "github.com/ivkhr/orderstream/services/api_service/internal/api/orders_service_api"
	ordersackconsumer "github.com/ivkhr/orderstream/services/api_service/internal/consumer/orders_ack_consumer"
	ordersapiswagger "github.com/ivkhr/orderstream/services/api_service/internal/pb/swagger/orders_api"
)

type Router interface {
	Route(pattern string, fn func(r chi.Router))
}

func RunHTTP(httpPort string, apiService *api.OrdersServiceAPI, ackConsumer *ordersackconsumer.Consumer) error {
	if ackConsumer != nil {
		go ackConsumer.Consume(context.Background())
	}

	r := chi.NewRouter()

	// Swagger берём строго из pb (сгенерированный orders.swagger.json).
	r.Get("/swagger.json", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write(ordersapiswagger.SwaggerJSON())
	})
	r.Get("/docs", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(docsHTML))
	})
	r.Get("/docs/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(docsHTML))
	})

	apiService.Routes(r)

	addr := fmt.Sprintf(":%s", httpPort)
	slog.Info("HTTP server listening on " + addr)
	return http.ListenAndServe(addr, r)
}

const docsHTML = `<!doctype html>
<html lang="en">
  <head>
    <meta charset="utf-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1"/>
    <title>api_service docs</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css"/>
    <style>html,body{margin:0;padding:0}#swagger-ui{height:100vh}</style>
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
      window.ui = SwaggerUIBundle({
        url: "/swagger.json",
        dom_id: "#swagger-ui"
      });
    </script>
  </body>
</html>`
