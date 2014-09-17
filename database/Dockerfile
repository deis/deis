FROM deis/common:0.12.0-dev

# install 9.3 from postgresql.org repository
RUN echo "deb http://apt.postgresql.org/pub/repos/apt/ trusty-pgdg main" > /etc/apt/sources.list.d/pgdg.list
RUN curl -sk https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
RUN apt-get update \
    && apt-get install -yq postgresql-9.3 \
    && /etc/init.d/postgresql stop \
    && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# define the execution environment
VOLUME ["/var/lib/postgresql"]
WORKDIR /app
CMD ["/app/bin/boot"]
EXPOSE 5432
ADD . /app
