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

.PHONY: pb
pb:
	protoc -I=./proto \
		--go_out=pkg/server \
		--go_opt=paths=source_relative \
		--go-grpc_out=pkg/server \
		--go-grpc_opt=paths=source_relative \
		./proto/*.proto

CLIENT_DIR := ./cmd/client
SERVER_DIR := ./cmd/server
DATE := $(shell date +%Y-%m-%d)
OUTPUT_DIR := bin

.PHONY: gen-mocks
gen-mocks:
	mockgen -source=internal/service/service.go -destination=internal/service/mocks/service_mock.gen.go -package=mocks

.PHONY: build
build:
	go build -ldflags '-X main.buildVersion=$(VERSION) -X main.buildDate=$(DATE)' -o $(CLIENT_DIR)/client $(CLIENT_DIR) && \
 	go build -o $(SERVER_DIR)/server $(SERVER_DIR)

.PHONY: run-server
run-server:
	$(SERVER_DIR)/server

.PHONY: run-client
run-client:
	$(CLIENT_DIR)/client

PLATFORMS = linux_amd64 windows_amd64 darwin_amd64 darwin_arm64
BIN_DIR = bin

.PHONY: build-all
build-all:
	@for platform in $(PLATFORMS); do \
		OSARCH=$$platform; \
		GOOS=$$(echo "$$OSARCH" | cut -d_ -f1); \
		GOARCH=$$(echo "$$OSARCH" | cut -d_ -f2); \
		echo "Building for $$GOOS/$$GOARCH..."; \
		go build -ldflags '-X main.buildVersion=$(version) -X main.buildDate=$(DATE)' -o $(BIN_DIR)/client_$$OSARCH $(CLIENT_DIR); \
	done && \
	echo "Building server..."; \
	go build -o $(BIN_DIR)/gophkeeper_server $(SERVER_DIR)