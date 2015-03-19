FROM alpine:3.1

RUN apk add --update redis && rm -rf /var/cache/apk/*

ADD . /app
WORKDIR /app

CMD ["/app/bin/boot"]
EXPOSE 6379

ENV DEIS_RELEASE 1.5.0-dev
