FROM alpine:3.1

ENV JAVA_HOME /jre
ENV PATH ${PATH}:${JAVA_HOME}/bin

EXPOSE 2181 2888 3888

ADD . /app

RUN DOCKER_BUILD=true /app/build.sh

WORKDIR /app

VOLUME ["/opt/zookeeper-data"]

ENTRYPOINT ["/app/bin/boot"]

