FROM alpine:3.1

# install curl in the image so it is possible to get the runtime
# profiling information without any additional package installation.
RUN apk add --update-cache curl && rm -rf /var/cache/apk/*

ADD bin/publisher /usr/local/bin/publisher
ENTRYPOINT ["/usr/local/bin/publisher"]

ENV DEIS_RELEASE 1.8.0-dev
