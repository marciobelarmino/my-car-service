package main

import (
	"encoding/json"
	"net/http"
)

// Car stores car information
type Car struct {
	Make     string
	Model    string
	Package  string
	Color    string
	Year     int
	Category string
	Mileage  int
	Price    int
	Id       string
}

const jsonContentType = "application/json"

// CarStore stores data of cars
type CarStore interface {
	GetAll() []Car
}

// CarServer is a HTTP interface for car information
type CarServer struct {
	store CarStore
	http.Handler
}

func (c *CarServer) carsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(c.store.GetAll())
}

// NewCarServer creates a CarServer with routing configured.
func NewCarServer(store CarStore) *CarServer {
	c := new(CarServer)

	c.store = store

	router := http.NewServeMux()
	router.Handle("/cars", http.HandlerFunc(c.carsHandler))

	c.Handler = router

	return c
}
