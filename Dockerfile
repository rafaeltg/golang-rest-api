FROM golang:latest

# Copy the local package files to the container’s workspace.
ADD . /go/src/github.com/golang-rest-api

# Install app's dependencies
RUN go get github.com/golang-rest-api

# Install app binary globally within container
RUN go install github.com/golang-rest-api

ENTRYPOINT /go/bin/golang-rest-api

# Expose default port (8000)
EXPOSE 8000