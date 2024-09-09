# First stage: Build the Go binary
FROM golang:1.22.5-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o client main.go

# Second stage: Create minimal runtime image
FROM alpine:latest

# Install certificates if your app makes HTTPS requests
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/client /client

# Run the binary
CMD ["/client"]
