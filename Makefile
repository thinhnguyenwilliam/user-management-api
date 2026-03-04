# user-management-api/Makefile
APP_NAME=user-management-api
MAIN_PATH=cmd/api/main.go
HOST=http://localhost:8086/api/v1
API_KEY=william-hehe

.PHONY: run build test clean dev rate-limit

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

rate-limit:
	hey -n 50 -c 20 \
		-H "X-API-Key: $(API_KEY)" \
		$(HOST)/users/1772594263538205928