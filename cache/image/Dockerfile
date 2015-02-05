FROM ubuntu-debootstrap:14.04

ADD . /app
WORKDIR /app

CMD ["/app/bin/boot"]
EXPOSE 6379

ENV DEIS_RELEASE 1.4.0-dev
