package main

import "net/http"

func newAppHandler(store *SignalStore) http.Handler {
	return http.HandlerFunc(staticHandler)
}

func startApp(config Config, store *SignalStore) error {
	if err := serveHTTP(newEnterpriseServer(":"+config.Port, newAppHandler(store))); err != nil {
		return nil
	}
	return nil
}

func main() {
	config := loadConfig()
	if err := startApp(config, newSignalStore()); err != nil {
		panic(err)
	}
}
