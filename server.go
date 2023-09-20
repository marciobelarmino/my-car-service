package main

import (
	"encoding/json"
	"io"
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
	Update(id string, car Car) (Car, error)
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

func (c *CarServer) carByIdHandler(w http.ResponseWriter, r *http.Request) {
	carId := strings.TrimPrefix(r.URL.Path, "/cars/")

	switch r.Method {
	case http.MethodGet:
		c.getCarByIdHandler(w, carId)
	case http.MethodPut:
		c.putCarByIdHandler(w, carId, r.Body)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (c *CarServer) putCarByIdHandler(w http.ResponseWriter, carId string, body io.Reader) {
	existingCar := c.store.Get(carId)

	if existingCar.Id == "" {
		w.WriteHeader(http.StatusNotFound)
	}

	var carToUpdate Car
	err := json.NewDecoder(body).Decode(&carToUpdate)
	if err != nil {
		http.Error(w, "Invalid JSON request body", http.StatusBadRequest)
		return
	}

	carUpdated, err := c.store.Update(carId, carToUpdate)
	if err != nil {
		http.Error(w, "Error update a car", http.StatusInternalServerError)
	}

	w.Header().Set("content-type", jsonContentType)
	json.NewEncoder(w).Encode(carUpdated)
}

func (c *CarServer) getCarByIdHandler(w http.ResponseWriter, carId string) {
	car := c.store.Get(carId)

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
	router.Handle("/cars/", http.HandlerFunc(c.carByIdHandler))

	c.Handler = router

	return c
}
