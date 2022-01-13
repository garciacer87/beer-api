package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/garciacer87/beer-api/internal/db"
	"github.com/garciacer87/beer-api/internal/service"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

//Server abstraction of a server
type Server interface {
	ListenAndServe() error
	Shutdown(context.Context) error
}

type server struct {
	httpPort   string
	httpServer *http.Server
	db         db.Database
	curClient  service.CurrencyClient
}

//NewServer creates a new server object
func NewServer(port string, db db.Database, curClient service.CurrencyClient) Server {
	r := mux.NewRouter()

	r.HandleFunc("/health", healthHandler).Methods(http.MethodGet)

	srv := &server{
		httpPort:  port,
		db:        db,
		curClient: curClient,
	}

	beersRouter := r.PathPrefix("/beers").Subrouter()
	beersRouter.HandleFunc("", validateBeer(srv.insertBeer)).Methods(http.MethodPost)
	beersRouter.HandleFunc("", srv.getBeers).Methods(http.MethodGet)
	beersRouter.HandleFunc("/{beerID}", validateExistence(db, srv.getBeer)).Methods(http.MethodGet)
	beersRouter.HandleFunc("/{beerID}/boxprice", validateExistence(db, srv.getBoxPrice)).Methods(http.MethodGet)

	srv.httpServer = &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%v", port),
		Handler: r,
	}

	return srv
}

//ListenAndServe starts the http server on the previously configurated port
func (s *server) ListenAndServe() error {
	logrus.Printf("serving on port %s\n", s.httpPort)
	return s.httpServer.ListenAndServe()
}

//Shutdown the http server
func (s *server) Shutdown(ctx context.Context) error {
	logrus.Infof("Shutting down API server")

	// close DB connection
	s.db.Close()

	return s.httpServer.Shutdown(ctx)
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("healthy"))
}
