REPO	:= registry.digitalocean.com/zeus-fyi
NAME    := zeus-cloud
GIT_SHA := $(shell git rev-parse HEAD)
IMG     := ${REPO}/${NAME}:${GIT_SHA}
LATEST  := ${REPO}/${NAME}:latest
SRC	    := ../../

GOOS 	:= linux
GOARCH  := amd64

docker.pubbuildx:
	@ docker buildx build -t ${IMG} -t ${LATEST} --platform=${GOOS}/${GOARCH} -f ../../docker/zeus_cloud/Dockerfile ../../ --push

docker.debug:
	docker run -it --entrypoint /bin/bash registry.digitalocean.com/zeus-fyi/zeus-cloud:latest
