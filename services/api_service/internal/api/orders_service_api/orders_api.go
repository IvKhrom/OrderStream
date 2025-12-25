package orders_service_api

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/ivkhr/orderstream/services/api_service/internal/models"
)

type OrdersService interface {
	UpsertOrder(ctx context.Context, orderID string, userID uuid.UUID, status string, payload json.RawMessage) (string, string, error)
	GetViewByID(ctx context.Context, id uuid.UUID) (*models.Order, error)
	GetViewByExternalID(ctx context.Context, externalID string, userID uuid.UUID) (*models.Order, error)
}

type OrdersServiceAPI struct {
	orders OrdersService
}

func NewOrdersServiceAPI(orders OrdersService) *OrdersServiceAPI {
	return &OrdersServiceAPI{orders: orders}
}

func (a *OrdersServiceAPI) Routes(r chi.Router) {
	r.Get("/health", a.health)
	r.Post("/orders", a.upsertOrder)
	r.Get("/orders/{order_id}", a.getOrderByID)
	r.Get("/orders/by-external/{external_id}", a.getOrderByExternalID)
}

func (a *OrdersServiceAPI) health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

type getOrderResp struct {
	Order *orderHTTPModel `json:"order,omitempty"`
}

type orderHTTPModel struct {
	OrderID     string          `json:"order_id"`
	UserID      string          `json:"user_id"`
	Amount      float64         `json:"amount,omitempty"`
	Status      string          `json:"status"`
	PayloadJSON json.RawMessage `json:"payload_json,omitempty"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
	Bucket      int32           `json:"bucket"`
}

func mapOrderToHTTP(o *models.Order) *orderHTTPModel {
	if o == nil {
		return nil
	}
	return &orderHTTPModel{
		OrderID:     o.OrderID.String(),
		UserID:      o.UserID.String(),
		Amount:      o.Amount,
		Status:      o.Status,
		PayloadJSON: json.RawMessage(o.Payload),
		CreatedAt:   o.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:   o.UpdatedAt.UTC().Format(time.RFC3339Nano),
		Bucket:      int32(o.Bucket),
	}
}

func (a *OrdersServiceAPI) getOrderByID(w http.ResponseWriter, r *http.Request) {
	oid, err := uuid.Parse(chi.URLParam(r, "order_id"))
	if err != nil {
		http.Error(w, "bad order_id", http.StatusBadRequest)
		return
	}
	o, err := a.orders.GetViewByID(r.Context(), oid)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(getOrderResp{Order: mapOrderToHTTP(o)})
}

func (a *OrdersServiceAPI) getOrderByExternalID(w http.ResponseWriter, r *http.Request) {
	externalID := chi.URLParam(r, "external_id")
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "bad user_id", http.StatusBadRequest)
		return
	}
	o, err := a.orders.GetViewByExternalID(r.Context(), externalID, userID)
	if err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(getOrderResp{Order: mapOrderToHTTP(o)})
}
