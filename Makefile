# user-management-api/Makefile
include .env
export
APP_NAME=user-management-api
MAIN_PATH=cmd/api/main.go
HOST=http://localhost:8086/api/v1
API_KEY=william-hehe


DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

.PHONY: run build test clean dev rate-limit export_sql import_sql create_db \
		create_migration migrate_up migrate_down

migrate_up:
	migrate -path internal/db/migrations -database "$(DB_URL)" up

# make migrate_down step=1
migrate_down:
	migrate -path internal/db/migrations -database "$(DB_URL)" down $(step)

# make migrate_force version=1
migrate_force:
	migrate -path internal/db/migrations -database "$(DB_URL)" force $(version)

# make create_migration name=create_profiles_table
# make create_migration name=add_phone_to_users
# make create_migration name=rename_phone_column
# make create_migration name=add_password_to_users
create_migration:
	migrate create -ext sql -dir internal/db/migrations -seq $(name)

export_sql:
	pg_dump -U $(DB_USER) -h $(DB_HOST) -p $(DB_PORT) $(DB_NAME) > backup.sql

import_sql:
	psql -U $(DB_USER) -h $(DB_HOST) -p $(DB_PORT) $(DB_NAME) < backup.sql
create_db:
	psql -U $(DB_USER) -h $(DB_HOST) -p $(DB_PORT) -d postgres -c 'CREATE DATABASE "$(DB_NAME)";'

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