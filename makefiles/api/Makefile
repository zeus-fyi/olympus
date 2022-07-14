GIT_REF ?= $(shell git rev-parse HEAD)

VERSION=latest-2

docker.build:
	docker buildx build --platform=linux/amd64 --tag registry.digitalocean.com/zeus-fyi/zeus-indexer:${VERSION} -f ../../docker/Dockerfile ../../

docker.push:
	docker push registry.digitalocean.com/zeus-fyi/zeus-indexer:${VERSION}

docker.run:
	docker run -it -p 9000:9000 registry.digitalocean.com/zeus-fyi/zeus-indexer sh -c "olympus --postgres-conn-str='' --beacon-endpoint=''"

docker.debug:
	docker run -it --entrypoint /bin/bash registry.digitalocean.com/zeus-fyi/zeus-indexer

k8s.upgrade:
	kubectl set image deployment.v1.apps/indexer-deployment zeus-indexer=registry.digitalocean.com/zeus-fyi/zeus-indexer:${VERSION}

k8s.deploy:
	kubectl apply -f ../../configs/kubernetes/indexer-deployment.yaml

k8s.pod:
	kubectl apply -f ../../configs/kubernetes/indexer-pod.yaml

k8s.delete:
	kubectl delete deployment indexer-deployment
