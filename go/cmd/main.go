package main

import (
	"log"
	"net/http"

	"github.com/1602077/webscraper/go/pkg/server"
)

func main() {
	router := server.NewRouter()
	log.Fatal(http.ListenAndServe(":8999", router))
}
