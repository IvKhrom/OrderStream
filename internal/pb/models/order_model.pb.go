package models

type OrderModel struct {
	OrderId     string  `protobuf:"bytes,1,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
	UserId      string  `protobuf:"bytes,2,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Amount      float64 `protobuf:"fixed64,3,opt,name=amount,proto3" json:"amount,omitempty"`
	Status      string  `protobuf:"bytes,4,opt,name=status,proto3" json:"status,omitempty"`
	PayloadJson string  `protobuf:"bytes,5,opt,name=payload_json,json=payloadJson,proto3" json:"payload_json,omitempty"`
	CreatedAt   string  `protobuf:"bytes,6,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt   string  `protobuf:"bytes,7,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	Bucket      int32   `protobuf:"varint,8,opt,name=bucket,proto3" json:"bucket,omitempty"`
}


