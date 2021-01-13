# Copyright (C) 2020 ConsenSys Software Inc
# Docker file used to build code.

FROM golang

# Add code to be run.
# Also grab dependancies from source directory, to improve build speed.
ADD pkg /go/src/github.com/ConsenSys/fc-retrieval-client/pkg
ADD cmd /go/src/github.com/ConsenSys/fc-retrieval-client/cmd
ADD internal /go/src/github.com/ConsenSys/fc-retrieval-client/internal

# Add the runtime dockerfile into the context as Dockerfile
COPY Dockerfile.run /go/bin/Dockerfile

# Add the settings file needed at runtime to the bin directory so it
# can be accessed by the runtime Dockerfile.
COPY settings.json /go/bin/settings.json

# Remove any cached dependancies. TODO is this really needed?
RUN go clean -modcache
# Get all dependancies, but don't install.
RUN go get -d -v github.com/ConsenSys/fc-retrieval-client/cmd/client
# Do a full compile of app and dependancies, forcing static linking.
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o /go/bin/client github.com/ConsenSys/fc-retrieval-client/cmd/client
# Don't do install, as build now done. 
#RUN go install github.com/ConsenSys/fc-retrieval-client/cmd/client


# Set the workdir to be /go/bin which is where the binaries are built
WORKDIR /go/bin
# Export the WORKDIR as a tar stream
CMD tar -cf - .

