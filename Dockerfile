FROM golang:1.22.4-alpine

WORKDIR /app

COPY . /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-gs-ping

EXPOSE 8080

CMD ["go", "run", "."]
# CMD ["/docker-gs-ping"]