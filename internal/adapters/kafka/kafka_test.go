package kafka

import (
	"context"
	"testing"
	"time"
)

func TestNewProducerAndClose(t *testing.T) {
	p := NewProducer([]string{"127.0.0.1:0"}, "test-topic")
	if p == nil {
		t.Fatalf("expected producer, got nil")
	}
	if err := p.Close(); err != nil {
		// Close may return error if writer wasn't connected; that's fine
		t.Logf("producer close returned error (ok): %v", err)
	}
}

func TestPublishRawTimesOut(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	err := PublishRaw(ctx, []string{"127.0.0.1:0"}, "topic", []byte("hello"))
	if err == nil {
		t.Fatalf("expected error when publishing to invalid broker")
	}
}
