package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

// CarsInitialData stores the initial data for InMemoryCarStore
var CarsInitialData map[string]Car = map[string]Car{
	"JHk290Xj": {"Ford", "F10", "Base", "Silver", 2010, "Truck", 120123, 1999900, "JHk290Xj"},
	"fWl37la":  {"Toyota", "Camry", "SE", "White", 2019, "Sedan", 3999, 2899000, "fWl37la"},
}

type StubCarStore struct {
	store map[string]Car
}

func (s *StubCarStore) GetAll() []Car {
	var cars []Car
	for _, car := range s.store {
		cars = append(cars, car)
	}
	return cars
}

func (s *StubCarStore) Get(id string) Car {
	return s.store[id]
}

func (s *StubCarStore) Create(car Car) (Car, error) {
	if car.Id == "" {
		return Car{}, ErrCarCreationMessage
	}

	s.store[car.Id] = car
	return s.store[car.Id], nil
}

func (s *StubCarStore) Update(id string, car Car) (Car, error) {
	carToUpdate := s.Get(id)

	// car not exists
	if carToUpdate.Id == "" {
		return Car{}, ErrCarUpdatingMessage
	}

	updateCarFromTo(&car, &carToUpdate)

	s.store[id] = carToUpdate
	return s.store[id], nil
}

func TestGETCars(t *testing.T) {

	t.Run("it returns 404 on missing cars", func(t *testing.T) {
		store := StubCarStore{CarsInitialData}
		server := NewCarServer(&store)

		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/cars/%s", "missing-key"), nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})

	t.Run("it retrieve an existing car", func(t *testing.T) {
		key := "JHk290Xj"
		want := CarsInitialData[key]

		store := StubCarStore{CarsInitialData}
		server := NewCarServer(&store)

		request, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/cars/%s", key), nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getCarFromResponse(t, response.Body)

		assertStatus(t, response.Code, http.StatusOK)
		assertCar(t, got, want)
		assertContentType(t, response, jsonContentType)
	})

	t.Run("it retrieve the list of cars", func(t *testing.T) {
		want := []Car{
			{"Ford", "F10", "Base", "Silver", 2010, "Truck", 120123, 1999900, "JHk290Xj"},
			{"Toyota", "Camry", "SE", "White", 2019, "Sedan", 3999, 2899000, "fWl37la"},
		}

		store := StubCarStore{CarsInitialData}
		server := NewCarServer(&store)

		request, _ := http.NewRequest(http.MethodGet, "/cars", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getCarsFromResponse(t, response.Body)

		assertStatus(t, response.Code, http.StatusOK)
		assertCars(t, got, want)
		assertContentType(t, response, jsonContentType)
	})
}

func TestPOSTCars(t *testing.T) {

	store := StubCarStore{map[string]Car{}}
	server := NewCarServer(&store)

	t.Run("it creates a new car when POST", func(t *testing.T) {

		want := Car{
			Id:    "Xyz123",
			Make:  "Toyota",
			Model: "Camry",
			Year:  2022,
		}

		var requestBody bytes.Buffer
		err := json.NewEncoder(&requestBody).Encode(want)
		if err != nil {
			t.Fatal(err)
		}

		request, _ := http.NewRequest(http.MethodPost, "/cars", &requestBody)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getCarFromResponse(t, response.Body)

		assertStatus(t, response.Code, http.StatusCreated)
		assertCar(t, got, want)
		assertContentType(t, response, jsonContentType)
	})

	t.Run("it returns a bad request when request data is malformed", func(t *testing.T) {
		malformedJSON := `{malformed...string}`

		request, _ := http.NewRequest(http.MethodPost, "/cars", strings.NewReader(malformedJSON))
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatus(t, response.Code, http.StatusBadRequest)
	})

}

func TestPUTCars(t *testing.T) {
	existingCar := map[string]Car{
		"Xyz1234": {
			Id:       "Xyz1234",
			Make:     "Toyota",
			Model:    "Camry",
			Year:     2022,
			Color:    "Silver",
			Category: "Sedan",
			Package:  "SE",
			Mileage:  1000,
			Price:    2999000,
		},
	}

	store := StubCarStore{existingCar}
	server := NewCarServer(&store)

	t.Run("it updates an existent car when PUT", func(t *testing.T) {
		want := Car{
			Id:       "Xyz1234",
			Make:     "Toyota",
			Model:    "Camry",
			Year:     2022,
			Color:    "Gold",
			Category: "Sedan",
			Package:  "SE",
			Mileage:  3000,
			Price:    2999000,
		}

		var requestBody bytes.Buffer
		err := json.NewEncoder(&requestBody).Encode(want)
		if err != nil {
			t.Fatal(err)
		}

		request, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("/cars/%s", want.Id), &requestBody)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		got := getCarFromResponse(t, response.Body)

		assertStatus(t, response.Code, http.StatusOK)
		assertCar(t, got, want)
		assertContentType(t, response, jsonContentType)
	})
}

func getCarFromResponse(t testing.TB, body io.Reader) (car Car) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&car)

	if err != nil {
		t.Fatalf("Unable to parse response from server %q into a Car object, '%v'", body, err)
	}

	return
}

func assertContentType(t testing.TB, response *httptest.ResponseRecorder, want string) {
	t.Helper()
	if response.Result().Header.Get("content-type") != want {
		t.Errorf("response did not have content-type of %s, got %v", want, response.Result().Header)
	}
}

func getCarsFromResponse(t testing.TB, body io.Reader) (cars []Car) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&cars)

	if err != nil {
		t.Fatalf("Unable to parse response from server %q into slice of Car, '%v'", body, err)
	}

	return
}

func assertCar(t testing.TB, got, want Car) {
	t.Helper()
	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertCars(t testing.TB, got, want []Car) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}
