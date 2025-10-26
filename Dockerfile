FROM golang:alpine AS builder

# Install git to fetch dependencies and gcc for cgo
RUN apk add --no-cache git gcc musl-dev

# Add GitHub token secret for private repo access
RUN --mount=type=secret,id=GITHUB_TOKEN \
    GITHUB_TOKEN=$(cat /run/secrets/GITHUB_TOKEN) && \
    git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/".insteadOf "https://github.com/"

# Set the Current Working Directory inside the container and add source files
WORKDIR /app
COPY . /app

# Make sure the scripts and .env file have LF line endings
RUN apk add --no-cache dos2unix
RUN if [ -f /app/.env ]; then dos2unix /app/.env; fi
RUN dos2unix /app/compile.sh /app/serve.sh /app/generate_docs.sh

# Make sure the compile and serve scripts are executables
RUN chmod +x /app/compile.sh /app/serve.sh /app/generate_docs.sh

# Install Go Swagger for API documentation generation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Go modules download
RUN go mod download

# Generate API documentation
ENV PATH="/root/go/bin:${PATH}"
RUN /app/generate_docs.sh

# Compile the Go app
RUN /app/compile.sh

FROM alpine AS server

# Set the Current Working Directory inside the container and copy the compiled Go app from the builder stage
WORKDIR /app
COPY --from=builder /app .

# Expose port to the outside world
EXPOSE 8080

# Run the Go app when the container launches
ENTRYPOINT ["/app/serve.sh"]