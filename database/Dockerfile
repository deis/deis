FROM deis/base:latest
MAINTAINER Gabriel Monroy <gabriel@opdemand.com>

# install 9.3 from postgresql.org repository
RUN echo "deb http://apt.postgresql.org/pub/repos/apt/ precise-pgdg main" > /etc/apt/sources.list.d/pgdg.list
RUN wget --quiet --no-check-certificate -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
RUN apt-get update && apt-get install -yq libxslt-dev libxml2-dev libpam-dev libedit-dev openssh-client rsync postgresql-9.3 postgresql-server-dev-9.3

# install repmgr 2
ADD https://github.com/2ndQuadrant/repmgr/archive/v2.0.tar.gz /tmp/repmgr-v2.0.tar.gz
WORKDIR /tmp
RUN tar -xzf repmgr-v2.0.tar.gz
WORKDIR /tmp/repmgr-2.0
RUN make USE_PGXS=1
RUN make USE_PGXS=1 install
RUN update-alternatives --install /usr/bin/repmgr repmgr /usr/lib/postgresql/9.3/bin/repmgr 10
RUN update-alternatives --install /usr/bin/repmgrd repmgrd /usr/lib/postgresql/9.3/bin/repmgrd 10

RUN mkdir -p /home/postgres/.ssh
RUN chown -R postgres:postgres /home/postgres/.ssh
RUN chmod 0700 /home/postgres/.ssh

# add the current build context to /app
ADD . /app

# define the execution environment
WORKDIR /app
CMD ["/app/bin/boot"]
EXPOSE 5432
