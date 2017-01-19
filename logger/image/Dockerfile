FROM alpine:3.4

ENTRYPOINT ["/bin/logger"]
CMD ["--enable-publish"]
EXPOSE 514
EXPOSE 8088

ADD . /

ENV DEIS_RELEASE 1.13.4
