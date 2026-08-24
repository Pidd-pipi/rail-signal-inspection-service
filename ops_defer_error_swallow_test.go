package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func seedService004() *OpsService {
	return newOpsService([]OpsRecord{
		{ID: "r1", Subject: "signal block", Owner: "op-a", Priority: "high", Labels: map[string]string{"site": "s1"}},
	})
}

func TestTransitionNotFoundErrorPreserved(t *testing.T) {
	svc := seedService004()
	_, err := svc.Transition(context.Background(), "missing", 0, OpsStatusActive, "op-a")
	if !errors.Is(err, ErrOpsNotFound) {
		t.Fatalf("not-found sentinel must survive transition: %v", err)
	}
}

func TestTransitionConflictErrorPreserved(t *testing.T) {
	svc := seedService004()
	_, err := svc.Transition(context.Background(), "r1", 999, OpsStatusActive, "op-a")
	if !errors.Is(err, ErrOpsConflict) {
		t.Fatalf("conflict sentinel must survive transition: %v", err)
	}
}

func TestInspectionValidationErrorDetailPreserved(t *testing.T) {
	err := validateInspection("bogus")
	if err == nil {
		t.Fatal("expected an error for unsupported inspection")
	}
	if !strings.Contains(err.Error(), "bogus") {
		t.Fatalf("validation error must keep the offending value: %v", err)
	}
}

func TestInspectionValidationHTTPShowsDetail(t *testing.T) {
	ts := httptest.NewServer(newServer(newSignalStore()))
	defer ts.Close()
	resp, err := http.Post(ts.URL+"/api/signals/SIG-21/inspection", "application/json", strings.NewReader(`{"inspection":"bogus"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
	var payload map[string]string
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("invalid JSON body: %v (%s)", err, body)
	}
	if !strings.Contains(payload["error"], "bogus") {
		t.Fatalf("HTTP error must include the rejected value: %v", payload["error"])
	}
}
