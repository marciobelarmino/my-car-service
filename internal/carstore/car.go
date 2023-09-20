package carstore

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

// CarStore stores data of cars
type CarStore interface {
	Get(id string) Car
	GetAll() []Car
	Create(car Car) (Car, error)
	Update(id string, car Car) (Car, error)
}
