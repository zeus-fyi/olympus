
FOLDER_NAME=redis_docker
VOLUME_NAME=${FOLDER_NAME}_redis
CONTAINER_NAME=redis
ROOT_PATH=../../

redis.fresh:
	docker volume rm ${VOLUME_NAME}
	docker-compose -f ${ROOT_PATH}pkg/datastores/redis_docker/docker-compose-redis.yml up

redis.start:
	docker-compose -f ${ROOT_PATH}pkg/datastores/redis_docker/docker-compose-redis.yml up

redis.strop:
	docker-compose -f ${ROOT_PATH}pkg/datastores/redis_docker/docker-compose-redis.yml down

docker.volume.rm:
	docker volume rm ${VOLUME_NAME}

docker.volume.ls:
	docker volume ls