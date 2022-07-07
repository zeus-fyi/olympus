GIT_REF ?= $(shell git rev-parse HEAD)

docker.build:
	docker build --tag registry.digitalocean.com/zeus-fyi/zeus-indexer -f ../../docker/Dockerfile ../../

docker.push:
	docker push registry.digitalocean.com/zeus-fyi/zeus-indexer

docker.run:
	docker run -it -p 9000:9000 registry.digitalocean.com/zeus-fyi/zeus-indexer sh -c "eth_indexer"

docker.debug:
	docker run -it --entrypoint /bin/bash registry.digitalocean.com/zeus-fyi/zeus-indexer
