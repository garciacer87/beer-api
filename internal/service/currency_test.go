package service

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/garciacer87/beer-api/internal/contract"
)

func getMockBeer() contract.Beer {
	return contract.Beer{
		ID:       1,
		Name:     "golden",
		Brewery:  "kross",
		Country:  "Chile",
		Price:    1000,
		Currency: "CLP",
	}
}

func okResponseHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"CLP_USD":2.2}`))
}

func badResponseHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}

func emptyResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestGetBoxPrice(t *testing.T) {
	tests := map[string]struct {
		srv           *httptest.Server
		apiToken      string
		beer          contract.Beer
		currentyTo    string
		qty           int
		totalExpected float64
		errExpected   bool
	}{
		"#1: bad response": {
			srv:         httptest.NewServer(http.HandlerFunc(badResponseHandler)),
			apiToken:    "asdf",
			beer:        getMockBeer(),
			currentyTo:  "asdf",
			qty:         1,
			errExpected: true,
		},
		"#2: empty response": {srv: httptest.NewServer(http.HandlerFunc(emptyResponse)),
			beer:        getMockBeer(),
			apiToken:    "asdf",
			currentyTo:  "asdf",
			qty:         1,
			errExpected: true,
		},
		"#3: empty token": {srv: httptest.NewServer(http.HandlerFunc(emptyResponse)),
			apiToken:    "",
			errExpected: true,
		},
		"#3: ok response": {srv: httptest.NewServer(http.HandlerFunc(okResponseHandler)),
			beer:          getMockBeer(),
			apiToken:      "asdf",
			currentyTo:    "usd",
			qty:           2,
			totalExpected: 2.2,
			errExpected:   false,
		},
	}

	for desc, tc := range tests {
		defer tc.srv.Close()

		cc := NewCurrencyClient(tc.apiToken, tc.srv.URL)
		total, err := cc.GetBoxPrice(tc.beer, tc.currentyTo, tc.qty)
		isErr := err != nil

		if isErr != tc.errExpected {
			t.Errorf("%s: error expected? %v\n error got? %v", desc, tc.errExpected, isErr)
		}

		if isErr != false && total != tc.totalExpected {
			t.Errorf("%s: total expected: %v\n total got: %v", desc, tc.totalExpected, total)
		}
	}

}
