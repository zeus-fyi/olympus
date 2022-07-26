
FOLDER_NAME=redis
VOLUME_NAME=${FOLDER_NAME}_redis
CONTAINER_NAME=redis
ROOT_PATH=../../

docker.compose.up:
	docker-compose -f ${ROOT_PATH}pkg/datastores/redis/docker-compose-redis.yml up

docker.compose.down:
	docker-compose -f ${ROOT_PATH}pkg/datastores/redis/docker-compose-redis.yml down

docker.volume.rm:
	docker volume rm ${VOLUME_NAME}

docker.volume.ls:
	docker volume ls