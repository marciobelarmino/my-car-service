package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCar(t *testing.T) {
	t.Run("it returns 200 status code", func(t *testing.T) {

		server := NewCarServer()

		request, _ := http.NewRequest(http.MethodGet, "/cars", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)

	})
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}
