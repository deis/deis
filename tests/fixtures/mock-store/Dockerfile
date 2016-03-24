FROM alpine:3.3

# install common packages
RUN apk add --no-cache curl bash sudo

WORKDIR /app
EXPOSE 8888
CMD ["/app/bin/boot"]
ADD bin/boot /app/bin/boot

ADD mock-s3-patch.diff /tmp/mock-s3-patch.diff
ADD build.sh /tmp/build.sh
RUN DOCKER_BUILD=true /tmp/build.sh
