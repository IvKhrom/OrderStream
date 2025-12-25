package ordersackprocessor

type Notifier interface {
	Notify(orderID string)
}

type Processor struct {
	notifier Notifier
}

func New(notifier Notifier) *Processor {
	return &Processor{notifier: notifier}
}
