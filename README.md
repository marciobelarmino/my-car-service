# Car Microservice

This is a basic microservice for managing information about cars. It provides RESTful API endpoints for retrieving, creating, and updating car information. The implementation is written in Go and uses standard libraries including net/http. It also includes some basic observability features and automated testing.

## Endpoints

The application hosts the following endpoints:

1. **GET** endpoint to retrieve an existing car by ID
2. **GET** endpoint to retrieve a list of all cars
3. **POST** endpoint to create a new car
4. **PUT** endpoint to update an existing car by ID

You can access the API documentation and interact with these endpoints using Swagger UI by visiting [http://localhost:8585/swagger-ui/#/default](http://localhost:8585/swagger-ui/#/default).

## Prerequisites

Before running the microservice, make sure you have the following prerequisites installed:

- Go (Golang)
- cURL (for testing the endpoints)

## Getting Started

Follow these steps to build, run, and test the microservice:

1. **Build the Project**:

   ```bash
   make build

2. **Start the Microservice**:

   ```bash
   make start

3. **Create a new car**:

   ```bash
   curl -X POST http://localhost:8585/cars -d '{"make": "Toyota", "model": "Camry", "package": "SE", "color": "White", "year": 2010, "category": "Sedan", "mileage": 3999, "price": 2899000, "id": "fWl37la"}'

4. **Get an existent car**:

   ```bash
    curl -X GET http://localhost:8585/cars/fWl37la

4. **Update an existent car**:

   ```bash
    curl -X PUT http://localhost:8585/cars/fWl37la -d '{"make": "Toyota", "model": "Camry", "package": "SE", "color": "Gold", "year": 2015, "category": "Sedan", "mileage": 2000, "price": 3899000, "id": "fWl37la"}'

4. **List all cars**:

   ```bash
    curl -X GET http://localhost:8585/cars
