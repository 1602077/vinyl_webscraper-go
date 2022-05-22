SHELL=/bin/zsh
APP=webscraper

.DEFAULT_GOAL=help

docker_run: docker_build # run application (webscraper & pg db) inside docker instances
	cd .. && docker compose up

docker_build: # builds go & pg containers
	cd .. && docker compose build

build: clean # compiles go binaries
	go build -o bin/${APP} cmd/main.go

run: # runs webscraper locally (requires local pg instance)
	go run cmd/main.go -env="../.env.testing"

test: # tests all packages recursively (requires local pg instance)
	go test -v ./... | sed ''/PASS/s//$$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/''

clean: # cleans out legacy binaries
	go clean

help: # shows help message
	@egrep -h '\s#\s' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?# "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
