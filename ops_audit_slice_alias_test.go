package main

import (
	"testing"
	"time"
)

func seedAudit006() *OpsAudit {
	a := newOpsAudit()
	a.Add("r1", "created", "op-a")
	a.Add("r1", "status_changed", "op-b")
	a.Add("r2", "created", "op-a")
	a.Add("r2", "status_changed", "op-c")
	return a
}

func seedAuditWithTimes() *OpsAudit {
	a := newOpsAudit()
	a.mu.Lock()
	a.events = append(a.events,
		OpsEvent{ID: "e1", RecordID: "r1", Type: "created", At: "2026-01-01T00:00:00Z"},
		OpsEvent{ID: "e2", RecordID: "r1", Type: "status_changed", At: "2026-01-02T00:00:00Z"},
		OpsEvent{ID: "e3", RecordID: "r2", Type: "created", At: "2026-01-03T00:00:00Z"},
	)
	a.mu.Unlock()
	return a
}

func TestAuditForResultsNotOverwritten(t *testing.T) {
	a := seedAudit006()
	first := a.For("r1")
	if len(first) != 2 {
		t.Fatalf("expected 2 events for r1, got %d", len(first))
	}
	_ = a.For("r2")
	if len(first) != 2 {
		t.Fatalf("first query result was overwritten by a later query: len=%d", len(first))
	}
	for _, ev := range first {
		if ev.RecordID != "r1" {
			t.Fatalf("first query result contaminated with record %q", ev.RecordID)
		}
	}
}

func TestAuditSinceResultsNotOverwritten(t *testing.T) {
	a := seedAuditWithTimes()
	first := a.Since(time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC))
	if len(first) != 3 {
		t.Fatalf("expected 3 events since 2026-01-01, got %d", len(first))
	}
	_ = a.Since(time.Date(2026, 1, 3, 0, 0, 0, 0, time.UTC))
	if len(first) != 3 {
		t.Fatalf("first since-result was overwritten by a later query: len=%d", len(first))
	}
	if first[0].ID != "e1" {
		t.Fatalf("first since-result contaminated at head, got event %q", first[0].ID)
	}
}

func TestAuditLatestDetailsIsolated(t *testing.T) {
	a := newOpsAudit()
	a.mu.Lock()
	a.events = append(a.events, OpsEvent{
		ID: "e9", RecordID: "r1", Type: "created", At: "2026-01-01T00:00:00Z",
		Details: map[string]string{"site": "s1"},
	})
	a.mu.Unlock()
	latest, ok := a.Latest()
	if !ok {
		t.Fatal("expected a latest event")
	}
	latest.Details["site"] = "changed"
	latest.Details["extra"] = "x"
	stored, _ := a.Latest()
	if stored.Details["site"] != "s1" {
		t.Fatalf("mutating the returned latest event leaked into the audit store: %#v", stored.Details)
	}
	if _, ok := stored.Details["extra"]; ok {
		t.Fatalf("added detail leaked into the audit store: %#v", stored.Details)
	}
}

func TestOpsClonePageIsolated(t *testing.T) {
	page := OpsPage{Items: []OpsRecord{{ID: "r1"}, {ID: "r2"}}, Page: 1, PageSize: 10, Total: 2}
	clone := opsClonePage(page)
	clone.Items[0].ID = "changed"
	clone.Items[1].ID = "changed2"
	if page.Items[0].ID != "r1" || page.Items[1].ID != "r2" {
		t.Fatalf("mutating the cloned page leaked into the original: %#v", page.Items)
	}
}

func TestOpsPageLastIDEmptySafe(t *testing.T) {
	if got := opsLastID(OpsPage{}); got != "" {
		t.Fatalf("expected empty last id for empty page, got %q", got)
	}
}
