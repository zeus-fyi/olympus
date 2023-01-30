REPO	:= registry.digitalocean.com/zeus-fyi
NAME    := hydra
GIT_SHA := $(shell git rev-parse HEAD)
IMG     := ${REPO}/${NAME}:${GIT_SHA}
LATEST  := ${REPO}/${NAME}:latest
SRC	    := ../../
GOMODCACHE := $(shell go env GOMODCACHE)
GOCACHE := $(shell go env GOCACHE)
GOOS 	:= linux
GOARCH  := amd64

docker.pubbuildx:
	@ docker buildx build -t ${IMG} -t ${LATEST} --build-arg GOMODCACHE=${GOMODCACHE} --build-arg GOCACHE=${GOCACHE} --platform=linux/amd64 -f ../../docker/hydra/Dockerfile ../../ --push

CHOREOGRAPHY_NAME    := hydrachoreography
CHOREOGRAPHY_IMG     := ${REPO}/${NAME}:${GIT_SHA}
CHOREOGRAPHY_LATEST  := ${REPO}/${NAME}:latest

docker.cho.pubbuildx:
	@ docker buildx build -t ${CHOREOGRAPHY_IMG} -t ${CHOREOGRAPHY_LATEST} --build-arg GOMODCACHE=${GOMODCACHE} --build-arg GOCACHE=${GOCACHE} --platform=linux/amd64 -f ../../docker/hydra_choreography/Dockerfile ../../ --push
