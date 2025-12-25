package orders_service_api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type upsertReq struct {
	OrderID     string          `json:"order_id"`
	UserID      string          `json:"user_id"`
	Status      string          `json:"status"`
	PayloadJSON json.RawMessage `json:"payload_json"`
	Payload     any             `json:"payload"`
}

type upsertResp struct {
	OrderID string `json:"order_id"`
	Status  string `json:"status"`
}

func (a *OrdersServiceAPI) upsertOrder(w http.ResponseWriter, r *http.Request) {
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

	id, st, err := a.orders.UpsertOrder(r.Context(), req.OrderID, userID, req.Status, payload)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(upsertResp{OrderID: id, Status: st})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(upsertResp{OrderID: id, Status: st})
}


