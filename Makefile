include .env

debug:
	go run ${ENGINE} service --authport ${RPC_AUTH_PORT} --movieport ${RPC_MOVIE_PORT} --gwport ${GATEWAY_PORT}
.PHONY: debug

build:
	@echo "Building app"
	go build -o ${BUILD_DIR}/app ${ENGINE}
	@echo "Success build app. Your app is ready to use in 'build/' directory."
.PHONY: build

proto-gen:	clean-proto
	@echo "Generating the stubs"
	./scripts/proto-gen.sh
	@echo "Success generate stubs. All stubs created are in the 'stubs/' directory"
	@echo "Generating the Swagger UI"
	./scripts/swagger-ui-gen.sh
	@echo "Success generate Swagger UI. If you want to change Swagger UI to previous version copy the previous version from './cache/swagger-ui' directory"
	@echo "You can try swagger-ui with command 'make debug'"
	@echo "DO NOT EDIT ANY FILES STUBS!"
.PHONY: proto-gen

ssl-gen:
	@echo "Generating ssl configuration"
	./scripts/ssl-gen.sh
	@echo "Success generate ssl configuration. All SSL Configuration created in the 'ssl/' directory"
	@echo "DO NOT EXPOSE SSL DIRECTORY!"
.PHONY: ssl-gen

dependency:
	@echo "Downloading all Go dependencies needed"
	go mod download
	go mod verify
	go mod tidy
	@echo "All Go dependencies was downloaded. you can run 'make debug' to compile locally or 'make build' to build app."
.PHONY: dependency

clean-proto:
	@echo "Delete all previous stubs ..."
	rm -rf stubs/auth/* stubs/movie/*
	@echo "All stubs successfully deleted"
.PHONY: clean-proto

tidy:
	@echo "Synchronize dependency"
	go mod tidy
	@echo "Finish Synchronize dependency"
.PHONY: tidy

lint:
	golangci-lint run ./...
.PHONY: lint

migrate-up:
	@echo "Starting migrations up"
	migrate -database ${CONN_STRING} -path ${MIGRATION_PATH} up
	@echo "Finish migration up"
.PHONY: migrate-up

migrate-down:
	@echo "Starting migrations down"
	migrate -database ${CONN_STRING} -path ${MIGRATION_PATH} down
	@echo "Finish migrations down"
.PHONY: migrate-down

migrate-down-all:
	@echo "Starting migrations down"
	migrate -database ${CONN_STRING} -path ${MIGRATION_PATH} down -all
	@echo "Finish migrations down"
.PHONY: migrate-down-all

seed-up:
	@echo "Starting seeders up"
	migrate -database ${CONN_STRING} -path ${SEEDER_PATH} up
	@echo "Finish migration up"
.PHONY: seed-up

seed-down:
	@echo "Starting seeders down"
	migrate -database ${CONN_STRING} -path ${SEEDER_PATH} down
	@echo "Finish seeders down"
.PHONY: seed-down

seed-down-all:
	@echo "Starting seeders down"
	migrate -database ${CONN_STRING} -path ${SEEDER_PATH} down -all
	@echo "Finish seeders down"
.PHONY: seed-down-all

generate:
	@echo "Starting generate all tags go:generate"
	go generate ./...
	@echo "Finish generate all tags"
.PHONY: generate