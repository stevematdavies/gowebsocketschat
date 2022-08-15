package main

import (
	"log"
	"net/http"
	"stevematdavies/websockets/chat/internal/handlers"
)

func main() {
	mux := routes()

	log.Println("Starting channel listener")
	go handlers.ListenToWsChannel()
	log.Println("Starting Webserver on port 8080")
	_ = http.ListenAndServe(":8080", mux)
}
