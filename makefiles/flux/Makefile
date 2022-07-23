
PROD_FLUX_PATH=../../devops/apps/production
NAMESPACE=eth-indexer
BEACON_SECRET_NAME=beacon-auth
PG_SECRET_NAME=postgres-auth
BEACON_SECRET=
PG_SECRET=

secret.pg:
	kubectl -n ${NAMESPACE} create secret generic ${PG_SECRET_NAME} \
    --from-literal=postgres-conn-str=${PG_SECRET} \
    --dry-run=client -o yaml | tee ${PG_SECRET_NAME}.yaml

encrypt.pg:
	sops -e ${PG_SECRET_NAME}.yaml | tee ${PROD_FLUX_PATH}/${PG_SECRET_NAME}.yaml

secret.beacon:
	kubectl -n ${NAMESPACE} create secret generic ${BEACON_SECRET_NAME} \
    --from-literal=beacon-conn-str=${BEACON_SECRET} \
    --dry-run=client -o yaml | tee ${BEACON_SECRET_NAME}.yaml

encrypt.beacon:
	sops -e ${BEACON_SECRET_NAME}.yaml | tee ${PROD_FLUX_PATH}/${BEACON_SECRET_NAME}.yaml

flux.age:
	cat age.agekey | \
    kubectl create secret generic sops-age \
    --namespace=flux-system \
    --from-file=../../configs/age.agekey=/dev/stdin
