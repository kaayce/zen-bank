# ZenBank - Full featured bank application backend [Golang + Postgres + Kubernetes + gRPC]

## Overview
Welcome to ZenBank, a sophisticated banking application developed to showcases various technologies and practices in the Go ecosystem. The application includes a robust backend implemented in Go, a Postgres database, and modern deployment techniques using Kubernetes and AWS.

## Technologies Used
- **Golang**: The primary programming language for the backend.
- **PostgreSQL**: The relational database system.
- **Docker**: Containerization tool for the database and application.
- **Kubernetes**: For container orchestration and deployment on AWS.
- **gRPC**: For high-performance communication between services.
- **Gin**: Web framework for the RESTful API.
- **Viper**: Configuration management.
- **AWS**: Cloud platform for hosting the application.
- **GitHub Actions**: CI/CD pipeline for automated testing and deployment.
- **Redis**: Used for background tasks.
- **Cert-manager and Let's Encrypt**: For handling TLS certificates.

## Features
- **Database Schema Design**: Using dbdiagram.io to design and generate SQL code.
- **CRUD Operations**: Generate CRUD operations from SQL using sqlc.
- **Database Transactions**: Implementation of database transactions and handling deadlocks.
- **RESTful API**: Built using Gin to handle HTTP requests.
- **Authentication**: Secure user authentication using PASETO tokens.
- **Role-Based Access Control (RBAC)**: To manage user permissions.
- **gRPC API**: For communication between microservices.
- **Docker and Kubernetes**: Containerization and orchestration for scalable deployment.
- **Automated CI/CD**: Using GitHub Actions for continuous integration and deployment.
- **Async Tasks**: Integration with Redis for background processing.

## Getting Started
To run this application locally, follow the steps below:

### Prerequisites
- Docker
- Go (> v1.15)
- Sqlc
- Migrate

### Setup Development Environment
1. **Clone the Repository**
   ```sh
   git clone https://github.com/kaayce/zen-bank.git
   cd zenbank
   ```

2. **Environment Variables**
   Create a `.env` file in the project root with the following content:
   ```env
   POSTGRES_USER=your_postgres_user
   POSTGRES_PASSWORD=your_postgres_password
   POSTGRES_DB=zen_bank
   POSTGRES_CONTAINER_NAME=zenbank_postgres
   POSTGRES_IMAGE=postgres:latest
   DB_URL=postgres://your_postgres_user:your_postgres_password@localhost:5434/zen_bank?sslmode=disable
   SCHEMA_DIR=path/to/migrations
   ```

3. **Start PostgreSQL Container**
   ```sh
   make startdb
   ```

4. **Create Database**
   ```sh
   make createdb
   ```

5. **Run Database Migrations**
   ```sh
   make migrate-up
   ```

6. **Generate SQL Code**
   ```sh
   make sqlc
   ```

### Running the Application
1. **Build and Run the Application**
   ```sh
   make run-app
   ```

2. **Running Tests**
   ```sh
   make test
   ```

### Stopping and Removing Docker Container
To stop and remove the PostgreSQL container (only in development or local environment):
```sh
make reset
```

## Deployment
For deployment, ZenBank uses Kubernetes on AWS. Follow the instructions provided in the course to set up AWS EKS and deploy the application.

For any questions or issues, please feel free to contact me at [pat.nzediegwu@gmail.com].