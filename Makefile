# Determine which environment file to use
ifneq ("$(wildcard .env)", "")
	ENV_FILE=.env
else
	ENV_FILE=app.env
endif

# Source Environment variables
include $(ENV_FILE)
export $(shell sed 's/=.*//' $(ENV_FILE))

# -- Docker Network --
network:
	docker network create $(NETWORK_NAME)

# -- DB Targets --
stop-postgres:
	docker stop $(POSTGRES_CONTAINER_NAME)

run-postgres:
	docker run --name $(POSTGRES_CONTAINER_NAME) --network $(NETWORK_NAME) -p 5432:5432 \
	-e POSTGRES_USER=$(POSTGRES_USER) \
	-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
	-d $(POSTGRES_IMAGE)
	@echo "Waiting for PostgreSQL to start..."
	@until docker exec $(POSTGRES_CONTAINER_NAME) pg_isready --username=$(POSTGRES_USER); do \
		sleep 1; \
		echo "Waiting for PostgreSQL to be ready..."; \
	done

# Create the zen_bank database
createdb:
	@echo "Creating database..."
	@until docker exec -it $(POSTGRES_CONTAINER_NAME) \
		createdb --username=$(POSTGRES_USER) --owner=$(POSTGRES_USER) $(POSTGRES_DB); do \
		sleep 1; \
	done
	@echo "Database created ✅"

verify:
	docker exec $(POSTGRES_CONTAINER_NAME) psql -U $(POSTGRES_USER) -c '\l'

migrate-up:
	migrate -path $(SCHEMA_DIR) -database $(DB_SOURCE) -verbose up

migrate-down:
	migrate -path $(SCHEMA_DIR) -database $(DB_SOURCE) -verbose down

# create new migration file, accepts name var
migrate-create:
	@if [ -z "$(name)" ]; then \
		echo "Error: 'name' variable is required. Usage: make migrate-create name=<filename>"; \
		exit 1; \
	else \
		migrate create -ext sql -dir db/migration -seq $(name).sql; \
	fi

startdb: network run-postgres createdb migrate-up
	@echo "Migrations complete ✅"

dropdb:
	docker exec $(POSTGRES_CONTAINER_NAME) dropdb $(POSTGRES_DB)

# sqlc
sqlc:
	sqlc generate

test:
	go test ./...

test-cov:
	go test -v -cover -short ./...

test-cov-mem:
	go test -v -cover ./... -gcflags '-m -l'

test-bench:
	go test -bench=. -benchmem -benchtime=10s

mock:
	mockgen -package mockdb -destination db/mocks/store.go github.com/kaayce/zen-bank/db/sqlc Store

mod:
	go mod tidy && go mod vendor

build-mem:
	go build -gcflags '-m -l'

server:
	go run main.go

start:
	startdb test server

run-app:
	go build && chmod +x zen-bank && ./zen-bank

reset:
	@if [ "$(GIN_MODE)" = "test" ] || [ "$(GIN_MODE)" = "debug" ]; then \
		echo "Dropping database $(POSTGRES_DB)..."; \
		docker exec $(POSTGRES_CONTAINER_NAME) dropdb $(POSTGRES_DB); \
		echo "Dropping Network $(NETWORK_NAME)..."; \
		docker rm $(NETWORK_NAME); \
		echo "Stopping and removing container $(POSTGRES_CONTAINER_NAME)..."; \
		docker stop $(POSTGRES_CONTAINER_NAME) && docker rm $(POSTGRES_CONTAINER_NAME); \
	else \
		echo "Not allowed in production environment"; \
	fi

# docker
build-docker-app:
	docker build -t zenbank:latest .

run-docker-app:
	docker run --name zenbank --network zenbank-network -p 8080:8080 -e GIN_MODE=release \
	-e DB_SOURCE="postgresql://pat:secret@postgres12:5432/zen_bank?sslmode=disable" \
	zenbank:latest

.PHONY: network startdb dropdb migrate-up migrate-down migrate-create sqlc test server reset stop-postgres run-postgres verify mock start
