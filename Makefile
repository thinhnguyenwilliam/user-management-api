# user-management-api/Makefile
APP_NAME=user-management-api
MAIN_PATH=cmd/api/main.go

.PHONY: run build test clean dev

dev:
	air
run:
	go run $(MAIN_PATH)

build:
	go build -o bin/$(APP_NAME) $(MAIN_PATH)

test:
	go test ./...

clean:
	rm -rf bin/