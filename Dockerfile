# Use the official Go image as the builder
FROM golang:1.20 AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app
RUN go build -o myapp .

# Use a smaller image for the final build
FROM gcr.io/distroless/base

# Copy the binary from the builder stage
COPY --from=builder /app/myapp /myapp

# Command to run the binary
CMD ["/myapp"]
