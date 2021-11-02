export GOPRIVATE=direct
export GO111MODULE=on
export GOSUMDB=off

BIN_NAME := $(or $(PROJECT_NAME),'api')
PKG_PATH := $(or $(PKG),'.')
PKG_LIST := $(shell go list ${PKG_PATH}/... | grep -v /vendor/)
GOLINT := golangci-lint

POSTGRES_MASTER_NAME := $(or $(POSTGRES_MASTER_NAME), micro-mrr-server-db)
POSTGRES_MASTER_USER := $(or $(POSTGRES_MASTER_USER), postgres)
POSTGRES_MASTER_PASSWORD := $(or $(POSTGRES_MASTER_PASSWORD), 12345)

POSTGRES_TEST_NAME := $(or $(POSTGRES_TEST_NAME), micro-mrr-server-test)
POSTGRES_TEST_USER := $(or $(POSTGRES_TEST_USER), postgres)
POSTGRES_TEST_PASSWORD := $(or $(POSTGRES_TEST_PASSWORD), 12345)

MIGRATE=migrate -path sql/migrations -database postgres://${POSTGRES_MASTER_USER}:${POSTGRES_MASTER_PASSWORD}@localhost:5432/${POSTGRES_MASTER_NAME}?sslmode=disable
TEST_MIGRATE=migrate -path sql/migrations -database postgres://${POSTGRES_TEST_USER}:${POSTGRES_TEST_PASSWORD}@localhost:5432/${POSTGRES_TEST_NAME}?sslmode=disable

check-lint:
	@which $(GOLINT) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.41.1

dep: # Download required dependencies
	go mod tidy
	go mod download
	go mod vendor

lint: dep check-lint ## Lint the files local env
	$(GOLINT) run -c .golangci.yml --timeout 5m

test: dep test-db-prepare ## Run unit tests
	go test -race -tags musl,unit -count=1 -short ./...

build: dep ## Build the binary file
	CGO_ENABLED=1 go build -mod=vendor -o ./bin/${BIN_NAME} -a ./src

clean: ## Remove previous build
	rm -f src/bin/$(BIN_NAME)

run: ## run application
	go run src/main.go

check-swagger:
	@which swagger || (GO111MODULE=off go get -u github.com/go-swagger/go-swagger/cmd/swagger)

swagger: check-swagger
	swagger generate spec -o src/server/http/static/swagger.yaml  -w ./ --scan-models
	swagger generate spec -o src/server/http/static/swagger.json  -w ./ --scan-models

fmt: ## format source files
	go fmt github.com/erp/api/src/...

migrate-macos-install: ## Install migration tool on MacOS
	brew install golang-migrate

migrate-linux-install: ## Install migration tool on Linux Debian
	curl -s https://packagecloud.io/install/repositories/golang-migrate/migrate/script.deb.sh | sudo bash
	sudo apt-get install migrate=4.4.0

migrate-up: submodule submodule-update ## Run migrations
	$(MIGRATE) up

migrate-down: ## Rollback migrations
	$(MIGRATE) down

test-db-prepare: ## cleanup test db
	docker exec -u postgres postgres-erp dropdb api-db-test || true
	docker exec -u postgres postgres-erp createdb api-db-test
	$(TEST_MIGRATE) up