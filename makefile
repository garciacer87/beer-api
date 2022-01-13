PRINT= printf "%b"
BLUE=\033[0;94m
NC=\033[0m

create-postgres:
	@$(PRINT) "$(BLUE)Starting PostgreSQL$(NC)\n"
	docker run --name beer-db -e POSTGRES_USER=beerapi -e POSTGRES_PASSWORD=beerapi -d -p 5432:5432 postgres:13-alpine

drop-db:
	@$(PRINT) "$(BLUE)Dropping DB$(NC)\n"
	docker exec beer-db dropdb --username=beerapi beerapi

create-test-unit-db:
	@$(PRINT) "$(BLUE)Creating DB for tests$(NC)\n"
	docker exec beer-db createdb --username=beerapi --owner=beerapi beerapitest

drop-test-unit-db:
	@$(PRINT) "$(BLUE)Dropping tests DB$(NC)\n"
	docker exec beer-db dropdb --username=beerapi beerapitest

drop-postgres:
	@$(PRINT) "$(BLUE)Stopping and removing postgresql container$(NC)\n"
	docker stop beer-db
	docker container rm beer-db

migrate-up:
	@$(PRINT) "$(BLUE)Running Database migration...$(NC)\n"
	migrate -path sql/postgresql/ -database "postgresql://beerapi:beerapi@localhost:5432/beerapi?sslmode=disable" -verbose up

migrate-down:
	@$(PRINT) "$(BLUE)Undoing Database migration...$(NC)\n"
	migrate -path sql/postgresql/ -database "postgresql://beerapi:beerapi@localhost:5432/beerapi?sslmode=disable" -verbose down -all

migrate-drop:
	@$(PRINT) "$(BLUE)Dropping migration...$(NC)\n"
	migrate -path sql/postgresql/ -database "postgresql://beerapi:beerapi@localhost:5432/beerapi?sslmode=disable" -verbose drop -f

vet:
	@$(PRINT) "$(BLUE)Vetting the source code...$(NC)\n"
	go vet ./...

revive:
	@$(PRINT) "$(BLUE)Running revive...$(NC)\n"
	revive -formatter friendly ./...

staticcheck:
	@$(PRINT) "$(BLUE)Running Static Code Analysys...$(NC)\n"
	staticcheck ./...

check: vet revive staticcheck

go-build:
	go build -o ./target/beer-api ./cmd/api

docker-build:
	GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ./target/beer-api ./cmd/api
	docker build -f build/docker/dockerfile -t beer-api:1.0.0 .

run:
	./target/beer-api

test:
	go test -cover ./...

.PHONY: build run check