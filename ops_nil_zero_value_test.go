package main

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNormalizeOpsRecordInitializesLabels(t *testing.T) {
	rec := normalizeOpsRecord(OpsRecord{ID: "R1", Subject: "  signal block ", Owner: " op-a ", Priority: "high"})
	if rec.Labels == nil {
		t.Fatal("labels must be initialized for zero-value records")
	}
	if rec.Labels["normalized"] != "1" {
		t.Fatalf("normalize marker must be set, got %#v", rec.Labels)
	}
}

func TestOpsRecordCloneDeepCopiesLabels(t *testing.T) {
	orig := OpsRecord{ID: "r1", Subject: "signal", Owner: "op", Priority: "high", Labels: map[string]string{"site": "s1"}}
	clone := orig.Clone()
	clone.Labels["site"] = "s2"
	clone.Labels["extra"] = "x"
	if orig.Labels["site"] != "s1" {
		t.Fatalf("clone must not share the labels map, orig site=%q", orig.Labels["site"])
	}
	if _, ok := orig.Labels["extra"]; ok {
		t.Fatalf("clone additions leaked into original: %#v", orig.Labels)
	}
}

func TestNewOpsStoreZeroValueSeedNoPanic(t *testing.T) {
	store := newOpsStore([]OpsRecord{{ID: "r1", Subject: "signal", Owner: "op", Priority: "high"}})
	rec, err := store.Get(context.Background(), "r1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Labels["normalized"] != "1" {
		t.Fatalf("stored record must carry normalize marker, got %#v", rec.Labels)
	}
}

func TestSignalListEmptyIsArray(t *testing.T) {
	store := &SignalStore{signals: map[string]Signal{}}
	ts := httptest.NewServer(newServer(store))
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/api/signals")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if strings.Contains(string(body), `"signals":null`) {
		t.Fatalf("empty signal list must serialize as [], got: %s", body)
	}
	if !strings.Contains(string(body), `"signals":[]`) {
		t.Fatalf("expected empty array for signals, got: %s", body)
	}
}
