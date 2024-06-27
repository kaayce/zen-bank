# Source Environment variables
include .env
export $(shell sed 's/=.*//' .env)

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

# do not run in prod, dev/local only - stops and removes container
reset:
	docker stop postgres12_zenbank && docker rm postgres12_zenbank

.PHONY: startdb