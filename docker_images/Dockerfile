FROM golang:1.22.4-alpine3.20

RUN apk add --no-progress --no-cache gcc musl-dev

WORKDIR /app

# Copy go.mod and go.sum to the WORKDIR
COPY go.mod go.sum ./

# Copies everything from your root directory into /app
COPY . .

# RUN go get -tags musl github.com/confluentinc/confluent-kafka-go/kafka

# Installs Go dependencies
RUN go mod download


# Builds your app with optional configuration
RUN go build -tags musl -ldflags '-extldflags "-static"' -o /main

FROM 533267397749.dkr.ecr.ap-south-1.amazonaws.com/legal-referral:latestqclr