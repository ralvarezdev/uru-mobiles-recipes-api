FROM golang:alpine AS builder

# Install git to fetch dependencies and protoc for compiling .proto files
RUN apk add --no-cache git

# Install Go for building the application
RUN apk add --no-cache go

# Add port argument
ARG PORT=8080
ENV PORT=${PORT}

# Add GitHub token secret for private repo access
RUN --mount=type=secret,id=GITHUB_TOKEN \
    GITHUB_TOKEN=$(cat /run/secrets/GITHUB_TOKEN) && \
    git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/".insteadOf "https://github.com/" \

# Set the Current Working Directory inside the container and add source files
WORKDIR /app
COPY . /app

# Pull rebase to get on the latest commit
RUN git pull --rebase

# Make sure the scripts and .env file have LF line endings
RUN apk add --no-cache dos2unix
RUN if [ -f /app/.env ]; then dos2unix /app/.env; fi
RUN dos2unix /app/compile.sh
RUN dos2unix /app/serve.sh
RUN dos2unix /app/generate_docs.sh

# Make sure the compile and serve scripts are executables
RUN chmod +x /app/compile.sh
RUN chmod +x /app/serve.sh
RUN chmod +x /app/generate_docs.sh

# Install Go Swagger for API documentation generation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Generate API documentation
RUN /app/generate_docs.sh

# Compile the Go app
RUN /app/compile.sh

FROM alpine AS server

# Set the Current Working Directory inside the container and copy the compiled Go app from the builder stage
WORKDIR /app
COPY --from=builder /app .

# Expose port to the outside world
EXPOSE ${PORT}

# Run the Go app when the container launches
ENTRYPOINT ["/app/serve.sh"]