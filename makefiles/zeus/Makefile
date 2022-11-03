
REPO	:= registry.digitalocean.com/zeus-fyi
NAME    := zeus
GIT_SHA := $(shell git rev-parse HEAD)
IMG     := ${REPO}/${NAME}:${GIT_SHA}
LATEST  := ${REPO}/${NAME}:latest

app.indexer:
	go run ../../apps/zeus

docker.build:
	@ docker build -t ${IMG} -t ${LATEST} -f ../../docker/zeus/Dockerfile ../../

docker.pubbuildx:
	@ docker buildx build -t ${IMG} -t ${LATEST} --platform=linux/amd64 -f ../../docker/zeus/Dockerfile ../../ --push
