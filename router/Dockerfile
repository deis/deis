FROM        deis/base:latest
MAINTAINER  Matthew Fisher <matthewf@opdemand.com>

# Install nginx dependencies
RUN apt-get update
RUN apt-get install -qy build-essential \
                        autoconf \
                        libssl-dev \
                        curl \
                        git \
                        libcurl4-gnutls-dev \
                        zlib1g \
                        zlib1g-dev \
                        libxml2 \
                        libxml2-dev \
                        libxslt-dev \
                        libreadline6-dev \
                        libpcre3 \
                        libpcre3-dev

RUN apt-get install -qy lua5.1 \
                        liblua5.1-0 \
                        liblua5.1-0-dev \
                        liblua5.1-socket2

# Make sure to symlink lua so the linker can find it
RUN ln -s /usr/lib/x86_64-linux-gnu/liblua5.1.so /usr/lib/liblua.so

# Download and extract all files
RUN curl -L http://nginx.org/download/nginx-1.4.7.tar.gz | tar zxv
RUN curl -L http://zlib.net/zlib-1.2.8.tar.gz | tar zxv
RUN curl -L http://www.openssl.org/source/openssl-1.0.1c.tar.gz | tar zxv
RUN curl -L https://github.com/simpl/ngx_devel_kit/archive/v0.2.19.tar.gz | tar zxv
RUN curl -L https://github.com/chaoslawful/lua-nginx-module/archive/v0.9.6.tar.gz | tar zxv
RUN git clone https://github.com/yaoweibin/nginx_tcp_proxy_module.git

# Patch nginx with TCP proxy support
RUN cd nginx-1.4.7 && patch -p1 < ../nginx_tcp_proxy_module/tcp.patch

# Install nginx
ENV LUA_LIB /usr/lib
ENV LUA_INC /usr/include/lua5.1
RUN cd nginx-1.4.7 && ./configure  --with-http_gzip_static_module \
                                    --with-http_ssl_module \
                                    --with-http_stub_status_module \
                                    --with-zlib=../zlib-1.2.8 \
                                    --with-openssl=../openssl-1.0.1c \
                                    --add-module=../ngx_devel_kit-0.2.19 \
                                    --add-module=../lua-nginx-module-0.9.6 \
                                    --add-module=../nginx_tcp_proxy_module
RUN cd nginx-1.4.7 && make && make install

# remove default config
RUN rm /usr/local/nginx/conf/nginx.conf

# Install redis library
RUN mkdir -p /usr/local/lib/lua/5.1
RUN cd /usr/local/lib/lua/5.1 && wget https://raw.github.com/nrk/redis-lua/version-2.0/src/redis.lua

# Add the current context to /app
ADD . /app

WORKDIR /usr/local/nginx

EXPOSE  80

CMD ["/app/bin/boot"]
