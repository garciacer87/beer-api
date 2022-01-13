package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/garciacer87/beer-api/internal/api"
	"github.com/garciacer87/beer-api/internal/db"
	"github.com/garciacer87/beer-api/internal/service"
	"github.com/sirupsen/logrus"
)

func main() {
	port, ok := os.LookupEnv("PORT")
	if !ok || port == "" {
		port = "8080"
		logrus.Println("env variable PORT not defined. Using default port 8080")
	}

	dbURI := os.Getenv("DATABASE_URI")
	if dbURI == "" {
		logrus.Panicf("DATABASE_URI environment variable not defined")
	}

	db, err := db.NewPostgreSQLDB(dbURI)
	if err != nil {
		logrus.Panicf("could not initialize database: %v", err)
	}

	curHelper := service.NewCurrencyClient(os.Getenv("CURRENCY_API_TOKEN"), "https://free.currconv.com/api/v7/convert")

	srv := api.NewServer(port, db, curHelper)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logrus.Error(err)
		}
	}()

	s := <-signalChan
	logrus.Infof("Signal triggered: %v", s)

	srv.Shutdown(context.Background())

	os.Exit(0)
}
