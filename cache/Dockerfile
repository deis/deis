FROM deis/base:latest
MAINTAINER Gabriel Monroy <gabriel@opdemand.com>

# install redis from OS package
RUN apt-get update && apt-get install -yq python-software-properties
RUN add-apt-repository ppa:chris-lea/redis-server -y
RUN apt-get update
RUN apt-get install -yq redis-server

# install latest etcdctl including no-sync options
RUN wget -q https://s3-us-west-2.amazonaws.com/deis/etcdctl.no-sync -O /usr/local/bin/etcdctl
RUN chmod +x /usr/local/bin/etcdctl

# add the current build context to /app
ADD . /app

# define the execution environment
WORKDIR /app
CMD ["/app/bin/boot"]
EXPOSE 6379
