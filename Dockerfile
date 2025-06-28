# Stage 1: Build the Go application
FROM golang:1.24 AS builder

RUN apt-get update && apt-get install git curl bash procps qrencode -y

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

CMD ["tail", "-f", "/dev/null"]

# Build the Go application
# This will build the main.go file. If your main file has a different name, change it here.
# RUN CGO_ENABLED=0 GOOS=linux go build -o /main main.go

# # Stage 2: Create a minimal final image
# FROM alpine:latest

# WORKDIR /root/

# # Copy the pre-built binary from the builder stage
# COPY --from=builder /main .

# # Expose port 8080 to the outside world
# EXPOSE 8080

# # Command to run the executable
# CMD ["./main"]
