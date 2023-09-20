package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

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

func TestCar(t *testing.T) {

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
