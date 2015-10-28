FROM alpine:3.2

# install common packages
RUN apk add --update-cache \
	bash \
	curl \
	geoip \
	libssl1.0 \
	openssl \
	pcre \
	sudo \
	&& rm -rf /var/cache/apk/*

# install confd
RUN curl -sSL -o /usr/local/bin/confd https://s3-us-west-2.amazonaws.com/opdemand/confd-git-73f7489 \
  && chmod +x /usr/local/bin/confd

# add nginx user
RUN addgroup -S nginx && \
  adduser -S -G nginx -H -h /opt/nginx -s /sbin/nologin -D nginx

COPY rootfs /

# compile nginx from source
RUN build

CMD ["boot"]
EXPOSE 80 2222 9090

ENV DEIS_RELEASE 1.12.0
