BIN_CALENDAR := "./bin/calendar"

MAIN_CALENDAR := "./cmd/calendar"

CONFIG_CALENDAR := "./configs/calendar/config.toml"

DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.50.0

lint: install-lint-deps
	golangci-lint run ./...

generate:
	protoc --plugin=protoc-gen-go=/Users/yu.zadorozhnaya/go/bin/protoc-gen-go --go_out=internal/server/grpc --plugin=protoc-gen-go-grpc=/Users/yu.zadorozhnaya/go/bin/protoc-gen-go-grpc --go-grpc_out=internal/server/grpc internal/server/grpc/EventService.proto

build:
	go build -v -o $(BIN_CALENDAR) -ldflags "$(LDFLAGS)" $(MAIN_CALENDAR)

build-calendar:
	go build -v -o $(BIN_CALENDAR) -ldflags "$(LDFLAGS)" $(MAIN_CALENDAR)

run-calendar: build-calendar
	$(BIN_CALENDAR) -config $(CONFIG_CALENDAR)

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build-calendar
	$(BIN_CALENDAR) version

migrate-up:
	docker exec -it postgresdb createdb --username postgres calendardb
	migrate -path migrations -database "postgres://postgres:1234512345@localhost:5432/calendardb?sslmode=disable" -verbose up

.PHONY: up down \
		generate rabbitmq postgre \
		install-lint-deps lint test \
		build build-calendar build-scheduler build-sender \
 		run-calendar run-scheduler run-sender \
 		build-img run-img version migrateup