package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestBucketFromUUID(t *testing.T) {
	var u uuid.UUID
	u[0] = 5
	b := BucketFromUUID(u, 4)
	if b != 1 {
		t.Fatalf("expected bucket 1, got %d", b)
	}
}
