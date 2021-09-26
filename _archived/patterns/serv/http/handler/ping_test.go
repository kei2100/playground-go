package handler_test

import (
	"testing"
)

func TestPing(t *testing.T) {
	rec := servAndRecord(t, "GET", "/ping", nil, nil)
	assertResponseCode(t, rec.Code, 200)
	assertResponseJSON(t, rec.Body, map[string]interface{}{
		"message": "ok",
	})
}
