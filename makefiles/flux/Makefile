
NAMESPACE=eth-indexer
PG_SECRET_NAME=postgres-auth
secret.pg:
	kubectl -n ${NAMESPACE} create secret generic ${PG_SECRET_NAME} \
    --from-literal=postgres-conn-str=change-me \
    --dry-run=client \
    -o yaml > ${PG_SECRET_NAME}.yaml

BEACON_SECRET_NAME=beacon-auth
secret.beacon:
	kubectl -n ${NAMESPACE} create secret generic ${BEACON_SECRET_NAME} \
    --from-literal=beacon-conn-str=change-me \
    --dry-run=client \
    -o yaml > ${BEACON_SECRET_NAME}.yaml

flux.age:
	cat age.agekey | \
    kubectl create secret generic sops-age \
    --namespace=flux-system \
    --from-file=../../configs/age.agekey=/dev/stdin

AGE_PUBKEY=age14h7yvfgkcdat8kdk2wl7gmxjw6a9j32yn2zhz2e4fgercv5gnf0s4svmu8
encrypt.pg:
	sops --age=${AGE_PUBKEY} \
    --encrypt --encrypted-regex '^(data|stringData)$' --in-place ${PG_SECRET_NAME}.yaml

encrypt.beacon:
	sops --age=${AGE_PUBKEY} \
    --encrypt --encrypted-regex '^(data|stringData)$' --in-place ${BEACON_SECRET_NAME}.yaml