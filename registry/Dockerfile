FROM ubuntu:14.04

ENV DEBIAN_FRONTEND noninteractive

# install common packages
RUN apt-get update && apt-get install -y curl

# install etcdctl
RUN curl -sSL -o /usr/local/bin/etcdctl https://s3-us-west-2.amazonaws.com/opdemand/etcdctl-v0.4.6 \
    && chmod +x /usr/local/bin/etcdctl

# install confd
RUN curl -sSL -o /usr/local/bin/confd https://s3-us-west-2.amazonaws.com/opdemand/confd-v0.5.0-json \
    && chmod +x /usr/local/bin/confd

# install required packages (copied from dotcloud/docker-registry Dockerfile)
RUN sed -i 's/main$/main universe/' /etc/apt/sources.list
RUN apt-get update && apt-get install -y git-core build-essential python-dev \
    libevent-dev python-openssl liblzma-dev

# install pip
RUN curl -sSL https://raw.githubusercontent.com/pypa/pip/1.5.6/contrib/get-pip.py | python -

# create a registry user
RUN useradd -s /bin/bash registry

# add the docker registry source from github
RUN git clone https://github.com/deis/docker-registry /docker-registry && \
    cd /docker-registry && \
    git checkout 94e7707 && \
    chown -R registry:registry /docker-registry

# install boto configuration
RUN cp /docker-registry/config/boto.cfg /etc/boto.cfg
RUN cd /docker-registry && pip install -r requirements/main.txt

# Install core
RUN pip install /docker-registry/depends/docker-registry-core

# Install registry
RUN pip install file:///docker-registry#egg=docker-registry[bugsnag,newrelic,cors]

ENV DOCKER_REGISTRY_CONFIG /docker-registry/config/config.yml
ENV SETTINGS_FLAVOR deis

# define the execution environment
WORKDIR /app
CMD ["/app/bin/boot"]
EXPOSE 5000
ADD . /app
