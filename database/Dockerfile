FROM gliderlabs/alpine:3.1

# install common packages
RUN apk-install curl bash sudo

# install etcdctl
RUN curl -sSL -o /usr/local/bin/etcdctl https://s3-us-west-2.amazonaws.com/opdemand/etcdctl-v0.4.6 \
    && chmod +x /usr/local/bin/etcdctl

# install confd
RUN curl -sSL -o /usr/local/bin/confd https://s3-us-west-2.amazonaws.com/opdemand/confd-git-b8e693c \
    && chmod +x /usr/local/bin/confd

ADD build.sh /tmp/build.sh
RUN DOCKER_BUILD=true /tmp/build.sh

# define the execution environment
WORKDIR /app
CMD ["/app/bin/boot"]
EXPOSE 5432
ADD . /app

ENV DEIS_RELEASE 1.5.0-dev

