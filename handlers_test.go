package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSignalHTTPFlow(t *testing.T) {
	ts := httptest.NewServer(newServer(newSignalStore()))
	defer ts.Close()
	for _, path := range []string{"/healthz", "/api/signals", "/", "/app.js"} {
		response, err := http.Get(ts.URL + path)
		if err != nil {
			t.Fatal(err)
		}
		response.Body.Close()
		if response.StatusCode != http.StatusOK {
			t.Fatalf("GET %s: got %d", path, response.StatusCode)
		}
	}
	post := func(id, body string) int {
		response, err := http.Post(ts.URL+"/api/signals/"+id+"/inspection", "application/json", strings.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		defer response.Body.Close()
		return response.StatusCode
	}
	if got := post("SIG-21", `{"inspection":"clear"}`); got != http.StatusOK {
		t.Fatalf("valid inspection: got %d", got)
	}
	if got := post("SIG-21", `{"inspection":"cleared"}`); got != http.StatusConflict {
		t.Fatalf("invalid transition: got %d", got)
	}
	if got := post("SIG-22", `{"inspection":"cleared"}`); got != http.StatusOK {
		t.Fatalf("attention cleared: got %d", got)
	}
	if got := post("missing", `{"inspection":"clear"}`); got != http.StatusNotFound {
		t.Fatalf("missing signal: got %d", got)
	}
}
