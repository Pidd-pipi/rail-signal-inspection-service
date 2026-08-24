package main

import (
	"sync"
	"testing"
)

func TestSignalStoreConcurrentListNoRace(t *testing.T) {
	store := newSignalStore()
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 20; j++ {
				_ = store.List()
			}
		}()
	}
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 20; j++ {
				store.RecordInspection("SIG-21", "clear")
			}
		}()
	}
	close(start)
	wg.Wait()
}

func TestSignalInspectionLostUpdate(t *testing.T) {
	for attempt := 0; attempt < 80; attempt++ {
		store := newSignalStore()
		start := make(chan struct{})
		results := make(chan bool, 2)
		var wg sync.WaitGroup
		for _, target := range []string{"clear", "needs_attention"} {
			wg.Add(1)
			go func(tgt string) {
				defer wg.Done()
				<-start
				_, _, changed := store.RecordInspection("SIG-21", tgt)
				results <- changed
			}(target)
		}
		close(start)
		wg.Wait()
		close(results)
		wins := 0
		for changed := range results {
			if changed {
				wins++
			}
		}
		if wins != 1 {
			t.Fatalf("attempt %d: expected exactly 1 accepted transition, got %d", attempt, wins)
		}
	}
}
