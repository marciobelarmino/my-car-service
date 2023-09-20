package main

import (
	"log"
	"net/http"

	"github.com/marciobelarmino/my-car-service/internal/carstore"
	"github.com/marciobelarmino/my-car-service/internal/server"
)

func main() {
	server := server.NewCarServer(carstore.NewInMemoryCarStore())
	log.Fatal(http.ListenAndServe(":8080", server))
}
