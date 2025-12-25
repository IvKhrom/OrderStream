package orders

import (
	"context"
	"time"
)

// AckCoordinator инкапсулирует паттерн:
// 1) зарегистрироваться на ACK
// 2) выполнить публикацию
// 3) дождаться ACK/таймаута/ctx.Done
type AckCoordinator interface {
	ExecuteAndWait(ctx context.Context, orderID string, publish func() error) error
}

type ackCoordinator struct {
	reg     AckWaitRegistry
	timeout time.Duration
}

func NewAckCoordinator(reg AckWaitRegistry, timeout time.Duration) AckCoordinator {
	if reg == nil {
		return nil
	}
	return &ackCoordinator{reg: reg, timeout: timeout}
}

func (c *ackCoordinator) ExecuteAndWait(ctx context.Context, orderID string, publish func() error) error {
	if c.timeout <= 0 {
		return publish()
	}

	ch, cleanup := c.reg.Register(orderID)
	defer cleanup()

	if err := publish(); err != nil {
		return err
	}

	select {
	case <-ch:
		return nil
	case <-time.After(c.timeout):
		return context.DeadlineExceeded
	case <-ctx.Done():
		return ctx.Err()
	}
}


