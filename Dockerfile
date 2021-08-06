FROM golang:alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

RUN apk update \
    && apk add --no-cache ca-certificates \
    && apk add --update gcc musl-dev \
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
RUN go build -tags=fts5 -o main .

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/main .

# Build a small image
FROM scratch

COPY --from=builder /dist/main /

ENV PORT=8000
ENV HOST=""
ENV COUNT=50
ENV VERSION=""
ENV TOKEN=""

# Command to run
CMD main --port ${PORT} --host ${HOST} --count ${HOST} --version ${VERSION} --token ${TOKEN}