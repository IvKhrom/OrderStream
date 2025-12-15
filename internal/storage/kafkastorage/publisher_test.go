package kafkastorage

import (
	"testing"
)

func TestNewPublisherAndClose(t *testing.T) {
	p := NewPublisher([]string{"127.0.0.1:0"}, "topic")
	if p == nil {
		t.Fatalf("expected publisher, got nil")
	}
	if err := p.Close(); err != nil {
		t.Logf("close returned error (ok): %v", err)
	}
}
