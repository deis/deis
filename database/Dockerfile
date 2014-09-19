FROM ubuntu:14.04

ENV DEBIAN_FRONTEND noninteractive

# install common packages
RUN apt-get update && apt-get install -y curl

# install etcdctl
RUN curl -sSL -o /usr/local/bin/etcdctl https://s3-us-west-2.amazonaws.com/opdemand/etcdctl-v0.4.6 \
    && chmod +x /usr/local/bin/etcdctl

# install confd
RUN curl -sSL -o /usr/local/bin/confd https://s3-us-west-2.amazonaws.com/opdemand/confd-v0.5.0-json \
    && chmod +x /usr/local/bin/confd

# install 9.3 from postgresql.org repository
RUN echo "deb http://apt.postgresql.org/pub/repos/apt/ trusty-pgdg main" > /etc/apt/sources.list.d/pgdg.list
RUN curl -sk https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
RUN apt-get update && apt-get install -yq postgresql-9.3 && /etc/init.d/postgresql stop

# define the execution environment
VOLUME ["/var/lib/postgresql"]
WORKDIR /app
CMD ["/app/bin/boot"]
EXPOSE 5432
ADD . /app
