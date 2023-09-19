package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type StubCarStore struct {
	store []Car
}

func (s *StubCarStore) GetAll() []Car {
	return s.store
}

func TestCar(t *testing.T) {
	t.Run("it returns the list of cars", func(t *testing.T) {
		wantedCars := []Car{
			{"Ford", "F10", "Base", "Silver", 2010, "Truck", 120123, 1999900, "JHk290Xj"},
			{"Toyota", "Camry", "SE", "White", 2019, "Sedan", 3999, 2899000, "fWl37la"},
		}

		store := StubCarStore{wantedCars}
		server := NewCarServer(&store)

		request, _ := http.NewRequest(http.MethodGet, "/cars", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		var cars []Car

		err := json.NewDecoder(response.Body).Decode(&cars)
		if err != nil {
			t.Fatalf("Unable to parse response from server %q into slice of Car, '%v'", response.Body, err)
		}

		assertStatus(t, response.Code, http.StatusOK)

		if !reflect.DeepEqual(cars, wantedCars) {
			t.Errorf("got %v want %v", cars, wantedCars)
		}

		if response.Result().Header.Get("content-type") != jsonContentType {
			t.Errorf("response did not have content-type of %s, got %v", jsonContentType, response.Result().Header)
		}
	})
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}
