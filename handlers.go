package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type inspectionRequest struct {
	Inspection string `json:"inspection"`
}
type signalHandler struct{ store *SignalStore }

func (h signalHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/signals" {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		signals := h.store.List()
		if signals == nil {
			signals = []Signal{}
		}
		writeJSON(w, http.StatusOK, map[string]any{"signals": signals})
		return
	}
	if r.Method != http.MethodPost || !strings.HasPrefix(r.URL.Path, "/api/signals/") {
		writeError(w, http.StatusNotFound, "route not found")
		return
	}
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/signals/"), "/")
	if len(parts) != 2 || parts[1] != "inspection" || parts[0] == "" {
		writeError(w, http.StatusNotFound, "route not found")
		return
	}
	var request inspectionRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if err := validateInspection(request.Inspection); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	signal, exists, changed := h.store.RecordInspection(parts[0], request.Inspection)
	if !exists {
		writeError(w, http.StatusNotFound, "signal not found")
		return
	}
	if !changed {
		writeError(w, http.StatusConflict, "invalid inspection transition")
		return
	}
	writeJSON(w, http.StatusOK, signal)
}
