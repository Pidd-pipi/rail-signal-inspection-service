package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestLoadConfigRefreshesEnv(t *testing.T) {
	t.Setenv("PORT", "18091")
	first := loadConfig()
	if first.Port != "18091" {
		t.Fatalf("expected first load to read PORT=18091, got %q", first.Port)
	}
	t.Setenv("PORT", "18092")
	second := loadConfig()
	if second.Port != "18092" {
		t.Fatalf("expected second load to reflect the updated PORT, got %q (stale cache?)", second.Port)
	}
}

func TestLoadConfigConcurrentNoRace(t *testing.T) {
	t.Setenv("PORT", "18095")
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 40; j++ {
				_ = loadConfig()
			}
		}()
	}
	close(start)
	wg.Wait()
}

func TestBuildServerUsesConfigPort(t *testing.T) {
	server := buildServer(Config{Port: "18093"}, newSignalStore())
	if server.Addr != ":18093" {
		t.Fatalf("expected server to listen on :18093, got %q", server.Addr)
	}
}

func TestBuildServerUsesProvidedStore(t *testing.T) {
	store := newSignalStore()
	store.RecordInspection("SIG-21", "clear")
	server := buildServer(Config{Port: "18094"}, store)
	ts := httptest.NewServer(server.Handler)
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/api/signals")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var payload struct {
		Signals []Signal `json:"signals"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("invalid JSON: %s", body)
	}
	for _, s := range payload.Signals {
		if s.ID == "SIG-21" && s.Inspection != "clear" {
			t.Fatalf("expected the provided store's inspection to be served, got %q", s.Inspection)
		}
	}
}
