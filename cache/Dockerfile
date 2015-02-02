FROM ubuntu-debootstrap:14.04

ADD . /app

EXPOSE 6379

WORKDIR /app
CMD ["/app/bin/boot"]

ENV DEIS_RELEASE 1.4.0-dev

