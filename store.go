package main

import "sync"

type SignalStore struct {
	mu      sync.RWMutex
	signals map[string]Signal
}

func newSignalStore() *SignalStore {
	return &SignalStore{signals: map[string]Signal{
		"SIG-21": {ID: "SIG-21", Block: "B-12", Aspect: "proceed", Inspection: "pending"},
		"SIG-22": {ID: "SIG-22", Block: "B-13", Aspect: "caution", Inspection: "needs_attention"},
	}}
}

func (s *SignalStore) List() []Signal {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Signal, 0, len(s.signals))
	for _, signal := range s.signals {
		result = append(result, signal)
	}
	return result
}

func (s *SignalStore) RecordInspection(id, inspection string) (Signal, bool, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	signal, exists := s.signals[id]
	if !exists {
		return Signal{}, false, false
	}
	if !inspectionTransitions[signal.Inspection][inspection] {
		return signal, true, false
	}
	signal.Inspection = inspection
	s.signals[id] = signal
	return signal, true, true
}
