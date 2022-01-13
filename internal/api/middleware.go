package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/garciacer87/beer-api/internal/contract"
	"github.com/garciacer87/beer-api/internal/db"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

//validates beer fields
func validateBeer(next http.HandlerFunc) http.HandlerFunc {
	beerValidator := newValidator()

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		bodyBytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			logrus.Errorf("could not decode the body %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		body := ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		beer := contract.Beer{}

		err = json.NewDecoder(body).Decode(&beer)
		if err != nil {
			logrus.Errorf("could not decode the body %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//validates beer fields from decoded body
		err = beerValidator.Struct(beer)
		if err != nil {
			errs := beerValidator.translate(err)
			logrus.Printf("Validation error(s):\n%s", strings.Join(errs, " | "))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		next(w, req)
	})
}

func validateExistence(dbConn db.Database, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		idStr := mux.Vars(req)["beerID"]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			logrus.Errorf("id is not valid")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err = dbConn.GetBeer(id)
		if err != nil {
			logrus.Errorf("could not get beer: %v", err)
			if _, ok := err.(*db.NotFoundError); ok {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		next(w, req)
	})
}
