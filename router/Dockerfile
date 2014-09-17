FROM deis/common:0.12.0-dev

# install build and runtime packages
RUN apt-get update && apt-get install -y \
        gcc make patch \
        libpcre3 libpcre3-dev libssl-dev libgeoip-dev \
        geoip-bin \
        --no-install-recommends

# build nginx+tcp_proxy from source
RUN curl -sSL http://nginx.org/download/nginx-1.6.1.tar.gz | tar -C /tmp -xz
RUN curl -sSL https://github.com/yaoweibin/nginx_tcp_proxy_module/archive/v0.4.5.tar.gz | tar -C /tmp -xz
RUN cd /tmp/nginx-1.6.1 \
    && patch -p1 < /tmp/nginx_tcp_proxy_module-0.4.5/tcp.patch \
    && ./configure \
        --prefix=/var/lib/nginx \
        --sbin-path=/usr/sbin/nginx \
        --conf-path=/etc/nginx/nginx.conf \
        --http-log-path=/var/log/nginx/access.log \
        --error-log-path=/var/log/nginx/error.log \
        --lock-path=/var/lock/nginx.lock \
        --pid-path=/run/nginx.pid \
        --http-client-body-temp-path=/var/lib/nginx/body \
        --http-fastcgi-temp-path=/var/lib/nginx/fastcgi \
        --http-proxy-temp-path=/var/lib/nginx/proxy \
        --http-scgi-temp-path=/var/lib/nginx/scgi \
        --http-uwsgi-temp-path=/var/lib/nginx/uwsgi \
        --with-debug \
        --with-pcre-jit \
        --with-ipv6 \
        --with-http_ssl_module \
        --with-http_stub_status_module \
        --with-http_realip_module \
        --with-http_auth_request_module \
        --with-http_addition_module \
        --with-http_dav_module \
        --with-http_geoip_module \
        --with-http_gzip_static_module \
        --with-http_spdy_module \
        --with-http_sub_module \
        --with-mail \
        --with-mail_ssl_module \
        --add-module=/tmp/nginx_tcp_proxy_module-0.4.5 \
    && make && make install

WORKDIR /app
EXPOSE 80 2222
CMD ["/app/bin/boot"]
ADD . /app
