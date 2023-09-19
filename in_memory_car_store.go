package main

// CarsInitialData stores the initial data for InMemoryCarStore
var CarsInitialData map[string]Car = map[string]Car{
	"JHk290Xj": {"Ford", "F10", "Base", "Silver", 2010, "Truck", 120123, 1999900, "JHk290Xj"},
	"fWl37la":  {"Toyota", "Camry", "SE", "White", 2019, "Sedan", 3999, 2899000, "fWl37la"},
}

// InMemoryCarStore collects data about cars in memory.
type InMemoryCarStore struct {
	store map[string]Car
}

// NewInMemoryCarStore initilizes an empty car store.
func NewInMemoryCarStore() *InMemoryCarStore {
	return &InMemoryCarStore{CarsInitialData}
}

// GetAll retrieves all cars from the store
func (i *InMemoryCarStore) GetAll() []Car {
	var cars []Car
	for _, car := range i.store {
		cars = append(cars, car)
	}
	return cars
}
