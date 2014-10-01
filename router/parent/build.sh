#!/usr/bin/env bash

cd /tmp

export VERSION_NGINX=nginx-1.6.2
export VERSION_TCP_PROXY=v0.4.5

# URLs to the source directories
export SOURCE_NGINX=http://nginx.org/download
export SOURCE_TCP_PROXY=https://github.com/yaoweibin/nginx_tcp_proxy_module/archive

export BPATH=`pwd`/build
export PREFIX=/deis
export STATICLIBSSL="$PREFIX"

rm -rf $PREFIX
mkdir $PREFIX

# install required packages to build
apt-get update \
  && apt-get install -y patch curl wget build-essential \
  libpcre3 libpcre3-dev libssl-dev libgeoip-dev zlib1g-dev

# where the installers are
mkdir $BPATH

# grab the source files
wget -P $BPATH ${SOURCE_NGINX}/${VERSION_NGINX}.tar.gz
wget -P $BPATH ${SOURCE_TCP_PROXY}/${VERSION_TCP_PROXY}.tar.gz


# expand the source files
cd $BPATH
tar xzf $VERSION_NGINX.tar.gz
tar xzf $VERSION_TCP_PROXY.tar.gz

# build nginx
cd $BPATH/$VERSION_NGINX

patch -p1 < /$BPATH/nginx_tcp_proxy_module-0.4.5/tcp.patch

./configure \
  --prefix=/nginx \
  --pid-path=/run/nginx.pid \
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
  --add-module=$BPATH/nginx_tcp_proxy_module-0.4.5 \
  && make && make install
