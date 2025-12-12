package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	w *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		w: kafka.NewWriter(kafka.WriterConfig{
			Brokers: brokers,
			Topic:   topic,
		}),
	}
}

func (p *Producer) Publish(ctx context.Context, value []byte) error {
	return p.w.WriteMessages(ctx, kafka.Message{
		Value: value,
		Time:  time.Now(),
	})
}

func (p *Producer) Close() error {
	return p.w.Close()
}

// PublishRaw публикует сообщение, создавая временный writer (без постоянного объекта Writer)
func PublishRaw(ctx context.Context, brokers []string, topic string, value []byte) error {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
	})
	defer w.Close()
	return w.WriteMessages(ctx, kafka.Message{Value: value})
}

// Обёртка для Consumer
type Consumer struct {
	r *kafka.Reader
}

func NewConsumer(brokers []string, topic, groupID string) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   topic,
	})
	return &Consumer{r: r}
}

func (c *Consumer) ReadMessage(ctx context.Context) (kafka.Message, error) {
	return c.r.ReadMessage(ctx)
}

func (c *Consumer) Close() error {
	return c.r.Close()
}
