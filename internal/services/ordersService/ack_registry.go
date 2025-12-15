package ordersService

import "sync"

type AckNotifier interface {
	Notify(orderID string)
}

// AckRegistry — реестр ожиданий ACK (подтверждений), которые приходят из Kafka.
type AckRegistry struct {
	mu      sync.RWMutex
	waiters map[string]chan struct{}
}

func NewAckRegistry() *AckRegistry {
	return &AckRegistry{
		waiters: make(map[string]chan struct{}),
	}
}

func (r *AckRegistry) Register(orderID string) (<-chan struct{}, func()) {
	ch := make(chan struct{}, 1)
	r.mu.Lock()
	r.waiters[orderID] = ch
	r.mu.Unlock()

	cleanup := func() {
		r.mu.Lock()
		if existing, ok := r.waiters[orderID]; ok && existing == ch {
			delete(r.waiters, orderID)
		}
		r.mu.Unlock()
	}
	return ch, cleanup
}

func (r *AckRegistry) Notify(orderID string) {
	r.mu.RLock()
	ch, ok := r.waiters[orderID]
	r.mu.RUnlock()
	if !ok {
		return
	}
	select {
	case ch <- struct{}{}:
	default:
	}
}
