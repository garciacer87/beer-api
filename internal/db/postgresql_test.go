package db

import (
	"fmt"
	"testing"

	"github.com/garciacer87/beer-api/internal/contract"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var dbURI = "postgres://beerapi:beerapi@localhost:5432/beerapitest"

func initTestDB(t *testing.T) *migrate.Migrate {
	m, err := migrate.New(
		"file://../../sql/postgresql",
		fmt.Sprintf("%s?sslmode=disable", dbURI),
	)

	if err != nil {
		t.Fatalf("could not create new migrate object: %s", err)
	}
	if err := m.Up(); err != nil {
		t.Fatalf("could not up migrate %s", err)
	}

	return m
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

func TestNewPostgreSQL(t *testing.T) {
	tests := map[string]struct {
		dbURI       string
		errExpected bool
	}{
		"#1: invalid URI": {dbURI: "postgres://wrong:wrong@wrong:5432/wrongdb", errExpected: true},
		"#2: valid case":  {dbURI: dbURI, errExpected: false},
	}

	for desc, tc := range tests {
		db, err := NewPostgreSQLDB(tc.dbURI)
		isErr := err != nil

		if isErr != tc.errExpected {
			t.Errorf("%s:\n got Error? %v.\n Error expected? %v.\n Error: %v", desc, isErr, tc.errExpected, err)
		}

		if db != nil {
			db.Close()
		}
	}
}

func TestInsertBeer(t *testing.T) {
	m := initTestDB(t)

	defer func() {
		if err := m.Down(); err != nil {
			t.Fatalf("could not down migrate %s", err)
		}
	}()

	db, err := NewPostgreSQLDB(dbURI)
	if err != nil {
		t.Fatalf("could not init database connection: %s", err)
	}

	defer db.Close()

	db.InsertBeer(getMockBeer())

	tests := map[string]struct {
		beer        contract.Beer
		errExpected bool
	}{
		"#1: duplicated beer id": {
			beer:        contract.Beer{ID: 1, Name: "golden", Brewery: "kross", Country: "Chile", Price: 1000, Currency: "CLP"},
			errExpected: true,
		},

		"#2: valid case": {
			beer:        contract.Beer{ID: 2, Name: "gran torobayo", Brewery: "kunstmann", Country: "Chile", Price: 1000, Currency: "CLP"},
			errExpected: false,
		},
	}

	for desc, tc := range tests {
		err = db.InsertBeer(tc.beer)
		isErr := err != nil
		if isErr != tc.errExpected {
			t.Errorf("%s:\n got Error? %v.\n Error expected? %v.\n Error: %v", desc, isErr, tc.errExpected, err)
		}
	}
}

func TestGetBeers(t *testing.T) {
	m := initTestDB(t)
	defer func() {
		if err := m.Down(); err != nil {
			t.Fatalf("could not down migrate %s", err)
		}
	}()

	db, err := NewPostgreSQLDB(dbURI)
	if err != nil {
		t.Fatalf("could not init database connection: %s", err)
	}

	defer db.Close()

	beers, err := db.GetBeers()
	if err != nil {
		t.Fatal()
	}

	if len(beers) != 0 {
		t.Errorf("#1: no beers. Beer slice must be empty")
	}

	db.InsertBeer(getMockBeer())

	beers, err = db.GetBeers()
	if err != nil {
		t.Fatal()
	}

	if len(beers) != 1 {
		t.Errorf("#2: Should be one beer in the slice")
	}
}

func TestGetBeer(t *testing.T) {
	m := initTestDB(t)
	defer func() {
		if err := m.Down(); err != nil {
			t.Fatalf("could not down migrate %s", err)
		}
	}()

	db, err := NewPostgreSQLDB(dbURI)
	if err != nil {
		t.Fatalf("could not init database connection: %s", err)
	}

	defer db.Close()

	db.InsertBeer(getMockBeer())

	tests := map[string]struct {
		id          int
		errExpected error
	}{
		"#1: valid case": {
			id:          1,
			errExpected: nil,
		},
		"#2: beer not found": {
			id:          123,
			errExpected: &NotFoundError{},
		},
	}

	for desc, tc := range tests {
		beer, err := db.GetBeer(tc.id)

		if tc.errExpected != err {
			t.Errorf("%s:\n Error expected: %v\n Error got: %v", desc, tc.errExpected, err)
		}

		if beer != nil && beer.ID != 1 {
			t.Errorf("%s\n expected different beer ID", desc)
		}
	}
}
