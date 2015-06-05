FROM alpine:3.2

ENTRYPOINT ["/bin/logger"]
CMD ["--enable-publish"]
EXPOSE 514

ADD . /

ENV DEIS_RELEASE 1.7.0-dev
