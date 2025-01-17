# Build Stage
FROM golang:1.22-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz | tar xvz

# Run Stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
# Conditionally copy .env file if it exists
RUN if [ -f ".env" ]; then cp .env .; fi
COPY start.sh .
COPY --from=builder /app/migrate ./migrate
COPY db/migration ./db/migration

# Ensure scripts are executable
RUN chmod +x /app/start.sh

EXPOSE 8080 9090
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
