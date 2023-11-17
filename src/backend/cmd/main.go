package main

import (
	"log"
	"net/http"

	"github.com/AlexSkyBell/push-the-button/internal/server"
)

func main() {
	log.SetFlags(log.Lshortfile)

	server := server.NewServer()
	go server.Listen()

	log.Fatal(http.ListenAndServe(":18080", nil))
}
