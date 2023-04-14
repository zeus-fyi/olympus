REPO	:= registry.digitalocean.com/zeus-fyi
NAME    := hestia
GIT_SHA := $(shell git rev-parse HEAD)
IMG     := ${REPO}/${NAME}:${GIT_SHA}
LATEST  := ${REPO}/${NAME}:latest
SRC	    := ../../
GOMODCACHE := $(shell go env GOMODCACHE)
GOCACHE := $(shell go env GOCACHE)
GOOS 	:= linux
GOARCH  := amd64

docker.pubbuildx:
	@ docker buildx build \
				-t ${IMG} \
				-t ${LATEST} \
				--build-arg GOMODCACHE=${GOMODCACHE} \
				--build-arg GOCACHE=${GOCACHE} \
				--build-arg CA_CRT=./configs/ca-certificate.crt \
				--build-arg CA_PEM=./configs/zeus.fyi/ca.pem \
				--build-arg CA_KEY=./configs/zeus.fyi/ca.key \
				--platform=${GOOS}/${GOARCH} \
				-f ../../docker/hestia/Dockerfile \
				../../ --push
docker.debug:
	docker run -it --entrypoint /bin/bash registry.digitalocean.com/zeus-fyi/hestia:latest

docker.pull:
	docker pull registry.digitalocean.com/zeus-fyi/hestia:latest
