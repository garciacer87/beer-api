package db

import (
	"context"
	"fmt"

	"github.com/garciacer87/beer-api/internal/contract"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
)

//PostgreSQLDB implementation of postgresql database
type PostgreSQLDB struct {
	pool *pgxpool.Pool
}

// NewPostgreSQLDB retrieves a new PostgreSQLDB object
func NewPostgreSQLDB(dbURI string) (Database, error) {
	pool, err := pgxpool.Connect(context.Background(), dbURI)
	if err != nil {
		return nil, fmt.Errorf("could not create database connection: %v", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("database is unreachable: %v", err)
	}
	logrus.Info("Connection succesfully to database")

	return &PostgreSQLDB{pool}, nil
}

// Close close connections from the pool
func (db *PostgreSQLDB) Close() {
	logrus.Info("Closing database connections")
	db.pool.Close()
}

//InsertBeer creates a new beer in database
func (db *PostgreSQLDB) InsertBeer(beer contract.Beer) error {
	query := "INSERT INTO public.beer(id, name, brewery, country, price, currency) VALUES($1, $2, $3, $4, $5, $6)"

	_, err := db.pool.Exec(context.Background(), query, beer.ID, beer.Name, beer.Brewery, beer.Country, beer.Price, beer.Currency)
	if err != nil {
		logrus.Infof("type: %T\n", err)
		if pgerr, ok := err.(*pgconn.PgError); ok {
			logrus.Info(pgerr.ConstraintName)
			if pgerr.ConstraintName == "beer_pkey" {
				return &DuplicateKeyError{}
			}
		}
		return err
	}

	return nil
}

//GetBeers retrieves all the beers stored
func (db *PostgreSQLDB) GetBeers() ([]contract.Beer, error) {
	query := "SELECT id, name, brewery, country, price, currency FROM public.beer"

	var (
		id                               int
		name, brewery, country, currency string
		price                            float64
	)

	rows, err := db.pool.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	beers := make([]contract.Beer, 0)
	for rows.Next() {
		if err = rows.Scan(&id, &name, &brewery, &country, &price, &currency); err != nil {
			return nil, fmt.Errorf("could not get beers: %v", err)
		}
		beers = append(beers, contract.Beer{
			ID:       id,
			Name:     name,
			Brewery:  brewery,
			Country:  country,
			Currency: currency,
			Price:    price,
		})
	}

	return beers, nil
}

//GetBeer retrieves a beer by its ID
func (db *PostgreSQLDB) GetBeer(id int) (*contract.Beer, error) {
	query := "SELECT name, brewery, country, price, currency FROM public.beer WHERE id = $1"

	var (
		name, brewery, country, currency string
		price                            float64
	)

	row := db.pool.QueryRow(context.Background(), query, id)
	err := row.Scan(&name, &brewery, &country, &price, &currency)
	if err != nil {
		switch err.Error() {
		case "no rows in result set":
			return nil, &NotFoundError{}
		default:
			return nil, err
		}
	}

	return &contract.Beer{
		ID:       id,
		Name:     name,
		Brewery:  brewery,
		Country:  country,
		Price:    price,
		Currency: currency,
	}, nil
}
