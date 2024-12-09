/*
 * Music info
 * API version: 0.0.1
 */
package main

import (
	"log"
	"musicinfo/server"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("config/dbconf.env")
	if err != nil {
		log.Println("error loading.env file:", err)
		return
	}
	port := os.Getenv("MUSIC_INFO_PORT")
	port = ":" + port
	log.Printf("Server started on port: %s", port)

	router := server.NewRouter()

	log.Fatal(http.ListenAndServe(port, router))
}
