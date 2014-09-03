FROM deis/base:latest
MAINTAINER OpDemand <info@opdemand.com>

RUN wget -q -O- 'https://ceph.com/git/?p=ceph.git;a=blob_plain;f=keys/release.asc' | apt-key add -
RUN echo "deb http://ceph.com/debian-firefly precise main" > /etc/apt/sources.list.d/ceph.list
RUN DEBIAN_FRONTEND=noninteractive apt-get update && apt-get install -yq ceph

# Use modified confd with a fix for /etc/hosts - see https://github.com/kelseyhightower/confd/pull/123
ADD https://s3-us-west-2.amazonaws.com/opdemand/confd-git-0e563e5 /usr/local/bin/confd
RUN chmod +x /usr/local/bin/confd

# Add shared confd configuration
ADD . /app
