# syntax=docker/dockerfile:1

FROM golang:1.21

RUN go install github.com/magefile/mage@latest

# Copy files
COPY ./cmd /app/cmd
COPY ./internal /app/internal
COPY ./pkg /app/pkg
COPY ./magefile.go /app/magefile.go
COPY ./go.mod /app/go.mod
COPY ./go.sum /app/go.sum

WORKDIR /app

# Download dependencies
RUN go mod download

# Run verification process
RUN mage verify

# Build
RUN mage build

CMD [ "./bin/app" ]