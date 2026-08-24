package main

import "net/http"

func newServer(store *SignalStore) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthHandler)
	mux.Handle("/api/signals", signalHandler{store: store})
	mux.Handle("/api/signals/", signalHandler{store: store})
	mux.HandleFunc("/", staticHandler)
	return mux
}
