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

# add nginx user
RUN addgroup -S nginx && \
  adduser -S -G nginx -H -h /opt/nginx -s /sbin/nologin -D nginx

COPY rootfs /

# compile nginx from source
RUN build

# Re-copy nginx.conf since the build replaces it with its own
COPY rootfs/opt/nginx/conf/nginx.conf /opt/nginx/conf/nginx.conf

RUN ln -sf /dev/stdout /opt/nginx/logs/error.log
RUN ln -sf /dev/stdout /opt/nginx/logs/access.log

CMD ["/opt/nginx/sbin/nginx"]
EXPOSE 443

ENV DEIS_RELEASE 1.13.0-dev
