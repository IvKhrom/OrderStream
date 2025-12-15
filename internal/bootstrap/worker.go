package bootstrap

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	orderseventsconsumer "github.com/ivkhr/orderstream/internal/consumer/orders_events_consumer"
)

func WorkerRun(consumer *orderseventsconsumer.OrdersEventsConsumer) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		cancel()
	}()

	consumer.Consume(ctx)
}


