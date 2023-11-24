REPO	:= registry.digitalocean.com/zeus-fyi
NAME    := aegis-secrets
GIT_SHA := $(shell git rev-parse HEAD)
IMG     := ${REPO}/${NAME}:${GIT_SHA}
LATEST  := ${REPO}/${NAME}:latest
SRC	    := ../../
GOMODCACHE := $(shell go env GOMODCACHE)
GOCACHE := $(shell go env GOCACHE)

#docker.pubbuildx:
#	@ docker buildx build -t ${IMG} -t ${LATEST} --build-arg GOMODCACHE=${GOMODCACHE} --build-arg GOCACHE=${GOCACHE} --platform=linux/amd64 -f ../../docker/aegis/Dockerfile ../../ --push

docker.pubbuildx:
	@ docker buildx build -t ${IMG} -t ${LATEST} --platform=linux/amd64 -f ../../docker/aegis_secrets/Dockerfile ../../ --push


WEB3SIGNER_NAME 	:= web3signer
WEB3SIGNER_IMG     	:= ${REPO}/${WEB3SIGNER_NAME}:${GIT_SHA}
WEB3SIGNER_LATEST  	:= ${REPO}/${WEB3SIGNER_NAME}:latest
VERSION 			:= 22.11.0

w3.docker.pubbuildx:
	@ docker buildx build -t ${WEB3SIGNER_IMG} -t ${WEB3SIGNER_LATEST} --build-arg VERSION=${VERSION} --build-arg GOMODCACHE=${GOMODCACHE} --build-arg GOCACHE=${GOCACHE} --platform=linux/amd64 -f ../../docker/web3signer/Dockerfile ../../ --push

docker.debug:
	docker run -it --entrypoint /bin/bash registry.digitalocean.com/zeus-fyi/aegis-secrets:latest