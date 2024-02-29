SHELL=cmd.exe
API_SERVICE_BINARY=apiApp
USER_SERVICE_BINARY=userApp
LOGGING_SERVICE_BINARY=loggingApp
MAIL_SERVICE_BINARY=mailApp
CARD_QUIZ_SERVICE_BINARY=cardQuizApp

## up: starts all containers in the background without forcing build
up:
	@echo Starting Docker images...
	docker-compose up -d
	@echo Docker images started!

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_api build_card_quiz build_user build_mail build_logging
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
	@echo Building api-service binary...
	chdir api-service && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${API_SERVICE_BINARY} ./cmd/api
	@echo Done!

## build_card_quiz: builds the card-quiz-service binary as a linux executable
build_card_quiz:
	@echo Building card-quizzler-service binary...
	chdir card-quizzler-service && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${CARD_QUIZ_SERVICE_BINARY} ./cmd/api
	@echo Done!

## build_user: builds the user-service binary as a linux executable
build_user:
	@echo Building user-serivce binary...
	chdir user-service && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${USER_SERVICE_BINARY} ./cmd/api
	@echo Done!

## build_logging: builds the logging-service binary as a linux executable
build_logging:
	@echo Building logging-serivce binary...
	chdir logging-service && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${LOGGING_SERVICE_BINARY} ./cmd/api
	@echo Done!

## build_mail: builds the mail-service binary as a linux executable
build_mail:
	@echo Building mail-service binary...
	chdir mail-service && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${MAIL_SERVICE_BINARY} ./cmd/api
	@echo Done!
