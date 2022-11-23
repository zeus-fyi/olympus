REPO	:= registry.digitalocean.com/zeus-fyi
NAME    := artemis
GIT_SHA := $(shell git rev-parse HEAD)
IMG     := ${REPO}/${NAME}:${GIT_SHA}
LATEST  := ${REPO}/${NAME}:latest
SRC	    := ../../

docker.pubbuildx:
	@ docker buildx build -t ${IMG} -t ${LATEST} --platform=linux/amd64 -f ../../docker/artemis/Dockerfile ../../ --push

