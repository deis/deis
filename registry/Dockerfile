FROM deis/base:latest
MAINTAINER OpDemand <info@opdemand.com>

# install required packages (copied from dotcloud/docker-registry Dockerfile)
# install recent pip
# create a registry user
# add the docker registry source from github
# install boto configuration
# Install core
# Install registry
RUN sed -i 's/main$/main universe/' /etc/apt/sources.list && \
    apt-get update && apt-get install -y git-core build-essential python-dev \
    libevent-dev python-openssl liblzma-dev wget && \
    wget -qO- https://raw.githubusercontent.com/pypa/pip/1.5.5/contrib/get-pip.py | python - && \
    useradd -s /bin/bash registry && \
    git clone https://github.com/deis/docker-registry /docker-registry && \
    cd /docker-registry && \
    git checkout 3cf76a4 && \
    chown -R registry:registry /docker-registry && \
    cp /docker-registry/config/boto.cfg /etc/boto.cfg && \
    cd /docker-registry && pip install -r requirements/main.txt && \
    pip install /docker-registry/depends/docker-registry-core && \
    pip install file:///docker-registry#egg=docker-registry[bugsnag]

ENV DOCKER_REGISTRY_CONFIG /docker-registry/config/config.yml
ENV SETTINGS_FLAVOR deis

# create data volume
RUN mkdir -p /data/repositories && chown -R registry:registry /data
ENV STORAGE_PATH /data
VOLUME /data

# define the execution environment
WORKDIR /app
CMD ["/app/bin/boot"]
EXPOSE 5000
ADD . /app
