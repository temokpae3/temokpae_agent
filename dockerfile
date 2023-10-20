## First stage - Build stage
# Pull golang image from Dockerhub
FROM golang:alpine AS builder

# Set up the working directory
WORKDIR /temokpae_agent

# copy the source code and build
COPY go.mod .
COPY go.sum .
COPY main.go .

# Build a statically-linked Go binary for Linux
RUN CGO_ENABLED=0 GOOS=linux go build -a -o main .

## Second stage - Run stage
FROM alpine:latest

# Set up the working directory
WORKDIR /temokpae_agent

# Copy the executable binary file and env file from the last stage to the new stage
COPY --from=builder /temokpae_agent/main .

# Add environment variables
ENV LOGGLY_TOKEN="LOGGLY_TOKEN"

# Check results
RUN env

# Execute the build
CMD ["./main"]