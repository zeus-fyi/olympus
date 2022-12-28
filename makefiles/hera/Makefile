REPO	:= registry.digitalocean.com/zeus-fyi
NAME    := hera
GIT_SHA := $(shell git rev-parse HEAD)
IMG     := ${REPO}/${NAME}:${GIT_SHA}
LATEST  := ${REPO}/${NAME}:latest
SRC	    := ../../
GOMODCACHE := $(shell go env GOMODCACHE)
GOCACHE := $(shell go env GOCACHE)

docker.pubbuildx:
	@ docker buildx build -t ${IMG} -t ${LATEST} --build-arg GOMODCACHE=${GOMODCACHE} --build-arg GOCACHE=${GOCACHE} --platform=linux/amd64 -f ../../docker/hera/Dockerfile ../../ --push


CORE_DIR=../../pkg/codegen/cookbook

TYPE=func
DIR=${CORE_DIR}/${TYPE}
FN=${DIR}/${TYPE}_template_gen
NEW_FN=${DIR}/${TYPE}

fn.gen:
	TYPE=func
	tojen gen ${FN}.go ${NEW_FN}.go

struct.gen:
	TYPE=struct
	tojen gen ${FN}.go ${NEW_FN}.go

gen.test:
	TYPE=test
	tojen gen ${FN}.go ${NEW_FN}.go

SQL_HOST=localhost
SQL_USER=postgres
SQL_PW=postgres
SQL_PORT=5432
SQL_DB=postgres
TLS_CRT_CERT_PATH=./DigiCertGlobalRootCA.crt.pem

PKG_NAME=autogen_structs
PKG_DIR=../../datastores/postgres/apps/zeus/structs/autogen
# https://github.com/fraenky8/tables-to-go
# absolutely stupid that this currently only works in the jetbrains terminal wtf dude
pg.gen.structs:
	tables-to-go -h ${SQL_HOST} -u ${SQL_USER} -p ${SQL_PW} -pn ${PKG_NAME} -of ${PKG_DIR}
