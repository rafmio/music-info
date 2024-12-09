/*
 * Music info
 * API version: 0.0.1
 */
package main

import (
	// sw "getcode/swagger"
	"log"
	"musicinfo/server"
	"net/http"
)

func main() {
	log.Printf("Server started on port 8080")

	router := server.NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
