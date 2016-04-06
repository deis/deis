FROM alpine:3.3

# install curl in the image so it is possible to get the runtime
# profiling information without any additional package installation.
RUN apk add --no-cache curl

ADD bin/publisher /usr/local/bin/publisher
ENTRYPOINT ["/usr/local/bin/publisher"]

ENV DEIS_RELEASE 1.14.0-dev
