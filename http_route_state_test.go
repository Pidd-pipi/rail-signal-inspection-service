package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSignalListEndpointReturnsJSON(t *testing.T) {
	server := newEnterpriseServer(":18096", newAppHandler(newSignalStore()))
	ts := httptest.NewServer(server.Handler)
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/api/signals")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d (body=%s)", resp.StatusCode, body)
	}
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("list endpoint must return JSON, got: %s", body)
	}
	if _, ok := payload["signals"]; !ok {
		t.Fatalf("list payload missing signals field: %s", body)
	}
}

func TestSignalInspectionEndpointReturnsJSON(t *testing.T) {
	server := newEnterpriseServer(":18096", newAppHandler(newSignalStore()))
	ts := httptest.NewServer(server.Handler)
	defer ts.Close()
	resp, err := http.Post(ts.URL+"/api/signals/SIG-21/inspection", "application/json", strings.NewReader(`{"inspection":"clear"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d (body=%s)", resp.StatusCode, body)
	}
	var signal Signal
	if err := json.Unmarshal(body, &signal); err != nil {
		t.Fatalf("inspection endpoint must return a JSON signal, got: %s", body)
	}
	if signal.Inspection != "clear" {
		t.Fatalf("expected inspection clear, got %q", signal.Inspection)
	}
}

func TestAppJsServedAsJavaScript(t *testing.T) {
	server := newEnterpriseServer(":18096", newAppHandler(newSignalStore()))
	ts := httptest.NewServer(server.Handler)
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/app.js")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	ct := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "text/javascript") {
		t.Fatalf("expected text/javascript content type, got %q", ct)
	}
	if strings.Contains(string(body), "<!doctype") {
		t.Fatal("app.js must not be served as the index page")
	}
	if !strings.Contains(string(body), "loadSignals") {
		t.Fatal("app.js body does not look like JavaScript")
	}
}

func TestStaticResponseNotCached(t *testing.T) {
	server := newEnterpriseServer(":18096", newAppHandler(newSignalStore()))
	ts := httptest.NewServer(server.Handler)
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/app.js")
	if err != nil {
		t.Fatal(err)
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
	if got := resp.Header.Get("Cache-Control"); got != "no-cache" {
		t.Fatalf("expected Cache-Control no-cache on static assets, got %q", got)
	}
}

func TestStartAppReportsStartupError(t *testing.T) {
	err := startApp(Config{Port: "bad port"}, newSignalStore())
	if err == nil {
		t.Fatal("startup failure must be reported, not swallowed")
	}
}
