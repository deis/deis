FROM deis/base:latest
MAINTAINER OpDemand <info@opdemand.com>

RUN useradd ceph --groups sudo --shell /bin/bash --home /data/ceph

# We just install ceph-deploy here -- actual setup is handled in /app/bin/boot
RUN wget -q -O- 'https://ceph.com/git/?p=ceph.git;a=blob_plain;f=keys/release.asc' | apt-key add -
RUN echo "deb http://ceph.com/debian-firefly precise main" > /etc/apt/sources.list.d/ceph.list
RUN DEBIAN_FRONTEND=noninteractive apt-get update && apt-get install -yq openssh-server openssh-client keychain ceph-deploy

# configure ssh server
RUN rm /etc/ssh/ssh_host_*
RUN dpkg-reconfigure openssh-server
RUN mkdir -p /var/run/sshd

# Change the listen port for SSHd
RUN sed -i 's/Port 22/Port 2202/g' /etc/ssh/sshd_config

WORKDIR /app
CMD ["/app/bin/boot"]
EXPOSE 2202 6789 6800 6801 6802 6803 6804 6805 6806 6807 6808 6809 6810
ADD . /app
