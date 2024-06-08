BIN := "./bin/calendar"
MAIN := "./cmd"
CONFIG := "./config/config.toml"

LINT_EXCLUDE_DIRS="hw12_13_14_15_calendar/internal/server"

DOCKER_IMG="develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

test:
	go test -race ./internal/...

install-lint-deps:
	export GOROOT=$(go env GOROOT)
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.55.2

lint: install-lint-deps
	golangci-lint run ./... --exclude=".*internal/server/.*"

generate:
	protoc --go_out=internal/server/grpc/api --go-grpc_out=internal/server/grpc/ internal/server/grpc/EventService.proto

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" $(MAIN)

run: build
	$(BIN) -config $(CONFIG)

version: build
	$(BIN) version

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

postgres:
	docker run --name postgresdb --env POSTGRES_PASSWORD="1234512345" --publish "5436:5432" --detach --rm postgres

migrate-up:
	docker exec -it postgresdb createdb --username postgres calendardb
	migrate -path migrations -database "postgres://postgres:1234512345@localhost:5432/calendardb?sslmode=disable" -verbose up

run-img: build-img
	docker run $(DOCKER_IMG)

up:
	cd deployments && docker-compose -f docker-compose.yaml build && docker-compose -f docker-compose.yaml up

down:
	cd deployments && docker-compose down

.PHONY: up down \
		generate rabbitmq postgre \
		install-lint-deps lint test \
		build build-calendar build-scheduler build-sender \
 		run-calendar run-scheduler run-sender \
 		build-img run-img version migrateup