package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/garciacer87/beer-api/internal/contract"
)

func TestInsertBeer(t *testing.T) {
	mockDB := &mockDB{}
	mockCH := &mockCurrencyHelper{}
	srv := NewServer("8081", mockDB, mockCH)

	defer func(srv Server) {
		if err := srv.Shutdown(context.Background()); err != nil {
			t.Fatalf("could not shutdown the test server")
		}
	}(srv)

	go srv.ListenAndServe()

	tests := map[string]struct {
		beer           contract.Beer
		statusExpected int
	}{
		"#1: invalid id": {
			beer:           contract.Beer{Name: "golden", Brewery: "kross", Country: "Chile", Price: 1000, Currency: "CLP"},
			statusExpected: http.StatusBadRequest,
		},
		"#9: database error": {
			beer:           getMockBeer(),
			statusExpected: http.StatusInternalServerError,
		},
		"#10: valid case": {
			beer:           getMockBeer(),
			statusExpected: http.StatusCreated,
		},
	}

	for desc, tc := range tests {
		mockDB.throwError = tc.statusExpected == http.StatusInternalServerError

		body, _ := json.Marshal(&tc.beer)

		resp, err := http.Post("http://localhost:8081/beers", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Errorf("Error not expected: %v", err)
		}

		if resp.StatusCode != tc.statusExpected {
			t.Errorf("%s Status expected: %v Status got: %v", desc, tc.statusExpected, resp.StatusCode)
		}
	}
}

func TestGetBeers(t *testing.T) {
	mockCH := &mockCurrencyHelper{}

	tests := map[string]struct {
		srv            Server
		statusExpected int
		beersExpected  int
	}{
		"#1: valid case":            {srv: NewServer("8081", &mockDB{beerCount: 2}, mockCH), statusExpected: http.StatusOK, beersExpected: 2},
		"#2: internal server error": {srv: NewServer("8081", &mockDB{throwError: true}, mockCH), statusExpected: http.StatusInternalServerError},
		"#3: empty list":            {srv: NewServer("8081", &mockDB{}, mockCH), statusExpected: http.StatusNotFound},
	}

	for desc, tc := range tests {
		go tc.srv.ListenAndServe()
		resp, err := http.Get("http://localhost:8081/beers")
		if err != nil {
			t.Fatalf("Error not expected: %v", err)
		}

		if resp.StatusCode != tc.statusExpected {
			t.Errorf("%s:\n response code different than expected\n Got: %v\n Expected: %v", desc, resp.StatusCode, tc.statusExpected)
		}

		if resp.StatusCode == http.StatusOK {
			respBody, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("could not get response body %v", err)
			}

			var beers []contract.Beer
			err = json.Unmarshal(respBody, &beers)
			if err != nil {
				t.Fatalf("could not get response body %v", err)
			}

			if len(beers) != tc.beersExpected {
				t.Errorf("there must be %v beers in the slice", tc.beersExpected)
			}

			resp.Body.Close()
		}

		if err := tc.srv.Shutdown(context.Background()); err != nil {
			t.Fatalf("could not shutdown the test server")
		}
	}
}

func TestGetBeer(t *testing.T) {
	mockCH := &mockCurrencyHelper{}

	tests := map[string]struct {
		db             *mockDB
		id             int
		statusExpected int
	}{
		"#1: internal server error": {db: &mockDB{throwError: true, beerCount: 0}, id: 1, statusExpected: http.StatusInternalServerError},
		"#2: beer not found":        {db: &mockDB{beerCount: 0}, id: 2, statusExpected: http.StatusNotFound},
		"#3: valid case":            {db: &mockDB{beerCount: 1}, id: 3, statusExpected: http.StatusOK},
	}

	for desc, tc := range tests {
		srv := NewServer("8081", tc.db, mockCH)
		go srv.ListenAndServe()

		url := fmt.Sprintf("http://localhost:8081/beers/%v", tc.id)
		resp, err := http.Get(url)

		if err != nil {
			t.Fatalf("error not expected")
		}

		if resp.StatusCode != tc.statusExpected {
			t.Errorf("%s:\n Status code got: %v\n Status code expected: %v", desc, resp.StatusCode, tc.statusExpected)
		}

		if err := srv.Shutdown(context.Background()); err != nil {
			t.Fatalf("could not shutdown the test server")
		}
	}
}

func TestGetBoxPrice(t *testing.T) {
	mockCH := &mockCurrencyHelper{}

	tests := map[string]struct {
		db             *mockDB
		id             int
		currencyTo     string
		qty            int
		statusExpected int
		totalExpected  float64
	}{
		"#1: internal server error": {db: &mockDB{throwError: true}, id: 1, statusExpected: http.StatusInternalServerError},
		"#2: unknown currency":      {db: &mockDB{beerCount: 1}, id: 2, currencyTo: "unknown", statusExpected: http.StatusInternalServerError},
		"#3: valid case":            {db: &mockDB{beerCount: 1}, id: 3, currencyTo: "USD", qty: 4, statusExpected: http.StatusOK, totalExpected: 4.8},
	}

	for desc, tc := range tests {
		srv := NewServer("8081", tc.db, mockCH)
		go srv.ListenAndServe()

		url := fmt.Sprintf("http://localhost:8081/beers/%v/boxprice?currency=%s&quantity=%v", tc.id, tc.currencyTo, tc.qty)
		resp, err := http.Get(url)

		if err != nil {
			t.Fatalf("error not expected")
		}

		if resp.StatusCode != tc.statusExpected {
			t.Errorf("%s:\n Status code got: %v\n Status code expected: %v", desc, resp.StatusCode, tc.statusExpected)
		}

		if resp.StatusCode == http.StatusOK {
			bytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("could not get response body %v", err)
			}

			var body map[string]interface{}
			err = json.Unmarshal(bytes, &body)
			if err != nil {
				t.Fatalf("could not get response body %v", err)
			}

			total := body["totalPrice"].(float64)

			if total != tc.totalExpected {
				t.Errorf("%s:\n Total got: %v\n Total expected: %v", desc, total, tc.totalExpected)
			}

			resp.Body.Close()
		}

		if err := srv.Shutdown(context.Background()); err != nil {
			t.Fatalf("could not shutdown the test server")
		}
	}
}
