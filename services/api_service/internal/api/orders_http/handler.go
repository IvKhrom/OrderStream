package ordershttp

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Service interface {
	Upsert(ctx context.Context, orderID string, userID uuid.UUID, status string, payload json.RawMessage) (string, string, error)
}

type Handler struct {
	svc Service
}

func New(svc Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Routes(r chi.Router) {
	r.Post("/orders", h.upsert)
}

type upsertReq struct {
	OrderID    string          `json:"order_id"`
	UserID     string          `json:"user_id"`
	Status     string          `json:"status"`
	PayloadJSON json.RawMessage `json:"payload_json"`
	// payload как объект тоже поддержим
	Payload any `json:"payload"`
}

type upsertResp struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

func (h *Handler) upsert(w http.ResponseWriter, r *http.Request) {
	var req upsertReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		http.Error(w, "bad user_id", http.StatusBadRequest)
		return
	}

	payload := req.PayloadJSON
	if len(payload) == 0 && req.Payload != nil {
		b, _ := json.Marshal(req.Payload)
		payload = b
	}
	if len(payload) == 0 {
		http.Error(w, "payload_json required", http.StatusBadRequest)
		return
	}

	id, st, err := h.svc.Upsert(r.Context(), req.OrderID, userID, req.Status, payload)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(upsertResp{OrderID: id, Status: st})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(upsertResp{OrderID: id, Status: st})
}


