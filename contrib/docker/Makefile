all: build run

build:
	docker build -t deis/deis:latest .

run:
	docker run --privileged deis/deis
	#docker run --privileged -v /var/run/docker.sock:/var/run/docker.sock deis/deis

stop:
	-docker stop $(DISCOVERY_NAME)
	-docker stop $(LOGGER_NAME)
	-docker stop $(DATABASE_NAME)
	-docker stop $(CACHE_NAME)
	-docker stop $(CONTROLLER_NAME)

start:
	docker start $(DISCOVERY_NAME)
	docker start $(LOGGER_NAME)
	docker start $(DATABASE_NAME)
	docker start $(CACHE_NAME)
	docker start $(CONTROLLER_NAME)

restart:
	docker restart $(DISCOVERY_NAME)
	docker restart $(LOGGER_NAME)
	docker restart $(DATABASE_NAME)
	docker restart $(CACHE_NAME)
	docker restart $(CONTROLLER_NAME)

clean: stop
	-docker rm $(DISCOVERY_NAME)
	-docker rm $(LOGGER_NAME)
	-docker rm $(DATABASE_NAME)
	-docker rm $(CACHE_NAME)
	-docker rm $(CONTROLLER_NAME)

full-clean: clean
	docker rmi $(DISCOVERY_IMAGE)
	docker rmi $(LOGGER_IMAGE)
	docker rmi $(DATABASE_IMAGE)
	docker rmi $(CACHE_IMAGE)
	docker rmi $(CONTROLLER_IMAGE)
	docker rmi $(DEIS_IMAGE)
