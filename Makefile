# Copyright (C) 2020 ConsenSys Software Inc.

# Build the Filecoin Retrieval Gateway

# Usage:
#   [VERSION=v3] [REGISTRY="gcr.io/google_containers"] make build
VERSION?=dev
REGISTRY?=

# This target (the first target in the build file) is the one that is executed if no 
# command line args are specified.
release1: clean useremoteconfig utest build

# Use this target if you are using local packages, or if the build is via circle ci, 
# and the go.mod and go.sum file should not be updated
release: clean utest build

# builds a docker image that builds the app and packages it into a minimal docker image
build:
	docker build -t ${REGISTRY}fc-retrieval-client:${VERSION} .


# push the image to an registry
push:
#	gcloud docker -- push ${REGISTRY}/fc-retrieval-client:${VERSION}

detectlocal:
	cd scripts; bash detect-local-gateway-repo.sh

detectmisconfig:
	cd scripts; bash detect-pkg-misconfig.sh

useremoteconfig:
	cd scripts; bash use-remote-repos.sh

uselocal:
	cd scripts; bash use-local-repos.sh

utest:
	go test ./...

itest:
	docker images
	docker-compose down
	docker-compose up --abort-on-container-exit --exit-code-from client

stop:
	docker-compose down

# remove previous images and containers
clean:
#	rm -f /etc/client/
	docker rm -f ${REGISTRY}fc-retrieval-client-builder 2> /dev/null || true
	docker rmi -f ${REGISTRY}fc-retrieval-client-builder || true
	docker rmi -f "${REGISTRY}fc-retrieval-client:${VERSION}" || true

# Alays assume these targets are out of date.
.PHONY: clean itest utest build release push detectmisconfig detectlocal stop

