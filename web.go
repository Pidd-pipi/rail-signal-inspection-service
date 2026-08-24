package main

import (
	"net/http"
	"os"
)

func staticHandler(w http.ResponseWriter, r *http.Request) {
	name := "web/index.html"
	if r.URL.Path == "/app.js" {
		name = "web/app.js"
	}
	data, err := os.ReadFile(name)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if name == "web/app.js" {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	} else {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	}
	_, _ = w.Write(data)
}
