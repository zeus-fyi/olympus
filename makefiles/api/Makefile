
REPO	:= registry.digitalocean.com/zeus-fyi
NAME    := zeus-indexer
GIT_SHA := $(shell git rev-parse HEAD)
IMG     := ${REPO}/${NAME}:${GIT_SHA}
LATEST  := ${REPO}/${NAME}:latest

docker.buildx:
	@ docker buildx build -t ${IMG} -t ${LATEST} -f ../../docker/Dockerfile ../../ --push

docker.print:
	echo ${IMG} ${LATEST}

docker.build:
	@ docker build -t ${IMG} -t ${LATEST} -f ../../docker/Dockerfile ../../

docker.push:
	@ docker push ${IMG}
	@ docker push ${LATEST}

docker.publish:
	@ docker buildx build -t ${IMG} -t ${LATEST} --platform=linux/amd64 -f ../../docker/Dockerfile ../../

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

helm.add:
	helm repo add ../../devops/helm/eth-indexer/ chartmuseum

helm.push:
	helm cm-push ../../devops/helm/eth-indexer/ chartmuseum -f

helm.push.url:
	helm cm-push ../../devops/helm/eth-indexer/ http://localhost:8080

helm.list:
	helm list

helm.search:
	helm search repo

chart.pf:
	kubectl port-forward chartmuseum-7c75f78c87-dk959 8080:8080

helm.install:
	helm install eth-indexer ../../helm/eth-indexer/ --values ../../helm/eth-indexer/values.yaml