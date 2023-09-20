package main

import (
	"encoding/json"
	"net/http"
	"strings"
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
	Get(id string) Car
	GetAll() []Car
}

// CarServer is a HTTP interface for car information
type CarServer struct {
	store CarStore
	http.Handler
}

func (c *CarServer) listOfCarsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(c.store.GetAll())
}

func (c *CarServer) getCarByIdHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/cars/")
	car := c.store.Get(id)

	if car.Id == "" {
		w.WriteHeader(http.StatusNotFound)
	}

	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(car)
}

// NewCarServer creates a CarServer with routing configured.
func NewCarServer(store CarStore) *CarServer {
	c := new(CarServer)

	c.store = store

	router := http.NewServeMux()
	router.Handle("/cars", http.HandlerFunc(c.listOfCarsHandler))
	router.Handle("/cars/", http.HandlerFunc(c.getCarByIdHandler))

	c.Handler = router

	return c
}
