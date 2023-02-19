FROM golang:1.20-alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

RUN apk update \
    && apk add --no-cache ca-certificates \
    && update-ca-certificates

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o main .

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/main .

# Build a small image
FROM alpine:3.15.0

ENV PORT=8000 \
    HOST="" \
    COUNT=50 \
    VERSION="" \
    TOKEN=""
EXPOSE 8080

COPY --from=builder /dist/main /

# Command to run
CMD ["./main"]
