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

# install postgresql 9.3 from postgresql.org repository as well as requirements for building wal-e
RUN echo "deb http://apt.postgresql.org/pub/repos/apt/ trusty-pgdg main" > /etc/apt/sources.list.d/pgdg.list
RUN curl -sk https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
RUN apt-get update && apt-get install -yq \
                                      curl \
                                      daemontools \
                                      file \
                                      gcc \
                                      git \
                                      libxml2-dev \
                                      libxslt1-dev \
                                      lzop \
                                      postgresql-9.3 \
                                      pv \
                                      python-dev

RUN /etc/init.d/postgresql stop

# install pip
RUN curl -sSL https://raw.githubusercontent.com/pypa/pip/1.5.6/contrib/get-pip.py | python -

# install wal-e
WORKDIR /tmp
RUN git clone https://github.com/wal-e/wal-e.git
WORKDIR /tmp/wal-e
RUN git checkout c6f0d02
RUN pip install .
RUN mkdir -p /etc/wal-e.d/env
RUN chown -R root:postgres /etc/wal-e.d

# define the execution environment
WORKDIR /app
CMD ["/app/bin/boot"]
EXPOSE 5432
ADD . /app
