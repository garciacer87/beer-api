package db

import "github.com/garciacer87/beer-api/internal/contract"

//Database abstraction of database connection
type Database interface {
	InsertBeer(beer contract.Beer) error
	GetBeers() ([]contract.Beer, error)
	GetBeer(id int) (*contract.Beer, error)
	Close()
}
