package api

import (
	"fmt"

	"github.com/garciacer87/beer-api/internal/contract"
	"github.com/garciacer87/beer-api/internal/db"
)

type mockDB struct {
	throwError bool
	beerCount  int
}

func (mdb *mockDB) InsertBeer(beer contract.Beer) error {
	if mdb.throwError {
		return fmt.Errorf("mocked error")
	}
	return nil
}

func (mdb *mockDB) GetBeers() ([]contract.Beer, error) {
	if mdb.throwError {
		return nil, fmt.Errorf("mock error")
	}

	beers := []contract.Beer{}
	for i := 0; i < mdb.beerCount; i++ {
		beers = append(beers, getMockBeer())
	}

	return beers, nil
}

func (mdb *mockDB) GetBeer(_ int) (*contract.Beer, error) {
	if mdb.beerCount == 0 {
		if mdb.throwError {
			return nil, fmt.Errorf("mocked error")
		}
		return nil, &db.NotFoundError{}
	}

	beer := getMockBeer()

	return &beer, nil
}

func (mdb *mockDB) Close() {}

type mockCurrencyHelper struct{}

func (mch *mockCurrencyHelper) GetBoxPrice(beer contract.Beer, to string, quantity int) (float64, error) {
	if to == "unknown" {
		return 0, fmt.Errorf("invalid currency")
	}

	if quantity == 0 {
		quantity = 6
	}

	rate := 0.0012
	result := beer.Price * rate * float64(quantity)

	return result, nil
}

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
