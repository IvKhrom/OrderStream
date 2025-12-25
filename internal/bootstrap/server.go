package bootstrap

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	server "github.com/ivkhr/orderstream/internal/api/orders_service_api"
	"github.com/ivkhr/orderstream/internal/api/swagger"
	ordersackconsumer "github.com/ivkhr/orderstream/internal/consumer/orders_ack_consumer"
	"github.com/ivkhr/orderstream/internal/pb/orders_api"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"
)

func AppRun(api server.OrdersServiceAPI, ackConsumer *ordersackconsumer.OrdersAckConsumer, httpPort string) {
	if ackConsumer != nil {
		go ackConsumer.Consume(context.Background())
	}
	go func() {
		if err := runGRPCServer(api); err != nil {
			panic(fmt.Errorf("failed to run gRPC server: %v", err))
		}
	}()
	if err := runHTTPServer(api, httpPort); err != nil {
		panic(fmt.Errorf("failed to run gateway server: %v", err))
	}
}

func runGRPCServer(api server.OrdersServiceAPI) error {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	orders_api.RegisterOrdersServiceServer(s, &api)
	slog.Info("gRPC server listening on :50051")
	return s.Serve(lis)
}

func runHTTPServer(api server.OrdersServiceAPI, httpPort string) error {
	r := chi.NewRouter()

	sw := swagger.NewHTTP(swagger.NewEmbeddedProvider())
	sw.Register(r)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		resp, err := (&api).Health(r.Context(), &orders_api.HealthRequest{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	})
	r.Get("/health/", func(w http.ResponseWriter, r *http.Request) {
		resp, err := (&api).Health(r.Context(), &orders_api.HealthRequest{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	})

	r.Post("/orders", func(w http.ResponseWriter, r *http.Request) {
		var raw map[string]any
		if err := json.NewDecoder(r.Body).Decode(&raw); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		req := &orders_api.UpsertOrderRequest{}
		if v, ok := raw["order_id"].(string); ok {
			req.OrderId = v
		}
		if v, ok := raw["user_id"].(string); ok {
			req.UserId = v
		}
		if v, ok := raw["status"].(string); ok {
			req.Status = v
		}
		if v, ok := raw["payload_json"].(string); ok {
			req.PayloadJson = v
		} else if v, ok := raw["payload"]; ok {
			// payload как объект -> сериализуем в строку
			b, _ := json.Marshal(v)
			req.PayloadJson = string(b)
		}

		resp, err := (&api).UpsertOrder(r.Context(), req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	})

	r.Get("/orders/{order_id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "order_id")
		resp, err := (&api).GetOrderByID(r.Context(), &orders_api.GetOrderByIDRequest{OrderId: id})
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if resp == nil || resp.Order == nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		// Для HTTP-клиента удобнее отдавать payload как JSON-объект, а не строку с экранированием.
		var payload any = nil
		if resp.Order.PayloadJson != "" {
			var tmp any
			if err := json.Unmarshal([]byte(resp.Order.PayloadJson), &tmp); err == nil {
				payload = tmp
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"order": map[string]any{
				"order_id":   resp.Order.OrderId,
				"user_id":    resp.Order.UserId,
				"amount":     resp.Order.Amount,
				"status":     resp.Order.Status,
				"payload":    payload,
				"created_at": resp.Order.CreatedAt,
				"updated_at": resp.Order.UpdatedAt,
				"bucket":     resp.Order.Bucket,
			},
		})
	})

	r.Get("/orders/by-external/{external_id}", func(w http.ResponseWriter, r *http.Request) {
		ext := chi.URLParam(r, "external_id")
		userID := r.URL.Query().Get("user_id")
		resp, err := (&api).GetOrderByExternalID(r.Context(), &orders_api.GetOrderByExternalIDRequest{
			ExternalId: ext,
			UserId:     userID,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if resp == nil || resp.Order == nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		var payload any = nil
		if resp.Order.PayloadJson != "" {
			var tmp any
			if err := json.Unmarshal([]byte(resp.Order.PayloadJson), &tmp); err == nil {
				payload = tmp
			}
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"order": map[string]any{
				"order_id":   resp.Order.OrderId,
				"user_id":    resp.Order.UserId,
				"amount":     resp.Order.Amount,
				"status":     resp.Order.Status,
				"payload":    payload,
				"created_at": resp.Order.CreatedAt,
				"updated_at": resp.Order.UpdatedAt,
				"bucket":     resp.Order.Bucket,
			},
		})
	})

	addr := fmt.Sprintf(":%s", httpPort)
	slog.Info("HTTP server listening on " + addr)
	return http.ListenAndServe(addr, r)
}
