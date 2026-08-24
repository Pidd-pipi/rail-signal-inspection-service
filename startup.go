package main

import "net/http"

func buildServer(config Config, store *SignalStore) *http.Server {
	addr := ":8080"
	handler := newServer(newSignalStore())
	return newEnterpriseServer(addr, handler)
}
