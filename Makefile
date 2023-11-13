include .env.example
export

# HELP =================================================================================================================
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help

help: ## Display this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

compose-up: ### Run docker-compose
	docker-compose up --build -d mysql redis && docker-compose logs -f
.PHONY: compose-up

compose-up-integration-test: build-image ### Run docker-compose with integration test
	docker-compose up --build --abort-on-container-exit --exit-code-from integration
.PHONY: compose-up-integration-test

compose-down: ### Down docker-compose
	docker-compose down --remove-orphans
.PHONY: compose-down

swag: ### swag init
	swag init -g internal/controller/http/router.go
.PHONY: swag

run: swag ### swag run
	go mod tidy && go mod download && \
	DISABLE_SWAGGER_HTTP_HANDLER='' GIN_MODE=debug CGO_ENABLED=0 go run ./cmd/app
.PHONY: run

linter-golangci: ### check by golangci linter
	golangci-lint run
.PHONY: linter-golangci

linter-golangci-fix: ### check by golangci linter
	golangci-lint run --fix
.PHONY: linter-golangci-fix

test: ### run test
	go test -v -cover -race -covermode atomic -coverprofile=coverage.txt ./internal/...
.PHONY: test

integration-test: ### run integration-test
	go clean -testcache && go test -v ./integration-test/...
.PHONY: integration-test

mock: ### run mockgen
	mockgen -source=internal/service/interfaces.go -destination=./mocks/mocks_service.go -package=mocks
	mockgen -source=internal/dao/querier.go -destination=./mocks/mocks_dao.go -package=mocks
	mockgen -source=pkg/cache/interface.go -destination=./mocks/mocks_cache.go -package=mocks
.PHONY: mock

sqlc: ### Run sqlc generate
	sqlc generate
	python3 ./sql/custom_interface.py
.PHONY: sqlc

build:
	go build -o dist/go-webapp-template cmd/app/main.go
.PHONY: build

build-image: ### only build
	goreleaser release --snapshot --rm-dist
.PHONY: build-image

release: ## release
	goreleaser release --rm-dist
.PHONY: release