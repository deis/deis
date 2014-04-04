FROM deis/base:latest
MAINTAINER Gabriel Monroy <gabriel@opdemand.com>

# install 9.3 from postgresql.org repository
RUN apt-get -yq install wget ca-certificates
RUN echo "deb http://apt.postgresql.org/pub/repos/apt/ precise-pgdg main" > /etc/apt/sources.list.d/pgdg.list
RUN wget --quiet --no-check-certificate -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
RUN apt-get update
RUN apt-get install -yq postgresql-9.3 && /etc/init.d/postgresql stop

# install latest etcdctl including no-peers options
RUN wget -q https://s3-us-west-2.amazonaws.com/deis/etcdctl.no-sync -O /usr/local/bin/etcdctl
RUN chmod +x /usr/local/bin/etcdctl

# install latest etcdctl including no-sync options
RUN wget -q https://s3-us-west-2.amazonaws.com/deis/etcdctl.no-sync -O /usr/local/bin/etcdctl
RUN chmod +x /usr/local/bin/etcdctl

# debug to remove
RUN apt-get install -yq curl

# add the current build context to /app
ADD . /app

# define the execution environment
VOLUME ["/var/lib/postgresql"]
WORKDIR /app
CMD ["/app/bin/boot"]
EXPOSE 5432
