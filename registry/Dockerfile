FROM alpine:3.1

# install common packages
RUN apk add --update-cache curl bash sudo && rm -rf /var/cache/apk/*

# install etcdctl
RUN curl -sSL -o /usr/local/bin/etcdctl https://s3-us-west-2.amazonaws.com/get-deis/etcdctl-v0.4.9 \
    && chmod +x /usr/local/bin/etcdctl

# install confd
RUN curl -sSL -o /usr/local/bin/confd https://github.com/kelseyhightower/confd/releases/download/v0.9.0/confd-0.9.0-linux-amd64 \
    && chmod +x /usr/local/bin/confd

ENV DOCKER_REGISTRY_CONFIG /docker-registry/config/config.yml
ENV SETTINGS_FLAVOR deis

# define the execution environment
WORKDIR /app
CMD ["/app/bin/boot"]
EXPOSE 5000

ADD build.sh /app/build.sh

RUN DOCKER_BUILD=true /app/build.sh

ADD . /app

ENV DEIS_RELEASE 1.8.0-dev
