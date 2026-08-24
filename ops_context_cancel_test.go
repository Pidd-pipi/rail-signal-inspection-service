package main

import (
	"context"
	"errors"
	"testing"
	"time"
)

func seedStore002() *OpsStore {
	return newOpsStore([]OpsRecord{
		{ID: "r1", Subject: "signal block", Owner: "op-a", Priority: "high", Labels: map[string]string{"site": "s1"}},
	})
}

func TestOpsStoreGetHonorsCanceledContext(t *testing.T) {
	store := seedStore002()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := store.Get(ctx, "r1"); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestOpsStoreListHonorsCanceledContext(t *testing.T) {
	store := seedStore002()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := store.List(ctx); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestOpsStorePutHonorsCanceledContext(t *testing.T) {
	store := newOpsStore(nil)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	item := OpsRecord{ID: "r9", Subject: "x", Owner: "o", Priority: "low", Labels: map[string]string{"site": "s1"}}
	if err := store.Put(ctx, item); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if _, err := store.Get(context.Background(), "r9"); !errors.Is(err, ErrOpsNotFound) {
		t.Fatalf("canceled put must not persist the record, got %v", err)
	}
}

func TestOpsStoreUpdateHonorsCanceledContext(t *testing.T) {
	store := seedStore002()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	rec, _ := store.Get(context.Background(), "r1")
	if err := store.Update(ctx, rec, 0); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}

func TestOpsStoreDeleteHonorsCanceledContext(t *testing.T) {
	store := seedStore002()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := store.Delete(ctx, "r1"); !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if _, err := store.Get(context.Background(), "r1"); err != nil {
		t.Fatalf("canceled delete must keep the record, got %v", err)
	}
}

func TestOpsContextPropagatesParentDeadline(t *testing.T) {
	parent, cancel := context.WithTimeout(context.Background(), 40*time.Millisecond)
	defer cancel()
	child, childCancel := opsContext(parent, 5*time.Second)
	defer childCancel()
	select {
	case <-child.Done():
	case <-time.After(300 * time.Millisecond):
		t.Fatal("child context did not inherit the parent deadline")
	}
}
