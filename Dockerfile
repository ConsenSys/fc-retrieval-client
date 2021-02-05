# Copyright (C) 2020 ConsenSys Software Inc
FROM golang:1.15-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers git


WORKDIR /go/src/github.com/ConsenSys/fc-retrieval-client
COPY . .
# Remove any cached dependancies. TODO is this really needed?
RUN go clean -modcache
# Get all dependancies, but don't install.
RUN go get -d -v github.com/ConsenSys/fc-retrieval-client/cmd/client
# Do a full compile of app and dependancies, forcing static linking.
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o /go/bin/client github.com/ConsenSys/fc-retrieval-client/cmd/client

# Pull build gateway into a second stage deploy alpine container
FROM alpine:latest
COPY --from=builder /go/bin/client /client

# Run the binary when the container starts.
WORKDIR /
CMD ["/client"]
