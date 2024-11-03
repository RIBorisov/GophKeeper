.PHONY: lint
lint: _golangci-lint-rm-unformatted-report

.PHONY: _golangci-lint-reports-mkdir
_golangci-lint-reports-mkdir:
	mkdir -p ./golangci-lint

.PHONY: _golangci-lint-run
_golangci-lint-run: _golangci-lint-reports-mkdir
	-golangci-lint run -c .golangci.yml > ./golangci-lint/report-unformatted.json

.PHONY: _golangci-lint-format-report
_golangci-lint-format-report: _golangci-lint-run
	cat ./golangci-lint/report-unformatted.json | jq > ./golangci-lint/report.json

.PHONY: _golangci-lint-rm-unformatted-report
_golangci-lint-rm-unformatted-report: _golangci-lint-format-report
	rm ./golangci-lint/report-unformatted.json

.PHONY: golangci-lint-clean
golangci-lint-clean:
	sudo rm -rf ./golangci-lint

.PHONY: migration
migration: #  example: make migration name=add-smth
	docker run --rm \
    -v $(realpath ./internal/storage/migrations):/migrations \
    migrate/migrate:v4.16.2 \
        create \
        -dir /migrations \
        -ext .sql \
        -seq -digits 3 \
        $(name)

DSN="postgresql://admin:password@localhost:5432/gophkeeper?sslmode=disable"
.PHONY: migrate-up
migrate-up:
	docker run --rm \
    -v $(realpath ./internal/storage/migrations):/migrations \
    migrate/migrate:v4.16.2 \
        -path=/migrations \
        -database $(DSN) \
        up

.PHONY: migrate-down
migrate-down:
	docker run --rm \
    -v $(realpath ./internal/storage/migrations):/migrations \
    migrate/migrate:v4.16.2 \
        -path=/migrations \
        -database $(DSN) \
        down 1

CLIENT_DIR := cmd/client
SERVER_DIR := cmd/server
APP_NAME := gophkeeper
COMMIT_HASH := $(shell git rev-parse --short=8 HEAD)
DATE := $(shell date +%Y-%m-%d)

.PHONY: server
server:
	echo "I AM SERVER ECHO"

.PHONY: client
client:
	echo "I AM CLIENT ECHO"

.PHONY: build
build: client server
	echo "I AM BUILD ECHO AFTER CLIENT && SERVER"
#	cd $(DIR) && \
#	go build -ldflags "-X main.buildVersion=$(version) -X main.buildDate=$(DATE) -X main.buildCommit=$(COMMIT_HASH)" -o $(APP_NAME)
#	cd $(DIR) && ./$(APP_NAME)

.PHONY: pb
pb:
	protoc -I=./proto \
		--go_out=pkg/server \
		--go_opt=paths=source_relative \
		--go-grpc_out=pkg/server \
		--go-grpc_opt=paths=source_relative \
		./proto/*.proto
