package main

import (
	"encoding/json"
	"io"
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
		want := []Car{
			{"Ford", "F10", "Base", "Silver", 2010, "Truck", 120123, 1999900, "JHk290Xj"},
			{"Toyota", "Camry", "SE", "White", 2019, "Sedan", 3999, 2899000, "fWl37la"},
		}

		store := StubCarStore{want}
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
