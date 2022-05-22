SHELL=/bin/zsh
APP_NAME=webscraper

.DEFAULT_GOAL=help

build: clean # compiles go binaries locally
	go build -o ./bin/${APP_NAME} cmd/main.go

dbuild: # builds go & pg containers
	cd .. && docker compose build

run: # run application locally from compiled binary (requires local pg instance)
	cd bin/ && ./${APP_NAME} -env="../../.env.testing"

up: dbuild # run application (webscraper & pg db) inside docker
	cd .. && docker compose up

test: # tests all packages recursively (requires local pg instance)
	go test -v ./... | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/''

clean: # cleans out legacy binaries
	go clean

help: # shows help message
	@egrep -h '\s#\s' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?# "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
