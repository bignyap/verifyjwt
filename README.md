# JWT Token Verifier

A simple Go application for verifying JWT tokens.

## Prerequisites

- [Go](https://golang.org/doc/install) (version 1.22.4+ recommended)
- [Docker](https://docs.docker.com/get-docker/) (optional, for Docker-based deployment)

## Getting Started

### 1. Clone the Repository

```sh
git clone https://github.com/bignyap/verifyjwt.git
cd verifyjwt
```

### 2. Download Dependencies

Ensure you have all required Go modules by running:

```sh
go mod download
```

### 3. Run the Application

#### Method 1: Direct Execution

```sh
go run .
```

#### Method 2: Build and Run Executable

Build the executable:

```sh
go build -o verifyjwt
```

Run the executable:

```sh
./verifyjwt
```

#### Method 3: Build and Run with Docker

Build the Docker image:

```sh
docker build -t verifyjwt .
```

Run the Docker container:

```sh
docker run -it --rm verifyjwt
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any changes.
