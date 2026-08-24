package main

import (
	"net"
	"net/http"
	"net/http/httptest"
	"sync"
	"syscall"
	"testing"
	"time"
)

const shutdownTestPort = "127.0.0.1:18237"

func TestServeHTTPGracefulShutdownReturns(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	server := &http.Server{Addr: shutdownTestPort, Handler: handler}
	done := make(chan error, 1)
	go func() { done <- serveHTTP(server) }()

	deadline := time.Now().Add(3 * time.Second)
	started := false
	for time.Now().Before(deadline) {
		conn, err := net.Dial("tcp", shutdownTestPort)
		if err == nil {
			conn.Close()
			started = true
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	if !started {
		t.Fatal("server did not start listening")
	}
	if err := syscall.Kill(syscall.Getpid(), syscall.SIGTERM); err != nil {
		t.Fatalf("failed to send SIGTERM: %v", err)
	}
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("serveHTTP returned error after graceful shutdown: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("serveHTTP deadlocked after SIGTERM; graceful shutdown never returned")
	}
}

func TestRequestIDCounterNoRace(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	h := requestIDMiddleware(inner)
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 50; j++ {
				rec := httptest.NewRecorder()
				h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/x", nil))
			}
		}()
	}
	close(start)
	wg.Wait()
}

func TestNewServerReusable(t *testing.T) {
	_ = newServer(newSignalStore())
	defer func() {
		if recovered := recover(); recovered != nil {
			t.Fatalf("newServer panicked on a second call: %v", recovered)
		}
	}()
	server := newServer(newSignalStore())
	if server == nil {
		t.Fatal("expected a mux from newServer")
	}
}

func TestNewServerConcurrentNoRace(t *testing.T) {
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 5; j++ {
				_ = newServer(newSignalStore())
			}
		}()
	}
	close(start)
	wg.Wait()
}
