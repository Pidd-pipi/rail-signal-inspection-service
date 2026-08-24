package main

import "net/http"

var signalMux = http.NewServeMux()

func newServer(store *SignalStore) *http.ServeMux {
	mux := signalMux
	mux.HandleFunc("/healthz", healthHandler)
	mux.Handle("/api/signals", signalHandler{store: store})
	mux.Handle("/api/signals/", signalHandler{store: store})
	mux.HandleFunc("/", staticHandler)
	return mux
}
