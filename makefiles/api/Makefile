GIT_REF ?= $(shell git rev-parse HEAD)

docker.build:
	docker build --tag registry.digitalocean.com/zeus-fyi/zeus-indexer -f ../../docker/Dockerfile ../../

docker.push:
	docker push registry.digitalocean.com/zeus-fyi/zeus-indexer

docker.run:
	docker run -it -p 9000:9000 registry.digitalocean.com/zeus-fyi/zeus-indexer sh -c "olympus --postgres-conn-str='' --beacon-endpoint=''"

docker.debug:
	docker run -it --entrypoint /bin/bash registry.digitalocean.com/zeus-fyi/zeus-indexer

k8s.deploy:
	kubectl apply -f ../../configs/kubernetes/indexer_pod.yaml