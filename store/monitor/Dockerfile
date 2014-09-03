FROM deis/store-base:latest
MAINTAINER OpDemand <info@opdemand.com>

VOLUME ["/etc/ceph", "/var/lib/ceph"]
EXPOSE 6789
ENTRYPOINT ["/app/bin/boot"]
ADD bin/boot /app/bin/boot
