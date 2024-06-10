BIN_CALENDAR := "./bin/calendar"
BIN_SCHEDULER := "./bin/scheduler"
BIN_SENDER := "./bin/sender"

MAIN_CALENDAR := "./cmd/calendar"
MAIN_SCHEDULER := "./cmd/calendar_scheduler"
MAIN_SENDER := "./cmd/calendar_sender"

CONFIG_CALENDAR := "./config/calendar_config.toml"
CONFIG_SCHEDULER := "./config/scheduler_config.toml"
CONFIG_SENDER := "./config/sender_config.toml"

DOCKER_IMG_CALENDAR := "calendar-app"
DOCKER_IMG_SCHEDULER := "scheduler-app"
DOCKER_IMG_SENDER := "sender-app"

LINT_EXCLUDE_DIRS="hw12_13_14_15_calendar/internal/server"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

test:
	go test -race ./internal/...

integration-tests:
	go test ./integration_tests

install-lint-deps:
	export GOROOT=$(go env GOROOT)
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.55.2

lint: install-lint-deps
	golangci-lint run ./... --exclude=".*internal/server/.*"

generate:
	protoc --go_out=internal/server/grpc/api --go-grpc_out=internal/server/grpc/ internal/server/grpc/EventService.proto

build:
	go build -v -o $(BIN_CALENDAR) -ldflags "$(LDFLAGS)" $(MAIN_CALENDAR)
	go build -o $(BIN_SENDER) $(MAIN_SENDER)
	go build -o $(BIN_SCHEDULER) $(MAIN_SCHEDULER)

build-calendar:
	go build -v -o $(BIN_CALENDAR) -ldflags "$(LDFLAGS)" $(MAIN_CALENDAR)

build-scheduler:
	go build -o $(BIN_SCHEDULER) $(MAIN_SCHEDULER)

build-sender:
	go build -o $(BIN_SENDER) $(MAIN_SENDER)

run-calendar: build-calendar
	$(BIN_CALENDAR) -config $(CONFIG_CALENDAR)

run-scheduler: build-scheduler
	$(BIN_SCHEDULER) -config $(CONFIG_SCHEDULER)

run-sender: build-sender
	$(BIN_SENDER) -config $(CONFIG_SENDER)

version: build-calendar
	$(BIN_CALENDAR) version

build-img-calendar:
	docker build --build-arg LDFLAGS="$(LDFLAGS)" -t $(DOCKER_IMG_CALENDAR) -f build/calendar/Dockerfile .

build-img-scheduler:
	docker build --build-arg LDFLAGS="$(LDFLAGS)" -t $(DOCKER_IMG_SCHEDULER) -f build/calendar_scheduler/Dockerfile .

build-img-sender:
	docker build --build-arg LDFLAGS="$(LDFLAGS)" -t $(DOCKER_IMG_SENDER) -f build/calendar_sender/Dockerfile .

postgres:
	docker run --name postgresdb --env POSTGRES_PASSWORD="1234512345" --publish "5436:5432" --detach --rm postgres

rabbitmq:
	docker run --name rabbitmq --detach --rm -p 15672:15672 -p 5672:5672 rabbitmq:3-management

migrate-up:
	docker exec -it postgresdb createdb --username postgres calendardb
	migrate -path migrations -database "postgres://postgres:1234512345@localhost:5432/calendardb?sslmode=disable" -verbose up

migrate-down:
	migrate -path migrations -database "postgres://postgres:1234512345@localhost:5436/calendardb?sslmode=disable" -verbose down

run-img: build-img
	docker run $(DOCKER_IMG)

up:
	cd deployments && docker-compose up --build -d

down:
	cd deployments && docker-compose down

integration-tests:
	cd deployments && docker-compose -f docker-compose-test.yaml build && docker-compose -f docker-compose.test.yaml up

.PHONY: up down \
		generate rabbitmq postgre \
		install-lint-deps lint test \
		build build-calendar build-scheduler build-sender \
 		run-calendar run-scheduler run-sender \
 		build-img run-img version migrateup