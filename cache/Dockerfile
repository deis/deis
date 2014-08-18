FROM deis/base:latest
MAINTAINER OpDemand <info@opdemand.com>

# install redis from OS package
RUN apt-get update && apt-get install -yq python-software-properties
RUN add-apt-repository ppa:chris-lea/redis-server -y
RUN apt-get update && apt-get install -yq redis-server

WORKDIR /app
CMD ["/app/bin/boot"]
EXPOSE 6379
ADD . /app
ADD redis.conf /etc/redis/redis.conf
