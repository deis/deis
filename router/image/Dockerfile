FROM ubuntu-debootstrap:14.04

ENV DEBIAN_FRONTEND noninteractive

# install common packages
RUN apt-get update && apt-get install -y curl net-tools sudo

# install confd from https://github.com/deis/confd/tree/deis
RUN curl -sSL -o /usr/local/bin/confd https://s3-us-west-2.amazonaws.com/opdemand/confd-git-4c50136 \
    && chmod +x /usr/local/bin/confd

# install common packages
RUN apt-get update && apt-get install -y libgeoip1 cron

WORKDIR /app

EXPOSE 80 2222

CMD ["/app/bin/boot"]

ADD . /app

ADD nginx.tgz /opt/nginx

RUN rm nginx.tgz

ENV DEIS_RELEASE 1.5.2

