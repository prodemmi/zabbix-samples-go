package main

import (
	"log"
	"net/http"
)

func main() {
	userMetrics := RegisterUserMetrics()
	go SetupUserSeeder(userMetrics)

	pMux := http.NewServeMux()
	pMux.HandleFunc("/ping", func(r http.ResponseWriter, req *http.Request) {
		r.Write([]byte("PONG"))
	})

	// Alert Testing
	go func() {
		// time.Sleep(20 * time.Second)
		// os.Exit(0)
	}()

	log.Fatal(http.ListenAndServe(":3000", pMux))
}
