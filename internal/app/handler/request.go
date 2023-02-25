package handler

import (
	"encoding/json"
	"io"
	"net/http"
)

func readJSONBody(v any, r *http.Request) error {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, v)
}
