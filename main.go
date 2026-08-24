package main

func main() {
	config := loadConfig()
	server := newEnterpriseServer(":"+config.Port, newServer(newSignalStore()))
	if err := serveHTTP(server); err != nil {
		panic(err)
	}
}
