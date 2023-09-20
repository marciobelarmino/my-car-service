package main

import (
	"errors"
)

var ErrCarCreationMessage = errors.New("unable to create a car without id")
var ErrCarUpdatingMessage = errors.New("unable to update a car without id")

// InMemoryCarStore collects data about cars in memory.
type InMemoryCarStore struct {
	store map[string]Car
}

// NewInMemoryCarStore initilizes an empty car store.
func NewInMemoryCarStore() *InMemoryCarStore {
	return &InMemoryCarStore{map[string]Car{}}
}

// GetAll retrieves all cars from the store
func (i *InMemoryCarStore) GetAll() []Car {
	var cars []Car
	for _, car := range i.store {
		cars = append(cars, car)
	}
	return cars
}

// Get retrieve car from the store by id
func (i *InMemoryCarStore) Get(id string) Car {
	return i.store[id]
}

// Create a new car
func (i *InMemoryCarStore) Create(car Car) (Car, error) {

	if car.Id == "" {
		return Car{}, ErrCarCreationMessage
	}

	i.store[car.Id] = car
	return i.store[car.Id], nil
}

// Update a car informing an id
func (i *InMemoryCarStore) Update(id string, car Car) (Car, error) {
	carToUpdate := i.Get(id)

	// car not exists
	if carToUpdate.Id == "" {
		return Car{}, ErrCarUpdatingMessage
	}

	updateCarFromTo(&car, &carToUpdate)

	i.store[id] = carToUpdate
	return i.store[id], nil
}
