FROM alpine:3.2

# install common packages
RUN apk add --update-cache curl bash sudo && rm -rf /var/cache/apk/*

WORKDIR /app
EXPOSE 8888
CMD ["/app/bin/boot"]
ADD bin/boot /app/bin/boot

ADD build.sh /tmp/build.sh
RUN DOCKER_BUILD=true /tmp/build.sh
