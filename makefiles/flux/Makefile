
CONFIGS_PATH=../../configs
PROD_FLUX_PATH=../../devops/apps/production
NAMESPACE=eth-indexer

# service secrets
BEACON_SECRET_NAME=beacon-auth
PG_SECRET_NAME=postgres-auth
BEACON_SECRET=
PG_SECRET=
# git repo source, adding ssh key
REPO_SSH_URL=ssh://git@github.com/zeus-fyi/olympus
GIT_REPO_SSH_SECRET_NAME=eth-indexer-auth
BRANCH=main

secret.pg:
	kubectl -n ${NAMESPACE} create secret generic ${PG_SECRET_NAME} \
    --from-literal=postgres-conn-str=${PG_SECRET} \
    --dry-run=client -o yaml | tee ${PG_SECRET_NAME}.yaml

encrypt.pg:
	sops -e ${PG_SECRET_NAME}.yaml | tee ${PROD_FLUX_PATH}/${PG_SECRET_NAME}.yaml

flux.read.pg.secret:
	kubectl get secret ${PG_SECRET_NAME} -n ${NAMESPACE} -ojson \
        | jq -r '.data."postgres-conn-str"' | base64 -d

secret.beacon:
	kubectl -n ${NAMESPACE} create secret generic ${BEACON_SECRET_NAME} \
    --from-literal=beacon-conn-str=${BEACON_SECRET} \
    --dry-run=client -o yaml | tee ${BEACON_SECRET_NAME}.yaml

encrypt.beacon:
	sops -e ${BEACON_SECRET_NAME}.yaml | tee ${PROD_FLUX_PATH}/${BEACON_SECRET_NAME}.yaml

flux.read.beacon.secret:
	kubectl get secret ${BEACON_SECRET_NAME} -n ${NAMESPACE} -ojson \
        | jq -r '.data."beacon-conn-str"' | base64 -d

# also needs to be re-installed for new flux install and then put in github
flux.git.secret.create:
	flux create secret git ${GIT_REPO_SSH_SECRET_NAME} --url=${REPO_SSH_URL} --ssh-key-algorithm=ecdsa --ssh-ecdsa-curve=p521

flux.read.git.pubkey:
	kubectl get secret eth-indexer-auth -n flux-system -ojson \
        | jq -r '.data."identity.pub"' | base64 -d

flux.git.secret.deploy:
	flux create source git ${NAMESPACE} --url=${REPO_SSH_URL} --branch=${BRANCH} --secret-ref=${GIT_REPO_SSH_SECRET_NAME}

flux.git.repos:
	 kubectl get gitrepository

# needed if you re-install flux
flux.age:
	cat age.agekey | \
    kubectl create secret generic sops-age \
    --namespace=flux-system \
    --from-file=age.agekey=/dev/stdin

flux.age.namespace:
	cat age.agekey | \
    kubectl create secret generic sops-age \
    --namespace=${NAMESPACE} \
    --from-file=age.agekey=/dev/stdin

flux.watch:
	flux get kustomizations --watch

flux.check:
	flux check

flux.all.sources:
	flux get sources all -A

flux.not.ready:
	flux get all -A --status-selector ready=false

flux.all:
	flux get all

flux.rec.git:
	 flux reconcile source git ${NAMESPACE}
	 flux reconcile source git flux-system

flux.bootstrap.dev:
	flux bootstrap github \
        --owner=zeus-fyi \
        --repository=${GITHUB_REPO} \
        --branch=main \
        --personal \
        --path=devops/clusters/dev