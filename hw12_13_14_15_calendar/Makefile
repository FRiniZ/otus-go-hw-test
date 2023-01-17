BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

generate:
	go generate ./cmd/calendar/main.go

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) -config ./configs/config.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

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


run_postgres:
	docker run -d --rm --name postgres -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secretkey -e PGDATA=/var/lib/postgresql/data/pgdata -v psqldata:/var/lib/postgresql/data -p 5432:5432 postgres

create_db:
	docker exec -it postgres createdb --username=root --owner=root calendar

drop_db:
	docker exec -it postgres dropdb calendar

migrate_up:
	migrate -path migrations/db -database "postgresql://root:secretkey@localhost:5432/calendar?sslmode=disable" -verbose up

migrate_down:
	migrate -path migrations/db -database "postgresql://root:secretkey@localhost:5432/calendar?sslmode=disable" -verbose down

.PHONY: build run build-img run-img version test lint run_postgres createdb dropdb

