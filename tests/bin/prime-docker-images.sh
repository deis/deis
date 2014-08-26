#!/bin/sh
#
# WARNING: Don't run this script unless you understand that it will remove all Docker items.
#
# Purges *all* Docker containers and images from the local graph, then pulls down a set of
# images to help tests run faster.

# Remove all Dockernalia
docker kill `docker ps -q`
docker rm `docker ps -a -q`
docker rmi -f `docker images -q`

# Pull Deis testing essentials
docker pull deis/base:latest
docker pull deis/slugbuilder:latest
docker pull deis/slugrunner:latest
docker pull deis/test-etcd:latest
docker pull paintedfox/postgresql:latest
