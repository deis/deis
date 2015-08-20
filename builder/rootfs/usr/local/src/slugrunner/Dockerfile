FROM heroku/cedar:14

RUN mkdir /app
RUN addgroup --quiet --gid 2000 slug && \
    useradd slug --uid=2000 --gid=2000 --home-dir /app --no-create-home \
        --shell /bin/bash
RUN chown 2000:2000 /app
WORKDIR /app

# add default port to expose (can be overridden)
ENV PORT 5000
EXPOSE 5000

ADD ./runner /runner
RUN chown slug:slug /runner/init
USER slug
ENV HOME /app
ENTRYPOINT ["/runner/init"]

ENV DEIS_RELEASE 1.10.0-dev

ONBUILD RUN mkdir -p /app
ONBUILD WORKDIR /app
ONBUILD ADD slug.tgz /app
