FROM deis/base:latest
MAINTAINER Gabriel Monroy <gabriel@opdemand.com>

# install required packages (copied from dotcloud/docker-registry Dockerfile)
RUN sed -i 's/main$/main universe/' /etc/apt/sources.list
RUN apt-get update && apt-get install -y git-core build-essential python-dev \
    libevent-dev python-openssl liblzma-dev wget

# install recent pip
RUN wget -qO- https://raw.githubusercontent.com/pypa/pip/1.5.5/contrib/get-pip.py | python -

# create a registry user
RUN useradd -s /bin/bash registry

# add the docker registry source from github
RUN git clone https://github.com/deis/docker-registry /docker-registry
RUN chown -R registry:registry /docker-registry

# lock the registry version to a tag
RUN cd /docker-registry && git checkout repository-import

# install boto configuration
RUN cp /docker-registry/config/boto.cfg /etc/boto.cfg
RUN cd /docker-registry && pip install -r requirements/main.txt

# Install core
RUN pip install /docker-registry/depends/docker-registry-core

# Install registry
RUN pip install file:///docker-registry#egg=docker-registry[bugsnag]

ENV DOCKER_REGISTRY_CONFIG /docker-registry/config/config_sample.yml
ENV SETTINGS_FLAVOR dev

# create data volume
RUN mkdir -p /data/repositories && chown -R registry:registry /data
ENV STORAGE_PATH /data
VOLUME /data

# add the current build context to /app
ADD . /app

# define the execution environment
WORKDIR /app
CMD ["/app/bin/boot"]
EXPOSE 5000
