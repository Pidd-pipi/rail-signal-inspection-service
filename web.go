package main

import (
	"net/http"
	"os"
)

func staticHandler(w http.ResponseWriter, r *http.Request) {
	name := "web/index.html"
	data, err := os.ReadFile(name)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(data)
}
