FROM deis/base:latest
MAINTAINER OpDemand <info@opdemand.com>

# install required packages (copied from dotcloud/docker-registry Dockerfile)
RUN sed -i 's/main$/main universe/' /etc/apt/sources.list
RUN apt-get update && apt-get install -y git-core build-essential python-dev \
    libevent-dev python-openssl liblzma-dev wget

# install recent pip
RUN wget -qO- https://raw.githubusercontent.com/pypa/pip/1.5.5/contrib/get-pip.py | python -

# create a registry user
RUN useradd -s /bin/bash registry

# add the docker registry source from github
RUN git clone https://github.com/deis/docker-registry /docker-registry && \
    cd /docker-registry && \
    git checkout 3cf76a4 && \
    chown -R registry:registry /docker-registry

# install boto configuration
RUN cp /docker-registry/config/boto.cfg /etc/boto.cfg
RUN cd /docker-registry && pip install -r requirements/main.txt

# Install core
RUN pip install /docker-registry/depends/docker-registry-core

# Install registry
RUN pip install file:///docker-registry#egg=docker-registry[bugsnag]

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
