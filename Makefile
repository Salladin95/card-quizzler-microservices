SHELL=cmd.exe
API_SERVICE_BINARY=apiApp
USER_SERVICE_BINARY=userApp

## up: starts all containers in the background without forcing build
up:
	@echo Starting Docker images...
	docker-compose up -d
	@echo Docker images started!

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_api build_user
	@echo Stopping docker images (if running...)
	docker-compose down
	@echo Building (when required) and starting docker images...
	docker-compose up --build -d
	@echo Docker images built and started!

## down: stop docker compose
down:
	@echo Stopping docker compose...
	docker-compose down
	@echo Done!

## build_api: builds the api-service binary as a linux executable
build_api:
	@echo Building broker binary...
	chdir broker && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${API_SERVICE_BINARY} ./cmd/api
	@echo Done!

## build_user: builds the user-service binary as a linux executable
build_user:
	@echo Building auth binary...
	chdir auth && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${USER_SERVICE_BINARY} ./cmd/api
	@echo Done!
