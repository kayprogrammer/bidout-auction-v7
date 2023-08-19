FROM golang:1.20-alpine

RUN mkdir build

# We create folder named build for our app.
WORKDIR /build

COPY ./.env .
COPY go.mod go.sum ./

# Download dependencies
RUN go install github.com/cosmtrek/air@latest
RUN go mod download