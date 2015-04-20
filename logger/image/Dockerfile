FROM ubuntu-debootstrap:14.04

ENTRYPOINT ["/bin/logger"]
CMD ["--enable-publish"]
EXPOSE 514

ADD . /

ENV DEIS_RELEASE 1.6.0-dev

