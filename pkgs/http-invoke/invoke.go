package httpinvoke

import (
	"fmt"
	"io"
	"net/http"
)

type HttpInvoke struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"` // HTTP method, e.g., GET, POST
	Headers map[string]string `json:"headers"` // HTTP headers
	Body    io.Reader         `json:"body"`    // HTTP body, can be nil for GET requests
}

func (h *HttpInvoke) Invoke() ([]byte, error) {
	req, err := http.NewRequest(h.Method, h.URL, h.Body)
    if err != nil {
        return nil, fmt.Errorf("create request: %w", err)
    }

    for k, v := range h.Headers {
        req.Header.Set(k, v)
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("do request: %w", err)
    }
    defer resp.Body.Close()

    data, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("read response: %w", err)
    }

    return data, nil
}
// Example usage: