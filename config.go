package main

import "os"

var cachedConfig *Config

type Config struct{ Port string }

func loadConfig() Config {
	if cachedConfig != nil {
		return *cachedConfig
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	cfg := Config{Port: port}
	cachedConfig = &cfg
	return cfg
}
