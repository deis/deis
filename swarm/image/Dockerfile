FROM alpine:3.1

COPY bin /app/bin/
WORKDIR /app/bin

ENTRYPOINT ["/app/bin/swarm"]
CMD ["--help"]
ENV SWARM_HOST :2375
EXPOSE 2375
VOLUME $HOME/.swarm

ENV DEIS_RELEASE 1.8.0-dev
