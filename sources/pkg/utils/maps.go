package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
)

var pool = &sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

func CopyValueWithJSONTags(dst, src any) error {
	buf := pool.Get().(*bytes.Buffer)
	defer func() {
		buf.Reset()
		pool.Put(buf)
	}()

	if err := json.NewEncoder(buf).Encode(src); err != nil {
		return fmt.Errorf("encode map to json format: %w", err)
	}
	return json.NewDecoder(buf).Decode(dst)
}
