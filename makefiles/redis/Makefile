
FOLDER_NAME=redis_docker
VOLUME_NAME=${FOLDER_NAME}_redis
CONTAINER_NAME=redis
ROOT_PATH=../../

redis.fresh:
	docker volume rm ${VOLUME_NAME}
	docker-compose -f ${ROOT_PATH}datastores/redis_docker/docker-compose-redis.yml up

redis.start:
	docker-compose -f ${ROOT_PATH}datastores/redis_docker/docker-compose-redis.yml up

redis.stop:
	docker-compose -f ${ROOT_PATH}datastores/redis_docker/docker-compose-redis.yml down

redis.volume.rm:
	docker volume rm ${VOLUME_NAME}

redis.volume.ls:
	docker volume ls