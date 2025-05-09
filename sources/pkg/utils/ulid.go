package utils

import (
	"math/rand"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

var (
	ulonce  sync.Once
	ulmx    sync.Mutex
	entropy *rand.Rand
)

func GenerateULID() string {
	ulonce.Do(func() {
		entropy = rand.New(rand.NewSource(time.Now().UnixNano()))
	})
	ulmx.Lock()
	defer ulmx.Unlock()

	id, _ := ulid.New(ulid.Timestamp(time.Now()), entropy)
	return id.String()
}
