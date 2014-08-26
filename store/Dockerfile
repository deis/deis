FROM deis/base:latest
MAINTAINER OpDemand <info@opdemand.com>

RUN useradd ceph --groups sudo --shell /bin/bash --home /data/ceph

# We just install ceph-deploy here -- actual setup is handled in /app/bin/boot
RUN wget -q -O- 'https://ceph.com/git/?p=ceph.git;a=blob_plain;f=keys/release.asc' | apt-key add -
RUN echo "deb http://ceph.com/debian-firefly precise main" > /etc/apt/sources.list.d/ceph.list
RUN DEBIAN_FRONTEND=noninteractive apt-get update && apt-get install -yq openssh-server openssh-client keychain ceph-deploy

WORKDIR /app
CMD ["/app/bin/boot"]
EXPOSE 6789 6800 6801 6802 6803 6804 6805 6806 6807 6808 6809 6810
ADD . /app
