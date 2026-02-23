BIN_CALENDAR := "./bin/calendar"
BIN_SCHEDULER := "./bin/scheduler"
BIN_SENDER := "./bin/sender"

CALENDAR_CONFIG := "./configs/calendar_config.yaml"
SCHEDULER_CONFIG := "./configs/scheduler_config.yaml"
SENDER_CONFIG := "./configs/sender_config.yaml"

DOCKER_IMG_CALENDAR="calendar:develop"
DOCKER_IMG_SCHEDULER="scheduler:develop"
DOCKER_IMG_SENDER="sender:develop"
 
GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

#Postrges
POSTGRES_USER ?= postgres
POSTGRES_PASSWORD ?= password
POSTGRES_DB ?= backend
POSTGRES_PORT ?= 5435
POSTGRES_CONTAINER := postgres-calendar

#Rabbit
RABBIT_USER ?= rabbit
RABBIT_PASSWORD ?= password
RABBIT_PORT ?= 5672
RABBIT_PORT_UI ?= 15672
RABBIT_CONTAINER := rabbitmq


run-postgres:
	docker run -d --name $(POSTGRES_CONTAINER) \
	-e POSTGRES_USER=$(POSTGRES_USER) \
	-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
	-e POSTGRES_DB=$(POSTGRES_DB) \
	-p $(POSTGRES_PORT):5432 \
	-v postgres-data:/var/lib/postgresql/data \
	postgres:latest

run-rabbit:
	docker run -d --name $(RABBIT_CONTAINER) \
	-p $(RABBIT_PORT):5672 \
	-p $(RABBIT_PORT_UI):15672 \
	-e RABBITMQ_DEFAULT_USER=$(RABBIT_USER) \
	-e RABBITMQ_DEFAULT_PASS=$(RABBIT_PASSWORD) \
	rabbitmq:management

generate:
	protoc --go_out=./internal/server/grpc --go-grpc_out=./internal/server/grpc ./api/*.proto

build-calendar:
	go build -v -o $(BIN_CALENDAR) -ldflags "$(LDFLAGS)" ./cmd/calendar

build-scheduler:
	go build -v -o $(BIN_SCHEDULER) -ldflags "$(LDFLAGS)" ./cmd/scheduler

build-sender:
	go build -v -o $(BIN_SENDER) -ldflags "$(LDFLAGS)" ./cmd/sender

build: build-calendar build-scheduler build-sender

run-calendar:
	$(BIN_CALENDAR) -config ./configs/calendar_config.yaml

run-scheduler:
	$(BIN_SCHEDULER) -config ./configs/scheduler_config.yaml

run-sender:
	$(BIN_SENDER) -config ./configs/sender_config.yaml

run:
	@echo "Starting all services..."
	$(BIN_CALENDAR) --config $(CALENDAR_CONFIG) &
	$(BIN_SCHEDULER) --config $(SCHEDULER_CONFIG) &
	$(BIN_SENDER) --config $(SENDER_CONFIG) &

build-calendar-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG_CALENDAR) \
		-f build/Dockerfile.calendar .

build-scheduler-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG_SCHEDULER) \
		-f build/Dockerfile.scheduler .

build-sender-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG_SENDER) \
		-f build/Dockerfile.sender .

run-calendar-img: build-calendar-img
	docker run $(DOCKER_IMG_CALENDAR)

run-scheduler-img: build-scheduler-img
	docker run $(DOCKER_IMG_SCHEDULER)

run-sender-img: build-sender-img
	docker run $(DOCKER_IMG_SENDER)

docker-build: build-calendar-img build-scheduler-img build-sender-img

up:
	docker-compose -f deployments/docker-compose.yaml up --build -d

down:
	docker-compose -f deployments/docker-compose.yaml down -v --remove-orphans

integration-tests:
	@echo "Starting services..."
	@docker-compose -f test/integration/docker-compose.yaml up -d --build postgres rabbitmq goose calendar scheduler sender
	@echo "Running integration tests..."
	@docker-compose -f test/integration/docker-compose.yaml up --build --abort-on-container-exit test; \
		EXIT_CODE=$$?; \
	docker-compose -f test/integration/docker-compose.yaml down -v --remove-orphans; \
	exit $$EXIT_CODE

version: build
	$(BIN) version

test:
	go test -race -v ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.64.2

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: generate lint test version build run docker-build docker-run up down

