FROM deis/slugrunner
RUN mkdir -p /app
WORKDIR /app
ENTRYPOINT ["/runner/init"]
ADD slug.tgz /app
