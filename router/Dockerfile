FROM deis/base
MAINTAINER Gabriel Monroy <gabriel@opdemand.com>

# install nginx
RUN apt-get install -yq python-software-properties
RUN add-apt-repository ppa:chris-lea/redis-server -y
RUN add-apt-repository ppa:nginx/stable -y
RUN apt-get update
RUN apt-get install -yq nginx

# install latest etcdctl including no-sync options
RUN wget -q https://s3-us-west-2.amazonaws.com/deis/etcdctl.no-sync -O /usr/local/bin/etcdctl
RUN chmod +x /usr/local/bin/etcdctl

ADD . /app
WORKDIR /app
EXPOSE 80
CMD ["/app/bin/boot"]
