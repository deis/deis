FROM gliderlabs/alpine:3.1

# install common packages
RUN apk-install curl bash sudo

# install confd
RUN curl -sSL -o /usr/local/bin/confd https://s3-us-west-2.amazonaws.com/opdemand/confd-git-b8e693c \
    && chmod +x /usr/local/bin/confd

# install common packages
RUN apk-install geoip openssl libssl1.0 pcre

# add nginx user
RUN addgroup -S nginx && \
  adduser -S -G nginx -H -h /app/errors -s /sbin/nologin -D nginx

WORKDIR /app

EXPOSE 80 2222

CMD ["/app/bin/boot"]

ADD . /app

ADD nginx.tgz /opt/nginx

RUN rm nginx.tgz

ENV DEIS_RELEASE 1.5.0-dev

