package json

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestMarshal(t *testing.T) {
	m := map[string]interface{}{
		"foo": "foo",
		"bar": make(chan struct{}),
	}
	enc := json.NewEncoder(os.Stdout)
	err := enc.Encode(m)
	fmt.Printf("err: %+v", err) // err: json: unsupported type: chan struct {}
}
