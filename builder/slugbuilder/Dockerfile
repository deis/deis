FROM progrium/cedarish:cedar
MAINTAINER OpDemand <info@opdemand.com>

RUN useradd slugbuilder --home-dir /app

ADD ./builder/ /tmp/builder
RUN /tmp/builder/install-buildpacks
ENTRYPOINT ["/tmp/builder/build.sh"]
