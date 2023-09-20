package carstore

import (
	"github.com/marciobelarmino/my-car-service/errors"
)

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
		return Car{}, errors.ErrCarCreationMessage
	}

	i.store[car.Id] = car
	return i.store[car.Id], nil
}

// Update a car informing an id
func (i *InMemoryCarStore) Update(id string, car Car) (Car, error) {
	carToUpdate := i.Get(id)

	// car not exists
	if carToUpdate.Id == "" {
		return Car{}, errors.ErrCarUpdatingMessage
	}

	UpdateCarFromTo(&car, &carToUpdate)

	i.store[id] = carToUpdate
	return i.store[id], nil
}

func UpdateCarFromTo(from *Car, to *Car) {
	if from.Make != "" {
		to.Make = from.Make
	}

	if from.Model != "" {
		to.Model = from.Model
	}

	if from.Year != 0 {
		to.Year = from.Year
	}

	if from.Color != "" {
		to.Color = from.Color
	}

	if from.Category != "" {
		to.Category = from.Category
	}

	if from.Package != "" {
		to.Package = from.Package
	}

	if from.Mileage != 0 {
		to.Mileage = from.Mileage
	}

	if from.Price != 0 {
		to.Price = from.Price
	}

}
