REPO	:= registry.digitalocean.com/zeus-fyi
NAME    := tokenizer
GIT_SHA := $(shell git rev-parse HEAD)
IMG     := ${REPO}/${NAME}:${GIT_SHA}
LATEST  := ${REPO}/${NAME}:latest
SRC	    := ../../

GOOS 	:= linux
GOARCH  := amd64

docker.pubbuildx:
	@ docker buildx build -t ${IMG} -t ${LATEST} --platform=${GOOS}/${GOARCH} -f ../../docker/tokenizer/Dockerfile ../../ --push

docker.debug:
	docker run -p 3000:3000 registry.digitalocean.com/zeus-fyi/zeus-cloud:latest
