
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