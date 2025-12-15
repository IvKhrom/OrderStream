package models

import "github.com/google/uuid"

func BucketFromUUID(id uuid.UUID, buckets int) int {
	return int(id[0] % byte(buckets))
}


