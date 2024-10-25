package main

import (
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}
func main() {
	http.HandleFunc("/", handler)

	log.Println("Connected to server on localhost:8080")
	if err := http.ListenAndServe("localhost:8080", nil); err != nil {
		log.Fatal(err)
	}
}
