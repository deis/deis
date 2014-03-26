#
# Deis Makefile
#
DEIS_IMAGE=deis/deis:latest
DEIS_NAME=deis
DISCOVERY_PORT=4001
DISCOVERY_PEER_PORT=7001
DISCOVERY_IMAGE=deis/discovery:latest
DISCOVERY_NAME=deis-discovery
LOGGER_PORT=514
LOGGER_IMAGE=deis/logger:latest
LOGGER_NAME=deis-logger
DATABASE_PORT=5432
DATABASE_IMAGE=deis/database:latest
DATABASE_NAME=deis-database
CACHE_PORT=6379
CACHE_IMAGE=deis/cache:latest
CACHE_NAME=deis-cache
CONTROLLER_PORT=8000
CONTROLLER_IMAGE=deis/controller:latest
CONTROLLER_NAME=deis-controller
# uses docker0 bridge for inter-container communication
# this assumes the address is static
ETCD=172.17.42.1:4001

all: build run

build:
	docker build -t deis/deis:latest .

run:
	docker run --privileged -v /var/run/docker.sock:/var/run/docker.sock deis/deis

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
