BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) -config ./configs/config.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

pull-postgres:
	docker pull postgres

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

postgres:
	docker run --name postgresdb --env POSTGRES_PASSWORD="1234512345" --publish "5436:5432" --detach --rm postgres

createdb:
	docker exec -it postgresdb createdb --username postgres calendardb

migrateup:
	migrate -path migrations -database "postgresql://postgres:1234512345@localhost:5436/calendardb?sslmode=disable" -verbose up

test:
	go test -race ./internal/config ./internal/logger ./internal/storage/memory

migratedown:
	migrate -path migrations -database "postgresql://postgres:1234512345@localhost:5436/calendardb?sslmode=disable" -verbose down

dropdb:
	docker exec -it postgresdb dropdb calendardb

.PHONY: build run build-img run-img version lint pull-postgres postgres createdb migrateup test migratedown dropdb