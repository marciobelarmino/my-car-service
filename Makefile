test:
	go test ./internal/server -v

build:
	go build -o my-car-service ./cmd/my-car-service/main.go

start:
	go run ./cmd/my-car-service/main.go

create:
	curl -X POST http://localhost:8080/cars -d '{"make": "Toyota", "model": "Camry", "package": "SE", "color": "While", "year": 2010, "category": "Sedan", "mileage": 3999, "price": 2899000, "id": "fWl37la"}'

update:
	curl -X PUT http://localhost:8080/cars/fWl37la -d '{"make": "Toyota", "model": "Camry", "package": "SE", "color": "Gold", "year": 2015, "category": "Sedan", "mileage": 2000, "price": 3899000, "id": "fWl37la"}'

get:
	curl -X GET http://localhost:8080/cars/fWl37la

list:
	curl -X GET http://localhost:8080/cars
