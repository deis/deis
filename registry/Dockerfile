FROM deis/base:latest
MAINTAINER Gabriel Monroy <gabriel@opdemand.com>

# install required packages (copied from dotcloud/docker-registry Dockerfile)
RUN sed -i 's/main$/main universe/' /etc/apt/sources.list 
RUN apt-get update
RUN apt-get install -y git-core build-essential python-dev \
    libevent-dev python-openssl liblzma-dev wget

# install recent pip
RUN wget -qO- https://raw.githubusercontent.com/pypa/pip/1.5.5/contrib/get-pip.py | python -

# create a registry user
RUN useradd -s /bin/bash registry

# add the docker registry source from github
RUN git clone https://github.com/dotcloud/docker-registry /docker-registry

# lock the registry version to a tag
RUN cd /docker-registry && git checkout 0.6.8

# install boto configuration
RUN cp /docker-registry/config/boto.cfg /etc/boto.cfg
RUN cd /docker-registry && pip install -r requirements.txt

# install latest etcdctl including no-sync options
RUN wget -q https://s3-us-west-2.amazonaws.com/deis/etcdctl.no-sync -O /usr/local/bin/etcdctl
RUN chmod +x /usr/local/bin/etcdctl

# create data volume
RUN mkdir -p /data/repositories && chown -R registry:registry /data
VOLUME /data

# add the current build context to /app
ADD . /app

# define the execution environment
WORKDIR /app
CMD ["/app/bin/boot"]
EXPOSE 5000
