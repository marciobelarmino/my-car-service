package main

import (
	"log"
	"net/http"
)

func main() {
	server := NewCarServer(NewInMemoryCarStore())
	log.Fatal(http.ListenAndServe(":8080", server))
}
