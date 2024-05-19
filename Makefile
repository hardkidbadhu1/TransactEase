working_directory := $(shell pwd)

# Path: Makefile
PROJECT := $(shell basename $(working_directory))

BUILD_VENDOR := go mod vendor

Run:
	docker-compose up -d

install_deps:
	docker-compose -f infrastructure/build.yml --project-name $(PROJECT) \
	run --rm build-env /bin/sh -c "$(BUILD_VENDOR)"

vet:
	docker-compose -f infrastructure/build.yml --project-name $(PROJECT) \
	run --rm build-env /bin/sh -c "go vet -mod=vendor ./..."

test:
	docker-compose -f infrastructure/build.yml --project-name $(PROJECT) \
	run --rm build-env /bin/sh -c "go test -tags='!integration' -mod=vendor ./..."

db_up:
	docker-compose -f docker-compose.yml up db db_migrate -d

clean:
	docker-compose down

integration_test: db_up
	docker-compose -f infrastructure/build.yml --project-name $(PROJECT) \
	run --rm build-env /bin/sh -c "go test -tags='integration' -mod=vendor ./..."

coverage:
	docker-compose -f infrastructure/test.yml --project-name $(PROJECT) run --rm test-env -t=90 -a=-mod=vendor

generate_docs:
	 swag init -g router/router.go

generate_mocks:
	go generate ./... -mod=vendor

build:
	docker-compose -f infrastructure/build.yml --project-name $(PROJECT) \
	run --rm build-env /bin/sh -c "go build -mod=vendor -o bin/$(PROJECT) cmd/main.go"

dockerise:
	docker build -t $(PROJECT) -f infrastructure/Dockerfile .
	docker build -t $(PROJECT)-migrate -f infrastructure/Migrate.Dockerfile .


clean_all:
	chmod -R +w ./.gopath vendor