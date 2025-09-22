package formparse

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/caddyserver/caddy/v2"
)

func newReplacerFunc(r *http.Request, fk []string) (caddy.ReplacerFunc, error) {
	// Save the body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("reading body: %w", err)
	}

	// Restore the body for FormValue parsing
	r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	fv := make(map[string]string)
	for _, key := range fk {
		fv[key] = r.FormValue(key)
	}

	// Restore the body for the upstream
	r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	return func(key string) (interface{}, bool) {
		prefix := "form."
		if !strings.HasPrefix(key, prefix) {
			return nil, false
		}
		key = strings.TrimPrefix(key, prefix)

		if v, ok := fv[key]; ok {
			return v, true
		}
		return "", true
	}, nil
}
