FROM alpine:3.1

ENTRYPOINT ["/bin/logger"]
CMD ["--enable-publish"]
EXPOSE 514

ADD . /

ENV DEIS_RELEASE 1.8.0-dev
