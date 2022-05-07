package data

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid"
)

func NewULID() (string, time.Time) {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	return ulid.MustNew(ulid.Timestamp(t), entropy).String(), t
}

func TTL(d time.Duration) int64 {
	return time.Now().Add(d).Unix()
}
