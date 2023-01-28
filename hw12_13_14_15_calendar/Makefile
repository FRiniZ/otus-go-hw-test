BIN1 := "./bin/calendar"
BIN2 := "./bin/calendar_scheduler"
BIN3 := "./bin/calendar_sender"

DOCKER_IMG1="calendar:develop"
DOCKER_IMG2="calendar_scheduler:develop"
DOCKER_IMG3="calendar_sender:develop"
DOCKER_IMG4="integration_tests:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

RABBITMQNAME := rabbitmq

Red='\033[0;31m'
Green='\033[0;32m'
Color_Off='\033[0m'

help:
	@echo ${Red}"Please select a subcommand"${Color_Off}
	@echo ${Green}"make build"${Color_Off}" to build applications"
	@echo ${Green}"make run_calendar"${Color_Off}" to run calendar"
	@echo ${Green}"make run_scheduler"${Color_Off}" to run scheduler"
	@echo ${Green}"make run_sender"${Color_Off}" to run sender"
	@echo 
	@echo ${Green}"make generate"${Color_Off}" to generate stub-files from protobuf-files"



generate:
	go generate ./...

build:
	go build -v -o $(BIN1) -ldflags "$(LDFLAGS)" ./cmd/calendar
	go build -v -o $(BIN2) -ldflags "$(LDFLAGS)" ./cmd/scheduler
	go build -v -o $(BIN3) -ldflags "$(LDFLAGS)" ./cmd/sender

run_calendar: build
	$(BIN1) -config=./configs/calendar_config.toml
run_scheduler: build
	$(BIN2) -config=./configs/scheduler_config.toml
run_sender: build
	$(BIN3) -config=./configs/sender_config.toml



build-img-calendar:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG1) \
		-f build/calendar/Dockerfile .	

build-img-scheduler:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG2) \
		-f build/scheduler/Dockerfile .	

build-img-sender:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG3) \
		-f build/sender/Dockerfile .	

build-img-tests:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG4) \
		-f build/integration-tests/Dockerfile .	

build-img: build-img-calendar build-img-scheduler build-img-sender build-img-tests

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

stop_rabbitmq:
	docker stop $(RABBITMQNAME)


run_rabbitmq:
	docker run -d --rm --hostname rabbitmq --name $(RABBITMQNAME) -e RABBITMQ_DEFAULT_USER=root -e RABBITMQ_DEFAULT_PASS=secretkey -p 15672:15672 -p 5672:5672 rabbitmq:3-management

run_postgres:
	docker run -d --rm --hostname postgres --network calendar --name postgres -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secretkey -e PGDATA=/var/lib/postgresql/data/pgdata -v psqldata:/var/lib/postgresql/data -p 5432:5432 postgres

create_db:
	docker exec -it postgres createdb --username=root --owner=root calendar

drop_db:
	docker exec -it postgres dropdb calendar

migrate_up:
	migrate -path migrations/db -database "postgresql://root:secretkey@localhost:5432/calendar?sslmode=disable" -verbose up

migrate_down:
	migrate -path migrations/db -database "postgresql://root:secretkey@localhost:5432/calendar?sslmode=disable" -verbose down


up:
	@docker-compose -f ./deployments/docker-compose.yaml up -d --scale tests=0

down:
	@docker-compose -f ./deployments/docker-compose.yaml down

destroy:
	@docker-compose -f ./deployments/docker-compose.yaml down -v


integration-tests: up
	@docker-compose -f ./deployments/docker-compose.yaml up tests --exit-code-from tests  && { docker-compose -f ./deployments/docker-compose.yaml down -v; echo "Success!"; } || { docker-compose -f ./deployments/docker-compose.yaml down -v; echo "Failure!"; exit 1; }

.PHONY: build run build-img run-img version test lint run_postgres createdb dropdb

