#Build state
FROM golang:1.22.5-alpine3.20 AS builder

RUN apk add --no-progress --no-cache gcc musl-dev

WORKDIR /app

# Copies everything from your root directory into /app
COPY . .

# Builds your app with optional configuration
RUN go build -tags musl -ldflags '-extldflags "-static"' -o main main.go

# Run state
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY service-account-key.json .

# Tells Docker which network port your container listens on
EXPOSE 8080

# Specifies the executable command that runs when the container starts
CMD ["/app/main"]