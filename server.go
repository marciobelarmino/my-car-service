package main

import "net/http"

// CarServer is a HTTP interface for car information
type CarServer struct {
	http.Handler
}

func (c *CarServer) carsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// NewCarServer creates a CarServer with routing configured.
func NewCarServer() *CarServer {
	c := new(CarServer)

	router := http.NewServeMux()
	router.Handle("/cars", http.HandlerFunc(c.carsHandler))

	c.Handler = router

	return c
}
