FROM deis/init:latest
MAINTAINER OpDemand <info@opdemand.com>

# install redis from OS package
RUN apt-get update && apt-get install -yq python-software-properties
RUN add-apt-repository ppa:chris-lea/redis-server -y
RUN apt-get update && apt-get install -yq redis-server

WORKDIR /app
CMD ["/sbin/my_init"]
EXPOSE 6379
ADD . /app
ADD bin/boot /etc/service/deis-cache/run
