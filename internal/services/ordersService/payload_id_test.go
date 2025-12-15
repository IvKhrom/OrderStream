package ordersService

import (
	"encoding/json"
	"testing"
)

func TestPayloadID(t *testing.T) {
	if got := payloadID(json.RawMessage(`{"id":"ext-123"}`)); got != "ext-123" {
		t.Fatalf("expected ext-123, got %q", got)
	}
}


