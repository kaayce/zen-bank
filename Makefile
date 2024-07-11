# Source Environment variables
include .env
export $(shell sed 's/=.*//' .env)

# -- DB Targets --

# Start the PostgreSQL container
postgres:
	docker run --name $(POSTGRES_CONTAINER_NAME) -p 5434:5432 \
	-e POSTGRES_USER=$(POSTGRES_USER) \
	-e POSTGRES_PASSWORD=$(POSTGRES_PASSWORD) \
	-d $(POSTGRES_IMAGE)
	@echo "Waiting for PostgreSQL to start..."
	@until docker exec -it $(POSTGRES_CONTAINER_NAME) pg_isready --username=$(POSTGRES_USER); do \
		sleep 1; \
		echo "Waiting for PostgreSQL to be ready..."; \
	done

# Create the zen_bank database
createdb:
	docker exec -it $(POSTGRES_CONTAINER_NAME) \
	createdb --username=$(POSTGRES_USER) \
	--owner=$(POSTGRES_USER) $(POSTGRES_DB)

#  verify db actual exists
verify:
	docker exec -it $(POSTGRES_CONTAINER_NAME) psql -U $(POSTGRES_USER) -c '\l'

migrate-up:
	migrate -path $(SCHEMA_DIR) -database $(DB_URL) -verbose up

migrate-down:
	migrate -path $(SCHEMA_DIR) -database $(DB_URL) -verbose down

startdb: postgres createdb

dropdb:
	docker exec -it $(POSTGRES_CONTAINER_NAME) dropdb $(POSTGRES_DB)

# sqlc
sqlc:
	sqlc generate

test:
	go test ./...

test-cov:
	go test -v -cover -short ./...

# performance x mem allocation
test-cov-mem:
	go test -v -cover ./... -gcflags '-m -l'

# performance
test-bench:
	go test -bench=. -benchmem -benchtime=10s

# go mod
mod:
	go mod tidy && go mod vendor

# build x perf
build-mem:
	go build -gcflags '-m -l'

# run app
run-app:
	go build && chmod +x zen-bank && ./zen-bank

# do not run in prod, dev/local only - stops and removes container
reset:
	@if [ "$(ENV)" = "local" ] || [ "$(ENV)" = "dev" ]; then \
		echo "Stopping and removing container $(POSTGRES_CONTAINER_NAME)..."; \
		docker stop $(POSTGRES_CONTAINER_NAME) && docker rm $(POSTGRES_CONTAINER_NAME); \
	else \
		echo "Not allowed in production environment"; \
	fi

.PHONY: startdb dropdb migrate-up migrate-down sqlc test
