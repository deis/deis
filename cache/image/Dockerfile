FROM alpine:3.2

RUN apk add --update-cache redis && rm -rf /var/cache/apk/*

ADD . /app
WORKDIR /app

CMD ["/app/bin/boot"]
EXPOSE 6379

ENV DEIS_RELEASE 1.7.0-dev
