
REPO	:= registry.digitalocean.com/zeus-fyi
NAME    := autok8s
GIT_SHA := $(shell git rev-parse HEAD)
IMG     := ${REPO}/${NAME}:${GIT_SHA}
LATEST  := ${REPO}/${NAME}:latest

app.indexer:
	go run ../../apps/autok8s

docker.build:
	@ docker build -t ${IMG} -t ${LATEST} -f ../../docker/autok8s/Dockerfile ../../
