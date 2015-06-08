FROM alpine:3.1

ENV DEBIAN_FRONTEND noninteractive

# install common packages
RUN apk add --update-cache curl bash sudo && rm -rf /var/cache/apk/*

# install confd
RUN curl -sSL -o /usr/local/bin/confd https://s3-us-west-2.amazonaws.com/opdemand/confd-git-73f7489 \
  && chmod +x /usr/local/bin/confd

# install common packages
RUN apk add --update-cache \
    geoip \
    libssl1.0 \
    openssl \
    pcre \
    && rm -rf /var/cache/apk/*

# add nginx user
RUN addgroup -S nginx && \
  adduser -S -G nginx -H -h /app/errors -s /sbin/nologin -D nginx

WORKDIR /app

EXPOSE 80 2222

CMD ["/app/bin/boot"]

ADD . /app

ADD nginx.tgz /opt/nginx

RUN rm nginx.tgz

ENV DEIS_RELEASE 1.8.0-dev
