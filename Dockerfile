# Stage 1: Build the Go application
FROM golang:latest AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o app

# Stage 2: Pack LibreOffice with the Go application
FROM ubuntu:latest

# Install dependencies for LibreOffice
RUN apt-get update && apt-get install -y libreoffice

WORKDIR /app

# Copy the Go binary from the builder stage
COPY --from=builder /app/app /app/app

# Set the command to start the application
CMD ["./app"]
