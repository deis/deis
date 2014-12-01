#!/usr/bin/env bash

# fail on any command exiting non-zero
set -eo pipefail

if [[ -z $DOCKER_BUILD ]]; then
  echo
  echo "Note: this script is intended for use by the Dockerfile and not as a way to build the router locally"
  echo
  exit 1
fi

export VERSION_NGINX=nginx-1.6.2
export VERSION_TCP_PROXY=0.4.5
export VERSION_NAXSI=0d53a64ed856e694fcb4038748c8cf6d5551a603

export BUILD_PATH=/tmp/build

# nginx installation directory
export PREFIX=/opt/nginx

rm -rf $PREFIX
mkdir $PREFIX

mkdir $BUILD_PATH
cd $BUILD_PATH

# install required packages to build
apt-get update \
  && apt-get install -y patch curl build-essential \
  libpcre3 libpcre3-dev libssl-dev libgeoip-dev zlib1g-dev

# grab the source files
curl -sSL http://nginx.org/download/$VERSION_NGINX.tar.gz -o $BUILD_PATH/$VERSION_NGINX.tar.gz
curl -sSL https://github.com/yaoweibin/nginx_tcp_proxy_module/archive/v$VERSION_TCP_PROXY.tar.gz -o $BUILD_PATH/$VERSION_TCP_PROXY.tar.gz
curl -sSL https://github.com/nbs-system/naxsi/archive/$VERSION_NAXSI.tar.gz -o $BUILD_PATH/$VERSION_NAXSI.tar.gz

# expand the source files
tar xzf $VERSION_NGINX.tar.gz
tar xzf $VERSION_TCP_PROXY.tar.gz
tar xzf $VERSION_NAXSI.tar.gz

# build nginx
cd $BUILD_PATH/$VERSION_NGINX

patch -p1 < $BUILD_PATH/nginx_tcp_proxy_module-$VERSION_TCP_PROXY/tcp.patch

./configure \
  --prefix=$PREFIX \
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
  --add-module=$BUILD_PATH/nginx_tcp_proxy_module-$VERSION_TCP_PROXY \
  --add-module=$BUILD_PATH/naxsi-$VERSION_NAXSI/naxsi_src \
  && make && make install
  
mv /tmp/firewall /opt/nginx/firewall
