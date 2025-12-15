package models

import (
	"testing"

	"github.com/google/uuid"
)

func TestBucketFromUUID(t *testing.T) {
	id := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	b := BucketFromUUID(id, 4)
	if b < 0 || b >= 4 {
		t.Fatalf("ожидали bucket в диапазоне [0..3], получили %d", b)
	}
}


