FROM deis/slugrunner
RUN mkdir -p /app
WORKDIR /app
ENTRYPOINT ["/runner/init"]
USER slug
ADD slug.tgz /app
