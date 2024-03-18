# Use an official Golang runtime as a base image
FROM golang:1.21.5 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download Go module dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o wasteGo .

# Start a new stage from scratch
FROM alpine:latest  

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the previous stage
COPY --from=builder /app/wasteGo .

# Expose the port the application listens on
EXPOSE 8000

# Command to run the executable
CMD ["./wasteGo"]