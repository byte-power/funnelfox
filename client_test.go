package funnelfox

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDoRequestSetsStatusCodeOnAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"status":"error","req_id":"req_123","error":[{"type":"validation","msg":"failed"}]}`))
	}))
	defer server.Close()

	client := NewClientWithHTTPClient("org", "secret", server.Client(), nil)
	client.baseURL = server.URL

	err := client.doRequest("/test", nil, nil, false)
	if err == nil {
		t.Fatal("expected API error")
	}
	if err.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected status code %d, got %d", http.StatusBadRequest, err.StatusCode)
	}
}

func TestDoRequestSetsStatusCodeOnInvalidJSONHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte("bad gateway"))
	}))
	defer server.Close()

	client := NewClientWithHTTPClient("org", "secret", server.Client(), nil)
	client.baseURL = server.URL

	err := client.doRequest("/test", nil, nil, false)
	if err == nil {
		t.Fatal("expected unmarshal error")
	}
	if err.StatusCode != http.StatusBadGateway {
		t.Fatalf("expected status code %d, got %d", http.StatusBadGateway, err.StatusCode)
	}
}
