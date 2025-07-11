# Use official Golang image as a base image
FROM golang:1.24-alpine3.20 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Go Modules manifests
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Use smaller base image
FROM gcr.io/distroless/base-debian11

# Copy the pre-built binary from the builder
COPY --from=builder /app/main .

# Expose the port that the app will run on
EXPOSE 8080

# Command to run the executable
CMD ["/main"]
