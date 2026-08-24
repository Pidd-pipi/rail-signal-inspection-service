package main

func main() {
	config := loadConfig()
	if err := serveHTTP(buildServer(config, newSignalStore())); err != nil {
		panic(err)
	}
}
