# Build stage
FROM golang:1.25.0-alpine3.22 AS builder

# Install necessary tools
RUN apk add --no-progress --no-cache gcc musl-dev

WORKDIR /app

# Copy everything from the root directory into /app
COPY . .

# Builds your app with optional configuration
RUN go build -o main main.go

# Run state
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY service-account-key.json .

# Expose port 8080 for App Runner
EXPOSE 8080

# Specifies the executable command that runs when the container starts
CMD ["/app/main"]