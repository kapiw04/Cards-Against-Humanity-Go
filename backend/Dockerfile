# Build Stage
FROM golang:1.23.2-alpine AS builder

# Install gcc and libc-dev (for CGO) and sqlite-dev
RUN apk add gcc musl-dev sqlite-dev

WORKDIR /app

# Cache Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the database
COPY cah/cah_cards.db ./

# Copy the rest of the source code
COPY *.go .

# Build the Go binary with CGO enabled
RUN CGO_ENABLED=1 GOOS=linux go build -o ./docker-gs-ping
EXPOSE 8080

CMD ["./docker-gs-ping"]
