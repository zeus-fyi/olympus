
FOLDER_NAME=postgres_docker
VOLUME_NAME=${FOLDER_NAME}_postgres
CONTAINER_NAME=postgres
ROOT_PATH=../../

pg.fresh:
	docker volume rm ${VOLUME_NAME}
	docker-compose -f ${ROOT_PATH}pkg/datastores/postgres_docker/docker-compose-postgres.yml up

pg.start:
	docker-compose -f ${ROOT_PATH}pkg/datastores/postgres_docker/docker-compose-postgres.yml up

pg.stop:
	docker-compose -f ${ROOT_PATH}pkg/datastores/postgres_docker/docker-compose-postgres.yml down

docker.volume.rm:
	docker volume rm ${VOLUME_NAME}

docker.volume.ls:
	docker volume ls

#SQL_HOST=localhost
#SQL_USER=postgres
#SQL_PW=postgres
#SQL_PORT=5432

SQL_HOST=zeus-pg-dev.postgres.database.azure.com
SQL_USER=postgres
SQL_PW=
SQL_PORT=5432
TLS_CRT_CERT_PATH=./DigiCertGlobalRootCA.crt.pem

temporal.create:
	./temporal-sql-tool --ep ${SQL_HOST} -u ${SQL_USER} -pw ${SQL_PW} -p ${SQL_PORT} --tls --tls-cert-file ${TLS_CRT_CERT_PATH} --db temporal --plugin postgres create

temporal.vis.create:
	./temporal-sql-tool --ep ${SQL_HOST} -u ${SQL_USER} -pw ${SQL_PW} -p ${SQL_PORT} --tls --tls-cert-file ${TLS_CRT_CERT_PATH} --db temporal_visibility --plugin postgres create

temporal.schema.create:
	./temporal-sql-tool --ep ${SQL_HOST} -u ${SQL_USER} -pw ${SQL_PW} -p ${SQL_PORT} --tls --tls-cert-file ${TLS_CRT_CERT_PATH} --plugin postgres --db temporal setup-schema -v 0.0

temporal.schema.update:
	./temporal-sql-tool --ep ${SQL_HOST} -u ${SQL_USER} -pw ${SQL_PW} -p ${SQL_PORT} --tls --tls-cert-file ${TLS_CRT_CERT_PATH} --plugin postgres --db temporal update-schema -d ../../../temporal/schema/postgresql/v96/temporal/versioned

temporal.vis.schema.create:
	./temporal-sql-tool --ep ${SQL_HOST} -u ${SQL_USER} -pw ${SQL_PW} -p ${SQL_PORT} --tls --tls-cert-file ${TLS_CRT_CERT_PATH} --plugin postgres --db temporal_visibility setup-schema -v 0.0

temporal.vis.schema.update:
	./temporal-sql-tool --ep ${SQL_HOST} -u ${SQL_USER} -pw ${SQL_PW} -p ${SQL_PORT} --tls --tls-cert-file ${TLS_CRT_CERT_PATH} --plugin postgres --db temporal_visibility update-schema -d ../../../temporal/schema/postgresql/v96/visibility/versioned

temporal.start:
	cd ../../../docker-compose && docker-compose -f docker-compose-postgres.yml up

temporal.stop:
	cd ../../../docker-compose && docker-compose -f docker-compose-postgres.yml down

temporal.stop.rm:
	cd ../../../docker-compose && docker-compose -f docker-compose-postgres.yml down --volumes