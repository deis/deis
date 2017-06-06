FROM alpine:3.4

# install common packages
RUN apk add --no-cache \
	bash \
	curl \
	geoip \
	libssl1.0 \
	openssl \
	pcre \
	sudo

# install confd
RUN curl -sSL -o /usr/local/bin/confd https://github.com/kelseyhightower/confd/releases/download/v0.11.0/confd-0.11.0-linux-amd64 \
	&& chmod +x /usr/local/bin/confd

# add nginx user
RUN addgroup -S nginx && \
  adduser -S -G nginx -H -h /opt/nginx -s /sbin/nologin -D nginx

COPY rootfs /

# compile nginx from source
RUN build

CMD ["boot"]
EXPOSE 80 2222 9090

ENV DEIS_RELEASE 1.13.4
