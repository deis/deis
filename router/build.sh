#!/usr/bin/env bash

# fail on any command exiting non-zero
set -eo pipefail

if [[ -z $DOCKER_BUILD ]]; then
  echo
  echo "Note: this script is intended for use by the Dockerfile and not as a way to build the router locally"
  echo
  exit 1
fi

function get_src {
  hash="$1"
  url="$2"
  f=$(basename "$url")

  curl -sSL "$url" -o "$f"
  echo "$hash  $f" | sha256sum -c - || exit 10
  tar xzf "$f"
  rm "$f"
}

export VERSION_NGINX=nginx-1.9.2
export VERSION_NAXSI=0d53a64ed856e694fcb4038748c8cf6d5551a603
export VERSION_NDK=0.2.19
export VERSION_SETMISC=0.29

export BUILD_PATH=/tmp/build

# nginx installation directory
export PREFIX=/opt/nginx

rm -rf "$PREFIX"
mkdir "$PREFIX"

mkdir "$BUILD_PATH"
cd "$BUILD_PATH"

# install required packages to build
apk add --update-cache \
  build-base \
  curl \
  geoip-dev \
  libcrypto1.0 \
  libpcre32 \
  patch \
  pcre-dev \
  openssl-dev \
  zlib \
  zlib-dev

# download, verify and extract the source files
get_src 80b6425be14a005c8cb15115f3c775f4bc06bf798aa1affaee84ed9cf641ed78 \
        "http://nginx.org/download/$VERSION_NGINX.tar.gz"

get_src 128b56873eedbd3f240dc0f88a8b260d791321db92f14ba2fc5c49fc5307e04d \
        "https://github.com/nbs-system/naxsi/archive/$VERSION_NAXSI.tar.gz"

get_src 501f299abdb81b992a980bda182e5de5a4b2b3e275fbf72ee34dd7ae84c4b679 \
        "https://github.com/simpl/ngx_devel_kit/archive/v$VERSION_NDK.tar.gz"

get_src 8d280fc083420afb41dbe10df9a8ceec98f1d391bd2caa42ebae67d5bc9295d8 \
        "https://github.com/openresty/set-misc-nginx-module/archive/v$VERSION_SETMISC.tar.gz"

# build nginx
cd "$BUILD_PATH/$VERSION_NGINX"

./configure \
  --prefix="$PREFIX" \
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
  --with-stream \
  --add-module="$BUILD_PATH/naxsi-$VERSION_NAXSI/naxsi_src" \
  --add-module="$BUILD_PATH/ngx_devel_kit-$VERSION_NDK" \
  --add-module="$BUILD_PATH/set-misc-nginx-module-$VERSION_SETMISC" \
  && make && make install

mv /tmp/firewall /opt/nginx/firewall
