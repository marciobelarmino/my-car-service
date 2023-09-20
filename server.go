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
	Create(car Car) (Car, error)
}

// CarServer is a HTTP interface for car information
type CarServer struct {
	store CarStore
	http.Handler
}

func (c *CarServer) carHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		c.listOfCarsHandler(w, r)
	case http.MethodPost:
		c.createCarHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (c *CarServer) listOfCarsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(c.store.GetAll())
}

func (c *CarServer) createCarHandler(w http.ResponseWriter, r *http.Request) {
	var car Car

	err := json.NewDecoder(r.Body).Decode(&car)
	if err != nil {
		http.Error(w, "Invalid JSON request body", http.StatusBadRequest)
		return
	}

	carCreated, err := c.store.Create(car)
	if err != nil {
		http.Error(w, "Error creating a car", http.StatusBadRequest)
		return
	}

	w.Header().Set("content-type", jsonContentType)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(carCreated)
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
	router.Handle("/cars", http.HandlerFunc(c.carHandler))
	router.Handle("/cars/", http.HandlerFunc(c.getCarByIdHandler))

	c.Handler = router

	return c
}
