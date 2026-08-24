package main

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

func TestOpsConflictErrorChainPreserved(t *testing.T) {
	seed := []OpsRecord{{ID: "dup", Subject: "signal", Owner: "op", Priority: "high", Labels: map[string]string{"site": "s1"}}}
	svc := newOpsService(seed)
	record := OpsRecord{ID: "dup", Subject: "signal", Owner: "op", Priority: "high", Labels: map[string]string{"site": "s1"}}
	_, err := svc.Create(context.Background(), record)
	if err == nil {
		t.Fatal("expected a conflict error for duplicate id")
	}
	if !errors.Is(err, ErrOpsConflict) {
		t.Fatalf("conflict sentinel lost in chain: %v", err)
	}
	if !opsIsConflict(err) {
		t.Fatalf("opsIsConflict must hold for duplicate create: %v", err)
	}
}

func TestOpsTransitionErrorChainPreserved(t *testing.T) {
	sm := newOpsStateMachine()
	err := sm.Move(OpsStatusQueued, OpsStatusPaused, "invalid hop")
	if err == nil {
		t.Fatal("expected transition error for queued->paused")
	}
	if !errors.Is(err, ErrOpsTransition) {
		t.Fatalf("transition sentinel lost: %v", err)
	}
	if !opsIsTransition(err) {
		t.Fatalf("opsIsTransition must hold for transition errors: %v", err)
	}
}

func TestOpsCodeClassifiesBySentinel(t *testing.T) {
	if code := opsCode(fmt.Errorf("wrapped: %w", ErrOpsNotFound)); code != "not_found" {
		t.Fatalf("wrapped not-found must map to not_found, got %q", code)
	}
	if code := opsCode(errors.New("revision conflict detected")); code != "internal" {
		t.Fatalf("unrelated text must not be classified as conflict, got %q", code)
	}
	if code := opsCode(fmt.Errorf("boom: %w", ErrOpsPolicy)); code != "policy" {
		t.Fatalf("wrapped policy must map to policy, got %q", code)
	}
}

func TestOpsErrorUnwrapCarriesCause(t *testing.T) {
	err := wrapOps("update", "store.update", ErrOpsConflict)
	if !errors.Is(err, ErrOpsConflict) {
		t.Fatalf("wrapOps must preserve the wrapped cause: %v", err)
	}
	var typed *OpsError
	if !errors.As(err, &typed) {
		t.Fatalf("errors.As must find OpsError: %v", err)
	}
	if typed.Cause == nil {
		t.Fatal("OpsError.Cause must be retained")
	}
}
