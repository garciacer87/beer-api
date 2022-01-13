# Beer API

This is a simple API to manage basic REST operations. <br/>
The service was made in [Go](https://go.dev/) and the database was implemented in [Postgresql](https://www.postgresql.org/).
Also, this service connects to an [external API](https://free.currencyconverterapi.com/) to get conversion rates between currencies.
Follow this link to get documentation about the API: https://www.currencyconverterapi.com/docs

## Setting up dev environment
### Tooling

In order to test this app, you must install the following dependencies:
* [go](https://go.dev/)
* [migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)
* [docker](https://www.docker.com/)
  
For code analysis:
* revive: go install github.com/mgechev/revive@latest
* staticcheck: go install honnef.co/go/tools/cmd/staticcheck@latest

<br/>

### Environment Variables
The application requires the following environment variables:
* **PORT:** server port. e.g.: 8080
* **DATABASE_URI:** e.g.: postgres://beerapi:beerapi@localhost:5432/beerapi
* **CURRENCY_API_TOKEN:** token needed to establish connection to [currency converter API](https://free.currencyconverterapi.com/)

<br/>

## Build and run

### 1. Locally using Go
To build and run this app locally using native Go, you need a running postgresql database in your localhost. To achieve that, go to the root of the app and run the following command from available from the makefile:

```console
make create-postgres
```

It will download the official postgresql docker image based in alpine.Then, it will create and run the docker container. In a couple of seconds, there will be available a postgresql database running in a docker container with the following params:
* Database name: beerapi
* User: beerapi
* password: beerapi

Now, to create the structure of the database, run the following command:

```console
make migrate-up
```

Finally, run the following commands to build and run the API (Remember to set the environment variables mentioned before):

```console
make go-build
make run
```

<br/>

### 2. Using docker-compose 
With this approach, you will use the docker-compose file located in the root. This will create a container with the database ready to use, and also, it will run the API from a docker container. Run the following command: 

```console
docker-compose up
```
<br/>

## Endpoints

* Retrieves all the beer stored:

```bash
curl --request GET 'http://localhost:8080/beers'
```

* Creates a new beer:

```bash
curl --request POST 'http://localhost:8080/beers' \
--header 'Content-Type: application/json' \
--data-raw '{
    "id": 1,
    "name": "golden",
    "brewery": "kross",
    "country": "chile",
    "price": 1000.00,
    "currency": "CLP"
}'
```

* Get a beer by ID:

```bash
curl --request GET 'http://localhost:8080/beers/{beerID}'
```

* Get box price of a beer according to quantity and currency

```bash
curl --request GET 'http://localhost:8080/beers/1/boxprice?quantity=<quantity_value>&currency=<currency_value>'

<br/>

## Unit testing

To locally run the unit tests, first execute the following to create the portable postgresql database:

```console
make create-postgresql
make create-test-unit-db
```

Then, to run the unit tests:

```console
make test
```

