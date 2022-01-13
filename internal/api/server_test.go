package api

import (
	"context"
	"net/http"
	"testing"
)

func TestListenAndServe(t *testing.T) {
	var (
		ch  = make(chan error)
		mdb = &mockDB{}
		mch = &mockCurrencyHelper{}
		srv = NewServer("8081", mdb, mch)
	)

	defer srv.Shutdown(context.Background())

	go srv.ListenAndServe() //serving in port 8081

	resp, err := http.Get("http://localhost:8081/health")
	if err != nil {
		t.Fatalf("Error not expected: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected code response 200")
	}

	newSrv := NewServer("8081", mdb, mch)

	go func() {
		ch <- newSrv.ListenAndServe() //trying to serve in the already used port 8081
	}()

	err = <-ch
	if err == nil {
		t.Errorf("error expected. Cannot serve in a port already used")
	}
}
