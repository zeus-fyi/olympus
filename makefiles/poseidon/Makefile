REPO	:= registry.digitalocean.com/zeus-fyi
NAME    := poseidon
GIT_SHA := $(shell git rev-parse HEAD)
IMG     := ${REPO}/${NAME}:${GIT_SHA}
LATEST  := ${REPO}/${NAME}:latest
SRC	    := ../../
GITHUB_TOKEN=${GITHUB_TOKEN_OLYMPUS}
GOMODCACHE := $(shell go env GOMODCACHE)
GOCACHE := $(shell go env GOCACHE)

docker.pubbuildx:
	@ docker buildx build -t ${IMG} -t ${LATEST} --build-arg GOMODCACHE=${GOMODCACHE} --build-arg GOCACHE=${GOCACHE} --platform=linux/amd64 -f ../../docker/poseidon/Dockerfile ../../ --push

