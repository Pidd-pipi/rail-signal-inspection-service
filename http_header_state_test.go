package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEnterpriseResponseHeadersPresent(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"status":"ok"}`)
	})
	ts := httptest.NewServer(opsEnterpriseMiddleware(inner))
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/healthz")
	if err != nil {
		t.Fatal(err)
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
	if got := resp.Header.Get("X-Operations-Domain"); got != opsDomainName {
		t.Fatalf("expected X-Operations-Domain %q, got %q", opsDomainName, got)
	}
	if got := resp.Header.Get("X-Operations-Latency-Ms"); got == "" {
		t.Fatal("expected X-Operations-Latency-Ms header to be present")
	}
	if got := resp.Header.Get("X-Operations-Request"); got != "generated" {
		t.Fatalf("expected X-Operations-Request generated, got %q", got)
	}
}

func TestHealthJSONContentTypeAndReadyHeader(t *testing.T) {
	ts := httptest.NewServer(newServer(newSignalStore()))
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/healthz")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
	ct := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("expected application/json content type on healthz, got %q", ct)
	}
	if got := resp.Header.Get("X-Service-Ready"); got != "yes" {
		t.Fatalf("expected X-Service-Ready header, got %q", got)
	}
}

func TestOpsActorFromRequestWebFallback(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/signals", nil)
	if got := opsActorFromRequest(req); got != "web" {
		t.Fatalf("expected web fallback actor, got %q", got)
	}
	req2 := httptest.NewRequest(http.MethodGet, "/api/signals", nil)
	req2.Header.Set("X-Operator", "  lin   ")
	if got := opsActorFromRequest(req2); got != "lin" {
		t.Fatalf("expected trimmed actor, got %q", got)
	}
}
