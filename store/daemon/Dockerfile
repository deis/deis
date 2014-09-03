FROM deis/store-base:latest
MAINTAINER OpDemand <info@opdemand.com>

EXPOSE 6800

ADD bin/boot /app/bin/boot
ENTRYPOINT ["/app/bin/boot"]

CMD ["ceph-osd"]
