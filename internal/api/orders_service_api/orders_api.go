package orders_service_api

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"github.com/ivkhr/orderstream/internal/models"
	proto_models "github.com/ivkhr/orderstream/internal/pb/models"
	"github.com/ivkhr/orderstream/internal/pb/orders_api"
)

type ordersService interface {
	UpsertOrder(ctx context.Context, orderID string, userID uuid.UUID, status string, payload json.RawMessage) (string, string, error)
	GetOrderByID(ctx context.Context, id uuid.UUID) (*models.Order, error)
	GetOrderByExternalID(ctx context.Context, externalID string, userID uuid.UUID) (*models.Order, error)
}

type OrdersServiceAPI struct {
	orders_api.UnimplementedOrdersServiceServer
	ordersService ordersService
}

func NewOrdersServiceAPI(ordersService ordersService) *OrdersServiceAPI {
	return &OrdersServiceAPI{ordersService: ordersService}
}

func mapOrderToProto(o *models.Order) *proto_models.OrderModel {
	if o == nil {
		return nil
	}
	return &proto_models.OrderModel{
		OrderId:     o.OrderID.String(),
		UserId:      o.UserID.String(),
		Amount:      o.Amount,
		Status:      o.Status,
		PayloadJson: string(o.Payload),
		CreatedAt:   o.CreatedAt.UTC().Format(time.RFC3339Nano),
		UpdatedAt:   o.UpdatedAt.UTC().Format(time.RFC3339Nano),
		Bucket:      int32(o.Bucket),
	}
}


