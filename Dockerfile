
# Specify the base image
FROM golang:alpine AS build

# Set the working directory
WORKDIR /snippetbox

COPY go1.23.4/bin/go /usr/local/bin/go

RUN echo "/usr/local/bin" >> /etc/shells

ENV PATH="/usr/local/bin:$PATH"

# Copy the Go source code
COPY . .

# Install dependencies
RUN \
  go mod download
  # go mod vendor

# Build the Go binary (adjust command for your project)
RUN go build ./cmd/web

# Create a new image for running the application
FROM alpine AS runtime

# Copy the binary
#COPY --from=build ./cmd/web ./cmd/web

# Expose the port (adjust port number if needed)
EXPOSE 8080 80

# Set the entrypoint to run the binary
CMD ["go run", "./cmd/web"]
