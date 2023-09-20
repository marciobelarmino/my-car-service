package server

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/marciobelarmino/my-car-service/internal/carstore"
)

const JsonContentType = "application/json"

// CarServer is a HTTP interface for car information
type CarServer struct {
	store carstore.CarStore
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
	w.Header().Set("content-type", JsonContentType)
	json.NewEncoder(w).Encode(c.store.GetAll())
}

func (c *CarServer) createCarHandler(w http.ResponseWriter, r *http.Request) {
	var car carstore.Car

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

	w.Header().Set("content-type", JsonContentType)
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

	var carToUpdate carstore.Car
	err := json.NewDecoder(body).Decode(&carToUpdate)
	if err != nil {
		http.Error(w, "Invalid JSON request body", http.StatusBadRequest)
		return
	}

	carUpdated, err := c.store.Update(carId, carToUpdate)
	if err != nil {
		http.Error(w, "Error update a car", http.StatusInternalServerError)
	}

	w.Header().Set("content-type", JsonContentType)
	json.NewEncoder(w).Encode(carUpdated)
}

func (c *CarServer) getCarByIdHandler(w http.ResponseWriter, carId string) {
	car := c.store.Get(carId)

	if car.Id == "" {
		w.WriteHeader(http.StatusNotFound)
	}

	w.Header().Set("content-type", JsonContentType)
	json.NewEncoder(w).Encode(car)
}

// NewCarServer creates a CarServer with routing configured.
func NewCarServer(store carstore.CarStore) *CarServer {
	c := new(CarServer)

	c.store = store

	router := http.NewServeMux()

	// List cars and it allows create a new car
	router.Handle("/cars", http.HandlerFunc(c.carHandler))

	// Get car and it allows update a car
	router.Handle("/cars/", http.HandlerFunc(c.carByIdHandler))

	// Serve Swagger UI
	router.Handle("/swagger.json", http.FileServer(http.Dir(".")))
	router.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir("./swagger-ui"))))

	c.Handler = router

	return c
}
