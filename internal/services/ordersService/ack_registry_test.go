package ordersService

import (
	"testing"
	"time"
)

func TestAckRegistry_RegisterNotify(t *testing.T) {
	r := NewAckRegistry()

	ch, cleanup := r.Register("id-1")
	defer cleanup()

	r.Notify("id-1")

	select {
	case <-ch:
		// ok
	case <-time.After(200 * time.Millisecond):
		t.Fatalf("expected ACK notification")
	}
}


