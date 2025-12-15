package ordersackprocessor

type ackNotifier interface {
	Notify(orderID string)
}

type OrdersAckProcessor struct {
	notifier ackNotifier
}

func NewOrdersAckProcessor(notifier ackNotifier) *OrdersAckProcessor {
	return &OrdersAckProcessor{notifier: notifier}
}


