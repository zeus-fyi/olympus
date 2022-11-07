
REPO	:= registry.digitalocean.com/zeus-fyi
NAME    := zeus
GIT_SHA := $(shell git rev-parse HEAD)
IMG     := ${REPO}/${NAME}:${GIT_SHA}
LATEST  := ${REPO}/${NAME}:latest

scale.down:
	kubectl scale deployment/${DEPLOYMENT_NAME} --replicas=0

app.zeus:
	go run ../../apps/zeus

docker.build:
	@ docker build -t ${IMG} -t ${LATEST} -f ../../docker/zeus/Dockerfile ../../

docker.pubbuildx:
	@ docker buildx build -t ${IMG} -t ${LATEST} --platform=linux/amd64 -f ../../docker/zeus/Dockerfile ../../ --push

docker.buildx.prune:
	@ docker buildx prune

docker.debug:
	docker run -it --entrypoint /bin/bash registry.digitalocean.com/zeus-fyi/zeus:latest
