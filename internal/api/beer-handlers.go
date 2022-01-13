package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/garciacer87/beer-api/internal/contract"
	"github.com/garciacer87/beer-api/internal/db"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func (s *server) insertBeer(w http.ResponseWriter, req *http.Request) {
	beer := contract.Beer{}
	json.NewDecoder(req.Body).Decode(&beer)

	//inserts beer into database
	err := s.db.InsertBeer(beer)
	if err != nil {
		logrus.Errorf("could not insert beer: %v", err)
		if _, ok := err.(*db.DuplicateKeyError); ok {
			w.WriteHeader(http.StatusConflict)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	logrus.Infof("Beer %v inserted", beer.ID)
	w.WriteHeader(http.StatusCreated)
}

func (s *server) getBeer(w http.ResponseWriter, req *http.Request) {
	idStr := mux.Vars(req)["beerID"]

	id, _ := strconv.Atoi(idStr)

	beer, _ := s.db.GetBeer(id)

	bodyResp, _ := json.Marshal(&beer)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bodyResp)
}

func (s *server) getBeers(w http.ResponseWriter, req *http.Request) {
	beers, err := s.db.GetBeers()
	if err != nil {
		logrus.Errorf("could not get beers: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(beers) == 0 {
		logrus.Info("no beers found in database")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	bodyResp, _ := json.Marshal(&beers)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bodyResp)
}

func (s *server) getBoxPrice(w http.ResponseWriter, req *http.Request) {
	idStr := mux.Vars(req)["beerID"]
	id, _ := strconv.Atoi(idStr)

	beer, _ := s.db.GetBeer(id)

	qParams := req.URL.Query()

	curTo := qParams["currency"][0]

	qty := 6
	var err error
	if qParams["quantity"][0] != "" {
		qty, err = strconv.Atoi(qParams["quantity"][0])
		if err != nil {
			logrus.Errorf("invalid value for quantity: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	result, err := s.curClient.GetBoxPrice(*beer, curTo, qty)
	if err != nil {
		logrus.Errorf("could not get box price: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bodyResp, _ := json.Marshal(&contract.BoxPriceResp{Total: result})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bodyResp)
}
