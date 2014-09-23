FROM deis/base:latest
MAINTAINER OpDemand <info@opdemand.com>

# install postgresql 9.3 from postgresql.org repository as well as requirements for building wal-e
RUN echo "deb http://apt.postgresql.org/pub/repos/apt/ precise-pgdg main" > /etc/apt/sources.list.d/pgdg.list
RUN wget --quiet --no-check-certificate -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
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
RUN pip install .
RUN mkdir -p /etc/wal-e.d/env
ADD WALE_S3_ENDPOINT /etc/wal-e.d/env/WALE_S3_ENDPOINT
RUN chown -R root:postgres /etc/wal-e.d

# Use modified confd with a fix for /etc/hosts - see https://github.com/kelseyhightower/confd/pull/123
ADD https://s3-us-west-2.amazonaws.com/opdemand/confd-git-0e563e5 /usr/local/bin/confd
RUN chmod +x /usr/local/bin/confd

# define the execution environment
VOLUME ["/var/lib/postgresql"]
WORKDIR /app
CMD ["/app/bin/boot"]
EXPOSE 5432
ADD . /app
