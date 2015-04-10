FROM ubuntu-debootstrap:14.04

ENV DEBIAN_FRONTEND noninteractive

WORKDIR /app
EXPOSE 8888
CMD ["/app/bin/boot"]
ADD bin/boot /app/bin/boot

ADD build.sh /tmp/build.sh
RUN /tmp/build.sh
