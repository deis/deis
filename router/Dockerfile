FROM deis/base:latest
MAINTAINER OpDemand <info@opdemand.com>

RUN apt-get update && \
    apt-get install -yq patch libpcre3 libpcre3-dev libssl-dev libgeoip-dev

RUN wget -q http://nginx.org/download/nginx-1.6.2.tar.gz -O /tmp/nginx-1.6.2.tar.gz
RUN wget -q https://github.com/yaoweibin/nginx_tcp_proxy_module/archive/v0.4.5.tar.gz -O /tmp/tcp-proxy.tar.gz

WORKDIR /tmp
RUN tar -xzf nginx-1.6.2.tar.gz
RUN tar -xzf tcp-proxy.tar.gz

WORKDIR /tmp/nginx-1.6.2
RUN patch -p1 < /tmp/nginx_tcp_proxy_module-0.4.5/tcp.patch
RUN ./configure --prefix=/var/lib/nginx --sbin-path=/usr/sbin/nginx --conf-path=/etc/nginx/nginx.conf --http-log-path=/var/log/nginx/access.log --error-log-path=/var/log/nginx/error.log --lock-path=/var/lock/nginx.lock --pid-path=/run/nginx.pid --http-client-body-temp-path=/var/lib/nginx/body --http-fastcgi-temp-path=/var/lib/nginx/fastcgi --http-proxy-temp-path=/var/lib/nginx/proxy --http-scgi-temp-path=/var/lib/nginx/scgi --http-uwsgi-temp-path=/var/lib/nginx/uwsgi --with-debug --with-pcre-jit --with-ipv6 --with-http_ssl_module --with-http_stub_status_module --with-http_realip_module --with-http_auth_request_module --with-http_addition_module --with-http_dav_module --with-http_geoip_module --with-http_gzip_static_module --with-http_spdy_module --with-http_sub_module --with-mail --with-mail_ssl_module --add-module=/tmp/nginx_tcp_proxy_module-0.4.5
RUN make && make install

WORKDIR /app
EXPOSE 80 2222
CMD ["/app/bin/boot"]
ADD . /app
